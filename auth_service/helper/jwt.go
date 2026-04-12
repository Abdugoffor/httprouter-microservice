package helper

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(ENV("JWT_KEY"))

type Claims struct {
	UserID int `json:"user_id"`
	RoleID int `json:"role_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int, roleID int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: userID,
		RoleID: roleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}

func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims.(*Claims), nil
}
