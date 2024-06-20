package types

import (
	secp256k1 "github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/txscript"
)

// VerifyPsbtSignatures verifies the signatures of the given psbt
// Note: assume that the psbt is valid and all inputs are native segwit
func VerifyPsbtSignatures(p *psbt.Packet) bool {
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
		hashType := p.Inputs[i].SighashType

		witness := p.Inputs[i].FinalScriptWitness
		if len(witness) < 72+33 {
			return false
		}

		sigBytes := witness[0 : len(witness)-33]
		pkBytes := witness[len(witness)-33:]

		if sigBytes[len(sigBytes)-1] != byte(hashType) {
			return false
		}

		sig, err := ecdsa.ParseDERSignature(sigBytes[0 : len(sigBytes)-1])
		if err != nil {
			return false
		}

		pk, err := secp256k1.ParsePubKey(pkBytes)
		if err != nil {
			return false
		}

		sigHash, err := txscript.CalcWitnessSigHash(output.PkScript, txscript.NewTxSigHashes(p.UnsignedTx, prevOutputFetcher),
			hashType, p.UnsignedTx, i, output.Value)
		if err != nil {
			return false
		}

		if !sig.Verify(sigHash, pk) {
			return false
		}
	}

	return true
}
