package devearn_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "sidechain/testutil/keeper"
	"sidechain/x/devearn"
	"sidechain/x/devearn/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		AssetsList: []types.Assets{
			{
				Id: 0,
			},
			{
				Id: 1,
			},
		},
		AssetsCount: 2,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.DevearnKeeper(t)
	devearn.InitGenesis(ctx, *k, genesisState)
	got := devearn.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	require.ElementsMatch(t, genesisState.AssetsList, got.AssetsList)
	require.Equal(t, genesisState.AssetsCount, got.AssetsCount)
	// this line is used by starport scaffolding # genesis/test/assert
}
