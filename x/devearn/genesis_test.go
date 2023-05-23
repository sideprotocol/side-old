package devearn_test

import (
	"testing"

	keepertest "sidechain/testutil/keeper"
	"sidechain/x/devearn"
	"sidechain/x/devearn/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		AssetsList: []types.Assets{
			{
				Denom: "aside",
			},
			{
				Denom: "bside",
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.DevearnKeeper(t)
	devearn.InitGenesis(ctx, *k, genesisState)
	got := devearn.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	require.ElementsMatch(t, genesisState.AssetsList, got.AssetsList)
	// this line is used by starport scaffolding # genesis/test/assert
}
