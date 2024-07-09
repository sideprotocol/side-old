package types

import (
	"encoding/base64"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
)

// VerifyMerkleProof verifies a Merkle proof
func VerifyMerkleProof(proofs []string, txHash, root *chainhash.Hash) bool {
	current := txHash
	for _, proof := range proofs {

		bytes, err := base64.StdEncoding.DecodeString(proof)
		if err != nil {
			return false
		}
		position := bytes[0]
		p := current
		if len(bytes) > 1 {
			p, err = chainhash.NewHash(bytes[1:])
			if err != nil {
				return false
			}
		}

		var temp chainhash.Hash
		if position == 0 {
			temp = blockchain.HashMerkleBranches(current, p)
		} else {
			temp = blockchain.HashMerkleBranches(p, current)
		}
		current = &temp
	}

	return current.IsEqual(root)
}
