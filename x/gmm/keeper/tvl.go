package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/gmm/types"
)

func (k Keeper) ObserveTVLFromPool(ctx sdk.Context, poolID string, assets sdk.Coins) error {
	// Get Asset from pool
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPoolTVLPrefix)
	rawTVLStack := store.Get([]byte(poolID))
	tvlStack := types.NewVolumeStack()
	err := tvlStack.Decode(rawTVLStack)
	if err != nil {
		return err
	}

	tvlStack.Observe(ctx, poolID, assets)
	result, err := tvlStack.Encode()
	if err != nil {
		return err
	}
	store.Set([]byte(poolID), result)
	return nil
}

// func (k Keeper) GetVolume24(ctx sdk.Context, poolID string) []sdk.Coin {
// 	// Get Asset from pool
// 	_, found := k.GetPool(ctx, poolID)
// 	if !found {
// 		return []sdk.Coin{}
// 	}
// 	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPoolVolumePrefix)
// 	rawVolumeStack := store.Get([]byte(poolID))
// 	var volumeStack types.VolumeStack
// 	err := volumeStack.Decode(rawVolumeStack)
// 	if err != nil {
// 		return []sdk.Coin{}
// 	}
// 	return volumeStack.Calculate24HourVolume(ctx, poolID)
// }

// func (k Keeper) GetTotalVolume(ctx sdk.Context, poolID string) []sdk.Coin {
// 	// Get Asset from pool
// 	_, found := k.GetPool(ctx, poolID)
// 	if !found {
// 		return []sdk.Coin{}
// 	}
// 	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPoolVolumePrefix)
// 	rawVolumeStack := store.Get([]byte(poolID))
// 	var volumeStack types.VolumeStack
// 	err := volumeStack.Decode(rawVolumeStack)
// 	if err != nil {
// 		return []sdk.Coin{}
// 	}
// 	return volumeStack.GetTotalVolume()
// }
