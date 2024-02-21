package types

import (
	"strings"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgAddLiquidity = "add_liquidity"

var _ sdk.Msg = &MsgAddLiquidity{}

func NewMsgAddLiquidity(
	sender, poolID string,
	liquidity sdk.Coins,
) *MsgAddLiquidity {
	return &MsgAddLiquidity{
		Sender:    sender,
		PoolId:    poolID,
		Liquidity: liquidity.Sort(),
	}
}

func (msg *MsgAddLiquidity) Route() string {
	return RouterKey
}

func (msg *MsgAddLiquidity) Type() string {
	return TypeMsgAddLiquidity
}

func (msg *MsgAddLiquidity) GetSigners() []sdk.AccAddress {
	Sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{Sender}
}

func (msg *MsgAddLiquidity) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddLiquidity) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid Sender address (%s)", err)
	}
	if strings.TrimSpace(msg.PoolId) == "" {
		return sdkerrors.Wrap(ErrInvalidPoolID, "pool id cannot be empty")
	}
	if len(msg.Liquidity) > 2 || len(msg.Liquidity) == 0 {
		return sdkerrors.Wrap(ErrInvalidLiquidityInLength, "liquidity cannot be empty or cannot be more than 2")
	}

	for _, asset := range msg.Liquidity {
		if asset.Amount.IsZero() {
			return sdkerrors.Wrap(ErrInvalidLiquidityAmount, "liquidity amount cannot be zero")
		}
	}
	return nil
}
