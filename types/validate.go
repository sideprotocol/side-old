package types

import (
	errorsmod "cosmossdk.io/errors"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
)

func ValidateAddress(address string) error {
	if !common.IsHexAddress(address) {
		return errorsmod.Wrapf(
			errortypes.ErrInvalidAddress, "address '%s' is not a valid ethereum hex address",
			address,
		)
	}
	return nil
}
