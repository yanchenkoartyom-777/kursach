package repositories

import (
	"car-dealership-system/internal/models"
	"database/sql"
	"fmt"
	"time"
)

// PasswordResetRepository handles password reset token data operations
type PasswordResetRepository struct {
	db *sql.DB
}

// NewPasswordResetRepository creates a new password reset repository
func NewPasswordResetRepository(db *sql.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

// Create creates a new password reset token
func (r *PasswordResetRepository) Create(token *models.PasswordResetToken) error {
	query := `
		INSERT INTO password_reset_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(query, token.UserID, token.Token, token.ExpiresAt).Scan(
		&token.ID,
		&token.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create password reset token: %w", err)
	}

	return nil
}

// GetByToken retrieves a password reset token by token string
func (r *PasswordResetRepository) GetByToken(token string) (*models.PasswordResetToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, is_used, created_at
		FROM password_reset_tokens
		WHERE token = $1
	`

	var resetToken models.PasswordResetToken
	err := r.db.QueryRow(query, token).Scan(
		&resetToken.ID,
		&resetToken.UserID,
		&resetToken.Token,
		&resetToken.ExpiresAt,
		&resetToken.IsUsed,
		&resetToken.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("token not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	return &resetToken, nil
}

// MarkAsUsed marks a password reset token as used
func (r *PasswordResetRepository) MarkAsUsed(tokenID uint) error {
	query := `UPDATE password_reset_tokens SET is_used = true WHERE id = $1`

	result, err := r.db.Exec(query, tokenID)
	if err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("token not found")
	}

	return nil
}

// DeleteExpired deletes expired password reset tokens
func (r *PasswordResetRepository) DeleteExpired() error {
	query := `DELETE FROM password_reset_tokens WHERE expires_at < $1`

	_, err := r.db.Exec(query, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete expired tokens: %w", err)
	}

	return nil
}

// DeleteByUserID deletes all password reset tokens for a user
func (r *PasswordResetRepository) DeleteByUserID(userID uint) error {
	query := `DELETE FROM password_reset_tokens WHERE user_id = $1`

	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete tokens: %w", err)
	}

	return nil
}
