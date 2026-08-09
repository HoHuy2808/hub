package authenticate

import (
	"errors"
	"fmt"
	"os"
	"time"

	"hub/internal/database/entities"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Token string

const (
	AccessToken  Token = "AccessToken"
	RefreshToken Token = "RefreshToken"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	UserId uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
}

func GenerateAccessToken(user *entities.User) (string, error) {
	expire := time.Now().Add(1 * time.Hour)

	claims := TokenClaims{
		UserId: user.Id,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expire),
		},
	}

	// expire := time.Now().Add(1 * time.Hour).Unix()

	// claims := jwt.MapClaims{
	// 	"userName": user.UserName,
	// 	"email":    user.Email,
	// 	"expire":   expire,
	// }

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateRefreshToken(user *entities.User) (string, error) {
	expire := time.Now().Add(7 * 24 * time.Hour)

	claims := TokenClaims{
		UserId: user.Id,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expire),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_REFRESH_TOKEN_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(tokenString string) (*TokenClaims, error) {
	tokenClaims := &TokenClaims{}

	token, err := jwt.ParseWithClaims(tokenString, tokenClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return tokenClaims, nil
}
