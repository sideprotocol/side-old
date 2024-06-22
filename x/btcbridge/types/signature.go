package types

import (
	secp256k1 "github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/txscript"
)

// VerifyPsbtSignatures verifies the signatures of the given psbt
// Note: assume that the psbt is finalized and all inputs are native segwit
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
		output := p.Inputs[i].WitnessUtxo

		witness := signedTx.TxIn[i].Witness
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

		if sigBytes[len(sigBytes)-1] != byte(SigHashType) {
			return false
		}

		sigHash, err := txscript.CalcWitnessSigHash(output.PkScript, txscript.NewTxSigHashes(p.UnsignedTx, prevOutputFetcher),
			SigHashType, p.UnsignedTx, i, output.Value)
		if err != nil {
			return false
		}

		if !sig.Verify(sigHash, pk) {
			return false
		}
	}

	return true
}
