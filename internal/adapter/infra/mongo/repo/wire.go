package repo

import "github.com/google/wire"

// ProviderSet is the wire provider set for the MongoDB repository.
var ProviderSet = wire.NewSet(
	NewUserRepo,
)
