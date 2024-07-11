package keeper

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"lukechampine.com/uint128"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
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

// NewSigningRequest creates a new signing request
func (k Keeper) NewSigningRequest(ctx sdk.Context, sender string, coin sdk.Coin, feeRate int64) (*types.BitcoinSigningRequest, error) {
	p := k.GetParams(ctx)
	btcVault := types.SelectVaultByAssetType(p.Vaults, types.AssetType_ASSET_TYPE_BTC)

	switch types.AssetTypeFromDenom(coin.Denom, p) {
	case types.AssetType_ASSET_TYPE_BTC:
		return k.NewBtcSigningRequest(ctx, sender, coin, feeRate, btcVault.Address)

	case types.AssetType_ASSET_TYPE_RUNE:
		runesVault := types.SelectVaultByAssetType(p.Vaults, types.AssetType_ASSET_TYPE_RUNE)
		return k.NewRunesSigningRequest(ctx, sender, coin, feeRate, runesVault.Address, btcVault.Address)

	default:
		return nil, types.ErrAssetNotSupported
	}
}

// NewBtcSigningRequest creates a signing request for btc withdrawal
func (k Keeper) NewBtcSigningRequest(ctx sdk.Context, sender string, coin sdk.Coin, feeRate int64, vault string) (*types.BitcoinSigningRequest, error) {
	utxos := k.GetOrderedUTXOsByAddr(ctx, vault)
	if len(utxos) == 0 {
		return nil, types.ErrInsufficientUTXOs
	}

	psbt, selectedUTXOs, _, err := types.BuildPsbt(utxos, sender, coin.Amount.Int64(), feeRate, vault)
	if err != nil {
		return nil, err
	}

	changeUTXO, err := k.handleBtcTxFee(psbt, vault)
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
	k.saveUTXO(ctx, changeUTXO)
	k.addToMintHistory(ctx, psbt.UnsignedTx.TxHash().String())

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

// NewBtcSigningRequest creates a signing request for runes withdrawal
func (k Keeper) NewRunesSigningRequest(ctx sdk.Context, sender string, coin sdk.Coin, feeRate int64, vault string, btcVault string) (*types.BitcoinSigningRequest, error) {
	var runeId types.RuneId
	runeId.FromDenom(coin.Denom)

	amount := uint128.FromBig(coin.Amount.BigInt())

	runesUTXOs, amountDelta := k.GetTargetRunesUTXOs(ctx, vault, runeId.ToString(), amount)
	if len(runesUTXOs) == 0 {
		return nil, types.ErrInsufficientUTXOs
	}

	paymentUTXOs := k.GetOrderedUTXOsByAddr(ctx, btcVault)
	if len(paymentUTXOs) == 0 {
		return nil, types.ErrInsufficientUTXOs
	}

	psbt, selectedUTXOs, changeUTXO, runesChangeUTXO, err := types.BuildRunesPsbt(runesUTXOs, paymentUTXOs, sender, runeId.ToString(), amount, feeRate, amountDelta, vault, btcVault)
	if err != nil {
		return nil, err
	}

	if err := k.handleRunesTxFee(ctx, psbt, sender); err != nil {
		return nil, err
	}

	psbtB64, err := psbt.B64Encode()
	if err != nil {
		return nil, types.ErrFailToSerializePsbt
	}

	// lock the involved utxos
	_ = k.LockUTXOs(ctx, runesUTXOs)
	_ = k.LockUTXOs(ctx, selectedUTXOs)

	// save the change utxo and mark minted
	if changeUTXO != nil {
		k.saveUTXO(ctx, changeUTXO)
		k.addToMintHistory(ctx, psbt.UnsignedTx.TxHash().String())
	}

	// save the runes change utxo and mark minted
	if runesChangeUTXO != nil {
		k.saveUTXO(ctx, runesChangeUTXO)
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

	// spend the locked utxos
	k.spendUTXOs(ctx, uTx)

	// burn the locked asset
	if err := k.burnLockedAsset(ctx, txHash.String()); err != nil {
		return nil, err
	}

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

// handleTxFee performs the fee handling for the btc withdrawal tx
// Make sure that the given psbt is valid
// There are at most two outputs and the change output is the last one if any
func (k Keeper) handleBtcTxFee(p *psbt.Packet, changeAddr string) (*types.UTXO, error) {
	recipientOut := p.UnsignedTx.TxOut[0]

	changeOut := new(wire.TxOut)
	if len(p.UnsignedTx.TxOut) > 1 {
		changeOut = p.UnsignedTx.TxOut[1]
	} else {
		changeOut = wire.NewTxOut(0, types.MustPkScriptFromAddress(changeAddr))
		p.UnsignedTx.TxOut = append(p.UnsignedTx.TxOut, changeOut)
	}

	txFee, err := p.GetTxFee()
	if err != nil {
		return nil, err
	}

	recipientOut.Value -= int64(txFee)
	changeOut.Value += int64(txFee)

	if types.IsDustOut(recipientOut) || types.IsDustOut(changeOut) {
		return nil, types.ErrDustOutput
	}

	return &types.UTXO{
		Txid:         p.UnsignedTx.TxHash().String(),
		Vout:         1,
		Address:      changeAddr,
		Amount:       uint64(changeOut.Value),
		PubKeyScript: changeOut.PkScript,
	}, nil
}

// handleRunesTxFee performs the fee handling for the runes withdrawal tx
func (k Keeper) handleRunesTxFee(ctx sdk.Context, p *psbt.Packet, recipient string) error {
	txFee, err := p.GetTxFee()
	if err != nil {
		return err
	}

	feeCoin := sdk.NewCoin(k.GetParams(ctx).BtcVoucherDenom, sdk.NewInt(int64(txFee)))
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sdk.MustAccAddressFromBech32(recipient), types.ModuleName, sdk.NewCoins(feeCoin)); err != nil {
		return err
	}

	k.lockAsset(ctx, p.UnsignedTx.TxHash().String(), feeCoin)

	return nil
}

// lockAsset locks the given asset by the tx hash
// we can guarantee that the keys do not overlap
func (k Keeper) lockAsset(ctx sdk.Context, txHash string, coin sdk.Coin) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&coin)
	store.Set(types.BtcLockedAssetKey(txHash, bz), []byte{})
}

// burnLockedAsset burns the locked asset
func (k Keeper) burnLockedAsset(ctx sdk.Context, txHash string) error {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, append(types.BtcLockedAssetKeyPrefix, []byte(txHash)...))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()

		var lockedAsset sdk.Coin
		k.cdc.MustUnmarshal(key[1+len(txHash):], &lockedAsset)

		if err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(lockedAsset)); err != nil {
			return err
		}

		store.Delete(key)
	}

	return nil
}
