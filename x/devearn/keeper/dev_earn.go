package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"sidechain/x/devearn/types"
)

// GetAllDevEarnInfos - get  all registered DevEarnInfo from the identifier
func (k Keeper) GetAllDevEarnInfos(ctx sdk.Context) []types.DevEarnInfo {
	devEarnInfos := []types.DevEarnInfo{}

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixDevEarn)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var devEarnInfo types.DevEarnInfo
		k.cdc.MustUnmarshal(iterator.Value(), &devEarnInfo)
		devEarnInfos = append(devEarnInfos, devEarnInfo)
	}

	return devEarnInfos
}

// GetDevEarnInfo - get registered DevEarnInfo from the identifier
func (k Keeper) GetDevEarnInfo(
	ctx sdk.Context,
	contract common.Address,
) (types.DevEarnInfo, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixDevEarn)
	bz := store.Get(contract.Bytes())
	if len(bz) == 0 {
		return types.DevEarnInfo{}, false
	}

	var devEarnInfo types.DevEarnInfo
	k.cdc.MustUnmarshal(bz, &devEarnInfo)
	return devEarnInfo, true
}

// SetDevEarnInfo stores an dev_earn info
func (k Keeper) SetDevEarnInfo(ctx sdk.Context, devEarnInfo types.DevEarnInfo) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixDevEarn)
	key := common.HexToAddress(devEarnInfo.Contract)
	bz := k.cdc.MustMarshal(&devEarnInfo)
	store.Set(key.Bytes(), bz)
}

func (k Keeper) ResetGasMeter(ctx sdk.Context, contract common.Address) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixDevEarn)
	bz := store.Get(contract.Bytes())
	if len(bz) == 0 {
		return false
	}
	var devEarnInfo types.DevEarnInfo
	k.cdc.MustUnmarshal(bz, &devEarnInfo)
	devEarnInfo.GasMeter = 0
	bz = k.cdc.MustMarshal(&devEarnInfo)
	store.Set(contract.Bytes(), bz)
	return true
}

// DeletDevEarnInfo removes an DevEarnInfo
func (k Keeper) DeleteDevEarnInfo(ctx sdk.Context, devEarnInfo types.DevEarnInfo) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixDevEarn)
	key := common.HexToAddress(devEarnInfo.Contract)
	store.Delete(key.Bytes())
}

// IsDevEarnInfoRegistered - check if registered dev earn info is registered
func (k Keeper) IsDevEarnInfoRegistered(
	ctx sdk.Context,
	contract common.Address,
) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixDevEarn)
	return store.Has(contract.Bytes())
}

// IterateDevEarnInfos iterates over all registered `Incentives` and performs a
// callback.
func (k Keeper) IterateDevEarnInfos(
	ctx sdk.Context,
	handlerFn func(devEarnInfo types.DevEarnInfo) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixDevEarn)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var devEarnInfo types.DevEarnInfo
		k.cdc.MustUnmarshal(iterator.Value(), &devEarnInfo)

		if handlerFn(devEarnInfo) {
			break
		}
	}
}
