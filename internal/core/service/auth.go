package service

import (
	"context"
	"errors"
	"sync"

	"clean-arch/internal/core/domain"
	"clean-arch/internal/core/port"
	"clean-arch/internal/core/util"
)

var (
	auth     port.AuthService
	authOnce sync.Once
)

type AuthService struct {
	repo port.UserRepository
	ts   port.TokenService
}

func NewAuth(userRepo port.UserRepository, ts port.TokenService) port.AuthService {
	authOnce.Do(func() {
		auth = &AuthService{
			repo: userRepo,
			ts:   ts,
		}
	})

	return auth
}

func ResetAuth() {
	authOnce = sync.Once{}
}

func (a *AuthService) Login(ctx context.Context, email, password string) (token string, err error) {
	user, err := a.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrDataNotFound) {
			return "", domain.ErrDataNotFound
		}
		return "", domain.ErrInternal
	}

	err = util.ComparePassword(password, user.Password)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	token, err = a.ts.CreateToken(user)
	if err != nil {
		return "", domain.ErrTokenCreationFailed
	}

	return token, nil
}
