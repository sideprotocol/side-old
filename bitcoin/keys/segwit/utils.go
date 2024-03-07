package segwit

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

func BitCoinAddr(pubKey []byte) (string, error) {
	witnessProg := btcutil.Hash160(pubKey)
	bech32Address, err := btcutil.NewAddressWitnessPubKeyHash(witnessProg, &chaincfg.MainNetParams)
	return bech32Address.String(), err
}
