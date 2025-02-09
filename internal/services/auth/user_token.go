package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const tokenTTL = 45 * time.Minute

func (s *Service) CreateAuthToken(userID string) (*string, error) {

	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(tokenTTL)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(s.cfg.AuthSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return &signedToken, nil
}

func (s *Service) VerifyAuthToken(tokenStr string) (*string, error) {
	claims := jwt.RegisteredClaims{}

	_, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.AuthSecret), nil
	})
	// в ідеалі перевіряти на кокретний тип помилки який віддає жвт
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// отримуємо ідентифікатор користувача з токена
	return &claims.Subject, nil

}
