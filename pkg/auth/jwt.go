package auth

import (
	"time"


	"github.com/golang-jwt/jwt/v5"
)

type JWTCustomClaims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func NewClaims(name, email string, admin bool) *JWTCustomClaims {
	return &JWTCustomClaims{
		Name:  name,
		Email: email,
		Admin: admin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
}

func GenerateJWT(claims *JWTCustomClaims, jwtSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return t, nil
}
