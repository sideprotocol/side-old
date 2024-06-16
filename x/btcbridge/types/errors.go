package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/yield module sentinel errors
var (
	ErrSenderAddressNotAuthorized = errorsmod.Register(ModuleName, 1000, "sender address not authorized")
	ErrInvalidHeader              = errorsmod.Register(ModuleName, 1100, "invalid block header")
	ErrReorgFailed                = errorsmod.Register(ModuleName, 1101, "failed to reorg chain")
	ErrForkedBlockHeader          = errorsmod.Register(ModuleName, 1102, "Invalid forked block header")

	ErrInvalidSenders = errorsmod.Register(ModuleName, 2100, "invalid allowed senders")

	ErrInvalidBtcTransaction    = errorsmod.Register(ModuleName, 3100, "invalid bitcoin transaction")
	ErrBlockNotFound            = errorsmod.Register(ModuleName, 3101, "block not found")
	ErrTransactionNotIncluded   = errorsmod.Register(ModuleName, 3102, "transaction not included in block")
	ErrNotConfirmed             = errorsmod.Register(ModuleName, 3200, "transaction not confirmed")
	ErrExceedMaxAcceptanceDepth = errorsmod.Register(ModuleName, 3201, "exceed max acceptance block depth")
	ErrUnsupportedScriptType    = errorsmod.Register(ModuleName, 3202, "unsupported script type")

	ErrInvalidSignatures      = errorsmod.Register(ModuleName, 4200, "invalid signatures")
	ErrInsufficientBalance    = errorsmod.Register(ModuleName, 4201, "insufficient balance")
	ErrSigningRequestNotExist = errorsmod.Register(ModuleName, 4202, "signing request does not exist")

	ErrUTXODoesNotExist = errorsmod.Register(ModuleName, 5100, "utxo does not exist")
	ErrUTXOLocked       = errorsmod.Register(ModuleName, 5101, "utxo locked")
	ErrUTXOUnlocked     = errorsmod.Register(ModuleName, 5102, "utxo unlocked")

	ErrInvalidFeeRate         = errorsmod.Register(ModuleName, 6100, "invalid fee rate")
	ErrDustOutput             = errorsmod.Register(ModuleName, 6101, "dust output value")
	ErrInsufficientUTXOs      = errorsmod.Register(ModuleName, 6102, "insufficient utxos")
	ErrFailToBuildTransaction = errorsmod.Register(ModuleName, 6103, "failed to build transaction")
	ErrFailToSerializePsbt    = errorsmod.Register(ModuleName, 6104, "failed to serialize psbt")
)
