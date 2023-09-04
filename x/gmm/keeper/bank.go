package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/gmm/types"
)

func (k Keeper) LockTokens(ctx sdk.Context, poolId string, sender sdk.AccAddress, tokens sdk.Coins) error {
	escrow := types.GetEscrowAddress(poolId)
	// escrow source tokens. It fails if balance insufficient
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, sender, escrow.String(), tokens,
	); err != nil {
		return err
	}
	return nil
}

func (k Keeper) UnLockTokens(ctx sdk.Context, poolId string, receiver sdk.AccAddress, tokens sdk.Coins) error {
	escrow := types.GetEscrowAddress(poolId)
	// escrow source tokens. It fails if balance insufficient
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx, escrow.String(), receiver, tokens,
	); err != nil {
		return err
	}
	return nil
}

func (k Keeper) BurnTokens(ctx sdk.Context, sender sdk.AccAddress, tokens sdk.Coin) error {
	// transfer the coins to the module account and burn them
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(tokens)); err != nil {
		return err
	}
	if err := k.bankKeeper.BurnCoins(
		ctx, types.ModuleName, sdk.NewCoins(tokens),
	); err != nil {
		// NOTE: should not happen as the module account was
		// retrieved on the step above and it has enough balance
		// to burn.
		panic(fmt.Sprintf("cannot burn coins after a successful send to a module account: %v", err))
	}
	return nil
}

func (k Keeper) MintTokens(ctx sdk.Context, receiver sdk.AccAddress, tokens sdk.Coin) error {
	// mint new tokens if the source of the transfer is the same chain
	if err := k.bankKeeper.MintCoins(
		ctx, types.ModuleName, sdk.NewCoins(tokens),
	); err != nil {
		return err
	}
	// send to receiver
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(
		ctx, types.ModuleName, receiver, sdk.NewCoins(tokens),
	); err != nil {
		panic(fmt.Sprintf("unable to send coins from module to account despite previously minting coins to module account: %v", err))
	}
	return nil
}
