package handler

import (
	"net/http"
	"strconv"
	"strings"
	"ytb-video-sharing-app-be/internal/dto"
	"ytb-video-sharing-app-be/internal/entities"
	"ytb-video-sharing-app-be/internal/service"
	"ytb-video-sharing-app-be/internal/websock"
	"ytb-video-sharing-app-be/pkg"
	"ytb-video-sharing-app-be/utils"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	accountService service.AccountService
	router         *gin.Engine
	queue          pkg.Queue
	opts           websock.RetentionMap
}

func NewAccountHandler(accountService service.AccountService, router *gin.Engine, queue pkg.Queue, opts websock.RetentionMap) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
		router:         router,
		queue:          queue,
		opts:           opts,
	}
}

// Register godoc
//
//	@Summary		Register new account
//	@Tags			accounts
//	@Description	create new account based info request
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CreateAccountRequest	true	"Thông tin đăng ký"
//	@Success		200		{object}	dto.CreateAccountResponseDocs
//	@Failure		400		{object}	dto.ResponseError
//	@Failure		500		{object}	dto.ResponseError
//	@Router			/accounts/register [post]
func (h *AccountHandler) Register(ctx *gin.Context) {
	req, _ := ctx.Get("data")
	data := req.(dto.CreateAccountRequest)

	// call service to create new account
	res, err := h.accountService.CreateAccount(ctx, &entities.Account{
		Email:     data.Email,
		FullName:  data.FullName,
		AvatarURL: data.AvatarURL,
	}, &entities.AccountPassword{Password: data.Password})

	if err != nil {
		utils.ErrorResponse(ctx, err.Code, err.Error())
		return
	}

	newOTP := h.opts.NewOTP().Key

	response := &dto.CreateAccountResponseWithOTP{
		CreateAccountResponse: *res,
		OTP:                   newOTP,
	}

	utils.SuccessResponse(ctx, http.StatusOK, response)
}

// Login godoc
//
//	@Summary		Login account
//	@Tags			accounts
//	@Description	Authenticate user and return access token & refresh token
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.LoginRequest	true	"Login payload"
//	@Success		200		{object}	dto.LoginResponseDocs
//	@Failure		400		{object}	dto.ResponseError
//	@Failure		500		{object}	dto.ResponseError
//	@Router			/accounts/login [post]
func (h *AccountHandler) Login(ctx *gin.Context) {
	req, _ := ctx.Get("data")
	data := req.(dto.LoginRequest)

	// call service to login
	res, err := h.accountService.Login(ctx, data.Email, data.Password)

	if err != nil {
		utils.ErrorResponse(ctx, err.Code, err.Error())
		return
	}

	newOTP := h.opts.NewOTP().Key

	response := &dto.LoginResponseWithOTP{
		LoginResponse: *res,
		OTP:           newOTP,
	}

	utils.SuccessResponse(ctx, http.StatusOK, response)
}

// Logout godoc
//
//	@Summary		Logout account
//	@Tags			accounts
//	@Description	Logout user by deleting refresh token
//	@Accept			json
//	@Produce		json
//	@Param			accountID		path	int		true	"Account ID"
//	@Param			X-Authorization	header	string	true	"Refresh Token"
//
//	@Security		BearerAuth
//
//	@Success		200	{object}	dto.LogoutResponseDocs
//	@Failure		400	{object}	dto.ResponseError
//	@Router			/accounts/logout/{accountID} [post]
func (h *AccountHandler) Logout(ctx *gin.Context) {
	refreshToken := strings.Split(ctx.Request.Header.Get("X-Authorization"), " ")[1]
	accountID, err := strconv.Atoi(ctx.Param("accountID"))

	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid account id")
	}

	// delete in database
	_, errRe := h.accountService.Logout(ctx, int64(accountID), refreshToken)

	if errRe != nil {
		utils.ErrorResponse(ctx, errRe.Code, errRe.Message)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, &dto.LogoutResponse{})
}

// RefreshToken godoc
//
//	@Summary		Refresh Token
//	@Description	Refresh access token using a valid refresh token
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			accountID		path	int		true	"Account ID"
//	@Param			X-Authorization	header	string	true	"Refresh Token"
//
//	@Security		BearerAuth
//
//	@Success		200	{object}	dto.RefreshTokenResponseDocs
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/accounts/refresh-token [post]
func (h *AccountHandler) RefreshToken(ctx *gin.Context) {
	refreshTokenReq, _ := ctx.Get("refresh_token")
	refreshToken := refreshTokenReq.(string)

	accountIDReq, _ := ctx.Get("account_id")
	accountID, ok := accountIDReq.(int64)

	if !ok {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "invalid account id")
		return
	}

	res, errRe := h.accountService.RefreshToken(ctx, accountID, refreshToken)

	if errRe != nil {
		utils.ErrorResponse(ctx, errRe.Code, errRe.Message)
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, res)
}

// CheckToken godoc
//
//	@Summary		check access token
//	@Description	check token every user access web page
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//
//	@Security		BearerAuth
//
//	@Success		200	{object}	dto.CheckTokenResponseDocs
//	@Failure		400	{object}	dto.ErrorResponse
//	@Router			/accounts/check-token [get]
func (h *AccountHandler) CheckToken(ctx *gin.Context) {
	newOTP := h.opts.NewOTP().Key
	utils.SuccessResponse(ctx, http.StatusOK, &dto.CheckTokenResponse{
		OTP: newOTP,
	})
}
