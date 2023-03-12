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
	require.EqualValues(t, params.EnableDevEarn, k.GetParams(ctx).EnableDevEarn)
	require.EqualValues(t, params.RewardEpochIdentifier, k.GetParams(ctx).RewardEpochIdentifier)
	require.EqualValues(t, params.DevEarnInflation_APR, k.GetParams(ctx).DevEarnInflation_APR)
}
