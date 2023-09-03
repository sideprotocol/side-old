package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/gmm module sentinel errors
var (
	ErrSample             = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrAlreadyCreatedPool = sdkerrors.Register(ModuleName, 1101, "already exist pool in the system")
	ErrNotEnoughBalance   = sdkerrors.Register(ModuleName, 1102, "not enough balance")
)
