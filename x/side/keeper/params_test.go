package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "sidechain/testutil/keeper"
	"sidechain/x/side/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.SideKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
