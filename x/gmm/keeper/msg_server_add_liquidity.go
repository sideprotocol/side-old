package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/gmm/types"
)

func (k msgServer) AddLiquidity(goCtx context.Context, msg *types.MsgAddLiquidity) (*types.MsgAddLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	pool, found := k.GetPool(ctx, msg.PoolId)
	if !found {
		return nil, types.ErrNotFoundAssetInPool
	}

	poolCreator := sdk.MustAccAddressFromBech32(msg.Sender)

	// Move asset from sender to module account

	if err := k.LockTokens(ctx, msg.PoolId, poolCreator, types.GetLiquidityAsCoins(msg.Liquidity)); err != nil {
		return nil, err
	}

	// Mint share to creator
	share, err := pool.EstimateShare(msg.Liquidity)
	if err != nil {
		return nil, err
	}
	if err := k.MintPoolShareToAccount(ctx, poolCreator, share); err != nil {
		return nil, err
	}
	// Update pool status
	if err := pool.IncreaseLiquidity(msg.Liquidity); err != nil {
		return nil, err
	}
	pool.IncreaseShare(share.Amount)

	// Save update information
	k.SetPool(ctx, pool)

	// Emit events
	k.EmitEvent(
		ctx, types.EventValueActionDeposit,
		msg.PoolId,
		msg.Sender,
		types.GetEventAttrOfAsset(msg.Liquidity)...,
	)
	return &types.MsgAddLiquidityResponse{}, nil
}
