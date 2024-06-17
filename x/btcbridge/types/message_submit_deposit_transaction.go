package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSubmitDepositTransaction = "submit_deposit_transaction"

func NewMsgSubmitTransactionRequest(
	sender string,
	blockhash string,
	transaction string,
	proof []string,
) *MsgSubmitDepositTransactionRequest {
	return &MsgSubmitDepositTransactionRequest{
		Sender:    sender,
		Blockhash: blockhash,
		TxBytes:   transaction,
		Proof:     proof,
	}
}

func (msg *MsgSubmitDepositTransactionRequest) Route() string {
	return RouterKey
}

func (msg *MsgSubmitDepositTransactionRequest) Type() string {
	return TypeMsgSubmitDepositTransaction
}

func (msg *MsgSubmitDepositTransactionRequest) GetSigners() []sdk.AccAddress {
	Sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{Sender}
}

func (msg *MsgSubmitDepositTransactionRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSubmitDepositTransactionRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(err, "invalid Sender address (%s)", err)
	}

	if len(msg.Blockhash) == 0 {
		return sdkerrors.Wrap(ErrInvalidBtcTransaction, "blockhash cannot be empty")
	}

	if len(msg.PrevTxBytes) == 0 {
		return sdkerrors.Wrap(ErrInvalidBtcTransaction, "transaction cannot be empty")
	}

	if len(msg.TxBytes) == 0 {
		return sdkerrors.Wrap(ErrInvalidBtcTransaction, "transaction cannot be empty")
	}

	if len(msg.Proof) == 0 {
		return sdkerrors.Wrap(ErrInvalidBtcTransaction, "proof cannot be empty")
	}

	return nil
}
