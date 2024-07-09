package types

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

const (
	// maximum allowed number of the non-vault outputs for the btc deposit transaction
	MaxNonVaultOutNum = 1

	// maximum allowed number of the non-vault outputs for the runes deposit transaction
	RunesMaxNonVaultOutNum = 3

	// allowed number of edicts in the runes payload for the runes deposit transaction
	RunesEdictNum = 1
)

// ExtractRecipientAddr extracts the recipient address for minting voucher token by the type of the asset to be deposited
func ExtractRecipientAddr(tx *wire.MsgTx, prevTx *wire.MsgTx, vaults []*Vault, isRunes bool, chainCfg *chaincfg.Params) (btcutil.Address, error) {
	if isRunes {
		return ExtractRunesRecipientAddr(tx, prevTx, vaults, chainCfg)
	}

	return ExtractCommonRecipientAddr(tx, prevTx, vaults, chainCfg)
}

// ExtractCommonRecipientAddr extracts the recipient address for minting voucher token in the common case.
// First, extract the recipient from the tx out which is a non-vault address;
// Then fall back to the first input
func ExtractCommonRecipientAddr(tx *wire.MsgTx, prevTx *wire.MsgTx, vaults []*Vault, chainCfg *chaincfg.Params) (btcutil.Address, error) {
	var recipient btcutil.Address

	nonVaultOutCount := 0

	// extract from the tx out which is a non-vault address
	for _, out := range tx.TxOut {
		pkScript, err := txscript.ParsePkScript(out.PkScript)
		if err != nil {
			return nil, err
		}

		addr, err := pkScript.Address(chainCfg)
		if err != nil {
			return nil, err
		}

		vault := SelectVaultByBitcoinAddress(vaults, addr.EncodeAddress())
		if vault == nil {
			recipient = addr
			nonVaultOutCount++
		}
	}

	// exceed allowed non vault out number
	if nonVaultOutCount > MaxNonVaultOutNum {
		return nil, ErrInvalidDepositTransaction
	}

	if recipient != nil {
		return recipient, nil
	}

	// fall back to extract from the first input
	pkScript, err := txscript.ParsePkScript(prevTx.TxOut[tx.TxIn[0].PreviousOutPoint.Index].PkScript)
	if err != nil {
		return nil, err
	}

	return pkScript.Address(chainCfg)
}

// ExtractRunesRecipientAddr extracts the recipient address for minting runes voucher token.
// First, extract the recipient from the tx out which is a non-vault and non-OP_RETURN output;
// Then fall back to the first input
func ExtractRunesRecipientAddr(tx *wire.MsgTx, prevTx *wire.MsgTx, vaults []*Vault, chainCfg *chaincfg.Params) (btcutil.Address, error) {
	var recipient btcutil.Address

	nonVaultOutCount := 0

	// extract from the tx out which is a non-vault and non-OP_RETURN output
	for _, out := range tx.TxOut {
		if IsOpReturnOutput(out) {
			nonVaultOutCount++
			continue
		}

		pkScript, err := txscript.ParsePkScript(out.PkScript)
		if err != nil {
			return nil, err
		}

		addr, err := pkScript.Address(chainCfg)
		if err != nil {
			return nil, err
		}

		vault := SelectVaultByBitcoinAddress(vaults, addr.EncodeAddress())
		if vault == nil {
			recipient = addr
			nonVaultOutCount++
		}
	}

	// exceed allowed non vault out number
	if nonVaultOutCount > RunesMaxNonVaultOutNum {
		return nil, ErrInvalidDepositTransaction
	}

	if recipient != nil {
		return recipient, nil
	}

	// fall back to extract from the first input
	pkScript, err := txscript.ParsePkScript(prevTx.TxOut[tx.TxIn[0].PreviousOutPoint.Index].PkScript)
	if err != nil {
		return nil, err
	}

	return pkScript.Address(chainCfg)
}

// CheckRunesDepositTransaction checks if the given tx is valid runes deposit tx
func CheckRunesDepositTransaction(tx *wire.MsgTx, vaults []*Vault) (*Edict, error) {
	edicts, err := ParseRunes(tx)
	if err != nil {
		return nil, ErrInvalidDepositTransaction
	}

	if len(edicts) == 0 {
		return nil, nil
	}

	if len(edicts) != RunesEdictNum {
		return nil, ErrInvalidDepositTransaction
	}

	// even split is not supported
	if edicts[0].Output == uint32(len(tx.TxOut)) {
		return nil, ErrInvalidDepositTransaction
	}

	vault := SelectVaultByPkScript(vaults, tx.TxOut[edicts[0].Output].PkScript)
	if vault == nil || vault.AssetType != AssetType_ASSET_TYPE_RUNE {
		return nil, ErrInvalidDepositTransaction
	}

	return edicts[0], nil
}
