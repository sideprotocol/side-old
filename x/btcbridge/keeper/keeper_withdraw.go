package keeper

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sideprotocol/side/x/btcbridge/types"
)

// GetRequestSeqence returns the request sequence
func (k Keeper) GetRequestSeqence(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.SequenceKey)
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

// IncrementRequestSequence increments the request sequence and returns the new sequence
func (k Keeper) IncrementRequestSequence(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	seq := k.GetRequestSeqence(ctx) + 1
	store.Set(types.SequenceKey, sdk.Uint64ToBigEndian(seq))
	return seq
}

// New signing request
// sender: the address of the sender
// txBytes: the transaction bytes
// vault: the address of the vault, default is empty.
// If empty, the vault will be Bitcoin vault, otherwise it will be Ordinals or Runes vault
func (k Keeper) NewSigningRequest(ctx sdk.Context, sender string, coin sdk.Coin, feeRate int64, vault string) (*types.BitcoinSigningRequest, error) {
	if len(vault) == 0 {
		// default to the first vault in the params for now
		// TODO: select an appropriate vault according to the utxos
		p := k.GetParams(ctx)
		for i, v := range p.Vaults {
			if v.AssetType == types.AssetType_ASSET_TYPE_BTC {
				vault = p.Vaults[i].Address
				break
			}
		}
	}

	utxos := k.GetOrderedUTXOsByAddr(ctx, vault)
	if len(utxos) == 0 {
		return nil, types.ErrInsufficientUTXOs
	}

	psbt, selectedUTXOs, changeUTXO, err := types.BuildPsbt(utxos, sender, coin.Amount.Int64(), feeRate, vault)
	if err != nil {
		return nil, err
	}

	psbtB64, err := psbt.B64Encode()
	if err != nil {
		return nil, types.ErrFailToSerializePsbt
	}

	// lock the selected utxos
	_ = k.LockUTXOs(ctx, selectedUTXOs)

	// save the change utxo and mark minted
	if changeUTXO != nil {
		k.saveUTXO(ctx, changeUTXO)
		k.addToMintHistory(ctx, psbt.UnsignedTx.TxHash().String())
	}

	signingRequest := &types.BitcoinSigningRequest{
		Address:      sender,
		Txid:         psbt.UnsignedTx.TxHash().String(),
		Psbt:         psbtB64,
		Status:       types.SigningStatus_SIGNING_STATUS_CREATED,
		Sequence:     k.IncrementRequestSequence(ctx),
		VaultAddress: vault,
	}

	k.SetSigningRequest(ctx, signingRequest)

	return signingRequest, nil
}

// GetSigningRequest returns the signing request
func (k Keeper) HasSigningRequest(ctx sdk.Context, hash string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.BtcSigningRequestHashKey(hash))
}

// GetSigningRequest returns the signing request
func (k Keeper) GetSigningRequest(ctx sdk.Context, hash string) *types.BitcoinSigningRequest {
	store := ctx.KVStore(k.storeKey)
	var signingRequest types.BitcoinSigningRequest
	// TODO replace the key with the hash
	bz := store.Get(types.BtcSigningRequestHashKey(hash))
	k.cdc.MustUnmarshal(bz, &signingRequest)
	return &signingRequest
}

// SetSigningRequest sets the signing request
func (k Keeper) SetSigningRequest(ctx sdk.Context, signingRequest *types.BitcoinSigningRequest) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(signingRequest)
	// TODO replace the key with the hash
	store.Set(types.BtcSigningRequestHashKey(signingRequest.Txid), bz)
}

// IterateSigningRequests iterates through all signing requests
func (k Keeper) IterateSigningRequests(ctx sdk.Context, process func(signingRequest types.BitcoinSigningRequest) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.BtcSigningRequestPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var signingRequest types.BitcoinSigningRequest
		k.cdc.MustUnmarshal(iterator.Value(), &signingRequest)
		if process(signingRequest) {
			break
		}
	}
}

// filter SigningRequest by status with pagination
func (k Keeper) FilterSigningRequestsByStatus(ctx sdk.Context, req *types.QuerySigningRequestRequest) []*types.BitcoinSigningRequest {
	var signingRequests []*types.BitcoinSigningRequest
	k.IterateSigningRequests(ctx, func(signingRequest types.BitcoinSigningRequest) (stop bool) {
		if signingRequest.Status == req.Status {
			signingRequests = append(signingRequests, &signingRequest)
		}
		// pagination TODO: limit the number of signing requests
		if len(signingRequests) >= 100 {
			return true
		}
		return false
	})
	return signingRequests
}

// filter SigningRequest by address with pagination
func (k Keeper) FilterSigningRequestsByAddr(ctx sdk.Context, req *types.QuerySigningRequestByAddressRequest) []*types.BitcoinSigningRequest {
	var signingRequests []*types.BitcoinSigningRequest
	k.IterateSigningRequests(ctx, func(signingRequest types.BitcoinSigningRequest) (stop bool) {
		if signingRequest.Address == req.Address {
			signingRequests = append(signingRequests, &signingRequest)
		}
		// pagination TODO: limit the number of signing requests
		if len(signingRequests) >= 100 {
			return true
		}
		return false
	})
	return signingRequests
}

// Process Bitcoin Withdraw Transaction
func (k Keeper) ProcessBitcoinWithdrawTransaction(ctx sdk.Context, msg *types.MsgSubmitWithdrawTransactionRequest) (*chainhash.Hash, error) {
	ctx.Logger().Info("accept bitcoin withdraw tx", "blockhash", msg.Blockhash)

	param := k.GetParams(ctx)

	if !param.IsAuthorizedSender(msg.Sender) {
		return nil, types.ErrSenderAddressNotAuthorized
	}

	header := k.GetBlockHeader(ctx, msg.Blockhash)
	// Check if block confirmed
	if header == nil {
		return nil, types.ErrBlockNotFound
	}

	best := k.GetBestBlockHeader(ctx)
	// Check if the block is confirmed
	if best.Height-header.Height < uint64(param.Confirmations) {
		return nil, types.ErrNotConfirmed
	}
	// Check if the block is within the acceptable depth
	if best.Height-header.Height > param.MaxAcceptableBlockDepth {
		return nil, types.ErrExceedMaxAcceptanceDepth
	}

	// Decode the base64 transaction
	txBytes, err := base64.StdEncoding.DecodeString(msg.TxBytes)
	if err != nil {
		fmt.Println("Error decoding transaction from base64:", err)
		return nil, err
	}

	// Create a new transaction
	var tx wire.MsgTx
	err = tx.Deserialize(bytes.NewReader(txBytes))
	if err != nil {
		fmt.Println("Error deserializing transaction:", err)
		return nil, err
	}

	uTx := btcutil.NewTx(&tx)
	if len(uTx.MsgTx().TxIn) < 1 {
		return nil, types.ErrInvalidBtcTransaction
	}

	txHash := uTx.MsgTx().TxHash()

	if !k.HasSigningRequest(ctx, txHash.String()) {
		return nil, types.ErrSigningRequestNotExist
	}

	signingRequest := k.GetSigningRequest(ctx, txHash.String())
	// if signingRequest.Status != types.SigningStatus_SIGNING_STATUS_BROADCASTED || signingRequest.Status != types.SigningStatus_SIGNING_STATUS_SIGNED {
	// 	return types.ErrInvalidStatus
	// }
	signingRequest.Status = types.SigningStatus_SIGNING_STATUS_CONFIRMED
	k.SetSigningRequest(ctx, signingRequest)

	// Validate the transaction
	if err := blockchain.CheckTransactionSanity(uTx); err != nil {
		fmt.Println("Transaction is not valid:", err)
		return nil, err
	}

	if len(uTx.MsgTx().TxIn[0].Witness) != 2 {
		return nil, types.ErrInvalidSenders
	}

	senderPubKey := uTx.MsgTx().TxIn[0].Witness[1]

	// check if the first sender is one of the vault addresses
	vault := types.SelectVaultByPubKey(param.Vaults, hex.EncodeToString(senderPubKey))
	if vault == nil {
		return nil, types.ErrInvalidSenders
	}

	k.spendUTXOs(ctx, uTx)

	return &txHash, nil
}

// spendUTXOs spends locked utxos
func (k Keeper) spendUTXOs(ctx sdk.Context, uTx *btcutil.Tx) {
	for _, in := range uTx.MsgTx().TxIn {
		hash := in.PreviousOutPoint.Hash.String()
		vout := in.PreviousOutPoint.Index

		if k.IsUTXOLocked(ctx, hash, uint64(vout)) {
			_ = k.SpendUTXO(ctx, hash, uint64(vout))
		}
	}
}
