package repository

import (
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"ytb-video-sharing-app-be/internal/entities"
)

type accountPasswordConfig struct {
	testConfig
	repo AccountPasswordRepository
}

func SetupAccountPasswordConfig(t *testing.T) *accountPasswordConfig {
	testConf := SetupTest(t)

	return &accountPasswordConfig{
		testConfig: *testConf,
		repo:       NewAccountPasswordRepository(testConf.db),
	}
}

func TestGetAccountPasswordByID(t *testing.T) {
	t.Run("Should return account password if found", func(t *testing.T) {
		cfg := SetupAccountPasswordConfig(t)
		defer cfg.TearDownTest()

		ctx := context.Background()
		expected := &entities.AccountPassword{ID: 1, Password: "hashed_password"}

		cfg.db.EXPECT().QueryRow(ctx, gomock.Any(), expected.ID).Return(cfg.row)
		cfg.row.EXPECT().Scan(gomock.Any(), gomock.Any()).DoAndReturn(func(args ...interface{}) error {
			*args[0].(*int64) = expected.ID
			*args[1].(*string) = expected.Password
			return nil
		})

		rs := cfg.repo.GetAccountPasswordByID(ctx, expected.ID)

		assert.NotNil(t, rs)
		assert.Equal(t, expected, rs)
	})

	t.Run("Should return nil if no row is found", func(t *testing.T) {
		cfg := SetupAccountPasswordConfig(t)
		defer cfg.TearDownTest()

		ctx := context.Background()
		cfg.db.EXPECT().QueryRow(ctx, gomock.Any(), gomock.Any()).Return(cfg.row)

		err := errors.New("sql: no rows in result set")
		cfg.row.EXPECT().Scan(gomock.Any(), gomock.Any()).Return(err)

		result := cfg.repo.GetAccountPasswordByID(ctx, 1)

		assert.Nil(t, result)
	})

	t.Run("Should return error if scan fails", func(t *testing.T) {
		cfg := SetupAccountPasswordConfig(t)
		defer cfg.TearDownTest()

		ctx := context.Background()
		cfg.db.EXPECT().QueryRow(ctx, gomock.Any(), gomock.Any()).Return(cfg.row)

		err := errors.New("scan error")
		cfg.row.EXPECT().Scan(gomock.Any(), gomock.Any()).Return(err)

		result := cfg.repo.GetAccountPasswordByID(ctx, 1)

		assert.Nil(t, result)
	})
}

func TestCreateAccountPassword(t *testing.T) {
	t.Run("Should create account password successfully", func(t *testing.T) {
		cfg := SetupAccountPasswordConfig(t)
		defer cfg.TearDownTest()

		ctx := context.Background()
		accountPassword := &entities.AccountPassword{ID: 1, Password: "hashed_password"}

		cfg.tx.EXPECT().Exec(ctx, gomock.Any(), accountPassword.ID, accountPassword.Password).Return(nil)

		err := cfg.repo.CreateAccountPassword(ctx, cfg.tx, accountPassword)
		assert.NoError(t, err)
	})

	t.Run("Should return error if database execution fails", func(t *testing.T) {
		cfg := SetupAccountPasswordConfig(t)
		defer cfg.TearDownTest()

		ctx := context.Background()
		expectedErr := errors.New("database error")
		accountPassword := &entities.AccountPassword{ID: 1, Password: "hashed_password"}

		cfg.tx.EXPECT().Exec(ctx, gomock.Any(), accountPassword.ID, accountPassword.Password).Return(expectedErr)

		err := cfg.repo.CreateAccountPassword(ctx, cfg.tx, accountPassword)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}
