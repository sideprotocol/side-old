package types

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/mempool"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// default tx version
	TxVersion = 2

	// default minimum relay fee
	MinRelayFee = 1000

	// default hash type for signature
	SigHashType = txscript.SigHashAll
)

// BuildPsbt builds a bitcoin psbt from the given params.
// Assume that the utxo script type is native segwit.
func BuildPsbt(utxos []*UTXO, recipient string, amount int64, feeRate int64, change string) (*psbt.Packet, []*UTXO, *UTXO, error) {
	chaincfg := sdk.GetConfig().GetBtcChainCfg()
	recipientAddr, err := btcutil.DecodeAddress(recipient, chaincfg)
	if err != nil {
		return nil, nil, nil, err
	}

	recipientPkScript, err := txscript.PayToAddrScript(recipientAddr)
	if err != nil {
		return nil, nil, nil, err
	}

	changeAddr, err := btcutil.DecodeAddress(change, chaincfg)
	if err != nil {
		return nil, nil, nil, err
	}

	txOuts := make([]*wire.TxOut, 0)
	txOuts = append(txOuts, wire.NewTxOut(amount, recipientPkScript))

	unsignedTx, selectedUTXOs, changeUTXO, err := BuildUnsignedTransaction(utxos, txOuts, feeRate, changeAddr)
	if err != nil {
		return nil, nil, nil, err
	}

	p, err := psbt.NewFromUnsignedTx(unsignedTx)
	if err != nil {
		return nil, nil, nil, err
	}

	for i, utxo := range selectedUTXOs {
		p.Inputs[i].SighashType = txscript.SigHashAll
		p.Inputs[i].WitnessUtxo = wire.NewTxOut(int64(utxo.Amount), utxo.PubKeyScript)
	}

	return p, selectedUTXOs, changeUTXO, nil
}

// BuildUnsignedTransaction builds an unsigned tx from the given params.
func BuildUnsignedTransaction(utxos []*UTXO, txOuts []*wire.TxOut, feeRate int64, change btcutil.Address) (*wire.MsgTx, []*UTXO, *UTXO, error) {
	tx := wire.NewMsgTx(TxVersion)

	outAmount := int64(0)
	for _, txOut := range txOuts {
		if mempool.IsDust(txOut, MinRelayFee) {
			return nil, nil, nil, ErrDustOutput
		}

		tx.AddTxOut(txOut)
		outAmount += txOut.Value
	}

	changePkScript, err := txscript.PayToAddrScript(change)
	if err != nil {
		return nil, nil, nil, err
	}

	changeOut := wire.NewTxOut(0, changePkScript)

	selectedUTXOs, err := AddUTXOsToTx(tx, utxos, outAmount, changeOut, feeRate)
	if err != nil {
		return nil, nil, nil, err
	}

	var changeUTXO *UTXO
	if len(tx.TxOut) > len(txOuts) {
		changeOut := tx.TxOut[len(tx.TxOut)-1]
		changeUTXO = &UTXO{
			Txid:         tx.TxHash().String(),
			Vout:         uint64(len(tx.TxOut) - 1),
			Address:      change.EncodeAddress(),
			Amount:       uint64(changeOut.Value),
			PubKeyScript: changeOut.PkScript,
		}
	}

	return tx, selectedUTXOs, changeUTXO, nil
}

// AddUTXOsToTx adds the given utxos to the tx.
func AddUTXOsToTx(tx *wire.MsgTx, utxos []*UTXO, outAmount int64, changeOut *wire.TxOut, feeRate int64) ([]*UTXO, error) {
	selectedUTXOs := make([]*UTXO, 0)
	inputAmount := int64(0)

	for _, utxo := range utxos {
		txIn := new(wire.TxIn)

		hash, err := chainhash.NewHashFromStr(utxo.Txid)
		if err != nil {
			return nil, err
		}

		txIn.PreviousOutPoint = *wire.NewOutPoint(hash, uint32(utxo.Vout))

		tx.AddTxIn(txIn)
		tx.AddTxOut(changeOut)

		selectedUTXOs = append(selectedUTXOs, utxo)

		inputAmount += int64(utxo.Amount)
		fee := GetTxVirtualSize(tx, utxos) * feeRate

		changeValue := inputAmount - outAmount - fee
		if changeValue > 0 {
			tx.TxOut[len(tx.TxOut)-1].Value = changeValue
			if mempool.IsDust(tx.TxOut[len(tx.TxOut)-1], btcutil.Amount(MinRelayFee)) {
				tx.TxOut = tx.TxOut[0 : len(tx.TxOut)-1]
			}

			return selectedUTXOs, nil
		}

		tx.TxOut = tx.TxOut[0 : len(tx.TxOut)-1]

		if changeValue == 0 {
			return selectedUTXOs, nil
		}

		if changeValue < 0 {
			feeWithoutChange := GetTxVirtualSize(tx, selectedUTXOs) * feeRate
			if inputAmount-outAmount-feeWithoutChange >= 0 {
				return selectedUTXOs, nil
			}
		}
	}

	return nil, ErrInsufficientUTXOs
}

// GetTxVirtualSize gets the virtual size of the given tx.
// Assume that the utxo script type is p2tr, p2wpkh, p2sh-p2wpkh or p2pkh.
func GetTxVirtualSize(tx *wire.MsgTx, utxos []*UTXO) int64 {
	newTx := tx.Copy()

	for i, txIn := range newTx.TxIn {
		var dummySigScript []byte
		var dummyWitness []byte

		switch txscript.GetScriptClass(utxos[i].PubKeyScript) {
		case txscript.WitnessV1TaprootTy:
			dummyWitness = make([]byte, 65)

		case txscript.WitnessV0PubKeyHashTy:
			dummyWitness = make([]byte, 73+33)

		case txscript.ScriptHashTy:
			dummySigScript = make([]byte, 1+1+1+20)
			dummyWitness = make([]byte, 73+33)

		case txscript.PubKeyHashTy:
			dummySigScript = make([]byte, 1+73+1+33)

		default:
		}

		txIn.SignatureScript = dummySigScript
		txIn.Witness = wire.TxWitness{dummyWitness}
	}

	return mempool.GetTxVirtualSize(btcutil.NewTx(newTx))
}

// CheckOutput checks the given output
func CheckOutput(address string, amount int64) error {
	addr, err := btcutil.DecodeAddress(address, sdk.GetConfig().GetBtcChainCfg())
	if err != nil {
		return err
	}

	pkScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return err
	}

	if mempool.IsDust(&wire.TxOut{Value: amount, PkScript: pkScript}, MinRelayFee) {
		return ErrDustOutput
	}

	return nil
}
