package keeper_test

import (
	"testing"

	testkeeper "github.com/sideprotocol/sidechain/testutil/keeper"
	"github.com/sideprotocol/sidechain/x/devearn/types"
	"github.com/stretchr/testify/require"
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
