package keeper

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"slices"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/btcbridge/types"
)

type (
	Keeper struct {
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
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		bankKeeper: bankKeeper,
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

// Process Bitcoin Deposit Transaction
func (k Keeper) ProcessBitcoinDepositTransaction(ctx sdk.Context, msg *types.MsgSubmitTransactionRequest) error {

	ctx.Logger().Info("accept bitcoin deposit tx", "blockhash", msg.Blockhash)

	param := k.GetParams(ctx)
	header := k.GetBlockHeader(ctx, msg.Blockhash)
	// Check if block confirmed
	if header == nil {
		return types.ErrBlockNotFound
	}

	best := k.GetBestBlockHeader(ctx)
	// Check if the block is confirmed
	if best.Height-header.Height < uint64(param.Confirmations) {
		return types.ErrNotConfirmed
	}
	// Check if the block is within the acceptable depth
	if best.Height-header.Height > param.MaxAcceptableBlockDepth {
		return types.ErrExceedMaxAcceptanceDepth
	}

	// Decode the base64 transaction
	txBytes, err := base64.StdEncoding.DecodeString(msg.TxBytes)
	if err != nil {
		fmt.Println("Error decoding transaction from base64:", err)
		return err
	}

	// Create a new transaction
	var tx wire.MsgTx
	err = tx.Deserialize(bytes.NewReader(txBytes))
	if err != nil {
		fmt.Println("Error deserializing transaction:", err)
		return err
	}
	uTx := btcutil.NewTx(&tx)
	if len(uTx.MsgTx().TxIn) < 1 {
		return types.ErrInvalidBtcTransaction
	}

	// Validate the transaction
	if err := blockchain.CheckTransactionSanity(uTx); err != nil {
		fmt.Println("Transaction is not valid:", err)
		return err
	}

	// extract senders from the previous transaction
	prevTxBytes, err := base64.StdEncoding.DecodeString(msg.PrevTxBytes)
	if err != nil {
		fmt.Println("Error decoding transaction from base64:", err)
		return err
	}

	// Create a new transaction
	var prevMsgTx wire.MsgTx
	err = prevMsgTx.Deserialize(bytes.NewReader(prevTxBytes))
	if err != nil {
		fmt.Println("Error deserializing transaction:", err)
		return err
	}

	prevTx := btcutil.NewTx(&prevMsgTx)
	if len(prevTx.MsgTx().TxOut) < 1 {
		return types.ErrInvalidBtcTransaction
	}
	// Validate the transaction
	if err := blockchain.CheckTransactionSanity(prevTx); err != nil {
		fmt.Println("Transaction is not valid:", err)
		return err
	}

	// check if the output is a valid address
	// if there are multiple inputs, then the first input is considered as the sender
	// assumpe all inputs are from the same sender
	out := prevTx.MsgTx().TxOut[tx.TxIn[0].PreviousOutPoint.Index]
	// check if the output is a valid address
	pk, err := txscript.ParsePkScript(out.PkScript)
	if err != nil {
		return err
	}

	sender, err := pk.Address(types.ChainCfg)
	if err != nil {
		return err
	}

	// if pk.Class() != txscript.WitnessV1TaprootTy || pk.Class() != txscript.WitnessV0PubKeyHashTy || pk.Class() != txscript.WitnessV0ScriptHashTy {
	// 	ctx.Logger().Error("Unsupported script type", "script", pk.Class(), "address", sender.EncodeAddress())
	// 	return types.ErrUnsupportedScriptType
	// }

	// check if the proof is valid
	root, err := chainhash.NewHashFromStr(header.MerkleRoot)
	if err != nil {
		return err
	}
	if !types.VerifyMerkleProof(msg.Proof, uTx.Hash(), root) {
		return types.ErrTransactionNotIncluded
	}

	for i, out := range uTx.MsgTx().TxOut {
		// check if the output is a valid address
		pks, err := txscript.ParsePkScript(out.PkScript)
		if err != nil {
			return err
		}
		addr, err := pks.Address(types.ChainCfg)
		if err != nil {
			return err
		}

		// TODO remove the true
		// check if the receiver is one of the voucher addresses
		if true || slices.Contains(param.BtcVoucherAddress, addr.EncodeAddress()) {
			// mint the voucher token
			coins := sdk.NewCoins(sdk.NewCoin(param.BtcVoucherDenom, sdk.NewInt(out.Value)))
			senderAddr, err := sdk.AccAddressFromBech32(sender.EncodeAddress())
			if err != nil {
				return err
			}
			k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
			k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, senderAddr, coins)

			utxo := types.UTXO{
				Txid:         uTx.Hash().String(),
				Vout:         uint64(i),
				Amount:       uint64(out.Value),
				PubKeyScript: out.PkScript,
				Height:       header.Height,
				Address:      addr.EncodeAddress(),
				IsCoinbase:   false,
				IsLocked:     false,
			}

			println("save utxo", utxo.Txid, utxo.Vout)

			ctx.Logger().Info("Minted Bitcoin Voucher", "index", i, "address", addr.EncodeAddress(), "amount", out.Value, "sender", sender.EncodeAddress(), "senderAddr", senderAddr.String(), "coins", coins.String())

		}

	}

	return nil
}

func (k Keeper) SetUtxo(ctx sdk.Context, utxo types.UTXO) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&utxo)
	store.Set(types.BtcUtxoKey(utxo.Txid, utxo.Vout), bz)
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
