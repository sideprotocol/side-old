package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/devearn module sentinel errors
var (
	ErrInternalDevEarn = errorsmod.Register(ModuleName, 2, "internal incentives error")
	ErrSample          = errorsmod.Register(ModuleName, 1100, "sample error")
)
