package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/yield/types"
)

func (k msgServer) LiquidStake(goCtx context.Context, msg *types.MsgLiquidStake) (*types.MsgLiquidStakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the host chain from the base denom in the message (e.g. uatom)
	hostChain, err := k.GetHostChainFromHostDenom(ctx, msg.Denom)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidToken, "no host zone found for denom (%s)", msg.Denom)
	}

	// Get user and module account addresses
	liquidStakerAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "user's address is invalid")
	}
	// hostZoneDepositAddress, err := sdk.AccAddressFromBech32(hostChain.DepositAddress)
	// if err != nil {
	// 	return nil, errorsmod.Wrapf(err, "host zone address is invalid")
	// }

	// The tokens that are sent to the protocol are denominated in the ibc hash of the native token on side (e.g. ibc/xxx)
	nativeDenom := hostChain.IbcDenom
	nativeCoin := sdk.NewCoin(nativeDenom, msg.Amount)
	if !types.IsIBCToken(nativeDenom) {
		return nil, errorsmod.Wrapf(types.ErrInvalidToken, "denom is not an IBC token (%s)", nativeDenom)
	}

	// Confirm the user has a sufficient balance to execute the liquid stake
	balance := k.bankKeeper.GetBalance(ctx, liquidStakerAddress, nativeDenom)
	if balance.IsLT(nativeCoin) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds, "balance is lower than staking amount. staking amount: %v, balance: %v", msg.Amount, balance.Amount)
	}

	// // Determine the amount of sdTokens to mint using the redemption rate
	// stAmount := (sdk.NewDecFromInt(msg.Amount).Quo(hostZone.RedemptionRate)).TruncateInt()
	// if stAmount.IsZero() {
	// 	return nil, errorsmod.Wrapf(types.ErrInsufficientLiquidStake,
	// 		"Liquid stake of %s%s would return 0 sdTokens", msg.Amount.String(), hostZone.HostDenom)
	// }

	// // Transfer the native tokens from the user to module account
	// if err := k.bankKeeper.SendCoins(ctx, liquidStakerAddress, hostZoneDepositAddress, sdk.NewCoins(nativeCoin)); err != nil {
	// 	return nil, errorsmod.Wrap(err, "failed to send tokens from Account to Module")
	// }

	// // Mint the sdTokens and transfer them to the user
	// stDenom := types.StAssetDenomFromHostZoneDenom(msg.HostDenom)
	// stCoin := sdk.NewCoin(stDenom, stAmount)
	// if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(stCoin)); err != nil {
	// 	return nil, errorsmod.Wrapf(err, "Failed to mint coins")
	// }
	// if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, liquidStakerAddress, sdk.NewCoins(stCoin)); err != nil {
	// 	return nil, errorsmod.Wrapf(err, "Failed to send %s from module to account", stCoin.String())
	// }

	// // Update the liquid staked amount on the deposit record
	// depositRecord.Amount = depositRecord.Amount.Add(msg.Amount)
	// k.RecordsKeeper.SetDepositRecord(ctx, *depositRecord)

	// // Emit liquid stake event
	// EmitSuccessfulLiquidStakeEvent(ctx, msg, *hostZone, stAmount)

	// k.hooks.AfterLiquidStake(ctx, liquidStakerAddress)

	// TODO: Handling the message
	// - Check balance
	// - Transfer amount given in message
	// - Deduce host denom from ibc from (msg.denom)
	// - Check for host zone using hose denom
	// - Add deposit record with status, amount and id
	// - Transfer native tokens(1st step of process)
	_ = ctx

	return &types.MsgLiquidStakeResponse{}, nil
}
