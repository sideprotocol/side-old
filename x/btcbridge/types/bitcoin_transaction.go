package types

import (
	"lukechampine.com/uint128"

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

	// default sig hash type
	DefaultSigHashType = txscript.SigHashAll
)

// BuildPsbt builds a bitcoin psbt from the given params.
// Assume that the utxo script type is witness.
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

	unsignedTx, selectedUTXOs, changeUTXO, err := BuildUnsignedTransaction([]*UTXO{}, txOuts, utxos, feeRate, changeAddr)
	if err != nil {
		return nil, nil, nil, err
	}

	p, err := psbt.NewFromUnsignedTx(unsignedTx)
	if err != nil {
		return nil, nil, nil, err
	}

	for i, utxo := range selectedUTXOs {
		p.Inputs[i].SighashType = DefaultSigHashType
		p.Inputs[i].WitnessUtxo = wire.NewTxOut(int64(utxo.Amount), utxo.PubKeyScript)
	}

	return p, selectedUTXOs, changeUTXO, nil
}

// BuildRunesPsbt builds a bitcoin psbt for runes edict from the given params.
// Assume that the utxo script type is witness.
func BuildRunesPsbt(utxos []*UTXO, paymentUTXOs []*UTXO, recipient string, runeId string, amount uint128.Uint128, feeRate int64, runesChangeAmount uint128.Uint128, runesChange string, change string) (*psbt.Packet, []*UTXO, *UTXO, *UTXO, error) {
	chaincfg := sdk.GetConfig().GetBtcChainCfg()

	recipientAddr, err := btcutil.DecodeAddress(recipient, chaincfg)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	recipientPkScript, err := txscript.PayToAddrScript(recipientAddr)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	changeAddr, err := btcutil.DecodeAddress(change, chaincfg)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	runesChangeAddr, err := btcutil.DecodeAddress(runesChange, chaincfg)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	runesChangePkScript, err := txscript.PayToAddrScript(runesChangeAddr)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	txOuts := make([]*wire.TxOut, 0)

	// fill the runes protocol script with empty output script first
	txOuts = append(txOuts, wire.NewTxOut(0, []byte{}))

	var runesChangeUTXO *UTXO
	edictOutputIndex := uint32(1)

	if runesChangeAmount.Cmp64(0) > 0 {
		// we can guarantee that every runes UTXO only includes a single rune by the deposit policy
		runesChangeUTXO = GetRunesChangeUTXO(runeId, runesChangeAmount, runesChange, runesChangePkScript, 1)

		// allocate the remaining runes to the first non-OP_RETURN output by default
		txOuts = append(txOuts, wire.NewTxOut(RunesOutValue, runesChangePkScript))

		// advance the edict output index
		edictOutputIndex++
	}

	// edict output
	txOuts = append(txOuts, wire.NewTxOut(RunesOutValue, recipientPkScript))

	runesScript, err := BuildEdictScript(runeId, amount, edictOutputIndex)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// populate the runes protocol script
	txOuts[0].PkScript = runesScript

	unsignedTx, selectedUTXOs, changeUTXO, err := BuildUnsignedTransaction(utxos, txOuts, paymentUTXOs, feeRate, changeAddr)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if runesChangeUTXO != nil {
		runesChangeUTXO.Txid = unsignedTx.TxHash().String()
	}

	p, err := psbt.NewFromUnsignedTx(unsignedTx)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	for i, utxo := range utxos {
		p.Inputs[i].SighashType = DefaultSigHashType
		p.Inputs[i].WitnessUtxo = wire.NewTxOut(int64(utxo.Amount), utxo.PubKeyScript)
	}

	for i, utxo := range selectedUTXOs {
		p.Inputs[i+len(utxos)].SighashType = DefaultSigHashType
		p.Inputs[i+len(utxos)].WitnessUtxo = wire.NewTxOut(int64(utxo.Amount), utxo.PubKeyScript)
	}

	return p, selectedUTXOs, changeUTXO, runesChangeUTXO, nil
}

// BuildUnsignedTransaction builds an unsigned tx from the given params.
func BuildUnsignedTransaction(utxos []*UTXO, txOuts []*wire.TxOut, paymentUTXOs []*UTXO, feeRate int64, change btcutil.Address) (*wire.MsgTx, []*UTXO, *UTXO, error) {
	tx := wire.NewMsgTx(TxVersion)

	inAmount := int64(0)
	outAmount := int64(0)

	for _, utxo := range utxos {
		AddUTXOToTx(tx, utxo)
		inAmount += int64(utxo.Amount)
	}

	for _, txOut := range txOuts {
		if IsDustOut(txOut) {
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

	selectedUTXOs, err := AddPaymentUTXOsToTx(tx, utxos, inAmount-outAmount, paymentUTXOs, changeOut, feeRate)
	if err != nil {
		return nil, nil, nil, err
	}

	var changeUTXO *UTXO
	if len(tx.TxOut) > len(txOuts) {
		changeUTXO = GetChangeUTXO(tx, change.EncodeAddress())
	}

	return tx, selectedUTXOs, changeUTXO, nil
}

// AddPaymentUTXOsToTx adds the given payment utxos to the tx.
func AddPaymentUTXOsToTx(tx *wire.MsgTx, utxos []*UTXO, inOutDiff int64, paymentUtxos []*UTXO, changeOut *wire.TxOut, feeRate int64) ([]*UTXO, error) {
	selectedUTXOs := make([]*UTXO, 0)
	paymentValue := int64(0)

	for _, utxo := range paymentUtxos {
		AddUTXOToTx(tx, utxo)
		tx.AddTxOut(changeOut)

		utxos = append(utxos, utxo)
		selectedUTXOs = append(selectedUTXOs, utxo)

		paymentValue += int64(utxo.Amount)
		fee := GetTxVirtualSize(tx, utxos) * feeRate

		changeValue := paymentValue + inOutDiff - fee
		if changeValue > 0 {
			tx.TxOut[len(tx.TxOut)-1].Value = changeValue
			if IsDustOut(tx.TxOut[len(tx.TxOut)-1]) {
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
			if paymentValue+inOutDiff-feeWithoutChange >= 0 {
				return selectedUTXOs, nil
			}
		}
	}

	return nil, ErrInsufficientUTXOs
}

// AddUTXOToTx adds the given utxo to the specified tx
// Make sure the utxo is valid
func AddUTXOToTx(tx *wire.MsgTx, utxo *UTXO) {
	txIn := new(wire.TxIn)

	hash, err := chainhash.NewHashFromStr(utxo.Txid)
	if err != nil {
		panic(err)
	}

	txIn.PreviousOutPoint = *wire.NewOutPoint(hash, uint32(utxo.Vout))

	tx.AddTxIn(txIn)
}

// GetChangeUTXO returns the change output from the given tx
// Make sure that the tx is valid and the change output is the last output
func GetChangeUTXO(tx *wire.MsgTx, change string) *UTXO {
	changeOut := tx.TxOut[len(tx.TxOut)-1]

	return &UTXO{
		Txid:         tx.TxHash().String(),
		Vout:         uint64(len(tx.TxOut) - 1),
		Address:      change,
		Amount:       uint64(changeOut.Value),
		PubKeyScript: changeOut.PkScript,
	}
}

// GetRunesChangeUTXO gets the runes change utxo.
func GetRunesChangeUTXO(runeId string, changeAmount uint128.Uint128, change string, changePkScript []byte, outIndex uint32) *UTXO {
	return &UTXO{
		Vout:         uint64(outIndex),
		Address:      change,
		Amount:       RunesOutValue,
		PubKeyScript: changePkScript,
		Runes: []*RuneBalance{
			{
				Id:     runeId,
				Amount: changeAmount.String(),
			},
		},
	}
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

// IsDustOut returns true if the given output is dust, false otherwise
func IsDustOut(out *wire.TxOut) bool {
	return !IsOpReturnOutput(out) && mempool.IsDust(out, MinRelayFee)
}

// CheckOutputAmount checks if the given output amount is dust
func CheckOutputAmount(address string, amount int64) error {
	addr, err := btcutil.DecodeAddress(address, sdk.GetConfig().GetBtcChainCfg())
	if err != nil {
		return err
	}

	pkScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return err
	}

	if IsDustOut(&wire.TxOut{Value: amount, PkScript: pkScript}) {
		return ErrDustOutput
	}

	return nil
}

// IsOpReturnOutput returns true if the script of the given out starts with OP_RETURN
func IsOpReturnOutput(out *wire.TxOut) bool {
	return len(out.PkScript) > 0 && out.PkScript[0] == txscript.OP_RETURN
}
