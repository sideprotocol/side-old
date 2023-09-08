package keeper

import (
	"context"

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

	out, err := pool.EstimateSwap(msg.TokenIn, msg.DenomOut)
	if err != nil {
		return nil, err
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
	return &types.MsgSwapResponse{}, nil
}
