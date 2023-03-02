package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "sidechain/testutil/keeper"
	"sidechain/x/devearn/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.DevearnKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
	require.EqualValues(t, params.EnableDevEarn, k.EnableDevEarn(ctx))
	require.EqualValues(t, params.DevEarnEpoch, k.DevEarnEpoch(ctx))
}
