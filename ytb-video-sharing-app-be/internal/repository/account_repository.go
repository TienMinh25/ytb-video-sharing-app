package repository

import (
	"context"
	"ytb-video-sharing-app-be/internal/entities"
	"ytb-video-sharing-app-be/pkg"
)

type AccountRepository interface {
	// Get account by email.
	GetAccountByEmail(ctx context.Context, email string) *entities.Account

	// Get account by email.
	GetAccountByID(ctx context.Context, id int64) *entities.Account

	// Get account by email within trasaction.
	GetAccountByEmailX(ctx context.Context, tx pkg.Tx, email string) *entities.Account

	// Create new account.
	CreateAccount(ctx context.Context, tx pkg.Tx, payload *entities.Account) error

	// Begin transaction.
	BeginTransaction(ctx context.Context) (pkg.Tx, error)
}

type accountRepository struct {
	db pkg.Database
}

func NewAccountRepository(db pkg.Database) AccountRepository {
	return &accountRepository{
		db: db,
	}
}

// CreateAccount implements AccountRepository.
func (a *accountRepository) CreateAccount(ctx context.Context, tx pkg.Tx, payload *entities.Account) error {
	query := `INSERT INTO accounts(email, fullname, avatarURL)
		VALUES(?, ?, ?)			
	`

	if err := tx.Exec(ctx, query, payload.Email, payload.FullName, payload.AvatarURL); err != nil {
		return err
	}

	return nil
}

// GetAccountByEmail implements AccountRepository.
func (a *accountRepository) GetAccountByEmail(ctx context.Context, email string) *entities.Account {
	query := "SELECT * FROM accounts WHERE email = ?"

	res := a.db.QueryRow(ctx, query, email)

	account := new(entities.Account)
	err := res.Scan(&account.ID, &account.Email, &account.FullName, &account.AvatarURL)

	if err != nil {
		return nil
	}

	return account
}

// GetAccountByEmailX implements AccountRepository.
func (a *accountRepository) GetAccountByEmailX(ctx context.Context, tx pkg.Tx, email string) *entities.Account {
	query := "SELECT * FROM accounts WHERE email = ?"

	res := tx.QueryRow(ctx, query, email)

	account := new(entities.Account)
	err := res.Scan(&account.ID, &account.Email, &account.FullName, &account.AvatarURL)

	if err != nil {
		return nil
	}

	return account
}

// GetAccountByID implements AccountRepository.
func (a *accountRepository) GetAccountByID(ctx context.Context, id int64) *entities.Account {
	query := "SELECT * FROM accounts WHERE id = ?"

	res := a.db.QueryRow(ctx, query, id)

	account := new(entities.Account)
	err := res.Scan(&account.ID, &account.Email, &account.FullName, &account.AvatarURL)

	if err != nil {
		return nil
	}

	return account
}

// BeginTransaction implements AccountRepository.
func (a *accountRepository) BeginTransaction(ctx context.Context) (pkg.Tx, error) {
	return a.db.Begin(ctx)
}
