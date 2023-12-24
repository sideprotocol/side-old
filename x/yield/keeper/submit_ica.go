package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/gogoproto/proto"

	"github.com/Stride-Labs/stride/v16/utils"
	icacallbackstypes "github.com/Stride-Labs/stride/v16/x/icacallbacks/types"

	"github.com/sideprotocol/side/x/yield/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
)

// SubmitTxs submits an ICA transaction containing multiple messages
func (k Keeper) SubmitTxs(
	ctx sdk.Context,
	connectionId string,
	msgs []proto.Message,
	timeoutTimestamp uint64,
	callbackId string,
	callbackArgs []byte,
) (uint64, error) {
	chainId, err := k.GetChainID(ctx, connectionId)
	if err != nil {
		return 0, err
	}
	owner := types.FormatICAAccountOwner(chainId, "delegationAccount")
	portID, err := icatypes.NewControllerPortID(owner)
	if err != nil {
		return 0, err
	}

	k.Logger(ctx).Info(utils.LogWithHostZone(chainId, "  Submitting ICA Tx on %s, %s with TTL: %d", portID, connectionId, timeoutTimestamp))
	protoMsgs := []proto.Message{}
	for _, msg := range msgs {
		k.Logger(ctx).Info(utils.LogWithHostZone(chainId, "    Msg: %+v", msg))
		protoMsgs = append(protoMsgs, msg)
	}

	channelID, found := k.icaControllerKeeper.GetActiveChannelID(ctx, connectionId, portID)
	if !found {
		return 0, errorsmod.Wrapf(icatypes.ErrActiveChannelNotFound, "failed to retrieve active channel for port %s", portID)
	}

	data, err := icatypes.SerializeCosmosTx(k.cdc, protoMsgs)
	if err != nil {
		return 0, err
	}

	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
	}

	// Submit ICA tx
	msgServer := icacontrollerkeeper.NewMsgServerImpl(&k.icaControllerKeeper)
	relativeTimeoutOffset := timeoutTimestamp - uint64(ctx.BlockTime().UnixNano())
	msgSendTx := icacontrollertypes.NewMsgSendTx(owner, connectionId, relativeTimeoutOffset, packetData)
	res, err := msgServer.SendTx(ctx, msgSendTx)
	if err != nil {
		return 0, err
	}
	sequence := res.Sequence

	// Store the callback data
	if callbackId != "" && callbackArgs != nil {
		callback := icacallbackstypes.CallbackData{
			CallbackKey:  icacallbackstypes.PacketID(portID, channelID, sequence),
			PortId:       portID,
			ChannelId:    channelID,
			Sequence:     sequence,
			CallbackId:   callbackId,
			CallbackArgs: callbackArgs,
		}
		k.Logger(ctx).Info(utils.LogWithHostZone(chainId, "Storing callback data: %+v", callback))
		k.ICACallbacksKeeper.SetCallbackData(ctx, callback)
	}

	return sequence, nil
}
