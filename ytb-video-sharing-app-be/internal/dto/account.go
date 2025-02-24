package dto

type CreateAccountRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FullName  string `json:"fullname" binding:"required"`
	AvatarURL string `json:"avatar_url" binding:"omitempty,url"`
}

type CreateAccountResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccountResponse
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	AccessToken  string          `json:"access_token"`
	RefreshToken string          `json:"refresh_token"`
	Info         AccountResponse `json:"info"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LogoutResponse struct {
}

type AccountResponse struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	FullName  string `json:"fullname"`
	AvatarURL string `json:"avatar_url"`
}
