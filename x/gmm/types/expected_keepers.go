package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktype "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// Methods imported from account should be defined here
	NewAccount(sdk.Context, types.AccountI) types.AccountI

	SetAccount(ctx sdk.Context, acc types.AccountI)
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here

	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error

	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SetDenomMetaData(ctx sdk.Context, denomMetaData banktype.Metadata)

	MintCoins(ctx sdk.Context, moduleName string, amounts sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amounts sdk.Coins) error

	HasSupply(ctx sdk.Context, denom string) bool
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}
