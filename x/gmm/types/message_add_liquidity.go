package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgAddLiquidity = "add_liquidity"

var _ sdk.Msg = &MsgAddLiquidity{}

func NewMsgAddLiquidity(
	creator, poolID string,
	liquidity sdk.Coins,
) *MsgAddLiquidity {
	return &MsgAddLiquidity{
		Creator:   creator,
		PoolId:    poolID,
		Liquidity: liquidity,
	}
}

func (msg *MsgAddLiquidity) Route() string {
	return RouterKey
}

func (msg *MsgAddLiquidity) Type() string {
	return TypeMsgAddLiquidity
}

func (msg *MsgAddLiquidity) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddLiquidity) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddLiquidity) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
