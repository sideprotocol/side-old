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

	if msg.Creator != k.GetParams(ctx).Admin {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "This operation can only be performed by admin")
	}

	// Get ConnectionEnd (for counterparty connection)
	_, found := k.ibcKeeper.ConnectionKeeper.GetConnection(ctx, msg.ConnectionId)
	if !found {
		errMsg := fmt.Sprintf("invalid connection id, %s not found", msg.ConnectionId)
		k.Logger(ctx).Error(errMsg)
		return nil, errorsmod.Wrapf(types.ErrFailedToRegisterHostChain, errMsg)
	}

	// get the host chain id
	chainID, err := k.GetChainID(ctx, msg.ConnectionId)
	if err != nil {
		return nil, fmt.Errorf("chain id not found for connection \"%s\": \"%w\"", msg.ConnectionId, err)
	}

	// get zone
	_, found = k.GetHostChain(ctx, chainID)
	if found {
		errMsg := fmt.Sprintf("invalid chain id, zone for %s already registered", chainID)
		k.Logger(ctx).Error(errMsg)
		return nil, errorsmod.Wrapf(types.ErrFailedToRegisterHostChain, errMsg)
	}

	// check the denom is not already registered
	hostChains := k.GetAllHostChain(ctx)
	for _, hostChain := range hostChains {
		if hostChain.HostDenom == msg.HostDenom {
			errMsg := fmt.Sprintf("host denom %s already registered", msg.HostDenom)
			k.Logger(ctx).Error(errMsg)
			return nil, errorsmod.Wrapf(types.ErrFailedToRegisterHostChain, errMsg)
		}
		if hostChain.ConnectionId == msg.ConnectionId {
			errMsg := fmt.Sprintf("connectionId %s already registered", msg.ConnectionId)
			k.Logger(ctx).Error(errMsg)
			return nil, errorsmod.Wrapf(types.ErrFailedToRegisterHostChain, errMsg)
		}
		if hostChain.TransferChannelId == msg.TransferChannelId {
			errMsg := fmt.Sprintf("transfer channel %s already registered", msg.TransferChannelId)
			k.Logger(ctx).Error(errMsg)
			return nil, errorsmod.Wrapf(types.ErrFailedToRegisterHostChain, errMsg)
		}
		if hostChain.Bech32Prefix == msg.Bech32Prefix {
			errMsg := fmt.Sprintf("bech32prefix %s already registered", msg.Bech32Prefix)
			k.Logger(ctx).Error(errMsg)
			return nil, errorsmod.Wrapf(types.ErrFailedToRegisterHostChain, errMsg)
		}
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
	k.SetHostChain(ctx, hc)

	// register delegate ICA
	// Note: We only need one account for now, can add multiple accounts in future
	delegateAccount := types.FormatICAAccountOwner(hc.ChainId, "delegationAccount")
	if err = k.RegisterICAAccount(ctx, hc.ConnectionId, delegateAccount); err != nil {
		return nil, errorsmod.Wrapf(
			types.ErrFailedToRegisterHostChain,
			"error registering %s delegate ica: %s",
			chainID,
			err.Error(),
		)
	}

	return &types.MsgRegisterHostChainResponse{}, nil
}
