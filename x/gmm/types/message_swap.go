package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgSwap = "swap"

var _ sdk.Msg = &MsgSwap{}

func NewMsgSwap(
	creator, poolID string,
	tokenIn sdk.Coin,
	denomOut string,
) *MsgSwap {
	return &MsgSwap{
		Creator:  creator,
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
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSwap) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSwap) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
