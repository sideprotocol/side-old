package types

import (
	"strings"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkmath "cosmossdk.io/math"
)

const TypeMsgSwap = "swap"

var _ sdk.Msg = &MsgSwap{}

func NewMsgSwap(
	sender, poolID string,
	tokenIn sdk.Coin,
	tokenOut sdk.Coin,
) *MsgSwap {
	return &MsgSwap{
		Sender:   sender,
		PoolId:   poolID,
		TokenIn:  tokenIn,
		TokenOut: tokenOut,
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

	if msg.TokenOut.Amount.IsZero() {
		return sdkerrors.Wrap(ErrInvalidTokenAmount, "tokenOut amount cannot be zero")
	}

	if msg.Slippage.IsNegative() || msg.Slippage.GTE(sdkmath.NewInt(100)) {
		return sdkerrors.Wrap(ErrInvalidSlippage, "slippage should be ranged from 0 to 100")
	}
	return nil
}
