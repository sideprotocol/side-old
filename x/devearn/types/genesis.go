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
	// Check for duplicated ID in assets
	assetsIdMap := make(map[uint64]bool)
	assetsCount := gs.GetAssetsCount()
	for _, elem := range gs.AssetsList {
		if _, ok := assetsIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for assets")
		}
		if elem.Id >= assetsCount {
			return fmt.Errorf("assets id should be lower or equal than the last id")
		}
		assetsIdMap[elem.Id] = true
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
