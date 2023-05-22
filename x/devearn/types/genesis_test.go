package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"sidechain/x/devearn/types"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{

				AssetsList: []types.Assets{
					{
						Id: 0,
					},
					{
						Id: 1,
					},
				},
				AssetsCount: 2,
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated assets",
			genState: &types.GenesisState{
				AssetsList: []types.Assets{
					{
						Id: 0,
					},
					{
						Id: 0,
					},
				},
			},
			valid: false,
		},
		{
			desc: "invalid assets count",
			genState: &types.GenesisState{
				AssetsList: []types.Assets{
					{
						Id: 1,
					},
				},
				AssetsCount: 0,
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
