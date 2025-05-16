//go:build wireinject
// +build wireinject

//go:generate wire
package app

import (
	"context"

	"clean-arch/internal/adapter/auth/jwt"
	"clean-arch/internal/adapter/config"
	"clean-arch/internal/adapter/handler/http"
	"clean-arch/internal/adapter/infra/mongo"
	repo "clean-arch/internal/adapter/infra/mongo/repo"
	"clean-arch/internal/core/service"

	"github.com/google/wire"
)

// InitializeApplication wires up all dependencies and returns an *Application.
func InitializeApplication(ctx context.Context) (*Application, error) {
	wire.Build(
		New,         // app.New (constructor for *Application)
		config.Load, // *config.Config
		mongo.New,
		repo.NewUserRepo,    // port.UserRepository
		jwt.ProviderSet, // port.TokenService
		service.ProviderSet, // provides port.UserService, port.AuthService
		http.ProviderSet,
	)
	return nil, nil
}
