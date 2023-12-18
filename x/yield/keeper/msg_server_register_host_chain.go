package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/yield/types"
)

func (k msgServer) RegisterHostChain(goCtx context.Context, msg *types.MsgRegisterHostChain) (*types.MsgRegisterHostChainResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgRegisterHostChainResponse{}, nil
}
