package services

import (
	"car-dealership-system/internal/models"
	"car-dealership-system/internal/repositories"
	"car-dealership-system/pkg/utils"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

// TwoFactorService handles two-factor authentication operations
type TwoFactorService struct {
	repo           *repositories.TwoFactorRepository
	codeLength     int
	codeExpiration time.Duration
}

// NewTwoFactorService creates a new two-factor service
func NewTwoFactorService(repo *repositories.TwoFactorRepository, codeLength int, codeExpiration time.Duration) *TwoFactorService {
	return &TwoFactorService{
		repo:           repo,
		codeLength:     codeLength,
		codeExpiration: codeExpiration,
	}
}

// GenerateCode generates a new verification code for a user
func (s *TwoFactorService) GenerateCode(userID uint) (string, error) {
	// Generate random numeric code
	code, err := s.generateNumericCode(s.codeLength)
	if err != nil {
		return "", fmt.Errorf("failed to generate code: %w", err)
	}

	// Hash the code before storing
	hashedCode, err := utils.HashPassword(code)
	if err != nil {
		return "", fmt.Errorf("failed to hash code: %w", err)
	}

	// Create two-factor code record
	twoFactorCode := &models.TwoFactorCode{
		UserID:    userID,
		Code:      hashedCode,
		ExpiresAt: time.Now().Add(s.codeExpiration),
		IsUsed:    false,
	}

	if err := s.repo.Create(twoFactorCode); err != nil {
		return "", fmt.Errorf("failed to save code: %w", err)
	}

	// Return the plain code (to be sent via email)
	return code, nil
}

// ValidateCode validates a verification code for a user
func (s *TwoFactorService) ValidateCode(userID uint, code string) (bool, error) {
	// Get the latest unused code for the user
	storedCode, err := s.repo.GetByUserID(userID)
	if err != nil {
		return false, fmt.Errorf("invalid or expired code")
	}

	// Check if code has expired
	if time.Now().After(storedCode.ExpiresAt) {
		return false, fmt.Errorf("code has expired")
	}

	// Check if code has already been used
	if storedCode.IsUsed {
		return false, fmt.Errorf("code has already been used")
	}

	// Verify the code
	if !utils.CheckPasswordHash(code, storedCode.Code) {
		return false, fmt.Errorf("invalid code")
	}

	// Mark code as used
	if err := s.repo.MarkAsUsed(storedCode.ID); err != nil {
		return false, fmt.Errorf("failed to mark code as used: %w", err)
	}

	return true, nil
}

// CleanupExpiredCodes removes expired and used codes
func (s *TwoFactorService) CleanupExpiredCodes() error {
	if err := s.repo.DeleteExpired(); err != nil {
		return fmt.Errorf("failed to cleanup expired codes: %w", err)
	}
	return nil
}

// CleanupOldCodes removes codes older than specified duration
func (s *TwoFactorService) CleanupOldCodes(olderThan time.Duration) error {
	if err := s.repo.CleanupOldCodes(olderThan); err != nil {
		return fmt.Errorf("failed to cleanup old codes: %w", err)
	}
	return nil
}

// InvalidateUserCodes invalidates all codes for a user
func (s *TwoFactorService) InvalidateUserCodes(userID uint) error {
	if err := s.repo.DeleteByUserID(userID); err != nil {
		return fmt.Errorf("failed to invalidate user codes: %w", err)
	}
	return nil
}

// generateNumericCode generates a random numeric code of specified length
func (s *TwoFactorService) generateNumericCode(length int) (string, error) {
	const digits = "0123456789"
	code := make([]byte, length)

	for i := range code {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		code[i] = digits[num.Int64()]
	}

	return string(code), nil
}

// GetActiveCodeCount returns the number of active codes for a user
func (s *TwoFactorService) GetActiveCodeCount(userID uint) (int, error) {
	count, err := s.repo.CountActiveByUserID(userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get active code count: %w", err)
	}
	return count, nil
}
