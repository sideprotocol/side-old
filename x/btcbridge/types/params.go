package types

import (
	"bytes"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewParams creates a new Params instance
func NewParams(relayers []string) Params {
	return Params{
		AuthorizedRelayers:      relayers,
		Confirmations:           2,
		MaxAcceptableBlockDepth: 100,
		BtcVoucherDenom:         "sat",
		Vaults: []*Vault{{
			Address:   "",
			PubKey:    "",
			AssetType: AssetType_ASSET_TYPE_BTC,
		}, {
			Address:   "",
			PubKey:    "",
			AssetType: AssetType_ASSET_TYPE_RUNE,
		}},
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams([]string{})
}

// Validate validates the set of params
func (p Params) Validate() error {
	for _, sender := range p.AuthorizedRelayers {
		_, err := sdk.AccAddressFromBech32(sender)
		if err != nil {
			return err
		}
	}
	return nil
}

// checks if the given address is an authorized sender
func (p Params) IsAuthorizedSender(sender string) bool {
	for _, s := range p.AuthorizedRelayers {
		if s == sender {
			return true
		}
	}
	return false
}

// SelectVaultByBitcoinAddress returns the vault if the address is found
// returns the vault if the address is found
func SelectVaultByBitcoinAddress(vaults []*Vault, address string) *Vault {
	for _, v := range vaults {
		if v.Address == address {
			return v
		}
	}
	return nil
}

// SelectVaultByPubKey returns the vault if the public key is found
// returns the vault if the public key is found
func SelectVaultByPubKey(vaults []*Vault, pubKey string) *Vault {
	for _, v := range vaults {
		if v.PubKey == pubKey {
			return v
		}
	}

	return nil
}

// SelectVaultByAssetType returns the vault by the asset type
func SelectVaultByAssetType(vaults []*Vault, assetType AssetType) *Vault {
	for _, v := range vaults {
		if v.AssetType == assetType {
			return v
		}
	}

	return nil
}

// SelectVaultByPkScript returns the vault by the pk script
func SelectVaultByPkScript(vaults []*Vault, pkScript []byte) *Vault {
	chainCfg := sdk.GetConfig().GetBtcChainCfg()

	for _, v := range vaults {
		addr, err := btcutil.DecodeAddress(v.Address, chainCfg)
		if err != nil {
			continue
		}

		addrScript, err := txscript.PayToAddrScript(addr)
		if err != nil {
			continue
		}

		if bytes.Equal(addrScript, pkScript) {
			return v
		}
	}

	return nil
}
