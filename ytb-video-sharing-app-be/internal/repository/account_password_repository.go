package repository

import (
	"context"
	"ytb-video-sharing-app-be/internal/entities"
	"ytb-video-sharing-app-be/pkg"
)

type AccountPasswordRepository interface {
	// Get account by email.
	GetAccountPasswordByID(ctx context.Context, accountID int64) *entities.AccountPassword

	// Create account password.
	CreateAccountPassword(ctx context.Context, tx pkg.Tx, payload *entities.AccountPassword) error
}

type accountPasswordRepository struct {
	db pkg.Database
}

func NewAccountPasswordRepository(db pkg.Database) AccountPasswordRepository {
	return &accountPasswordRepository{
		db: db,
	}
}

// CreateAccountPassword implements AccountPasswordRepository.
func (a *accountPasswordRepository) CreateAccountPassword(ctx context.Context, tx pkg.Tx, payload *entities.AccountPassword) error {
	query := `INSERT INTO account_password(id, password)
		VALUES(?, ?)`

	if err := tx.Exec(ctx, query, payload.ID, payload.Password); err != nil {
		return err
	}

	return nil
}

// GetAccountPasswordByID implements AccountPasswordRepository.
func (a *accountPasswordRepository) GetAccountPasswordByID(ctx context.Context, accountID int64) *entities.AccountPassword {
	query := `SELECT * FROM account_password WHERE id = ?`

	res := a.db.QueryRow(ctx, query, accountID)
	accountPassword := new(entities.AccountPassword)

	if err := res.Scan(&accountPassword.ID, &accountPassword.Password); err != nil {
		return nil
	}

	return accountPassword
}
