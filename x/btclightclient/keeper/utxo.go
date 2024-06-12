package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sideprotocol/side/x/btclightclient/types"
)

type UTXOViewKeeper interface {
	HasUTXO(ctx sdk.Context, hash string, vout uint64) bool
	IsUTXOLocked(ctx sdk.Context, hash string, vout uint64) bool

	GetAllUTXOs(ctx sdk.Context) []*types.UTXO

	GetUTXOsByAddr(ctx sdk.Context, addr string) []*types.UTXO
	GetOrderedUTXOsByAddr(ctx sdk.Context, addr string) []*types.UTXO
}

type UTXOKeeper interface {
	UTXOViewKeeper

	LockUTXO(ctx sdk.Context, hash string, vout uint64)
	SpendUTXO(ctx sdk.Context, hash string, vout uint64)
}

var _ UTXOKeeper = (*BaseUTXOKeeper)(nil)

type BaseUTXOViewKeeper struct {
	k Keeper
}

func NewBaseUTXOViewKeeper(k Keeper) *BaseUTXOViewKeeper {
	return &BaseUTXOViewKeeper{k}
}

func (bvk *BaseUTXOViewKeeper) HasUTXO(ctx sdk.Context, hash string, vout uint64) bool {
	// TODO
	return true
}

func (bvk *BaseUTXOViewKeeper) IsUTXOLocked(ctx sdk.Context, hash string, vout uint64) bool {
	// TODO
	return true
}

func (bvk *BaseUTXOViewKeeper) GetAllUTXOs(ctx sdk.Context) []*types.UTXO {
	// TODO
	return nil
}

func (bvk *BaseUTXOViewKeeper) GetUTXOsByAddr(ctx sdk.Context, addr string) []*types.UTXO {
	// TODO
	return nil
}

func (bvk *BaseUTXOViewKeeper) GetOrderedUTXOsByAddr(ctx sdk.Context, addr string) []*types.UTXO {
	// TODO
	return nil
}

type BaseUTXOKeeper struct {
	BaseUTXOViewKeeper

	k Keeper
}

func NewBaseUTXOKeeper(k Keeper) *BaseUTXOKeeper {
	return &BaseUTXOKeeper{BaseUTXOViewKeeper: *NewBaseUTXOViewKeeper(k), k: k}
}

func (bk *BaseUTXOKeeper) LockUTXO(ctx sdk.Context, hash string, vout uint64) {
	// TODO
}

func (bk *BaseUTXOKeeper) SpendUTXO(ctx sdk.Context, hash string, vout uint64) {
	// TODO
}
