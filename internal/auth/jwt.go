package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"skymates-api/internal/types"
	"time"
)

// TokenExpiry is the duration of the token, 24 hours
const TokenExpiry = 24 * time.Hour

// Claims 是 JWT Payload 部分, 明文的, 不要存储敏感信息
// JWT 的结构: 1. Header: 描述签名算法(如 HS256)
// 2. Payload: 存储 Claims 信息
// 3. Signature: 签名, 用于验证 token 完整性
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJwtToken(user *types.User) (string, error) {
	expirationTime := time.Now().Add(TokenExpiry)

	claims := &Claims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "skymates",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	return token.SignedString(secretKey)
}

func ValidateJwtToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法是否属于 HMAC 家族, HS256, HS384 和 HS512 等所有 HMAC 算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
	)

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrSignatureInvalid):
			return nil, fmt.Errorf("invalid token signature")
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, fmt.Errorf("token expired")
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, fmt.Errorf("token not active yet")
		default:
			return nil, fmt.Errorf("invalid token: %w", err)
		}
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
