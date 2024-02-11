package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/gmm/types"
)

func (k msgServer) Swap(goCtx context.Context, msg *types.MsgSwap) (*types.MsgSwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	pool, found := k.GetPool(ctx, msg.PoolId)
	if !found {
		return nil, types.ErrNotFoundAssetInPool
	}

	out, err := pool.EstimateSwap(msg.TokenIn, msg.TokenOut.Denom)
	if err != nil {
		return nil, err
	}

	// Track APR
	fee := pool.TakeFees(msg.TokenIn.Amount)
	if err := k.ObserveFeeFromPool(ctx, pool.PoolId, sdk.NewCoin(msg.TokenIn.Denom, fee.RoundInt())); err != nil {
		return nil, err
	}
	// Calculate the absolute difference between the expected and actual token output amounts
	differ := sdkmath.NewInt(0)
	if out.Amount.LT(msg.TokenOut.Amount) {
		differ = msg.TokenOut.Amount.Sub(out.Amount)
	}
	// Calculate the expected slippage. Make sure msg.Slippage is in the correct unit (e.g., percentage).
	// Divide by 100 if msg.Slippage is a percentage.
	expectedDiffer := msg.TokenOut.Amount.Mul(msg.Slippage).Quo(sdkmath.NewInt(100))

	// Check if the actual slippage exceeds the expected slippage
	if differ.GT(expectedDiffer) {
		return nil, types.ErrNotMeetSlippage
	}

	// Move asset from sender to module account
	if err := k.LockTokens(ctx, msg.PoolId, sdk.MustAccAddressFromBech32(msg.Sender), sdk.NewCoins(
		msg.TokenIn,
	)); err != nil {
		return nil, err
	}

	// Send wanted token from pool to user
	if err := k.UnLockTokens(ctx, msg.PoolId, sdk.MustAccAddressFromBech32(msg.Sender), sdk.NewCoins(
		out,
	)); err != nil {
		return nil, err
	}

	// Update state.
	if err := pool.IncreaseLiquidity([]sdk.Coin{msg.TokenIn}); err != nil {
		return nil, err
	}
	if err := pool.DecreaseLiquidity([]sdk.Coin{out}); err != nil {
		return nil, err
	}

	// Save pool.
	k.SetPool(ctx, pool)
	// Observe volume
	err = k.ObserveVolumeFromPool(ctx, pool.PoolId, msg.TokenIn, out)
	return &types.MsgSwapResponse{}, err
}
