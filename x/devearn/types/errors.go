package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/devearn module sentinel errors
var (
	ErrInternalDevEarn  = errorsmod.Register(ModuleName, 2, "internal dev earn error")
	ErrContractNotFound = errorsmod.Register(ModuleName, 1100, "contract is not fond,please check your contract address")
)
