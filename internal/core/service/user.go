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
	user     *userServiceImpl
	userOnce sync.Once
)

type userServiceImpl struct {
	repo port.UserRepository
}

func NewUser(repo port.UserRepository) port.UserService {
	userOnce.Do(func() {
		user = &userServiceImpl{repo: repo}
	})

	return user
}

func ResetUser() {
	userOnce = sync.Once{}
}

func (s *userServiceImpl) Create(ctx context.Context, user *domain.User) error {
	existing, _ := s.repo.GetByEmail(ctx, user.Email)
	if existing != nil {
		return domain.ErrDuplicateEmail
	}

	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	if err := s.repo.Create(ctx, user); err != nil {
		return err
	}
	return nil
}

func (s *userServiceImpl) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userServiceImpl) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *userServiceImpl) List(ctx context.Context, page, size int64) ([]*domain.User, error) {
	return s.repo.List(ctx, page, size)
}

func (s *userServiceImpl) Update(ctx context.Context, user *domain.User) error {
	existingUser, err := s.repo.GetByID(ctx, user.ID)
	if err != nil {
		return domain.ErrDataNotFound
	}

	if existingUser.Email != user.Email {
		existing, _ := s.repo.GetByEmail(ctx, user.Email)
		if existing != nil {
			return errors.New("email already exists")
		}
	}

	return s.repo.Update(ctx, user)
}

func (s *userServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *userServiceImpl) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}
