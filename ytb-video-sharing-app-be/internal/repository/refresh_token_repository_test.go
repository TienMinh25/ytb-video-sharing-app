package repository

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
	"ytb-video-sharing-app-be/internal/entities"
	"ytb-video-sharing-app-be/mocks"
)

type testConfig struct {
	ctr  *gomock.Controller
	db   *mocks.MockDatabase
	tx   *mocks.MockTx
	row  *mocks.MockRow
	repo RefreshTokenRepository
}

func __SetupConfig(t *testing.T) *testConfig {
	ctr := gomock.NewController(t)
	db := mocks.NewMockDatabase(ctr)
	tx := mocks.NewMockTx(ctr)
	row := mocks.NewMockRow(ctr)

	repo := NewRefreshTokenRepository(db)

	return &testConfig{ctr: ctr, db: db, tx: tx, repo: repo, row: row}
}

func TestSave(t *testing.T) {
	t.Run("Should save a new token", func(t *testing.T) {
		cfg := __SetupConfig(t)
		defer cfg.ctr.Finish()

		ctx := context.Background()
		refreshToken := &entities.RefreshToken{
			AccountID: 1,
			Token:     "test_token",
			ExpiresAt: time.Now().Add(time.Hour),
		}

		cfg.tx.EXPECT().
			Exec(ctx, gomock.Any(), refreshToken.AccountID, refreshToken.Token, refreshToken.ExpiresAt).
			Return(nil)

		err := cfg.repo.Save(ctx, cfg.tx, refreshToken)

		assert.NoError(t, err)
	})

	t.Run("Should throw an error if pass wrong params", func(t *testing.T) {
		cfg := __SetupConfig(t)
		defer cfg.ctr.Finish()

		ctx := context.Background()
		err := errors.New("wrong params")

		cfg.tx.EXPECT().
			Exec(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(err)

		err = cfg.repo.Save(ctx, cfg.tx, &entities.RefreshToken{})

		assert.Error(t, err)
	})
}

func TestGetRefreshToken(t *testing.T) {
	t.Run("Shoud return token if found", func(t *testing.T) {
		cfg := __SetupConfig(t)
		defer cfg.ctr.Finish()

		ctx := context.Background()
		expectedToken := &entities.RefreshToken{
			ID:        1,
			AccountID: 1,
			Token:     "test_token",
			ExpiresAt: time.Now().Add(time.Hour),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		cfg.db.EXPECT().
			QueryRow(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
			Return(cfg.row)

		cfg.row.EXPECT().
			Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(args ...interface{}) error {
				*args[0].(*int64) = expectedToken.ID
				*args[1].(*int64) = expectedToken.AccountID
				*args[2].(*string) = expectedToken.Token
				*args[3].(*time.Time) = expectedToken.ExpiresAt
				*args[4].(*time.Time) = expectedToken.CreatedAt
				*args[5].(*time.Time) = expectedToken.UpdatedAt
				return nil
			})

		rs := cfg.repo.GetRefreshToken(ctx, expectedToken.AccountID, expectedToken.Token)

		assert.NotNil(t, rs)
		assert.Equal(t, expectedToken, rs)
	})

	t.Run("Should return nil if account_id or token is incorrect or both account_id and refresh_token are incorrect", func(t *testing.T) {
		cfg := __SetupConfig(t)

		defer cfg.ctr.Finish()
		ctx := context.Background()
		cfg.db.EXPECT().QueryRow(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
			Return(cfg.row)

		cfg.row.EXPECT().
			Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(errors.New("sql: no rows in result set"))

		rs := cfg.repo.GetRefreshToken(ctx, 1, "test_token")

		assert.Nil(t, rs)
	})

	t.Run("Should return nil if Scan encounters an error", func(t *testing.T) {
		cfg := __SetupConfig(t)
		defer cfg.ctr.Finish()

		ctx := context.Background()
		cfg.db.EXPECT().QueryRow(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(cfg.row)
		cfg.row.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("scan error"))

		result := cfg.repo.GetRefreshToken(ctx, 1, "test_token")
		assert.Nil(t, result)
	})
}

func TestDeleteRefreshToken(t *testing.T) {
	t.Run("Should delete an existing token", func(t *testing.T) {
		cfg := __SetupConfig(t)
		defer cfg.ctr.Finish()

		ctx := context.Background()
		cfg.db.EXPECT().
			Exec(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)

		err := cfg.repo.DeleteRefreshToken(ctx, 1, "test_token")

		assert.NoError(t, err)
	})

	t.Run("Should not return error if no rows are affected (account_id or token incorrect)", func(t *testing.T) {
		cfg := __SetupConfig(t)
		defer cfg.ctr.Finish()

		ctx := context.Background()
		cfg.db.EXPECT().
			Exec(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil)

		err := cfg.repo.DeleteRefreshToken(ctx, 1, "wrong_token")

		assert.NoError(t, err)
	})

	t.Run("Should return error if database execution fails", func(t *testing.T) {
		cfg := __SetupConfig(t)
		defer cfg.ctr.Finish()

		expectedErr := errors.New("database error")
		ctx := context.Background()
		cfg.db.EXPECT().
			Exec(ctx, gomock.Any(), gomock.Any(), gomock.Any()).
			Return(expectedErr)

		err := cfg.repo.DeleteRefreshToken(ctx, 1, "test_token")

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestUpdateRefreshToken(t *testing.T) {
	t.Run("Should update an existing token", func(t *testing.T) {
		cfg := __SetupConfig(t)
		defer cfg.ctr.Finish()
		ctx := context.Background()

		id := int64(1)
		newToken := "new_token"
		expiresAt := time.Now().Add(time.Hour)

		cfg.tx.EXPECT().
			Exec(ctx, gomock.Any(), expiresAt, newToken, id).
			Return(nil)

		err := cfg.repo.Update(ctx, cfg.tx, id, newToken, expiresAt)

		assert.NoError(t, err)
	})

	t.Run("Should not return error if no rows are affected (account_id or token incorrect)", func(t *testing.T) {
		cfg := __SetupConfig(t)
		defer cfg.ctr.Finish()
		ctx := context.Background()

		id := int64(999)
		newToken := "new_token"
		expiresAt := time.Now().Add(time.Hour)

		cfg.tx.EXPECT().
			Exec(ctx, gomock.Any(), expiresAt, newToken, id).
			Return(nil)

		err := cfg.repo.Update(ctx, cfg.tx, id, newToken, expiresAt)

		assert.NoError(t, err)
	})

	t.Run("Should return error if database execution fails", func(t *testing.T) {
		cfg := __SetupConfig(t)
		defer cfg.ctr.Finish()
		ctx := context.Background()

		id := int64(1)
		newToken := "new_token"
		expiresAt := time.Now().Add(time.Hour)
		expectedErr := errors.New("database execution failed")

		cfg.tx.EXPECT().
			Exec(ctx, gomock.Any(), expiresAt, newToken, id).
			Return(expectedErr)

		err := cfg.repo.Update(ctx, cfg.tx, id, newToken, expiresAt)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}
