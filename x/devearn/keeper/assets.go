package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sidechain/x/devearn/types"
)

// GetAssetsCount get the total number of assets
func (k Keeper) GetAssetsCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.AssetsCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetAssetsCount set the total number of assets
func (k Keeper) SetAssetsCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.AssetsCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendAssets appends a assets in the store with a new id and update the count
func (k Keeper) AppendAssets(
	ctx sdk.Context,
	assets types.Assets,
) uint64 {
	// Create the assets
	count := k.GetAssetsCount(ctx)

	// Set the ID of the appended value
	assets.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AssetsKey))
	appendedValue := k.cdc.MustMarshal(&assets)
	store.Set(GetAssetsIDBytes(assets.Id), appendedValue)

	// Update assets count
	k.SetAssetsCount(ctx, count+1)

	return count
}

// SetAssets set a specific assets in the store
func (k Keeper) SetAssets(ctx sdk.Context, assets types.Assets) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AssetsKey))
	b := k.cdc.MustMarshal(&assets)
	store.Set(GetAssetsIDBytes(assets.Id), b)
}

// GetAssets returns a assets from its id
func (k Keeper) GetAssets(ctx sdk.Context, id uint64) (val types.Assets, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AssetsKey))
	b := store.Get(GetAssetsIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveAssets removes a assets from the store
func (k Keeper) RemoveAssets(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AssetsKey))
	store.Delete(GetAssetsIDBytes(id))
}

// GetAllAssets returns all assets
func (k Keeper) GetAllAssets(ctx sdk.Context) (list []types.Assets) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AssetsKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Assets
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetAssetsIDBytes returns the byte representation of the ID
func GetAssetsIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetAssetsIDFromBytes returns ID in uint64 format from a byte array
func GetAssetsIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
