package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgWithdraw = "withdraw"

var _ sdk.Msg = &MsgWithdraw{}

func NewMsgWithdraw(
	creator, poolID, receiver string,
	share sdk.Coin,
) *MsgWithdraw {
	return &MsgWithdraw{
		Creator:  creator,
		Receiver: receiver,
		PoolId:   poolID,
		Share:    share,
	}
}

func (msg *MsgWithdraw) Route() string {
	return RouterKey
}

func (msg *MsgWithdraw) Type() string {
	return TypeMsgWithdraw
}

func (msg *MsgWithdraw) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgWithdraw) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWithdraw) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
