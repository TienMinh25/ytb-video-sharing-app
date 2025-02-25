package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	"ytb-video-sharing-app-be/internal/dto"
	"ytb-video-sharing-app-be/internal/entities"
	"ytb-video-sharing-app-be/internal/repository"
	"ytb-video-sharing-app-be/utils"
)

type AccountService interface {
	CreateAccount(context.Context, *entities.Account, *entities.AccountPassword) (*dto.CreateAccountResponse, *dto.ErrorResponse)

	Login(ctx context.Context, email string, password string) (*dto.LoginResponse, *dto.ErrorResponse)

	Logout(ctx context.Context, accountID int64, refreshToken string) (*dto.LogoutResponse, *dto.ErrorResponse)

	RefreshToken(ctx context.Context, accountID int64, refreshToken string) (*dto.RefreshTokenResponse, *dto.ErrorResponse)
}

type accountService struct {
	accountRepository         repository.AccountRepository
	accountPasswordRepository repository.AccountPasswordRepository
	refreshTokenRepository    repository.RefreshTokenRepository
	keyManager                *utils.KeyManager
}

func NewAccountService(acaccountRepository repository.AccountRepository,
	accountPasswordRepository repository.AccountPasswordRepository,
	keyManager *utils.KeyManager,
	refreshTokenRepository repository.RefreshTokenRepository) AccountService {
	return &accountService{
		accountRepository:         acaccountRepository,
		accountPasswordRepository: accountPasswordRepository,
		keyManager:                keyManager,
		refreshTokenRepository:    refreshTokenRepository,
	}
}

// CreateAccount implements AccountService.
func (a *accountService) CreateAccount(ctx context.Context, accountPayload *entities.Account, accountPasswordPayload *entities.AccountPassword) (*dto.CreateAccountResponse, *dto.ErrorResponse) {
	// check email
	if a.accountRepository.GetAccountByEmail(ctx, accountPayload.Email) != nil {
		return nil, &dto.ErrorResponse{Code: http.StatusBadRequest, Message: "Duplicate email, please try again!"}
	}

	// hash password
	hashedPassword, err := utils.HashPassword(accountPasswordPayload.Password)

	if err != nil {
		return nil, &dto.ErrorResponse{Code: http.StatusInternalServerError, Message: utils.INTERNAL_SERVER_ERROR}
	}

	// start transaction
	tx, err := a.accountRepository.BeginTransaction(ctx)

	if err != nil {
		return nil, &dto.ErrorResponse{Code: http.StatusInternalServerError, Message: utils.INTERNAL_SERVER_ERROR}
	}
	defer tx.Rollback(ctx)

	// create account
	if err = a.accountRepository.CreateAccount(ctx, tx, &entities.Account{
		Email:     accountPayload.Email,
		FullName:  accountPayload.FullName,
		AvatarURL: accountPayload.AvatarURL,
	}); err != nil {
		return nil, &dto.ErrorResponse{Code: http.StatusInternalServerError, Message: utils.INTERNAL_SERVER_ERROR}
	}

	// get account
	account := a.accountRepository.GetAccountByEmailX(ctx, tx, accountPayload.Email)

	// create account password
	if err = a.accountPasswordRepository.CreateAccountPassword(ctx, tx, &entities.AccountPassword{
		ID:       account.ID,
		Password: hashedPassword,
	}); err != nil {
		return nil, &dto.ErrorResponse{Code: http.StatusInternalServerError, Message: utils.INTERNAL_SERVER_ERROR}
	}

	// generate access token and refresh token
	accessToken, refreshToken, errRe := a.generateTokens(account)

	if errRe != nil {
		fmt.Println("error:", errRe)
		return nil, &dto.ErrorResponse{Code: errRe.Code, Message: errRe.Message}
	}

	// save refresh token into db
	if err = a.refreshTokenRepository.Save(ctx, tx, &entities.RefreshToken{
		AccountID: account.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Duration(a.getExpireTime("EXPIRE_TIME_REFRESH_TOKEN")) * time.Hour * 24),
	}); err != nil {
		return nil, &dto.ErrorResponse{Message: utils.INTERNAL_SERVER_ERROR, Code: http.StatusInternalServerError}
	}

	// end transaction
	if err = tx.Commit(ctx); err != nil {
		return nil, &dto.ErrorResponse{Code: http.StatusInternalServerError, Message: utils.INTERNAL_SERVER_ERROR}
	}

	return &dto.CreateAccountResponse{
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		AccountResponse: dto.AccountResponse(*account),
	}, nil
}

// Login implements AccountService.
func (a *accountService) Login(ctx context.Context, email string, password string) (*dto.LoginResponse, *dto.ErrorResponse) {
	// load account and account password
	account := a.accountRepository.GetAccountByEmail(ctx, email)

	if account == nil {
		return nil, &dto.ErrorResponse{Code: http.StatusBadRequest, Message: utils.LOGIN_FAIL}
	}

	accountPassword := a.accountPasswordRepository.GetAccountPasswordByID(ctx, account.ID)

	if accountPassword == nil {
		return nil, &dto.ErrorResponse{Code: http.StatusBadRequest, Message: utils.LOGIN_FAIL}
	}

	// check matching password
	if !utils.CheckPassword(accountPassword.Password, password) {
		return nil, &dto.ErrorResponse{Code: http.StatusBadRequest, Message: utils.LOGIN_FAIL}
	}

	// generate access token and refresh token
	accessToken, refreshToken, err := a.generateTokens(account)
	if err != nil {
		return nil, &dto.ErrorResponse{Code: err.Code, Message: err.Message}
	}

	// start transaction
	tx, errCommon := a.accountRepository.BeginTransaction(ctx)

	if errCommon != nil {
		return nil, &dto.ErrorResponse{Code: http.StatusInternalServerError, Message: utils.INTERNAL_SERVER_ERROR}
	}
	defer tx.Rollback(ctx)

	// save refresh token into db
	if errCommon = a.refreshTokenRepository.Save(ctx, tx, &entities.RefreshToken{
		AccountID: account.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Duration(a.getExpireTime("EXPIRE_TIME_REFRESH_TOKEN")) * time.Hour * 24),
	}); errCommon != nil {
		return nil, &dto.ErrorResponse{Message: utils.INTERNAL_SERVER_ERROR, Code: http.StatusInternalServerError}
	}

	// end transaction
	if errCommon = tx.Commit(ctx); errCommon != nil {
		return nil, &dto.ErrorResponse{Code: http.StatusInternalServerError, Message: utils.INTERNAL_SERVER_ERROR}
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		AccountResponse: dto.AccountResponse{
			ID:        account.ID,
			Email:     account.Email,
			FullName:  account.FullName,
			AvatarURL: account.AvatarURL,
		},
	}, nil
}

func (a *accountService) Logout(ctx context.Context, accountID int64, refreshToken string) (*dto.LogoutResponse, *dto.ErrorResponse) {
	if err := a.refreshTokenRepository.DeleteRefreshToken(ctx, accountID, refreshToken); err != nil {
		return nil, &dto.ErrorResponse{Code: http.StatusBadRequest, Message: "Refresh token is invalid"}
	}

	return &dto.LogoutResponse{}, nil
}

func (a *accountService) RefreshToken(ctx context.Context, accountID int64, refreshTokenStr string) (*dto.RefreshTokenResponse, *dto.ErrorResponse) {
	oldRefreshToken := a.refreshTokenRepository.GetRefreshToken(ctx, accountID, refreshTokenStr)

	if oldRefreshToken == nil {
		return nil, &dto.ErrorResponse{Code: http.StatusBadRequest, Message: "Refresh token is invalid"}
	}

	account := a.accountRepository.GetAccountByID(ctx, accountID)

	if account == nil {
		return nil, &dto.ErrorResponse{Code: http.StatusBadRequest, Message: "Account is not found"}
	}

	accessToken, newRefreshToken, err := a.generateTokens(account)

	if err != nil {
		return nil, &dto.ErrorResponse{Code: http.StatusInternalServerError, Message: utils.INTERNAL_SERVER_ERROR}
	}

	// start transaction
	tx, errCommon := a.accountRepository.BeginTransaction(ctx)

	if errCommon != nil {
		return nil, &dto.ErrorResponse{Code: http.StatusInternalServerError, Message: utils.INTERNAL_SERVER_ERROR}
	}
	defer tx.Rollback(ctx)

	// update refresh token db
	if errCommon = a.refreshTokenRepository.Update(ctx, tx, oldRefreshToken.ID, newRefreshToken,
		time.Now().Add(time.Duration(a.getExpireTime("EXPIRE_TIME_REFRESH_TOKEN"))*time.Hour*24)); errCommon != nil {
		return nil, &dto.ErrorResponse{Message: utils.INTERNAL_SERVER_ERROR, Code: http.StatusInternalServerError}
	}

	// end transaction
	if errCommon = tx.Commit(ctx); errCommon != nil {
		return nil, &dto.ErrorResponse{Code: http.StatusInternalServerError, Message: utils.INTERNAL_SERVER_ERROR}
	}

	return &dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (a *accountService) generateTokens(account *entities.Account) (string, string, *dto.ErrorResponse) {
	expireAccessToken := a.getExpireTime("EXPIRE_TIME_ACCESS_TOKEN")
	expireRefreshToken := a.getExpireTime("EXPIRE_TIME_REFRESH_TOKEN")
	return utils.GenerateToken(account, a.keyManager, expireAccessToken, expireRefreshToken)
}

func (a *accountService) getExpireTime(envVar string) int {
	expireTime, err := strconv.Atoi(os.Getenv(envVar))
	if err != nil {
		return 24 // default to 24 hours or minutes if not set or invalid
	}
	return expireTime
}
