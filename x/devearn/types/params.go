package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	epochstypes "github.com/evmos/evmos/v10/x/epochs/types"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

// ParamsKey params store key
var (
	DefaultEnableDevEarn              bool    = true
	DefaultRewardEpochIdentifier      string  = epochstypes.WeekEpochID
	DefaultDevEarnInflationPercentage sdk.Dec = sdk.NewDecWithPrec(5, 2)
)

var (
	ParamsKey                               = []byte("Params")
	ParamStoreKeyEnableDevEarn              = []byte("EnableDevEarn")
	ParamStoreKeyRewardEpochIdentifier      = []byte("EarnEpochIdentifier")
	ParamStoreKeyDevEarnInflationPercentage = []byte("EarnInflationPercentage")
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	enableDevEarn bool,
	rewardEpochIdentifier string,
	devEarnInflationAPR sdk.Dec,
) Params {
	return Params{
		EnableDevEarn:         enableDevEarn,
		RewardEpochIdentifier: rewardEpochIdentifier,
		DevEarnInflation_APR:  devEarnInflationAPR,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultEnableDevEarn,
		DefaultRewardEpochIdentifier,
		DefaultDevEarnInflationPercentage,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyEnableDevEarn, &p.EnableDevEarn, validateBool),
		paramtypes.NewParamSetPair(ParamStoreKeyRewardEpochIdentifier, &p.RewardEpochIdentifier, epochstypes.ValidateEpochIdentifierInterface),
		paramtypes.NewParamSetPair(ParamStoreKeyDevEarnInflationPercentage, &p.DevEarnInflation_APR, validatePercentage),
	}
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// validateEnableDevEarn validates the EnableDevEarn param
func validateBool(v interface{}) error {
	_, ok := v.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	return nil
}

// validateDevEarnEpoch validates the DevEarnEpoch param
func validatePercentage(v interface{}) error {
	dec, ok := v.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}
	if dec.IsNegative() {
		return fmt.Errorf("DevEarnInflationPercentage must be positive: %s", dec)
	}
	if dec.GT(sdk.OneDec()) {
		return fmt.Errorf("DevEarnInflationPercentage must <= 100: %s", dec)
	}

	return nil
}
func (p Params) Validate() error {
	if err := validateBool(p.EnableDevEarn); err != nil {
		return err
	}

	if err := validatePercentage(p.DevEarnInflation_APR); err != nil {
		return err
	}

	return epochstypes.ValidateEpochIdentifierString(p.RewardEpochIdentifier)
}
