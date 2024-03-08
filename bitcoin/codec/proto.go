package codec

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256r1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/sideprotocol/side/bitcoin/keys/segwit"
)

// RegisterInterfaces registers the sdk.Tx interface.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	var pk *cryptotypes.PubKey
	registry.RegisterImplementations(pk, &segwit.PubKey{})

	var priv *cryptotypes.PrivKey
	registry.RegisterImplementations(priv, &segwit.PrivKey{})
	secp256r1.RegisterInterfaces(registry)
}
