package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/gmm/types"
)

func (k msgServer) Withdraw(goCtx context.Context, msg *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	pool, found := k.GetPool(ctx, msg.PoolId)
	if !found {
		return nil, types.ErrNotFoundAssetInPool
	}

	outs, err := pool.EstimateWithdrawals(msg.Share)
	if err != nil {
		return nil, err
	}

	// Check creator has enough share or not
	bal := k.bankKeeper.GetBalance(ctx, sdk.MustAccAddressFromBech32(msg.Creator),
		msg.Share.Denom,
	)
	if bal.IsLT(msg.Share) {
		return nil, types.ErrInsufficientBalance
	}

	// Unlock asset from pool
	if err = k.UnLockTokens(ctx, pool.PoolId, sdk.MustAccAddressFromBech32(msg.Receiver), outs); err != nil {
		return nil, err
	}

	// Burn lp token
	if err = k.BurnTokens(ctx, sdk.MustAccAddressFromBech32(msg.Creator), msg.Share); err != nil {
		return nil, err
	}

	// Update pool statues.
	if err := pool.DecreaseLiquidity(outs); err != nil {
		return nil, err
	}
	pool.DecreaseShare(msg.Share.Amount)

	// Save pool
	k.SetPool(ctx, pool)
	return &types.MsgWithdrawResponse{}, nil
}
