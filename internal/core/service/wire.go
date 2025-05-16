package service

import (
	"github.com/google/wire"
)

// ProviderSet is the wire provider set for the service package.
var ProviderSet = wire.NewSet(
	NewUser,
	NewAuth,
)