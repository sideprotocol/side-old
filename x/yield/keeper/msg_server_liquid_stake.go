package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/yield/types"
)

func (k msgServer) LiquidStake(goCtx context.Context, msg *types.MsgLiquidStake) (*types.MsgLiquidStakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

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
