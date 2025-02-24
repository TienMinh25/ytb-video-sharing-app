package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"ytb-video-sharing-app-be/internal/entities"
)

type accountConfig struct {
	testConfig
	repo AccountRepository
}

func SetupAccountConfig(t *testing.T) *accountConfig {
	testConf := SetupTest(t)

	return &accountConfig{
		testConfig: *testConf,
		repo:       NewAccountRepository(testConf.db),
	}
}

func TestCreateAccount(t *testing.T) {
	t.Run("Should create an account successfully", func(t *testing.T) {
		cfg := SetupAccountConfig(t)
		defer cfg.TearDownTest()

		ctx := context.Background()
		account := &entities.Account{
			Email:     "test@example.com",
			FullName:  "Test User",
			AvatarURL: "https://avatar.url",
		}

		cfg.tx.EXPECT().
			Exec(ctx, gomock.Any(), account.Email, account.FullName, account.AvatarURL).
			Return(nil)

		err := cfg.repo.CreateAccount(ctx, cfg.tx, account)
		assert.NoError(t, err)
	})

	t.Run("Should return error when DB execution fails", func(t *testing.T) {
		cfg := SetupAccountConfig(t)
		defer cfg.TearDownTest()

		ctx := context.Background()
		account := &entities.Account{
			Email:     "test@example.com",
			FullName:  "Test User",
			AvatarURL: "https://avatar.url",
		}

		expectedErr := errors.New("db execution failed")
		cfg.tx.EXPECT().
			Exec(ctx, gomock.Any(), account.Email, account.FullName, account.AvatarURL).
			Return(expectedErr)

		err := cfg.repo.CreateAccount(ctx, cfg.tx, account)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestGetAccountByEmail(t *testing.T) {
	t.Run("Should return account when found", func(t *testing.T) {
		cfg := SetupAccountConfig(t)
		defer cfg.TearDownTest()

		ctx := context.Background()
		expectedAccount := &entities.Account{
			ID:        1,
			Email:     "test@example.com",
			FullName:  "Test User",
			AvatarURL: "https://avatar.url",
		}

		cfg.db.EXPECT().
			QueryRow(ctx, gomock.Any(), expectedAccount.Email).
			Return(cfg.row)

		cfg.row.EXPECT().
			Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(args ...interface{}) error {
				*args[0].(*int64) = expectedAccount.ID
				*args[1].(*string) = expectedAccount.Email
				*args[2].(*string) = expectedAccount.FullName
				*args[3].(*string) = expectedAccount.AvatarURL
				return nil
			})

		account := cfg.repo.GetAccountByEmail(ctx, expectedAccount.Email)
		assert.NotNil(t, account)
		assert.Equal(t, expectedAccount, account)
	})

	t.Run("Should return nil if account not found", func(t *testing.T) {
		cfg := SetupAccountConfig(t)
		defer cfg.TearDownTest()

		ctx := context.Background()
		cfg.db.EXPECT().
			QueryRow(ctx, gomock.Any(), gomock.Any()).
			Return(cfg.row)

		cfg.row.EXPECT().
			Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(sql.ErrNoRows)

		account := cfg.repo.GetAccountByEmail(ctx, "test@example.com")
		assert.Nil(t, account)
	})

	t.Run("Should return nil if scan fails", func(t *testing.T) {
		cfg := SetupAccountConfig(t)
		defer cfg.TearDownTest()

		ctx := context.Background()
		cfg.db.EXPECT().QueryRow(ctx, gomock.Any(), gomock.Any()).Return(cfg.row)

		err := errors.New("scan error")
		cfg.row.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(err)

		result := cfg.repo.GetAccountByID(ctx, 1)

		assert.Nil(t, result)
	})
}

func TestGetAccountByEmailX(t *testing.T) {
	t.Run("Should return account when found in transaction", func(t *testing.T) {
		cfg := SetupAccountConfig(t)
		defer cfg.TearDownTest()

		ctx := context.Background()
		expectedAccount := &entities.Account{
			ID:        1,
			Email:     "test@example.com",
			FullName:  "Test User",
			AvatarURL: "https://avatar.url",
		}

		cfg.tx.EXPECT().
			QueryRow(ctx, gomock.Any(), expectedAccount.Email).
			Return(cfg.row)

		cfg.row.EXPECT().
			Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(args ...interface{}) error {
				*args[0].(*int64) = expectedAccount.ID
				*args[1].(*string) = expectedAccount.Email
				*args[2].(*string) = expectedAccount.FullName
				*args[3].(*string) = expectedAccount.AvatarURL
				return nil
			})

		account := cfg.repo.GetAccountByEmailX(ctx, cfg.tx, expectedAccount.Email)
		assert.NotNil(t, account)
		assert.Equal(t, expectedAccount, account)
	})

	t.Run("Should return nil if account not found in transaction", func(t *testing.T) {
		cfg := SetupAccountConfig(t)
		defer cfg.TearDownTest()

		ctx := context.Background()
		cfg.tx.EXPECT().
			QueryRow(ctx, gomock.Any(), gomock.Any()).
			Return(cfg.row)

		cfg.row.EXPECT().
			Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(sql.ErrNoRows)

		account := cfg.repo.GetAccountByEmailX(ctx, cfg.tx, "test@example.com")
		assert.Nil(t, account)
	})

	t.Run("Should return nil if scan fails in transaction", func(t *testing.T) {
		cfg := SetupAccountConfig(t)
		defer cfg.TearDownTest()

		ctx := context.Background()
		cfg.tx.EXPECT().QueryRow(ctx, gomock.Any(), gomock.Any()).Return(cfg.row)

		err := errors.New("scan error")
		cfg.row.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(err)

		result := cfg.repo.GetAccountByEmailX(ctx, cfg.tx, "test@example.com")

		assert.Nil(t, result)
	})
}

func TestGetAccountByID(t *testing.T) {
	t.Run("Should return account when found", func(t *testing.T) {
		cfg := SetupAccountConfig(t)
		defer cfg.TearDownTest()

		ctx := context.Background()
		expectedAccount := &entities.Account{
			ID:        1,
			Email:     "test@example.com",
			FullName:  "Test User",
			AvatarURL: "https://avatar.url",
		}

		cfg.db.EXPECT().
			QueryRow(ctx, gomock.Any(), expectedAccount.ID).
			Return(cfg.row)

		cfg.row.EXPECT().
			Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(args ...interface{}) error {
				*args[0].(*int64) = expectedAccount.ID
				*args[1].(*string) = expectedAccount.Email
				*args[2].(*string) = expectedAccount.FullName
				*args[3].(*string) = expectedAccount.AvatarURL
				return nil
			})

		account := cfg.repo.GetAccountByID(ctx, expectedAccount.ID)
		assert.NotNil(t, account)
		assert.Equal(t, expectedAccount, account)
	})

	t.Run("Should return nil if account not found", func(t *testing.T) {
		cfg := SetupAccountConfig(t)
		defer cfg.TearDownTest()

		ctx := context.Background()
		cfg.db.EXPECT().
			QueryRow(ctx, gomock.Any(), gomock.Any()).
			Return(cfg.row)

		cfg.row.EXPECT().
			Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(sql.ErrNoRows)

		account := cfg.repo.GetAccountByID(ctx, 1)
		assert.Nil(t, account)
	})

	t.Run("Should return nil if scan fails", func(t *testing.T) {
		cfg := SetupAccountConfig(t)
		defer cfg.TearDownTest()

		ctx := context.Background()
		cfg.db.EXPECT().QueryRow(ctx, gomock.Any(), gomock.Any()).Return(cfg.row)

		err := errors.New("scan error")
		cfg.row.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(err)

		result := cfg.repo.GetAccountByID(ctx, 1)

		assert.Nil(t, result)
	})
}

func TestBeginTransaction(t *testing.T) {
	t.Run("Should begin a transaction successfully", func(t *testing.T) {
		cfg := SetupAccountConfig(t)
		defer cfg.TearDownTest()

		cfg.db.EXPECT().
			Begin(gomock.Any()).
			Return(cfg.tx, nil)

		tx, err := cfg.repo.BeginTransaction(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, tx)
	})

	t.Run("Should return error when transaction fails", func(t *testing.T) {
		cfg := SetupAccountConfig(t)
		defer cfg.TearDownTest()

		expectedErr := errors.New("failed to begin transaction")
		cfg.db.EXPECT().
			Begin(gomock.Any()).
			Return(nil, expectedErr)

		tx, err := cfg.repo.BeginTransaction(context.Background())

		assert.Error(t, err)
		assert.Nil(t, tx)
		assert.Equal(t, expectedErr, err)
	})
}
