package types

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgCreatePool = "create_pool"

var _ sdk.Msg = &MsgCreatePool{}

func NewMsgCreatePool(
	creator string,
	params PoolParams,
	liquidity []PoolAsset,
) *MsgCreatePool {
	return &MsgCreatePool{
		Creator:   creator,
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
		return sdkerrors.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
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
		Creator:     msg.Creator,
		PoolParams:  *msg.Params,
		Assets:      assets,
		TotalShares: sdk.NewCoin(poolShareBaseDenom, totalShares),
	}
	return pool
}
