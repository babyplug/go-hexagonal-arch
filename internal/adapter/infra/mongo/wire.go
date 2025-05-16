//go:build wireinject
// +build wireinject

//go:generate wire

package mongo

import (
	"clean-arch/internal/adapter/config"
	"context"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	New,
)

func Wire(ctx context.Context, cfg *config.Config) (Client, error) {
	wire.Build(
		New,
	)
	return &client{}, nil
}
