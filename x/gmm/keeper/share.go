package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/gmm/types"
)

// MintPoolShareToAccount attempts to mint shares of a GAMM denomination to the
// specified address returning an error upon failure. Shares are minted using
// the x/gamm module account.
func (k Keeper) MintPoolShareToAccount(ctx sdk.Context, addr sdk.AccAddress, share sdk.Coin) error {
	amt := sdk.NewCoins(share)
	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, amt)
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, amt)
	if err != nil {
		return err
	}

	return nil
}
