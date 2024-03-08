package segwit

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
)

// Sign creates an ECDSA signature on curve Secp256k1, using SHA256 on the msg.
func (privKey *PrivKey) Sign(msg []byte) ([]byte, error) {
	derivedKey := secp256k1.PrivKey{
		Key: privKey.Key,
	}
	return derivedKey.Sign(msg)
}

// VerifySignature validates the signature.
// The msg will be hashed prior to signature verification.
func (pubKey *PubKey) VerifySignature(msg []byte, sigStr []byte) bool {
	derivedPubKey := secp256k1.PubKey{
		Key: pubKey.Key,
	}
	return derivedPubKey.VerifySignature(msg, sigStr)
}
