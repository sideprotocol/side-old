package keeper

import (
	"context"
	"fmt"

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
	chainID, err := k.GetChainID(ctx, msg.ConnectionId)
	if err != nil {
		return nil, fmt.Errorf("chain id not found for connection \"%s\": \"%w\"", msg.ConnectionId, err)
	}

	hc := &types.HostChain{
		ChainId:           chainID,
		Bech32Prefix:      msg.Bech32Prefix,
		ConnectionId:      msg.ConnectionId,
		TransferChannelId: msg.TransferChannelId,
		IbcDenom:          msg.IbcDenom,
		HostDenom:         msg.HostDenom,
	}

	// save the host chain
	// k.SetHostChain(ctx, hc)

	// register delegate ICA
	if err = k.RegisterICAAccount(ctx, hc.ConnectionId, hc.DelegationAccount.Owner); err != nil {
		return nil, errorsmod.Wrapf(
			types.ErrFailedToRegisterHostChain,
			"error registering %s delegate ica: %s",
			chainID,
			err.Error(),
		)
	}
	_ = ctx

	return &types.MsgRegisterHostChainResponse{}, nil
}
