package http

import "github.com/google/wire"

// ProviderSet is the wire provider set for the HTTP handler.
var ProviderSet = wire.NewSet(
	NewUserHandler,
	NewAuthHandler,

	NewRouter,
)
