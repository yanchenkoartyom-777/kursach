package repositories

import (
	"car-dealership-system/internal/models"
	"database/sql"
	"fmt"
	"time"
)

// TwoFactorRepository handles two-factor authentication code operations
type TwoFactorRepository struct {
	db *sql.DB
}

// NewTwoFactorRepository creates a new two-factor repository
func NewTwoFactorRepository(db *sql.DB) *TwoFactorRepository {
	return &TwoFactorRepository{db: db}
}

// Create creates a new two-factor code
func (r *TwoFactorRepository) Create(code *models.TwoFactorCode) error {
	query := `
		INSERT INTO two_factor_codes (user_id, code, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		query,
		code.UserID,
		code.Code,
		code.ExpiresAt,
	).Scan(&code.ID, &code.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create two-factor code: %w", err)
	}

	return nil
}

// GetByUserID retrieves the latest unused code for a user
func (r *TwoFactorRepository) GetByUserID(userID uint) (*models.TwoFactorCode, error) {
	query := `
		SELECT id, user_id, code, expires_at, is_used, created_at
		FROM two_factor_codes
		WHERE user_id = $1 AND is_used = false AND expires_at > CURRENT_TIMESTAMP
		ORDER BY created_at DESC
		LIMIT 1
	`

	var code models.TwoFactorCode
	err := r.db.QueryRow(query, userID).Scan(
		&code.ID,
		&code.UserID,
		&code.Code,
		&code.ExpiresAt,
		&code.IsUsed,
		&code.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no valid code found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get two-factor code: %w", err)
	}

	return &code, nil
}

// MarkAsUsed marks a code as used
func (r *TwoFactorRepository) MarkAsUsed(codeID uint) error {
	query := `
		UPDATE two_factor_codes
		SET is_used = true
		WHERE id = $1
	`

	result, err := r.db.Exec(query, codeID)
	if err != nil {
		return fmt.Errorf("failed to mark code as used: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("code not found")
	}

	return nil
}

// DeleteExpired deletes all expired codes
func (r *TwoFactorRepository) DeleteExpired() error {
	query := `
		DELETE FROM two_factor_codes
		WHERE expires_at < CURRENT_TIMESTAMP OR is_used = true
	`

	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to delete expired codes: %w", err)
	}

	return nil
}

// DeleteByUserID deletes all codes for a specific user
func (r *TwoFactorRepository) DeleteByUserID(userID uint) error {
	query := `DELETE FROM two_factor_codes WHERE user_id = $1`

	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user codes: %w", err)
	}

	return nil
}

// CountActiveByUserID counts active codes for a user
func (r *TwoFactorRepository) CountActiveByUserID(userID uint) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM two_factor_codes
		WHERE user_id = $1 AND is_used = false AND expires_at > CURRENT_TIMESTAMP
	`

	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count active codes: %w", err)
	}

	return count, nil
}

// CleanupOldCodes deletes codes older than specified duration
func (r *TwoFactorRepository) CleanupOldCodes(olderThan time.Duration) error {
	query := `
		DELETE FROM two_factor_codes
		WHERE created_at < $1
	`

	cutoffTime := time.Now().Add(-olderThan)
	_, err := r.db.Exec(query, cutoffTime)
	if err != nil {
		return fmt.Errorf("failed to cleanup old codes: %w", err)
	}

	return nil
}
