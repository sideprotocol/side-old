package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/yield module sentinel errors
var (
	ErrSample                    = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrFailedToRegisterHostChain = sdkerrors.Register(ModuleName, 1101, "failed to register host chain")
)
