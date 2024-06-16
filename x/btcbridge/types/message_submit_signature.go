package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSubmitSignature = "submit_signature"

func NewMsgSubmitWithdrawSignaturesRequest(
	sender string,
	txid string,
	pbst string,
) *MsgSubmitWithdrawSignaturesRequest {
	return &MsgSubmitWithdrawSignaturesRequest{
		Sender: sender,
		Txid:   txid,
		Psbt:   pbst,
	}
}

func (msg *MsgSubmitWithdrawSignaturesRequest) Route() string {
	return RouterKey
}

func (msg *MsgSubmitWithdrawSignaturesRequest) Type() string {
	return TypeMsgSubmitSignature
}

func (msg *MsgSubmitWithdrawSignaturesRequest) GetSigners() []sdk.AccAddress {
	Sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{Sender}
}

func (msg *MsgSubmitWithdrawSignaturesRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSubmitWithdrawSignaturesRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(err, "invalid Sender address (%s)", err)
	}

	if len(msg.Txid) == 0 {
		return sdkerrors.Wrap(ErrInvalidSignatures, "txid cannot be empty")
	}

	if len(msg.Psbt) == 0 {
		return sdkerrors.Wrap(ErrInvalidSignatures, "sigatures cannot be empty")
	}

	return nil
}
