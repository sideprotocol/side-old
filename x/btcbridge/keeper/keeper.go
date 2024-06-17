package keeper

import (
	"fmt"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/btcbridge/types"
)

type (
	Keeper struct {
		BaseUTXOKeeper

		cdc      codec.BinaryCodec
		storeKey storetypes.StoreKey
		memKey   storetypes.StoreKey

		bankKeeper types.BankKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,

	bankKeeper types.BankKeeper,
) *Keeper {
	return &Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		memKey:         memKey,
		bankKeeper:     bankKeeper,
		BaseUTXOKeeper: *NewBaseUTXOKeeper(cdc, storeKey),
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsStoreKey, bz)
}

func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	store := ctx.KVStore(k.storeKey)
	var params types.Params
	bz := store.Get(types.ParamsStoreKey)
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

func (k Keeper) GetBestBlockHeader(ctx sdk.Context) *types.BlockHeader {
	store := ctx.KVStore(k.storeKey)
	var blockHeader types.BlockHeader
	bz := store.Get(types.BtcBestBlockHeaderKey)
	k.cdc.MustUnmarshal(bz, &blockHeader)
	return &blockHeader
}

func (k Keeper) SetBestBlockHeader(ctx sdk.Context, header *types.BlockHeader) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(header)
	store.Set(types.BtcBestBlockHeaderKey, bz)
}

func (k Keeper) SetBlockHeaders(ctx sdk.Context, blockHeader []*types.BlockHeader) error {
	store := ctx.KVStore(k.storeKey)
	// check if the previous block header exists
	best := k.GetBestBlockHeader(ctx)
	for _, header := range blockHeader {

		// check the block header sanity
		err := blockchain.CheckBlockHeaderSanity(
			HeaderConvert(header),
			chaincfg.MainNetParams.PowLimit,
			blockchain.NewMedianTime(),
			blockchain.BFNone,
		)
		if err != nil {
			return err
		}

		// check whether it's next block header or not
		if best.Hash != header.PreviousBlockHash {
			// check if the block header already exists
			// if exists, then it is a forked block header
			if !store.Has(types.BtcBlockHeaderHeightKey(header.Height)) {
				return types.ErrInvalidHeader
			}

			// a forked block header is detected
			// check if the new block header has more work than the old one
			oldNode := k.GetBlockHeaderByHeight(ctx, header.Height)
			worksOld := blockchain.CalcWork(BitsToTargetUint32(oldNode.Bits))
			worksNew := blockchain.CalcWork(BitsToTargetUint32(header.Bits))
			if worksNew.Cmp(worksOld) <= 0 {
				return types.ErrForkedBlockHeader
			}

			// remove the block headers after the forked block header
			// and consider the forked block header as the best block header
			for i := header.Height; i <= best.Height; i++ {
				ctx.Logger().Info("Removing block header: ", i)
				thash := k.GetBlockHashByHeight(ctx, i)
				store.Delete(types.BtcBlockHeaderHashKey(thash))
				store.Delete(types.BtcBlockHeaderHeightKey(i))
			}
		}

		// store the block header
		bz := k.cdc.MustMarshal(header)
		store.Set(types.BtcBlockHeaderHashKey(header.Hash), bz)
		// store the height to hash mapping
		store.Set(types.BtcBlockHeaderHeightKey(header.Height), []byte(header.Hash))
		// update the best block header
		best = header
	}

	if len(blockHeader) > 0 {
		// set the best block header
		k.SetBestBlockHeader(ctx, best)
	}

	return nil
}

func (k Keeper) GetBlockHeader(ctx sdk.Context, hash string) *types.BlockHeader {
	store := ctx.KVStore(k.storeKey)
	var blockHeader types.BlockHeader
	bz := store.Get(types.BtcBlockHeaderHashKey(hash))
	k.cdc.MustUnmarshal(bz, &blockHeader)
	return &blockHeader
}

func (k Keeper) GetBlockHashByHeight(ctx sdk.Context, height uint64) string {
	store := ctx.KVStore(k.storeKey)
	hash := store.Get(types.BtcBlockHeaderHeightKey(height))
	return string(hash)
}

func (k Keeper) GetBlockHeaderByHeight(ctx sdk.Context, height uint64) *types.BlockHeader {
	store := ctx.KVStore(k.storeKey)
	hash := store.Get(types.BtcBlockHeaderHeightKey(height))
	return k.GetBlockHeader(ctx, string(hash))
}

// GetAllBlockHeaders returns all block headers
func (k Keeper) GetAllBlockHeaders(ctx sdk.Context) []*types.BlockHeader {
	var headers []*types.BlockHeader
	k.IterateBlockHeaders(ctx, func(header types.BlockHeader) (stop bool) {
		headers = append(headers, &header)
		return false
	})
	return headers
}

// IterateBlockHeaders iterates through all block headers
func (k Keeper) IterateBlockHeaders(ctx sdk.Context, process func(header types.BlockHeader) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.BtcBlockHeaderHashPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var header types.BlockHeader
		k.cdc.MustUnmarshal(iterator.Value(), &header)
		if process(header) {
			break
		}
	}
}
