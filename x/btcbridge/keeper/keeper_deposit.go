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
func (k Keeper) ProcessBitcoinDepositTransaction(ctx sdk.Context, msg *types.MsgSubmitDepositTransactionRequest) (*chainhash.Hash, btcutil.Address, error) {
	ctx.Logger().Info("accept bitcoin deposit tx", "blockhash", msg.Blockhash)

	params := k.GetParams(ctx)
	if !params.IsAuthorizedSender(msg.Sender) {
		return nil, nil, types.ErrSenderAddressNotAuthorized
	}

	tx, prevTx, err := k.ValidateDepositTransaction(ctx, msg.TxBytes, msg.PrevTxBytes, msg.Blockhash, msg.Proof)
	if err != nil {
		return nil, nil, err
	}

	recipient, err := k.Mint(ctx, tx, prevTx, k.GetBlockHeader(ctx, msg.Blockhash).Height)
	if err != nil {
		return nil, nil, err
	}

	return tx.Hash(), recipient, nil
}

// validateDepositTransaction validates the deposit transaction
func (k Keeper) ValidateDepositTransaction(ctx sdk.Context, txBytes string, prevTxBytes string, blockHash string, proof []string) (*btcutil.Tx, *btcutil.Tx, error) {
	params := k.GetParams(ctx)

	header := k.GetBlockHeader(ctx, blockHash)
	// Check if block confirmed
	if header == nil || header.Height == 0 {
		return nil, nil, types.ErrBlockNotFound
	}

	best := k.GetBestBlockHeader(ctx)
	// Check if the block is confirmed
	if best.Height-header.Height < uint64(params.Confirmations) {
		return nil, nil, types.ErrNotConfirmed
	}
	// Check if the block is within the acceptable depth
	// if best.Height-header.Height > param.MaxAcceptableBlockDepth {
	// 	return types.ErrExceedMaxAcceptanceDepth
	// }

	// Decode the base64 transaction
	rawTx, err := base64.StdEncoding.DecodeString(txBytes)
	if err != nil {
		fmt.Println("Error decoding transaction from base64:", err)
		return nil, nil, err
	}

	// Create a new transaction
	var tx wire.MsgTx
	err = tx.Deserialize(bytes.NewReader(rawTx))
	if err != nil {
		fmt.Println("Error deserializing transaction:", err)
		return nil, nil, err
	}

	uTx := btcutil.NewTx(&tx)

	// Validate the transaction
	if err := blockchain.CheckTransactionSanity(uTx); err != nil {
		fmt.Println("Transaction is not valid:", err)
		return nil, nil, err
	}

	// Decode the previous transaction
	rawPrevTx, err := base64.StdEncoding.DecodeString(prevTxBytes)
	if err != nil {
		fmt.Println("Error decoding transaction from base64:", err)
		return nil, nil, err
	}

	// Create a new transaction
	var prevMsgTx wire.MsgTx
	err = prevMsgTx.Deserialize(bytes.NewReader(rawPrevTx))
	if err != nil {
		fmt.Println("Error deserializing transaction:", err)
		return nil, nil, err
	}

	prevTx := btcutil.NewTx(&prevMsgTx)

	// Validate the transaction
	if err := blockchain.CheckTransactionSanity(prevTx); err != nil {
		fmt.Println("Transaction is not valid:", err)
		return nil, nil, err
	}

	if uTx.MsgTx().TxIn[0].PreviousOutPoint.Hash.String() != prevTx.Hash().String() {
		return nil, nil, types.ErrInvalidBtcTransaction
	}

	// check if the proof is valid
	root, err := chainhash.NewHashFromStr(header.MerkleRoot)
	if err != nil {
		return nil, nil, err
	}

	if !types.VerifyMerkleProof(proof, uTx.Hash(), root) {
		k.Logger(ctx).Error("Invalid merkle proof", "txhash", tx, "root", root, "proof", proof)
		return nil, nil, types.ErrTransactionNotIncluded
	}

	return uTx, prevTx, nil
}

// mint performs the minting operation of the voucher token
func (k Keeper) Mint(ctx sdk.Context, tx *btcutil.Tx, prevTx *btcutil.Tx, height uint64) (btcutil.Address, error) {
	params := k.GetParams(ctx)
	chainCfg := sdk.GetConfig().GetBtcChainCfg()

	// check if this is a valid runes deposit tx
	// if any error encountered, this tx is illegal runes deposit
	// if the edict is not nil, it indicates that this is a legal runes deposit tx
	edict, err := types.CheckRunesDepositTransaction(tx.MsgTx(), params.Vaults)
	if err != nil {
		return nil, err
	}

	isRunes := edict != nil

	// extract the recipient for minting voucher token
	recipient, err := types.ExtractRecipientAddr(tx.MsgTx(), prevTx.MsgTx(), params.Vaults, isRunes, chainCfg)
	if err != nil {
		return nil, err
	}

	// mint voucher token and save utxo if the receiver is a vault address
	for i, out := range tx.MsgTx().TxOut {
		if types.IsOpReturnOutput(out) {
			continue
		}

		// check if the output is a valid address
		pks, err := txscript.ParsePkScript(out.PkScript)
		if err != nil {
			return nil, err
		}
		addr, err := pks.Address(chainCfg)
		if err != nil {
			return nil, err
		}

		// check if the receiver is one of the voucher addresses
		vault := types.SelectVaultByBitcoinAddress(params.Vaults, addr.EncodeAddress())
		if vault == nil {
			continue
		}

		// mint the voucher token by asset type and save utxos
		// skip if the asset type of the sender address is unspecified
		switch vault.AssetType {
		case types.AssetType_ASSET_TYPE_BTC:
			err := k.mintBTC(ctx, tx, height, recipient.EncodeAddress(), vault, out, i, params.BtcVoucherDenom)
			if err != nil {
				return nil, err
			}

		case types.AssetType_ASSET_TYPE_RUNE:
			if isRunes && edict.Output == uint32(i) {
				if err := k.mintRunes(ctx, tx, height, recipient.EncodeAddress(), vault, out, i, edict.Id, edict.Amount); err != nil {
					return nil, err
				}
			}
		}
	}

	return recipient, nil
}

func (k Keeper) mintBTC(ctx sdk.Context, tx *btcutil.Tx, height uint64, sender string, vault *types.Vault, out *wire.TxOut, vout int, denom string) error {
	// save the hash of the transaction to prevent double minting
	hash := tx.Hash().String()
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
		Txid:         hash,
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

func (k Keeper) mintRunes(ctx sdk.Context, tx *btcutil.Tx, height uint64, recipient string, vault *types.Vault, out *wire.TxOut, vout int, id *types.RuneId, amount string) error {
	// save the hash of the transaction to prevent double minting
	hash := tx.Hash().String()
	if k.existsInHistory(ctx, hash) {
		return types.ErrTransactionAlreadyMinted
	}
	k.addToMintHistory(ctx, hash)

	coins := sdk.NewCoins(sdk.NewCoin(id.Denom(), sdk.NewIntFromBigInt(types.RuneAmountFromString(amount).Big())))

	receipientAddr, err := sdk.AccAddressFromBech32(recipient)
	if err != nil {
		return err
	}

	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return err
	}

	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receipientAddr, coins); err != nil {
		return err
	}

	utxo := types.UTXO{
		Txid:         hash,
		Vout:         uint64(vout),
		Amount:       uint64(out.Value),
		PubKeyScript: out.PkScript,
		Height:       height,
		Address:      vault.Address,
		IsCoinbase:   false,
		IsLocked:     false,
		Runes: []*types.RuneBalance{{
			Id:     id.ToString(),
			Amount: amount,
		}},
	}

	k.saveUTXO(ctx, &utxo)

	return nil
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
