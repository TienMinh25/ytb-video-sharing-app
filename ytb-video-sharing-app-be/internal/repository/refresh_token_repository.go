//go:generate mockgen -source=../../pkg/db.go -destination=../../mocks/db_mock.go -package=mocks
package repository

import (
	"context"
	"time"
	"ytb-video-sharing-app-be/internal/entities"
	"ytb-video-sharing-app-be/pkg"
)

type RefreshTokenRepository interface {
	// Save refresh token.
	Save(ctx context.Context, tx pkg.Tx, payload *entities.RefreshToken) error

	// GetRefreshToken Get refresh token.
	GetRefreshToken(ctx context.Context, accountID int64, refreshTokenStr string) *entities.RefreshToken

	// DeleteRefreshToken Delete refresh token.
	DeleteRefreshToken(ctx context.Context, accountID int64, refreshToken string) error

	// Update update refresh token
	Update(ctx context.Context, tx pkg.Tx, id int64, newToken string, expiresAt time.Time) error
}

type refreshTokenRepository struct {
	db pkg.Database
}

func NewRefreshTokenRepository(db pkg.Database) RefreshTokenRepository {
	return &refreshTokenRepository{
		db: db,
	}
}

// GetRefreshToken implements RefreshTokenRepository.
func (r *refreshTokenRepository) GetRefreshToken(ctx context.Context, accountID int64, refreshTokenStr string) *entities.RefreshToken {
	query := `SELECT * FROM refresh_token WHERE account_id = ? AND token = ?`

	refreshToken := new(entities.RefreshToken)
	row := r.db.QueryRow(ctx, query, accountID, refreshTokenStr)

	if err := row.Scan(&refreshToken.ID, &refreshToken.AccountID,
		&refreshToken.Token, &refreshToken.ExpiresAt,
		&refreshToken.CreatedAt, &refreshToken.UpdatedAt); err != nil {
		return nil
	}

	return refreshToken
}

// Save implements RefreshTokenRepository.
func (r *refreshTokenRepository) Save(ctx context.Context, tx pkg.Tx, payload *entities.RefreshToken) error {
	query := `INSERT INTO refresh_token (account_id, token, expires_at)
				VALUES(?, ?, ?)`

	return tx.Exec(ctx, query, payload.AccountID, payload.Token, payload.ExpiresAt)
}

// DeleteRefreshToken implements RefreshTokenRepository
func (r *refreshTokenRepository) DeleteRefreshToken(ctx context.Context, accountID int64, refreshToken string) error {
	query := "DELETE FROM refresh_token WHERE account_id = ? AND token = ?"

	return r.db.Exec(ctx, query, accountID, refreshToken)
}

// Update implements RefreshTokenRepository
func (r *refreshTokenRepository) Update(ctx context.Context, tx pkg.Tx, id int64, newToken string, expiresAt time.Time) error {
	query := `UPDATE refresh_token SET expires_at = ?, token = ? WHERE id = ?`

	return tx.Exec(ctx, query, expiresAt, newToken, id)
}
