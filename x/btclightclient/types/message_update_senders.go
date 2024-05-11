package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgUpdateSenders = "update_senders"

func NewMsgUpdateSendersRequest(
	sender string,
	senders []string,
) *MsgUpdateSendersRequest {
	return &MsgUpdateSendersRequest{
		Sender:  sender,
		Senders: senders,
	}
}

func (msg *MsgUpdateSendersRequest) Route() string {
	return RouterKey
}

func (msg *MsgUpdateSendersRequest) Type() string {
	return TypeMsgUpdateSenders
}

func (msg *MsgUpdateSendersRequest) GetSigners() []sdk.AccAddress {
	Sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{Sender}
}

func (msg *MsgUpdateSendersRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateSendersRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(err, "invalid Sender address (%s)", err)
	}

	if len(msg.Senders) == 0 {
		return sdkerrors.Wrap(ErrInvalidSenders, "senders cannot be empty")
	}

	for _, sender := range msg.Senders {
		_, err := sdk.AccAddressFromBech32(sender)
		if err != nil {
			return sdkerrors.Wrapf(ErrInvalidSenders, "address (%s) is invalid", sender)
		}
	}

	return nil
}
