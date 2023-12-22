package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"

	"github.com/Stride-Labs/stride/v16/utils"
	icacallbackstypes "github.com/Stride-Labs/stride/v16/x/icacallbacks/types"

	"github.com/sideprotocol/side/x/yield/types"
)

// Transfers native tokens, to the cosmos hub.
// TODO: Then from cosmos hub to stride
func (k Keeper) IBCTransferNativeTokens(ctx sdk.Context, msg *transfertypes.MsgTransfer, depositRecord types.DepositRecord) error {
	// Submit IBC transfer
	msgTransferResponse, err := k.TransferKeeper.Transfer(sdk.WrapSDKContext(ctx), msg)
	if err != nil {
		return err
	}

	// Build callback data
	transferCallback := types.TransferCallback{
		DepositRecordId: depositRecord.Id,
	}
	k.Logger(ctx).Info(utils.LogWithHostZone(depositRecord.HostChainId, "Marshaling TransferCallback args: %+v", transferCallback))
	marshalledCallbackArgs, err := k.MarshalTransferCallbackArgs(ctx, transferCallback)
	if err != nil {
		return err
	}

	// Store the callback data
	sequence := msgTransferResponse.Sequence
	callback := icacallbackstypes.CallbackData{
		CallbackKey:  icacallbackstypes.PacketID(msg.SourcePort, msg.SourceChannel, sequence),
		PortId:       msg.SourcePort,
		ChannelId:    msg.SourceChannel,
		Sequence:     sequence,
		CallbackId:   IBCCallbacksIDNativeTransfer,
		CallbackArgs: marshalledCallbackArgs,
	}
	k.Logger(ctx).Info(utils.LogWithHostZone(depositRecord.HostChainId, "Storing callback data: %+v", callback))
	k.ICACallbacksKeeper.SetCallbackData(ctx, callback)

	// update the record state to TRANSFER_IN_PROGRESS
	// TODO: Add ibc transfer tests
	depositRecord.Status = types.DepositRecord_TRANSFER_FIRST_IN_PROGRESS
	k.SetDepositRecord(ctx, depositRecord)

	return nil
}
