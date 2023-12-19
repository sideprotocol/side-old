package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/sideprotocol/side/x/yield/types"
)

// SetHostChain set a specific hostChain in the store
func (k Keeper) SetHostChain(ctx sdk.Context, hostChain types.HostChain) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.HostChainKey))
	b := k.cdc.MustMarshal(&hostChain)
	store.Set([]byte(hostChain.ChainId), b)
}

// GetHostChain returns a hostChain from its id
func (k Keeper) GetHostChain(ctx sdk.Context, chainID string) (val types.HostChain, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.HostChainKey))
	b := store.Get([]byte(chainID))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetHostChainFromHostDenom returns a HostChain from a HostDenom
func (k Keeper) GetHostChainFromHostDenom(ctx sdk.Context, denom string) (*types.HostChain, error) {
	var matchChain types.HostChain
	k.IterateHostChains(ctx, func(ctx sdk.Context, index int64, chainInfo types.HostChain) error {
		if chainInfo.HostDenom == denom {
			matchChain = chainInfo
			return nil
		}
		return nil
	})
	if matchChain.ChainId != "" {
		return &matchChain, nil
	}
	return nil, errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "No HostChain for %s found", denom)
}

// GetHostChainFromTransferChannelID returns a HostChain from a transfer channel ID
func (k Keeper) GetHostChainFromTransferChannelID(ctx sdk.Context, channelID string) (hostChain types.HostChain, found bool) {
	for _, hostChain := range k.GetAllHostChain(ctx) {
		if hostChain.TransferChannelId == channelID {
			return hostChain, true
		}
	}
	return types.HostChain{}, false
}

// RemoveHostChain removes a hostChain from the store
func (k Keeper) RemoveHostChain(ctx sdk.Context, chainID string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.HostChainKey))
	store.Delete([]byte(chainID))
}

// GetAllHostChain returns all hostChain
func (k Keeper) GetAllHostChain(ctx sdk.Context) (list []types.HostChain) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.HostChainKey))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.HostChain
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// IterateHostChains iterates chains
func (k Keeper) IterateHostChains(ctx sdk.Context, fn func(ctx sdk.Context, index int64, zoneInfo types.HostChain) error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.HostChainKey))

	iterator := sdk.KVStorePrefixIterator(store, nil)
	defer iterator.Close()

	i := int64(0)

	for ; iterator.Valid(); iterator.Next() {
		k.Logger(ctx).Debug(fmt.Sprintf("Iterating hostChain %d", i))
		zone := types.HostChain{}
		k.cdc.MustUnmarshal(iterator.Value(), &zone)

		error := fn(ctx, i, zone)

		if error != nil {
			break
		}
		i++
	}
}
