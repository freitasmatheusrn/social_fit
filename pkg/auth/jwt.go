package auth

import (
	"time"

	"github.com/freitasmatheusrn/social-fit/internal/user"
	"github.com/golang-jwt/jwt/v5"
)

type JWTCustomClaims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func NewClaims(user user.SigninResponse) *JWTCustomClaims {
	return &JWTCustomClaims{
		Name:  user.Name,
		Email: user.Email,
		Admin: user.Admin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
}

func GenerateJWT(claims *JWTCustomClaims, jwtSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return t, nil
}
