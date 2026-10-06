package session

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"
)

const (
	CookieName = "session"
)

var (
	ErrInvalidSession = errors.New("invalid session")
	ErrExpiredSession = errors.New("session expired")
)

// SessionData represents the data stored in a session
type SessionData struct {
	UserID    uint      `json:"user_id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	Role      string    `json:"role"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Manager handles session encryption and decryption
type Manager struct {
	secretKey []byte
	maxAge    int
}

// NewManager creates a new session manager
func NewManager(secretKey string, maxAge int) (*Manager, error) {
	if len(secretKey) != 32 {
		return nil, fmt.Errorf("secret key must be exactly 32 bytes, got %d", len(secretKey))
	}

	return &Manager{
		secretKey: []byte(secretKey),
		maxAge:    maxAge,
	}, nil
}

// CreateSession creates an encrypted session cookie value
func (m *Manager) CreateSession(userID uint, email, firstName, role string) (string, error) {
	data := SessionData{
		UserID:    userID,
		Email:     email,
		FirstName: firstName,
		Role:      role,
		ExpiresAt: time.Now().Add(time.Duration(m.maxAge) * time.Second),
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal session data: %w", err)
	}

	// Encrypt
	encrypted, err := m.encrypt(jsonData)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt session: %w", err)
	}

	// Encode to base64
	return base64.URLEncoding.EncodeToString(encrypted), nil
}

// ValidateSession decrypts and validates a session cookie value
func (m *Manager) ValidateSession(cookieValue string) (*SessionData, error) {
	if cookieValue == "" {
		return nil, ErrInvalidSession
	}

	// Decode from base64
	encrypted, err := base64.URLEncoding.DecodeString(cookieValue)
	if err != nil {
		return nil, ErrInvalidSession
	}

	// Decrypt
	decrypted, err := m.decrypt(encrypted)
	if err != nil {
		return nil, ErrInvalidSession
	}

	// Unmarshal JSON
	var data SessionData
	if err := json.Unmarshal(decrypted, &data); err != nil {
		return nil, ErrInvalidSession
	}

	// Check expiration
	if time.Now().After(data.ExpiresAt) {
		return nil, ErrExpiredSession
	}

	return &data, nil
}

// encrypt encrypts data using AES-256-GCM
func (m *Manager) encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(m.secretKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// decrypt decrypts data using AES-256-GCM
func (m *Manager) decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(m.secretKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// GetMaxAge returns the session max age in seconds
func (m *Manager) GetMaxAge() int {
	return m.maxAge
}
