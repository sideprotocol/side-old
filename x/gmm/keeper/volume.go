package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/gmm/types"
)

func (k Keeper) ObserveVolumeFromPool(ctx sdk.Context, poolID string, tokenIn, tokenOut sdk.Coin) error {
	// Get Asset from pool
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPoolVolumePrefix)
	rawVolumeStack := store.Get([]byte(poolID))
	volumeStack := types.NewVolumeStack()
	err := volumeStack.Decode(rawVolumeStack)
	if err != nil {
		return err
	}
	volumeStack.Observe(ctx, poolID, sdk.NewCoins(tokenIn, tokenOut))
	result, err := volumeStack.Encode()
	if err != nil {
		return err
	}
	store.Set([]byte(poolID), result)
	return nil
}

func (k Keeper) GetVolume24(ctx sdk.Context, poolID string) []sdk.Coin {
	// Get Asset from pool
	_, found := k.GetPool(ctx, poolID)
	if !found {
		return []sdk.Coin{}
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPoolVolumePrefix)
	rawVolumeStack := store.Get([]byte(poolID))
	var volumeStack types.VolumeStack
	err := volumeStack.Decode(rawVolumeStack)
	if err != nil {
		return []sdk.Coin{}
	}
	return volumeStack.Calculate24HourVolume(ctx, poolID)
}

func (k Keeper) GetTotalVolume(ctx sdk.Context, poolID string) []sdk.Coin {
	// Get Asset from pool
	_, found := k.GetPool(ctx, poolID)
	if !found {
		return []sdk.Coin{}
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPoolVolumePrefix)
	rawVolumeStack := store.Get([]byte(poolID))
	var volumeStack types.VolumeStack
	err := volumeStack.Decode(rawVolumeStack)
	if err != nil {
		return []sdk.Coin{}
	}
	return volumeStack.GetTotalVolume()
}
