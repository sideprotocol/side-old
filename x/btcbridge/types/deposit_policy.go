package types

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

const (
	// maximum allowed number of the non-vault outputs for the deposit transaction
	MaxNonVaultOutNum = 1
)

// ExtractRecipientAddr extracts the recipient address for minting voucher token.
// First, extract the recipient from the tx out which is a non-vault address;
// Then fallback to the first input
func ExtractRecipientAddr(tx *wire.MsgTx, prevTx *wire.MsgTx, vaults []*Vault, chainCfg *chaincfg.Params) (btcutil.Address, error) {
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

	// fallback to extract from the first input
	pkScript, err := txscript.ParsePkScript(prevTx.TxOut[tx.TxIn[0].PreviousOutPoint.Index].PkScript)
	if err != nil {
		return nil, err
	}

	return pkScript.Address(chainCfg)
}
