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

	if uTx.MsgTx().TxIn[0].PreviousOutPoint.Hash.String() != prevTx.Hash().String() {
		return types.ErrInvalidBtcTransaction
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

	// mint voucher token and save utxo if the receiver is a vault address
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
		// check if the receiver is one of the voucher addresses
		vault := types.SelectVaultByBitcoinAddress(param.Vaults, addr.EncodeAddress())
		if vault == nil {
			continue
		}

		// mint the voucher token by asset type and save utxos
		// skip if the asset type of the sender address is unspecified
		switch vault.AssetType {
		case types.AssetType_ASSET_TYPE_BTC:
			k.mintBTC(ctx, uTx, header.Height, sender.EncodeAddress(), vault, out, i, param.BtcVoucherDenom)
		case types.AssetType_ASSET_TYPE_RUNE:
			k.mintRUNE(ctx, uTx, header.Height, sender.EncodeAddress(), vault, out, i, "rune")
		}
	}

	return nil
}

func (k Keeper) mintBTC(ctx sdk.Context, uTx *btcutil.Tx, height uint64, sender string, vault *types.Vault, out *wire.TxOut, vout int, denom string) {

	// save the hash of the transaction to prevent double minting
	hash := uTx.Hash().String()
	if k.hasMintedTxHash(ctx, hash) {
		return
	}
	k.saveMintedTxHash(ctx, hash)

	// mint the voucher token
	coins := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(out.Value)))

	receipient, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return
	}

	k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
	k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receipient, coins)

	utxo := types.UTXO{
		Txid:         uTx.Hash().String(),
		Vout:         uint64(vout),
		Amount:       uint64(out.Value),
		PubKeyScript: out.PkScript,
		Height:       height,
		Address:      vault.GetAddressOnBitcoin(),
		IsCoinbase:   false,
		IsLocked:     false,
	}

	k.SetUTXO(ctx, &utxo)
	k.SetOwnerUTXO(ctx, &utxo)
}

func (k Keeper) mintRUNE(ctx sdk.Context, uTx *btcutil.Tx, height uint64, sender string, vault *types.Vault, out *wire.TxOut, vout int, denom string) {
}

func (k Keeper) hasMintedTxHash(ctx sdk.Context, txHash string) bool {
	store := ctx.KVStore(k.storeKey)

	return store.Has(types.BtcMintedTxHashKey(txHash))
}

func (k Keeper) saveMintedTxHash(ctx sdk.Context, txHash string) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.BtcMintedTxHashKey(txHash), nil)
}
