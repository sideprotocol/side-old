package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgSubmitBlockHeaderRequest{}, "btcbridge/MsgSubmitBlockHeaderRequest", nil)
	cdc.RegisterConcrete(&MsgUpdateQualifiedRelayersRequest{}, "btcbridge/MsgUpdateQualifiedRelayersRequest", nil)
	cdc.RegisterConcrete(&MsgSubmitDepositTransactionRequest{}, "btcbridge/MsgSubmitDepositTransactionRequest", nil)
	cdc.RegisterConcrete(&MsgSubmitWithdrawSignaturesRequest{}, "btcbridge/MsgSubmitWithdrawSignaturesRequest", nil)
	cdc.RegisterConcrete(&MsgSubmitWithdrawTransactionRequest{}, "btcbridge/MsgSubmitWithdrawTransactionRequest", nil)
	cdc.RegisterConcrete(&MsgWithdrawBitcoinRequest{}, "btcbridge/MsgWithdrawBitcoinRequest", nil)
	cdc.RegisterConcrete(&MsgSubmitWithdrawStatusRequest{}, "btcbridge/MsgSubmitWithdrawStatusRequest", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgSubmitBlockHeaderRequest{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgUpdateQualifiedRelayersRequest{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgSubmitDepositTransactionRequest{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgSubmitWithdrawSignaturesRequest{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgSubmitWithdrawTransactionRequest{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgWithdrawBitcoinRequest{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgSubmitWithdrawStatusRequest{})
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
