package services

import (
	"car-dealership-system/internal/models"
	"car-dealership-system/internal/repositories"
	"car-dealership-system/pkg/utils"
	"fmt"
	"sync"
	"time"
)

// LoginAttempt tracks login attempts for rate limiting
type LoginAttempt struct {
	Count       int
	LastAttempt time.Time
}

// AuthService handles authentication operations
type AuthService struct {
	userRepo          *repositories.UserRepository
	passwordResetRepo *repositories.PasswordResetRepository
	twoFactorService  *TwoFactorService
	emailService      *EmailService
	loginAttempts     map[string]*LoginAttempt
	mu                sync.RWMutex
	twoFactorEnabled  bool
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo *repositories.UserRepository,
	passwordResetRepo *repositories.PasswordResetRepository,
	twoFactorService *TwoFactorService,
	emailService *EmailService,
	twoFactorEnabled bool,
) *AuthService {
	service := &AuthService{
		userRepo:          userRepo,
		passwordResetRepo: passwordResetRepo,
		twoFactorService:  twoFactorService,
		emailService:      emailService,
		loginAttempts:     make(map[string]*LoginAttempt),
		twoFactorEnabled:  twoFactorEnabled,
	}

	// Start cleanup goroutine for login attempts
	go service.cleanupLoginAttempts()

	return service
}

// Register registers a new user
func (s *AuthService) Register(req *models.RegisterRequest) error {
	// Check if email already exists
	exists, err := s.userRepo.EmailExists(req.Email)
	if err != nil {
		return fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return fmt.Errorf("email already registered")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         "client",
	}

	if err := s.userRepo.Create(user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Send welcome email (non-blocking)
	go func() {
		_ = s.emailService.SendWelcomeEmail(user.Email, user.FirstName)
	}()

	return nil
}

// Login authenticates a user and sends a verification code
func (s *AuthService) Login(email, password string) (*models.LoginResponse, error) {
	// Check rate limiting
	if err := s.checkRateLimit(email); err != nil {
		return nil, err
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		s.recordFailedAttempt(email)
		return nil, fmt.Errorf("invalid email or password")
	}

	// Verify password
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		s.recordFailedAttempt(email)
		return nil, fmt.Errorf("invalid email or password")
	}

	// Reset login attempts on successful password verification
	s.resetLoginAttempts(email)

	// If 2FA is disabled, return user directly
	if !s.twoFactorEnabled {
		return &models.LoginResponse{
			Message: "Login successful",
			Email:   email,
			User:    user,
			Skip2FA: true,
		}, nil
	}

	// Generate verification code
	code, err := s.twoFactorService.GenerateCode(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification code: %w", err)
	}

	// Send verification code via email
	if err := s.emailService.SendVerificationCode(user.Email, code); err != nil {
		return nil, fmt.Errorf("failed to send verification code: %w", err)
	}

	return &models.LoginResponse{
		Message: "Verification code sent to your email",
		Email:   email,
	}, nil
}

// VerifyTwoFactor verifies the two-factor code and completes authentication
func (s *AuthService) VerifyTwoFactor(email, code string) (*models.User, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Validate code
	valid, err := s.twoFactorService.ValidateCode(user.ID, code)
	if err != nil || !valid {
		return nil, fmt.Errorf("invalid or expired verification code")
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// checkRateLimit checks if the user has exceeded login attempts
func (s *AuthService) checkRateLimit(email string) error {
	s.mu.RLock()
	attempt, exists := s.loginAttempts[email]
	s.mu.RUnlock()

	if !exists {
		return nil
	}

	// Check if 15 minutes have passed since last attempt
	if time.Since(attempt.LastAttempt) > 15*time.Minute {
		s.resetLoginAttempts(email)
		return nil
	}

	// Check if exceeded 5 attempts
	if attempt.Count >= 5 {
		return fmt.Errorf("too many login attempts, please try again in 15 minutes")
	}

	return nil
}

// recordFailedAttempt records a failed login attempt
func (s *AuthService) recordFailedAttempt(email string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	attempt, exists := s.loginAttempts[email]
	if !exists {
		s.loginAttempts[email] = &LoginAttempt{
			Count:       1,
			LastAttempt: time.Now(),
		}
		return
	}

	attempt.Count++
	attempt.LastAttempt = time.Now()
}

// resetLoginAttempts resets login attempts for an email
func (s *AuthService) resetLoginAttempts(email string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.loginAttempts, email)
}

// cleanupLoginAttempts periodically cleans up old login attempts
func (s *AuthService) cleanupLoginAttempts() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for email, attempt := range s.loginAttempts {
			if now.Sub(attempt.LastAttempt) > 15*time.Minute {
				delete(s.loginAttempts, email)
			}
		}
		s.mu.Unlock()
	}
}

// ChangePassword changes a user's password
func (s *AuthService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Verify old password
	if !utils.CheckPasswordHash(oldPassword, user.PasswordHash) {
		return fmt.Errorf("invalid current password")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update password
	if err := s.userRepo.UpdatePassword(userID, hashedPassword); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Invalidate all two-factor codes
	_ = s.twoFactorService.InvalidateUserCodes(userID)

	return nil
}

// UpdateProfile updates user profile information
func (s *AuthService) UpdateProfile(user *models.User) error {
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update profile: %w", err)
	}
	return nil
}

// RequestPasswordReset creates a password reset token and sends email
func (s *AuthService) RequestPasswordReset(email string) error {
	// Find user by email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		// Don't reveal if user exists or not for security
		return nil
	}

	// Delete any existing tokens for this user
	_ = s.passwordResetRepo.DeleteByUserID(user.ID)

	// Generate reset token
	token := utils.GenerateRandomString(32)

	// Create password reset token (valid for 1 hour)
	resetToken := &models.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		IsUsed:    false,
	}

	if err := s.passwordResetRepo.Create(resetToken); err != nil {
		return fmt.Errorf("failed to create reset token: %w", err)
	}

	// Send password reset email
	if err := s.emailService.SendPasswordResetEmail(user.Email, user.FirstName, token); err != nil {
		return fmt.Errorf("failed to send reset email: %w", err)
	}

	return nil
}

// ResetPassword resets user password using a valid token
func (s *AuthService) ResetPassword(token, newPassword string) error {
	// Get token from database
	resetToken, err := s.passwordResetRepo.GetByToken(token)
	if err != nil {
		return fmt.Errorf("invalid or expired token")
	}

	// Check if token is already used
	if resetToken.IsUsed {
		return fmt.Errorf("token has already been used")
	}

	// Check if token is expired
	if time.Now().After(resetToken.ExpiresAt) {
		return fmt.Errorf("token has expired")
	}

	// Get user
	user, err := s.userRepo.GetByID(resetToken.UserID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update user password
	user.PasswordHash = hashedPassword
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Mark token as used
	if err := s.passwordResetRepo.MarkAsUsed(resetToken.ID); err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	return nil
}
