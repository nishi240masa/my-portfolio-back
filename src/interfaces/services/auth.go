package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


type AuthService interface {
    GenerateToken(userID string) (string, error)
    ValidateToken(token string) (string, error)
}

type authService struct {
	secretKey string
}

func NewAuthService() AuthService {
	return &authService{
		secretKey: os.Getenv("SECRET_KEY"),
	}
}

// tokenの生成
func (a *authService) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.secretKey))
}

// tokenの検証
func (a *authService) ValidateToken(token string) (string, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(a.secretKey), nil
	})

	if err != nil || !parsedToken.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil {
		return "", errors.New("invalid claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("invalid user_id type")
	}

	return userID, nil
}
