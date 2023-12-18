package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/sideprotocol/side/x/yield/types"
)

func (k msgServer) RegisterHostChain(goCtx context.Context, msg *types.MsgRegisterHostChain) (*types.MsgRegisterHostChainResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	if msg.Creator != k.GetParams(ctx).Admin {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "This operation can only be performed by admin")
	}

	// get the host chain id
	// chainID, err := k.GetChainID(ctx, msg.ConnectionId)
	// if err != nil {
	// 	return nil, fmt.Errorf("chain id not found for connection \"%s\": \"%w\"", msg.ConnectionId, err)
	// }
	_ = ctx

	return &types.MsgRegisterHostChainResponse{}, nil
}
