package types

import (
	"fmt"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		AssetsList: []Assets{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated denom in assets
	assetsIdMap := make(map[string]bool)
	for _, elem := range gs.AssetsList {
		if _, ok := assetsIdMap[elem.Denom]; ok {
			return fmt.Errorf("duplicated id for assets")
		}
		assetsIdMap[elem.Denom] = true
	}
	// this line is used by starport scaffolding # genesis/types/validate
	return gs.Params.Validate()
}
