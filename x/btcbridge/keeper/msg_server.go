package keeper

import (
	"bytes"
	"context"
	"encoding/base64"
	"strconv"

	"github.com/btcsuite/btcd/btcutil/psbt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/btcbridge/types"
)

type msgServer struct {
	Keeper
}

// SubmitBlockHeaders implements types.MsgServer.
func (m msgServer) SubmitBlockHeaders(goCtx context.Context, msg *types.MsgSubmitBlockHeaderRequest) (*types.MsgSubmitBlockHeadersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// check if the sender is one of the authorized senders
	param := m.GetParams(ctx)
	if !param.IsAuthorizedSender(msg.Sender) {
		return nil, types.ErrSenderAddressNotAuthorized
	}

	// Set block headers
	err := m.SetBlockHeaders(ctx, msg.BlockHeaders)
	if err != nil {
		return nil, err
	}

	// Emit events
	// m.EmitEvent(
	// 	ctx,
	// 	msg.Sender,
	// 	sdk.Attribute{
	// 		Key:   types.AttributeKeyPoolCreator,
	// 		Value: msg.Sender,
	// 	},
	// )
	return &types.MsgSubmitBlockHeadersResponse{}, nil
}

// SubmitTransaction implements types.MsgServer.
// No Permission check required for this message
// Since everyone can submit a transaction to mint voucher tokens
// This message is usually sent by relayers
func (m msgServer) SubmitDepositTransaction(goCtx context.Context, msg *types.MsgSubmitDepositTransactionRequest) (*types.MsgSubmitDepositTransactionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		ctx.Logger().Error("Error validating basic", "error", err)
		return nil, err
	}

	txHash, recipient, err := m.ProcessBitcoinDepositTransaction(ctx, msg)
	if err != nil {
		ctx.Logger().Error("Error processing bitcoin deposit transaction", "error", err)
		return nil, err
	}

	// Emit Events
	m.EmitEvent(ctx, msg.Sender,
		sdk.NewAttribute("blockhash", msg.Blockhash),
		sdk.NewAttribute("txBytes", msg.TxBytes),
		sdk.NewAttribute("txid", txHash.String()),
		sdk.NewAttribute("recipient", recipient.EncodeAddress()),
	)

	return &types.MsgSubmitDepositTransactionResponse{}, nil
}

// SubmitTransaction implements types.MsgServer.
// No Permission check required for this message
// Since everyone can submit a transaction to mint voucher tokens
// This message is usually sent by relayers
func (m msgServer) SubmitWithdrawTransaction(goCtx context.Context, msg *types.MsgSubmitWithdrawTransactionRequest) (*types.MsgSubmitWithdrawTransactionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		ctx.Logger().Error("Error validating basic", "error", err)
		return nil, err
	}

	txHash, err := m.ProcessBitcoinWithdrawTransaction(ctx, msg)
	if err != nil {
		ctx.Logger().Error("Error processing bitcoin withdraw transaction", "error", err)
		return nil, err
	}

	// Emit Events
	m.EmitEvent(ctx, msg.Sender,
		sdk.NewAttribute("blockhash", msg.Blockhash),
		sdk.NewAttribute("txBytes", msg.TxBytes),
		sdk.NewAttribute("txid", txHash.String()),
	)

	return &types.MsgSubmitWithdrawTransactionResponse{}, nil
}

// UpdateSenders implements types.MsgServer.
func (m msgServer) UpdateQualifiedRelayers(goCtx context.Context, msg *types.MsgUpdateQualifiedRelayersRequest) (*types.MsgUpdateQualifiedRelayersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// check if the sender is one of the authorized senders
	param := m.GetParams(ctx)
	if !param.IsAuthorizedSender(msg.Sender) {
		return nil, types.ErrSenderAddressNotAuthorized
	}

	// Set block headers
	m.SetParams(ctx, types.NewParams(msg.Relayers))

	// Emit events

	return &types.MsgUpdateQualifiedRelayersResponse{}, nil
}

func (m msgServer) WithdrawBitcoin(goCtx context.Context, msg *types.MsgWithdrawBitcoinRequest) (*types.MsgWithdrawBitcoinResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender := sdk.MustAccAddressFromBech32(msg.Sender)

	coin, err := sdk.ParseCoinNormalized(msg.Amount)
	if err != nil {
		return nil, err
	}

	if err = m.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		return nil, err
	}

	feeRate, err := strconv.ParseInt(msg.FeeRate, 10, 64)
	if err != nil {
		return nil, err
	}

	req, err := m.Keeper.NewSigningRequest(ctx, msg.Sender, coin, feeRate, "")
	if err != nil {
		return nil, err
	}

	// Emit events
	m.EmitEvent(ctx, msg.Sender,
		sdk.NewAttribute("amount", msg.Amount),
		sdk.NewAttribute("txid", req.Txid),
	)

	return &types.MsgWithdrawBitcoinResponse{}, nil
}

// SubmitWithdrawSignatures submits the signatures of the withdraw transaction.
func (m msgServer) SubmitWithdrawSignatures(goCtx context.Context, msg *types.MsgSubmitWithdrawSignaturesRequest) (*types.MsgSubmitWithdrawSignaturesResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	exist := m.HasSigningRequest(ctx, msg.Txid)
	if !exist {
		return nil, types.ErrSigningRequestNotExist
	}

	b, err := base64.StdEncoding.DecodeString(msg.Psbt)
	if err != nil {
		return nil, types.ErrInvalidSignatures
	}

	packet, err := psbt.NewFromRawBytes(bytes.NewReader(b), false)
	if err != nil {
		return nil, err
	}

	if packet.UnsignedTx.TxHash().String() != msg.Txid {
		return nil, types.ErrInvalidSignatures
	}

	if err = packet.SanityCheck(); err != nil {
		return nil, err
	}
	if !packet.IsComplete() {
		return nil, types.ErrInvalidSignatures
	}

	// verify the signatures
	if !types.VerifyPsbtSignatures(packet) {
		return nil, types.ErrInvalidSignatures
	}

	// Set the signing request status to signed
	request := m.GetSigningRequest(ctx, msg.Txid)
	request.Psbt = msg.Psbt
	request.Status = types.SigningStatus_SIGNING_STATUS_SIGNED
	m.SetSigningRequest(ctx, request)

	return &types.MsgSubmitWithdrawSignaturesResponse{}, nil
}

func (m msgServer) SubmitWithdrawStatus(goCtx context.Context, msg *types.MsgSubmitWithdrawStatusRequest) (*types.MsgSubmitWithdrawStatusResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	param := m.GetParams(sdk.UnwrapSDKContext(goCtx))
	if !param.IsAuthorizedSender(msg.Sender) {
		return nil, types.ErrSenderAddressNotAuthorized
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	exist := m.HasSigningRequest(ctx, msg.Txid)
	if !exist {
		return nil, types.ErrSigningRequestNotExist
	}

	request := m.GetSigningRequest(ctx, msg.Txid)
	request.Status = msg.Status
	m.SetSigningRequest(ctx, request)

	return &types.MsgSubmitWithdrawStatusResponse{}, nil
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}
