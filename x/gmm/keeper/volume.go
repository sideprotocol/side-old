package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/gmm/types"
)

func (k Keeper) ObserveVolumeFromPool(ctx sdk.Context, pool types.Pool) error {
	// Get Asset from pool
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPoolVolumePrefix)
	rawVolumeStack := store.Get([]byte(pool.PoolId))
	var volumeStack types.VolumeStack
	err := volumeStack.Decode(rawVolumeStack)
	if err != nil {
		return err
	}

	assets := pool.GetAssetList()
	for _, asset := range assets {
		// update the volume information.
		volumeStack.PushOrUpdate(ctx, pool.PoolId, asset.Token.Denom, asset.Token.Amount)
	}
	result, err := volumeStack.Encode()
	if err != nil {
		return err
	}
	store.Set([]byte(pool.PoolId), result)
	return nil
}

func (k Keeper) GetVolume24(ctx sdk.Context, poolID string) error {
	// Get Asset from pool
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPoolVolumePrefix)
	rawVolumeStack := store.Get([]byte(poolID))
	var volumeStack types.VolumeStack
	err := volumeStack.Decode(rawVolumeStack)
	if err != nil {
		return err
	}
	if len(volumeStack.Data) == 0 {
		return nil
	}
	
	if volumeStack.Data[len(volumeStack.Data)-1].BlockHeight <
	return nil
}
