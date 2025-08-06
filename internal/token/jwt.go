package token

import (
	"T/internal/config"
	"crypto/sha512"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

type JwtManager struct {
	secretKey []byte
}

type AccessToken struct {
	Token Token
}

type RefreshToken struct {
	Token Token
}

func NewJWTManager(secret string) *JwtManager {
	hashed := sha512.Sum512([]byte(secret))
	return &JwtManager{secretKey: hashed[:]}
}

func NewAccessToken(cfg *config.Config) *AccessToken {
	return &AccessToken{
		Token: NewJWTManager(cfg.JWT.Secret),
	}
}

func NewRefreshToken(cfg *config.Config) *RefreshToken {
	return &RefreshToken{
		Token: NewJWTManager(cfg.JWT.Secret + "_refresh"),
	}
}

func (j *JwtManager) CreateToken(userID, roleID, sessionID uuid.UUID, duration time.Duration) (string, error) {
	claims := &Claims{
		UserID:    userID,
		RoleID:    roleID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString(j.secretKey)
}

func (j *JwtManager) VerifyToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS512 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
