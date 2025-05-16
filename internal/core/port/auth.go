//go:generate mockgen -source=auth.go -destination=mock/auth.go -package=mock
package port

import (
	"context"

	"clean-arch/internal/core/domain"
)

// TokenService is an interface for interacting with token-related business logic
type TokenService interface {
	// CreateToken creates a new token for a given user
	CreateToken(user *domain.User) (string, error)
	// VerifyToken verifies the token and returns the payload
	VerifyToken(token string) (*domain.TokenPayload, error)
}

type AuthService interface {
	// AuthenticateUser authenticates a user with the given email and password.
	// It returns the authenticated user and an error if authentication fails.
	Login(ctx context.Context, email, password string) (string, error)
}
