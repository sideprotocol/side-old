package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSubmitBlockHeader = "submit_block_header"

func NewMsgSubmitBlockHeaderRequest(
	sender string,
	headers []*BlockHeader,
) *MsgSubmitBlockHeaderRequest {
	return &MsgSubmitBlockHeaderRequest{
		Sender:       sender,
		BlockHeaders: headers,
	}
}

func (msg *MsgSubmitBlockHeaderRequest) Route() string {
	return RouterKey
}

func (msg *MsgSubmitBlockHeaderRequest) Type() string {
	return TypeMsgSubmitBlockHeader
}

func (msg *MsgSubmitBlockHeaderRequest) GetSigners() []sdk.AccAddress {
	Sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{Sender}
}

func (msg *MsgSubmitBlockHeaderRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSubmitBlockHeaderRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(err, "invalid Sender address (%s)", err)
	}

	if len(msg.BlockHeaders) == 0 {
		return sdkerrors.Wrap(ErrInvalidHeader, "block headers cannot be empty")
	}

	return nil
}
