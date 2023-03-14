package types

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	sidetypes "sidechain/types"
	"time"
)

// NewIncentive returns an instance of Incentive
func NewDevEarn(
	contract common.Address,
	gasMeter uint64,
	epochs uint32,
	ownerAddr string,
) DevEarnInfo {
	return DevEarnInfo{
		Contract:     contract.String(),
		GasMeter:     gasMeter,
		StartTime:    time.Time{},
		OwnerAddress: ownerAddr,
		Epochs:       epochs,
	}
}

// Validate performs a stateless validation of a EarnInfo
func (d DevEarnInfo) Validate() error {
	if err := sidetypes.ValidateAddress(d.Contract); err != nil {
		return err
	}

	if d.Epochs == 0 {
		return fmt.Errorf("epoch cannot be 0")
	}
	return nil
}

// IsActive returns true if the Incentive has remaining Epochs
func (d DevEarnInfo) IsActive() bool {
	return d.Epochs > 0
}
