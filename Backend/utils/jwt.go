package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID uint, role string) (string, error) {
	secret, err := jwtSecret()
	if err != nil {
		return "", err
	}
	canonicalRole := CanonicalRole(role)
	if userID == 0 || canonicalRole == "" {
		return "", errors.New("claim JWT tidak valid")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func ParseJWT(tokenString string) (uint, string, error) {
	secret, err := jwtSecret()
	if err != nil {
		return 0, "", err
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("metode token tidak valid")
		}

		return secret, nil
	})

	if err != nil || !token.Valid {
		return 0, "", errors.New("token tidak valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.New("claims token tidak valid")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", errors.New("user_id token tidak valid")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return 0, "", errors.New("role token tidak valid")
	}

	if userIDFloat <= 0 || userIDFloat != float64(uint64(userIDFloat)) || CanonicalRole(role) == "" {
		return 0, "", errors.New("claim token tidak valid")
	}

	return uint(userIDFloat), CanonicalRole(role), nil
}

func jwtSecret() ([]byte, error) {
	secret := os.Getenv("JWT_SECRET")
	if len(secret) < 32 || strings.Contains(strings.ToLower(secret), "replace_with") || strings.Contains(strings.ToLower(secret), "change_me") {
		return nil, fmt.Errorf("JWT_SECRET wajib diisi dengan secret minimal 32 karakter")
	}
	return []byte(secret), nil
}

func ValidateJWTSecret() error {
	_, err := jwtSecret()
	return err
}
