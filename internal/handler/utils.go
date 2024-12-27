package handler

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"regexp"
	"skymates-api/internal/types"
	"time"
)

func isValidEmail(email string) bool {
	// ^[a-zA-Z0-9] 要求以字母或数字开头
	// [a-zA-Z0-9._%+-]* 允许中间部分有字母、数字、符号等
	pattern := `^[a-zA-Z0-9][a-zA-Z0-9._%+-]*@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	return regexp.MustCompile(pattern).MatchString(email)
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (h *UserHandler) generateJWT(user *types.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(os.Getenv("JWT_SECRET"))
}
