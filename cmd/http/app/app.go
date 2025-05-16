package app

import (
	"context"
	"log"
	"net/http"
	"sync"

	"clean-arch/internal/adapter/background"
	"clean-arch/internal/adapter/config"
	handler "clean-arch/internal/adapter/handler/http"
	"clean-arch/internal/adapter/infra/mongo"
	"clean-arch/internal/core/port"
)

type Application struct {
	Config      *config.Config
	MongoClient mongo.Client
	UserRepo    port.UserRepository
	Router      *handler.Router
	server      *http.Server
}

var (
	_app     = &Application{}
	_appOnce sync.Once
)

// New creates a new Application instance.
func New(
	ctx context.Context,
	cfg *config.Config,
	client mongo.Client,
	userRepo port.UserRepository,
	userHandler *handler.UserHandler,
	authHandler *handler.AuthHandler,
	ts port.TokenService,
) (*Application, error) {
	var err error
	_appOnce.Do(func() {
		// Init router
		router, err := handler.NewRouter(
			cfg,
			ts,
			userHandler,
			authHandler,
		)
		if err != nil {
			log.Fatalf("Failed to create router: %v", err)
		}

		srv := &http.Server{
			Addr:    cfg.Port,
			Handler: router.Handler(),
		}

		_app = &Application{
			Config:      cfg,
			UserRepo:    userRepo,
			Router:      router,
			server:      srv,
			MongoClient: client,
		}
	})

	return _app, err
}

// ListenAndServe starts the HTTP server.
func (a *Application) ListenAndServe() error {
	if err := a.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

// Shutdown gracefully shuts down the server.
func (a *Application) Shutdown(ctx context.Context) error {
	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

// StartBackgroundProcess starts the background process for logging, etc.
func (a *Application) StartBackgroundProcess(stopCh chan struct{}) {
	background.StartUserCountLogger(a.UserRepo, stopCh)
}
