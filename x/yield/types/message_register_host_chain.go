package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgRegisterHostChain = "register_host_chain"

var _ sdk.Msg = &MsgRegisterHostChain{}

func NewMsgRegisterHostChain(
	creator string,
	connectionId string,
	bech32prefix string,
	hostDenom string,
	ibcDenom string,
	transferChannelId string,
) *MsgRegisterHostChain {
	return &MsgRegisterHostChain{
		Creator:           creator,
		ConnectionId:      connectionId,
		Bech32Prefix:      bech32prefix,
		HostDenom:         hostDenom,
		IbcDenom:          ibcDenom,
		TransferChannelId: transferChannelId,
	}
}

func (msg *MsgRegisterHostChain) Route() string {
	return RouterKey
}

func (msg *MsgRegisterHostChain) Type() string {
	return TypeMsgRegisterHostChain
}

func (msg *MsgRegisterHostChain) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRegisterHostChain) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRegisterHostChain) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
