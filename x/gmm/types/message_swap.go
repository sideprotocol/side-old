package types

import (
	"strings"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSwap = "swap"

var _ sdk.Msg = &MsgSwap{}

func NewMsgSwap(
	Sender, poolID string,
	tokenIn sdk.Coin,
	denomOut string,
) *MsgSwap {
	return &MsgSwap{
		Sender:  Sender,
		PoolId:   poolID,
		TokenIn:  tokenIn,
		DenomOut: denomOut,
	}
}

func (msg *MsgSwap) Route() string {
	return RouterKey
}

func (msg *MsgSwap) Type() string {
	return TypeMsgSwap
}

func (msg *MsgSwap) GetSigners() []sdk.AccAddress {
	Sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{Sender}
}

func (msg *MsgSwap) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSwap) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid Sender address (%s)", err)
	}

	if strings.TrimSpace(msg.PoolId) == "" {
		return sdkerrors.Wrap(ErrInvalidPoolID, "pool id cannot be empty")
	}

	if msg.TokenIn.Amount.IsZero() {
		return sdkerrors.Wrap(ErrInvalidTokenAmount, "tokenIn amount cannot be zero")
	}

	if strings.TrimSpace(msg.DenomOut) == "" {
		return sdkerrors.Wrap(ErrEmptyDenom, "denom_out cannot be empty")
	}

	return nil
}
