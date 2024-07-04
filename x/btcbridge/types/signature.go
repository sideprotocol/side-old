package types

import (
	secp256k1 "github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// VerifyPsbtSignatures verifies the signatures of the given psbt
// Note: assume that the psbt is finalized and all inputs are witness type
func VerifyPsbtSignatures(p *psbt.Packet) bool {
	// extract signed tx
	signedTx, err := psbt.Extract(p)
	if err != nil {
		return false
	}

	// build previous output fetcher
	prevOutputFetcher := txscript.NewMultiPrevOutFetcher(nil)

	for i, txIn := range p.UnsignedTx.TxIn {
		prevOutput := p.Inputs[i].WitnessUtxo
		if prevOutput == nil {
			return false
		}

		prevOutputFetcher.AddPrevOut(txIn.PreviousOutPoint, prevOutput)
	}

	// verify signatures
	for i := range p.Inputs {
		prevOutput := p.Inputs[i].WitnessUtxo
		hashType := DefaultSigHashType

		switch {
		case txscript.IsPayToWitnessPubKeyHash(prevOutput.PkScript):
			if !verifyWitnessSignature(signedTx, i, prevOutput, prevOutputFetcher, hashType) {
				return false
			}

		case txscript.IsPayToTaproot(prevOutput.PkScript):
			if !verifyTaprootSignature(signedTx, i, prevOutput, prevOutputFetcher, hashType) {
				return false
			}

		default:
			return false
		}
	}

	return true
}

// verifyWitnessSignature verifies the signature of the witness v0 input
func verifyWitnessSignature(tx *wire.MsgTx, idx int, prevOutput *wire.TxOut, prevOutputFetcher txscript.PrevOutputFetcher, hashType txscript.SigHashType) bool {
	witness := tx.TxIn[idx].Witness
	if len(witness) != 2 {
		return false
	}

	sigBytes := witness[0]
	pkBytes := witness[1]

	sig, err := ecdsa.ParseDERSignature(sigBytes)
	if err != nil {
		return false
	}

	pk, err := secp256k1.ParsePubKey(pkBytes)
	if err != nil {
		return false
	}

	if sigBytes[len(sigBytes)-1] != byte(hashType) {
		return false
	}

	sigHash, err := txscript.CalcWitnessSigHash(prevOutput.PkScript, txscript.NewTxSigHashes(tx, prevOutputFetcher),
		hashType, tx, idx, prevOutput.Value)
	if err != nil {
		return false
	}

	return sig.Verify(sigHash, pk)
}

// verifyTaprootSignature verifies the signature of the taproot input
func verifyTaprootSignature(tx *wire.MsgTx, idx int, prevOutput *wire.TxOut, prevOutputFetcher txscript.PrevOutputFetcher, hashType txscript.SigHashType) bool {
	witness := tx.TxIn[idx].Witness
	if len(witness) != 1 || len(witness[0]) == 0 {
		return false
	}

	sigBytes := witness[0]

	if hashType != txscript.SigHashDefault {
		if sigBytes[len(sigBytes)-1] != byte(hashType) {
			return false
		}

		sigBytes = sigBytes[0 : len(sigBytes)-1]
	}

	sig, err := schnorr.ParseSignature(sigBytes)
	if err != nil {
		return false
	}

	pk, err := schnorr.ParsePubKey(prevOutput.PkScript[2:34])
	if err != nil {
		return false
	}

	sigHash, err := txscript.CalcTaprootSignatureHash(txscript.NewTxSigHashes(tx, prevOutputFetcher),
		hashType, tx, idx, prevOutputFetcher)
	if err != nil {
		return false
	}

	return sig.Verify(sigHash, pk)
}
