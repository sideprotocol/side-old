package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSubmitTransaction = "submit_transaction"

func NewMsgSubmitTransactionRequest(
	sender string,
	transaction string,
	proof string,
) *MsgSubmitTransactionRequest {
	return &MsgSubmitTransactionRequest{
		Sender: sender,
		Tx:     transaction,
		Proof:  proof,
	}
}

func (msg *MsgSubmitTransactionRequest) Route() string {
	return RouterKey
}

func (msg *MsgSubmitTransactionRequest) Type() string {
	return TypeMsgSubmitTransaction
}

func (msg *MsgSubmitTransactionRequest) GetSigners() []sdk.AccAddress {
	Sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{Sender}
}

func (msg *MsgSubmitTransactionRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSubmitTransactionRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(err, "invalid Sender address (%s)", err)
	}

	if len(msg.Tx) == 0 {
		return sdkerrors.Wrap(ErrInvalidBtcTransaction, "transaction cannot be empty")
	}

	if len(msg.Proof) == 0 {
		return sdkerrors.Wrap(ErrInvalidBtcTransaction, "proof cannot be empty")
	}

	return nil
}
