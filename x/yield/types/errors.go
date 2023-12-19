package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/yield module sentinel errors
var (
	ErrSample                    = errorsmod.Register(ModuleName, 1100, "sample error")
	ErrFailedToRegisterHostChain = errorsmod.Register(ModuleName, 1101, "failed to register host chain")
)
