package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgUpdateSenders = "update_qualifed_relayers"

func NewMsgUpdateSendersRequest(
	sender string,
	relayers []string,
) *MsgUpdateQualifiedRelayersRequest {
	return &MsgUpdateQualifiedRelayersRequest{
		Sender:   sender,
		Relayers: relayers,
	}
}

func (msg *MsgUpdateQualifiedRelayersRequest) Route() string {
	return RouterKey
}

func (msg *MsgUpdateQualifiedRelayersRequest) Type() string {
	return TypeMsgUpdateSenders
}

func (msg *MsgUpdateQualifiedRelayersRequest) GetSigners() []sdk.AccAddress {
	Sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{Sender}
}

func (msg *MsgUpdateQualifiedRelayersRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateQualifiedRelayersRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(err, "invalid sender address (%s)", err)
	}

	if len(msg.Relayers) == 0 {
		return sdkerrors.Wrap(ErrInvalidSenders, "relayers cannot be empty")
	}

	for _, sender := range msg.Relayers {
		_, err := sdk.AccAddressFromBech32(sender)
		if err != nil {
			return sdkerrors.Wrapf(ErrInvalidSenders, "address (%s) is invalid", sender)
		}
	}

	return nil
}
