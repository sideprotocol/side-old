package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/yield module sentinel errors
var (
	ErrSample                    = errorsmod.Register(ModuleName, 1100, "sample error")
	ErrFailedToRegisterHostChain = errorsmod.Register(ModuleName, 1101, "failed to register host chain")
	ErrUnmarshalFailure          = errorsmod.Register(ModuleName, 1102, "cannot unmarshal")
	ErrUnknownDepositRecord      = errorsmod.Register(ModuleName, 1103, "unknown deposit record")
)
