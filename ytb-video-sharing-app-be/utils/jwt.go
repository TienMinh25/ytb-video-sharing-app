package utils

import (
	"crypto/rsa"
	"fmt"
	"net/http"
	"os"
	"time"
	"ytb-video-sharing-app-be/internal/dto"
	"ytb-video-sharing-app-be/internal/entities"

	"github.com/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// KeyManager manages private key and public key
type KeyManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

type UserClaims struct {
	AccountID int64  `json:"id"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}

func newUserClaims(id int64, email string, duration time.Duration) (*UserClaims, *dto.ErrorResponse) {
	tokenID, err := uuid.NewRandom()

	if err != nil {
		return nil, &dto.ErrorResponse{Message: INTERNAL_SERVER_ERROR, Code: http.StatusInternalServerError}
	}

	return &UserClaims{
		AccountID: id,
		Email:     email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			Subject:   email,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}, nil
}

// LoadKeys read jwtRSA256.key file and jwtRSA256.key.pub file
func LoadKeys() (*KeyManager, error) {
	privateKeyData, err := os.ReadFile(os.Getenv("PRIVATE_KEY_PATH"))

	if err != nil {
		return nil, errors.Wrap(err, "os.ReadFile")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)

	if err != nil {
		return nil, errors.Wrap(err, "jwt.ParseRSAPrivateKeyFromPEM")
	}

	publicKeyData, err := os.ReadFile(os.Getenv("PUBLIC_KEY_PATH"))

	if err != nil {
		return nil, errors.Wrap(err, "os.ReadFile")
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)

	if err != nil {
		return nil, errors.Wrap(err, "jwt.ParseRSAPublicKeyFromPEM")
	}

	fmt.Println("✅ Keys loaded successfully!")
	return &KeyManager{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func GenerateToken(payload *entities.Account, k *KeyManager, expireAccessToken, expireRefreshToken int) (string, string, *dto.ErrorResponse) {
	claimsAccessToken, errClaims := newUserClaims(payload.ID, payload.Email, time.Duration(expireAccessToken)*time.Minute)

	if errClaims != nil {
		return "", "", errClaims
	}

	claimsRefreshToken, errClaims := newUserClaims(payload.ID, payload.Email, time.Duration(expireRefreshToken)*24*time.Hour)

	if errClaims != nil {
		return "", "", errClaims
	}

	accessTokenJWT := jwt.NewWithClaims(jwt.SigningMethodRS256, claimsAccessToken)
	refreshTokenJWT := jwt.NewWithClaims(jwt.SigningMethodRS256, claimsRefreshToken)

	accessToken, err := accessTokenJWT.SignedString(k.privateKey)

	if err != nil {
		return "", "", &dto.ErrorResponse{Message: INTERNAL_SERVER_ERROR, Code: http.StatusInternalServerError}
	}

	refreshToken, err := refreshTokenJWT.SignedString(k.privateKey)

	if err != nil {
		return "", "", &dto.ErrorResponse{Message: INTERNAL_SERVER_ERROR, Code: http.StatusInternalServerError}
	}

	return accessToken, refreshToken, nil
}

func ValidateToken(tokenStr string, k *KeyManager) (*UserClaims, *dto.ErrorResponse) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Verify the signing method
		_, ok := t.Method.(*jwt.SigningMethodRSA)

		if !ok {
			return nil, errors.New("invalid signing method")
		}

		return k.publicKey, nil
	})

	if err != nil {
		return nil, &dto.ErrorResponse{Message: err.Error(), Code: http.StatusUnauthorized}
	}

	if !token.Valid {
		return nil, &dto.ErrorResponse{Message: "invalid token", Code: http.StatusUnauthorized}
	}

	claims, _ := token.Claims.(*UserClaims)

	return claims, nil
}
