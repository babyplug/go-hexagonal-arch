package jwt

import (
	"errors"
	"time"

	"clean-arch/internal/adapter/config"
	"clean-arch/internal/core/domain"
	"clean-arch/internal/core/port"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenAuth struct {
	secret   string
	duration time.Duration
}

func New(cfg *config.Config) (port.TokenService, error) {
	duration, err := time.ParseDuration(cfg.Duration)
	if err != nil {
		return nil, errors.New("invalid duration format")
	}

	return &TokenAuth{secret: cfg.JWTSecret, duration: duration}, nil
}

func (ta *TokenAuth) CreateToken(user *domain.User) (string, error) {
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(ta.duration)

	claims := jwt.MapClaims{
		"sub": user.ID,
		"iat": issuedAt.Unix(),
		"exp": expiredAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(ta.secret))
}

func (ta *TokenAuth) VerifyToken(token string) (*domain.TokenPayload, error) {
	tk, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenMalformed
		}
		return []byte(ta.secret), nil
	})

	if err != nil || !tk.Valid {
		return nil, errors.New("invalid token")
	}
	claims := tk.Claims.(jwt.MapClaims)
	return &domain.TokenPayload{
		ID: claims["sub"].(uuid.UUID),
	}, nil
}
