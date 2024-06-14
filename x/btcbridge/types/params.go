package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// NewParams creates a new Params instance
func NewParams(senders []string) Params {
	return Params{
		Senders:                 senders,
		Confirmations:           2,
		MaxAcceptableBlockDepth: 100,
		BtcVoucherDenom:         "sat",
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams([]string{})
}

// Validate validates the set of params
func (p Params) Validate() error {
	for _, sender := range p.Senders {
		_, err := sdk.AccAddressFromBech32(sender)
		if err != nil {
			return err
		}
	}
	return nil
}

// checks if the given address is an authorized sender
func (p Params) IsAuthorizedSender(sender string) bool {
	for _, s := range p.Senders {
		if s == sender {
			return true
		}
	}
	return false
}
