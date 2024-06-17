package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// NewParams creates a new Params instance
func NewParams(senders []string) Params {
	return Params{
		QualifiedRelayers:       senders,
		Confirmations:           2,
		MaxAcceptableBlockDepth: 100,
		Vaults: []*Vault{{
			AddressOnBitcoin: "",
			AssetType:        AssetType_ASSET_TYPE_BTC,
		}},
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams([]string{})
}

// Validate validates the set of params
func (p Params) Validate() error {
	for _, sender := range p.QualifiedRelayers {
		_, err := sdk.AccAddressFromBech32(sender)
		if err != nil {
			return err
		}
	}
	return nil
}

// checks if the given address is an authorized sender
func (p Params) IsAuthorizedSender(sender string) bool {
	for _, s := range p.QualifiedRelayers {
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
		if v.AddressOnBitcoin == address {
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
