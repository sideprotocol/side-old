package keeper

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/btcbridge/types"
)

// Process Bitcoin Deposit Transaction
func (k Keeper) ProcessBitcoinDepositTransaction(ctx sdk.Context, msg *types.MsgSubmitDepositTransactionRequest) error {

	ctx.Logger().Info("accept bitcoin deposit tx", "blockhash", msg.Blockhash)

	param := k.GetParams(ctx)
	header := k.GetBlockHeader(ctx, msg.Blockhash)
	// Check if block confirmed
	if header == nil || header.Height == 0 {
		return types.ErrBlockNotFound
	}

	best := k.GetBestBlockHeader(ctx)
	// Check if the block is confirmed
	if best.Height-header.Height < uint64(param.Confirmations) {
		return types.ErrNotConfirmed
	}
	// Check if the block is within the acceptable depth
	// if best.Height-header.Height > param.MaxAcceptableBlockDepth {
	// 	return types.ErrExceedMaxAcceptanceDepth
	// }

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

	// Decode the previous transaction
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

	if uTx.MsgTx().TxIn[0].PreviousOutPoint.Hash.String() != prevTx.Hash().String() {
		return types.ErrInvalidBtcTransaction
	}

	chainCfg := sdk.GetConfig().GetBtcChainCfg()

	// Extract the recipient address
	recipient, err := types.ExtractRecipientAddr(&tx, &prevMsgTx, param.Vaults, chainCfg)
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

	txhash := uTx.MsgTx().TxHash()
	if !types.VerifyMerkleProof(msg.Proof, &txhash, root) {
		k.Logger(ctx).Error("Invalid merkle proof", "txhash", tx, "root", root, "proof", msg.Proof)
		return types.ErrTransactionNotIncluded
	}

	// mint voucher token and save utxo if the receiver is a vault address
	for i, out := range uTx.MsgTx().TxOut {
		// check if the output is a valid address
		pks, err := txscript.ParsePkScript(out.PkScript)
		if err != nil {
			return err
		}
		addr, err := pks.Address(chainCfg)
		if err != nil {
			return err
		}
		// check if the receiver is one of the voucher addresses
		vault := types.SelectVaultByBitcoinAddress(param.Vaults, addr.EncodeAddress())
		if vault == nil {
			continue
		}

		// mint the voucher token by asset type and save utxos
		// skip if the asset type of the sender address is unspecified
		switch vault.AssetType {
		case types.AssetType_ASSET_TYPE_BTC:
			err := k.mintBTC(ctx, uTx, header.Height, recipient.EncodeAddress(), vault, out, i, param.BtcVoucherDenom)
			if err != nil {
				return err
			}
		case types.AssetType_ASSET_TYPE_RUNE:
			k.mintRUNE(ctx, uTx, header.Height, recipient.EncodeAddress(), vault, out, i, "rune")
		}
	}

	return nil
}

func (k Keeper) mintBTC(ctx sdk.Context, uTx *btcutil.Tx, height uint64, sender string, vault *types.Vault, out *wire.TxOut, vout int, denom string) error {

	// save the hash of the transaction to prevent double minting
	hash := uTx.Hash().String()
	if k.existsInHistory(ctx, hash) {
		return types.ErrTransactionAlreadyMinted
	}
	k.addToMintHistory(ctx, hash)

	// mint the voucher token
	if len(denom) == 0 {
		denom = "sat"
	}
	coins := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(out.Value)))

	receipient, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return err
	}

	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return err
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receipient, coins); err != nil {
		return err
	}

	utxo := types.UTXO{
		Txid:         uTx.Hash().String(),
		Vout:         uint64(vout),
		Amount:       uint64(out.Value),
		PubKeyScript: out.PkScript,
		Height:       height,
		Address:      vault.Address,
		IsCoinbase:   false,
		IsLocked:     false,
	}

	k.saveUTXO(ctx, &utxo)

	return nil
}

func (k Keeper) mintRUNE(ctx sdk.Context, uTx *btcutil.Tx, height uint64, sender string, vault *types.Vault, out *wire.TxOut, vout int, denom string) {
}

func (k Keeper) existsInHistory(ctx sdk.Context, txHash string) bool {
	store := ctx.KVStore(k.storeKey)

	return store.Has(types.BtcMintedTxHashKey(txHash))
}

func (k Keeper) addToMintHistory(ctx sdk.Context, txHash string) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.BtcMintedTxHashKey(txHash), []byte{1})
}

// need a query all history for exporting
