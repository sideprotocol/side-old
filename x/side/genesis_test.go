package side_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "sidechain/testutil/keeper"
	"sidechain/testutil/nullify"
	"sidechain/x/side"
	"sidechain/x/side/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.SideKeeper(t)
	side.InitGenesis(ctx, *k, genesisState)
	got := side.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
