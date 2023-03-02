package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	KeyEnableDevEarn = []byte("EnableDevEarn")
	// TODO: Determine the default value
	DefaultEnableDevEarn bool = false
)

var (
	KeyDevEarnEpoch = []byte("DevEarnEpoch")
	// TODO: Determine the default value
	DefaultDevEarnEpoch string = "dev_earn_epoch"
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	enableDevEarn bool,
	devEarnEpoch string,
) Params {
	return Params{
		EnableDevEarn: enableDevEarn,
		DevEarnEpoch:  devEarnEpoch,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(
		DefaultEnableDevEarn,
		DefaultDevEarnEpoch,
	)
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyEnableDevEarn, &p.EnableDevEarn, validateEnableDevEarn),
		paramtypes.NewParamSetPair(KeyDevEarnEpoch, &p.DevEarnEpoch, validateDevEarnEpoch),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateEnableDevEarn(p.EnableDevEarn); err != nil {
		return err
	}

	if err := validateDevEarnEpoch(p.DevEarnEpoch); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// validateEnableDevEarn validates the EnableDevEarn param
func validateEnableDevEarn(v interface{}) error {
	enableDevEarn, ok := v.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	// TODO implement validation
	_ = enableDevEarn

	return nil
}

// validateDevEarnEpoch validates the DevEarnEpoch param
func validateDevEarnEpoch(v interface{}) error {
	devEarnEpoch, ok := v.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", v)
	}

	// TODO implement validation
	_ = devEarnEpoch

	return nil
}
