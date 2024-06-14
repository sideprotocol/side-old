package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSubmitSignature = "submit_signature"

func NewMsgSubmitWithdrawSignaturesRequest(
	sender string,
	txid string,
	signatures []string,
) *MsgSubmitWithdrawSignaturesRequest {
	return &MsgSubmitWithdrawSignaturesRequest{
		Sender:     sender,
		Txid:       txid,
		Signatures: signatures,
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

	if len(msg.Signatures) == 0 {
		return sdkerrors.Wrap(ErrInvalidSignatures, "sigatures cannot be empty")
	}

	for _, signature := range msg.Signatures {
		if len(signature) == 0 {
			return sdkerrors.Wrap(ErrInvalidSignatures, "signature cannot be empty")

		}
	}

	return nil
}
