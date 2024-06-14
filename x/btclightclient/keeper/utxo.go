package keeper

import (
	"math/big"
	"sort"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sideprotocol/side/x/btclightclient/types"
)

type UTXOViewKeeper interface {
	HasUTXO(ctx sdk.Context, hash string, vout uint64) bool
	IsUTXOLocked(ctx sdk.Context, hash string, vout uint64) bool

	GetUTXO(ctx sdk.Context, hash string, vout uint64) *types.UTXO
	GetAllUTXOs(ctx sdk.Context) []*types.UTXO

	GetUTXOsByAddr(ctx sdk.Context, addr string) []*types.UTXO
	GetUnlockedUTXOsByAddr(ctx sdk.Context, addr string) []*types.UTXO
	GetOrderedUTXOsByAddr(ctx sdk.Context, addr string) []*types.UTXO

	IterateAllUTXOs(ctx sdk.Context, cb func(utxo *types.UTXO) (stop bool))
	IterateUTXOsByAddr(ctx sdk.Context, addr string, cb func(addr string, utxo *types.UTXO) (stop bool))
}

type UTXOKeeper interface {
	UTXOViewKeeper

	LockUTXO(ctx sdk.Context, hash string, vout uint64) error
	LockUTXOs(ctx sdk.Context, utxos []*types.UTXO) error

	UnlockUTXO(ctx sdk.Context, hash string, vout uint64) error
	UnlockUTXOs(ctx sdk.Context, utxos []*types.UTXO) error

	SpendUTXO(ctx sdk.Context, hash string, vout uint64) error
	SpendUTXOs(ctx sdk.Context, utxos []*types.UTXO) error
}

var _ UTXOKeeper = (*BaseUTXOKeeper)(nil)

type BaseUTXOViewKeeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
}

func NewBaseUTXOViewKeeper(cdc codec.BinaryCodec, storeKey storetypes.StoreKey) *BaseUTXOViewKeeper {
	return &BaseUTXOViewKeeper{
		cdc,
		storeKey,
	}
}

func (bvk *BaseUTXOViewKeeper) HasUTXO(ctx sdk.Context, hash string, vout uint64) bool {
	store := ctx.KVStore(bvk.storeKey)
	return store.Has(types.BtcUtxoKey(hash, vout))
}

// IsUTXOLocked returns true if the given utxo is locked, false otherwise.
// Note: it returns false if the given utxo does not exist.
func (bvk *BaseUTXOViewKeeper) IsUTXOLocked(ctx sdk.Context, hash string, vout uint64) bool {
	if !bvk.HasUTXO(ctx, hash, vout) {
		return false
	}

	utxo := bvk.GetUTXO(ctx, hash, vout)

	return utxo.IsLocked
}

func (bvk *BaseUTXOViewKeeper) GetUTXO(ctx sdk.Context, hash string, vout uint64) *types.UTXO {
	store := ctx.KVStore(bvk.storeKey)

	var utxo types.UTXO
	bz := store.Get(types.BtcUtxoKey(hash, vout))
	bvk.cdc.MustUnmarshal(bz, &utxo)

	return &utxo
}

func (bvk *BaseUTXOViewKeeper) GetAllUTXOs(ctx sdk.Context) []*types.UTXO {
	utxos := make([]*types.UTXO, 0)

	bvk.IterateAllUTXOs(ctx, func(utxo *types.UTXO) (stop bool) {
		utxos = append(utxos, utxo)
		return false
	})

	return utxos
}

func (bvk *BaseUTXOViewKeeper) GetUTXOsByAddr(ctx sdk.Context, addr string) []*types.UTXO {
	utxos := make([]*types.UTXO, 0)

	bvk.IterateUTXOsByAddr(ctx, addr, func(addr string, utxo *types.UTXO) (stop bool) {
		utxos = append(utxos, utxo)
		return false
	})

	return utxos
}

func (bvk *BaseUTXOViewKeeper) GetUnlockedUTXOsByAddr(ctx sdk.Context, addr string) []*types.UTXO {
	utxos := make([]*types.UTXO, 0)

	bvk.IterateUTXOsByAddr(ctx, addr, func(addr string, utxo *types.UTXO) (stop bool) {
		if !utxo.IsLocked {
			utxos = append(utxos, utxo)
		}

		return false
	})

	return utxos
}

// GetOrderedUTXOsByAddr gets all unlocked utxos of the given address in the descending order by amount
func (bvk *BaseUTXOViewKeeper) GetOrderedUTXOsByAddr(ctx sdk.Context, addr string) []*types.UTXO {
	utxos := make([]*types.UTXO, 0)

	bvk.IterateUTXOsByAddr(ctx, addr, func(addr string, utxo *types.UTXO) (stop bool) {
		if !utxo.IsLocked {
			utxos = append(utxos, utxo)
		}

		return false
	})

	// sort utxos in the descending order
	sort.SliceStable(utxos, func(i int, j int) bool {
		return utxos[i].Amount > utxos[j].Amount
	})

	return utxos
}

func (bvk *BaseUTXOViewKeeper) IterateAllUTXOs(ctx sdk.Context, cb func(utxo *types.UTXO) (stop bool)) {
	store := ctx.KVStore(bvk.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.BtcUtxoKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var utxo types.UTXO
		bvk.cdc.MustUnmarshal(iterator.Value(), &utxo)

		if cb(&utxo) {
			break
		}
	}
}

func (bvk *BaseUTXOViewKeeper) IterateUTXOsByAddr(ctx sdk.Context, addr string, cb func(addr string, utxo *types.UTXO) (stop bool)) {
	store := ctx.KVStore(bvk.storeKey)

	keyPrefix := append(types.BtcOwnerUtxoKeyPrefix, []byte(addr)...)
	iterator := sdk.KVStorePrefixIterator(store, keyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()

		hash := key[1+len(addr) : 1+len(addr)+64]
		vout := key[1+len(addr)+64:]

		utxo := bvk.GetUTXO(ctx, string(hash), new(big.Int).SetBytes(vout).Uint64())
		if cb(addr, utxo) {
			break
		}
	}
}

type BaseUTXOKeeper struct {
	BaseUTXOViewKeeper

	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
}

func NewBaseUTXOKeeper(cdc codec.BinaryCodec, storeKey storetypes.StoreKey) *BaseUTXOKeeper {
	return &BaseUTXOKeeper{
		BaseUTXOViewKeeper: *NewBaseUTXOViewKeeper(cdc, storeKey),
		cdc:                cdc,
		storeKey:           storeKey,
	}
}

func (bk *BaseUTXOKeeper) SetUTXO(ctx sdk.Context, utxo *types.UTXO) {
	store := ctx.KVStore(bk.storeKey)

	bz := bk.cdc.MustMarshal(utxo)
	store.Set(types.BtcUtxoKey(utxo.Txid, utxo.Vout), bz)
}

func (bk *BaseUTXOKeeper) SetOwnerUTXO(ctx sdk.Context, utxo *types.UTXO) {
	store := ctx.KVStore(bk.storeKey)

	store.Set(types.BtcOwnerUtxoKey(utxo.Address, utxo.Txid, utxo.Vout), nil)
}

func (bk *BaseUTXOKeeper) LockUTXO(ctx sdk.Context, hash string, vout uint64) error {
	if !bk.HasUTXO(ctx, hash, vout) {
		return types.ErrUTXODoesNotExist
	}

	utxo := bk.GetUTXO(ctx, hash, vout)
	if utxo.IsLocked {
		return types.ErrUTXOLocked
	}

	utxo.IsLocked = true
	bk.SetUTXO(ctx, utxo)

	return nil
}

func (bk *BaseUTXOKeeper) LockUTXOs(ctx sdk.Context, utxos []*types.UTXO) error {
	for _, utxo := range utxos {
		if err := bk.LockUTXO(ctx, utxo.Txid, utxo.Vout); err != nil {
			return err
		}
	}

	return nil
}

func (bk *BaseUTXOKeeper) UnlockUTXO(ctx sdk.Context, hash string, vout uint64) error {
	if !bk.HasUTXO(ctx, hash, vout) {
		return types.ErrUTXODoesNotExist
	}

	utxo := bk.GetUTXO(ctx, hash, vout)
	if !utxo.IsLocked {
		return types.ErrUTXOUnlocked
	}

	utxo.IsLocked = false
	bk.SetUTXO(ctx, utxo)

	return nil
}

func (bk *BaseUTXOKeeper) UnlockUTXOs(ctx sdk.Context, utxos []*types.UTXO) error {
	for _, utxo := range utxos {
		if err := bk.UnlockUTXO(ctx, utxo.Txid, utxo.Vout); err != nil {
			return err
		}
	}

	return nil
}

func (bk *BaseUTXOKeeper) SpendUTXO(ctx sdk.Context, hash string, vout uint64) error {
	if !bk.HasUTXO(ctx, hash, vout) {
		return types.ErrUTXODoesNotExist
	}

	bk.removeUTXO(ctx, hash, vout)

	return nil
}

func (bk *BaseUTXOKeeper) SpendUTXOs(ctx sdk.Context, utxos []*types.UTXO) error {
	for _, utxo := range utxos {
		if err := bk.SpendUTXO(ctx, utxo.Txid, utxo.Vout); err != nil {
			return err
		}
	}

	return nil
}

func (bk *BaseUTXOKeeper) removeUTXO(ctx sdk.Context, hash string, vout uint64) {
	store := ctx.KVStore(bk.storeKey)

	store.Delete(types.BtcUtxoKey(hash, vout))
}
