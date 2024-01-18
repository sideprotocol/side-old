package keeper_test

import (
	"testing"

	testkeeper "github.com/sideprotocol/side/testutil/keeper"
	"github.com/sideprotocol/side/x/router/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.RouterKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
