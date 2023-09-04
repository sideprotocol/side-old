package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/gmm module sentinel errors
var (
	ErrSample              = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrAlreadyCreatedPool  = sdkerrors.Register(ModuleName, 1101, "already exist pool in the system")
	ErrNotEnoughBalance    = sdkerrors.Register(ModuleName, 1102, "not enough balance")
	ErrInvalidPoolType     = sdkerrors.Register(ModuleName, 1103, "invalid pool type")
	ErrInvalidNumOfAssets  = sdkerrors.Register(ModuleName, 1104, "invalid number of assets")
	ErrNotFoundAssetInPool = sdkerrors.Register(ModuleName, 1105, "did not find asset")
	ErrNotFoundAPool       = sdkerrors.Register(ModuleName, 1106, "did not find pool")
	ErrOverflowShareAmount = sdkerrors.Register(ModuleName, 1107, "share amount is over than pool total supply")
)
