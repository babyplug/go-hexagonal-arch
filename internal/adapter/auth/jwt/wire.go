package jwt

import (
	"github.com/google/wire"
)

// ProviderSet is the wire provider set for the JWT package.
var ProviderSet = wire.NewSet(
	New,
)
