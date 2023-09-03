package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCreatePool = "create_pool"

var _ sdk.Msg = &MsgCreatePool{}

func NewMsgCreatePool(creator string) *MsgCreatePool {
	return &MsgCreatePool{
		Creator: creator,
	}
}

func (msg *MsgCreatePool) Route() string {
	return RouterKey
}

func (msg *MsgCreatePool) Type() string {
	return TypeMsgCreatePool
}

func (msg *MsgCreatePool) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreatePool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreatePool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

func (msg *MsgCreatePool) GetPoolType() PoolType {
	return msg.Params.Type
}

// The creator of the pool, who pays the PoolCreationFee, provides initial liquidity,
// and gets the initial LP shares.
func (msg *MsgCreatePool) PoolCreator() sdk.AccAddress {
	return sdk.AccAddress(msg.Creator)
}

// A stateful validation function.
func (msg *MsgCreatePool) Validate(ctx sdk.Context) error {
	return msg.ValidateBasic()
}

// Initial Liquidity for the pool that the sender is required to send to the pool account
func (msg *MsgCreatePool) InitialLiquidity() sdk.Coins {
	liquidity := sdk.Coins{}
	for _, asset := range msg.Liquidity {
		liquidity.Add(asset.Token)
	}
	return liquidity
}

// Return denom list of liquidity
func (msg *MsgCreatePool) GetAssetDenoms() []string {
	denoms := []string{}
	for _, asset := range msg.Liquidity {
		denoms = append(denoms, asset.Token.Denom)
	}
	return denoms
}
