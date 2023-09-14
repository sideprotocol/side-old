package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgCreatePool = "create_pool"

var _ sdk.Msg = &MsgCreatePool{}

func NewMsgCreatePool(
	sender string,
	params PoolParams,
	liquidity []PoolAsset,
) *MsgCreatePool {
	return &MsgCreatePool{
		Sender:    sender,
		Params:    &params,
		Liquidity: liquidity,
	}
}

func (msg *MsgCreatePool) Route() string {
	return RouterKey
}

func (msg *MsgCreatePool) Type() string {
	return TypeMsgCreatePool
}

func (msg *MsgCreatePool) GetSigners() []sdk.AccAddress {
	Sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{Sender}
}

func (msg *MsgCreatePool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreatePool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid Sender address (%s)", err)
	}
	if msg.Params == nil {
		return ErrInvalidPoolParams
	}
	if msg.Liquidity == nil {
		return ErrEmptyLiquidity
	}

	if len(msg.Liquidity) != 2 {
		return sdkerrors.Wrapf(ErrInvalidLiquidityInLength, "number of liquidity (%d)", len(msg.Liquidity))
	}

	if msg.Params.Amp.GT(sdkmath.NewInt(100)) {
		return sdkerrors.Wrapf(ErrInvalidAmp, "amp (%d) is out of range", msg.Params.Amp)
	}
	totalWeight := sdkmath.NewInt(0)
	for _, asset := range msg.Liquidity {
		totalWeight = totalWeight.Add(*asset.Weight)
	}
	return nil
}

func (msg *MsgCreatePool) GetPoolType() PoolType {
	return msg.Params.Type
}

// The Sender of the pool, who pays the PoolCreationFee, provides initial liquidity,
// and gets the initial LP shares.
func (msg *MsgCreatePool) PoolCreator() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(msg.Sender)
}

// Initial Liquidity for the pool that the sender is required to send to the pool account
func (msg *MsgCreatePool) InitialLiquidity() sdk.Coins {
	liquidity := sdk.NewCoins()
	for _, asset := range msg.Liquidity {
		liquidity = liquidity.Add(asset.Token)
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

// Return denom list of liquidity
func (msg *MsgCreatePool) CreatePool() Pool {
	// Extract denom list from Liquidity
	denoms := msg.GetAssetDenoms()

	assets := make(map[string]PoolAsset)
	totalShares := sdk.NewInt(0)
	for _, liquidity := range msg.Liquidity {
		assets[liquidity.Token.Denom] = liquidity
		totalShares = totalShares.Add(liquidity.Token.Amount)
	}

	// Generate new PoolId
	newPoolID := GetPoolID(denoms)
	poolShareBaseDenom := GetPoolShareDenom(newPoolID)
	pool := Pool{
		PoolId:      newPoolID,
		Sender:      msg.Sender,
		PoolParams:  *msg.Params,
		Assets:      assets,
		TotalShares: sdk.NewCoin(poolShareBaseDenom, totalShares),
	}
	return pool
}
