package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/gmm/types"
)

func (k Keeper) ObserveFeeFromPool(ctx sdk.Context, poolID string, fee sdk.Coin) error {
	// Get Asset from pool
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPoolAPRPrefix)
	rawAprData := store.Get([]byte(poolID))
	var aprData types.PoolAPR
	if err := aprData.Decode(ctx, rawAprData); err != nil {
		return err
	}
	aprData.Fees = aprData.Fees.Add(fee)
	data, err := aprData.Encode()
	if err != nil {
		return err
	}
	store.Set([]byte(poolID), data)
	return nil
}

func (k Keeper) GetAPR(ctx sdk.Context, poolID string) []sdk.Coin {
	// Get Asset from pool
	pool, found := k.GetPool(ctx, poolID)
	if !found {
		return []sdk.Coin{}
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPoolAPRPrefix)
	rawAPR := store.Get([]byte(poolID))
	var poolAPR types.PoolAPR
	err := poolAPR.Decode(ctx, rawAPR)
	if err != nil {
		return []sdk.Coin{}
	}
	return poolAPR.CalcAPR(ctx, pool.Assets)
}
