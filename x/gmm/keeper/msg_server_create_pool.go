package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/gmm/types"
)

func (k msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePool) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Initialize pool
	pooID, err := k.initializePool(ctx, msg)
	if err != nil {
		return nil, err
	}

	// Emit events
	k.EmitEvent(
		ctx, types.EventValueActionCreatePool, *pooID,
		msg.Creator,
		sdk.Attribute{
			Key:   types.AttributeKeyPoolCreator,
			Value: msg.Creator,
		},
	)

	return &types.MsgCreatePoolResponse{
		PoolId: *pooID,
	}, nil
}
