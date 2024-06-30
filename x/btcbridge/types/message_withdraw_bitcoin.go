package types

import (
	"strconv"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgWithdrawBitcoin = "withdraw_bitcoin"

func NewMsgWithdrawBitcoinRequest(
	sender string,
	amount string,
	feeRate string,
) *MsgWithdrawBitcoinRequest {
	return &MsgWithdrawBitcoinRequest{
		Sender:  sender,
		Amount:  amount,
		FeeRate: feeRate,
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

	_, err = sdk.ParseCoinNormalized(msg.Amount)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAmount, "invalid amount %s", msg.Amount)
	}

	feeRate, err := strconv.ParseInt(msg.FeeRate, 10, 64)
	if err != nil {
		return err
	}

	if feeRate <= 0 {
		return sdkerrors.Wrap(ErrInvalidFeeRate, "fee rate must be greater than zero")
	}

	return nil
}
