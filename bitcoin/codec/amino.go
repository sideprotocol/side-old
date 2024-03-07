package codec

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/sideprotocol/side/bitcoin/keys/segwit"
)

// RegisterCrypto registers all crypto dependency types with the provided Amino
// codec.
func RegisterCrypto(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*segwit.PubKey)(nil), nil)
	cdc.RegisterConcrete(segwit.PubKey{},
		segwit.PubKeyName, nil)

	cdc.RegisterInterface((*segwit.PrivKey)(nil), nil)
	cdc.RegisterConcrete(segwit.PrivKey{},
		segwit.PrivKeyName, nil)
}
