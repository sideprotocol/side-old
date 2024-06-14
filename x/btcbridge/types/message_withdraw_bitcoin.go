package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgWithdrawBitcoin = "withdraw_bitcoin"

func NewMsgWithdrawBitcoinRequest(
	sender string,
	amount string,
) *MsgWithdrawBitcoinRequest {
	return &MsgWithdrawBitcoinRequest{
		Sender: sender,
		Amount: amount,
	}
}

func (msg *MsgWithdrawBitcoinRequest) Route() string {
	return RouterKey
}

func (msg *MsgWithdrawBitcoinRequest) Type() string {
	return TypeMsgWithdrawBitcoin
}

func (msg *MsgWithdrawBitcoinRequest) GetSigners() []sdk.AccAddress {
	Sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{Sender}
}

func (msg *MsgWithdrawBitcoinRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWithdrawBitcoinRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(err, "invalid Sender address (%s)", err)
	}

	if len(msg.Amount) == 0 {
		return sdkerrors.Wrap(sdk.ErrInvalidLengthCoin, "amount cannot be empty")
	}

	return nil
}
