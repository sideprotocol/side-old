package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	keepertest "github.com/sideprotocol/side/testutil/keeper"
	"github.com/sideprotocol/side/testutil/nullify"
	"github.com/sideprotocol/side/x/yield/keeper"
	"github.com/sideprotocol/side/x/yield/types"
)

func createNHostChain(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.HostChain {
	items := make([]types.HostChain, n)
	for i := range items {
		items[i].ChainId = strconv.Itoa(i)
		keeper.SetHostChain(ctx, items[i])
	}
	return items
}

func TestHostChainGet(t *testing.T) {
	keeper, ctx := keepertest.YieldKeeper(t)
	items := createNHostChain(keeper, ctx, 10)
	for _, item := range items {
		got, found := keeper.GetHostChain(ctx, item.ChainId)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&got),
		)
	}
}

func TestHostChainRemove(t *testing.T) {
	keeper, ctx := keepertest.YieldKeeper(t)
	items := createNHostChain(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveHostChain(ctx, item.ChainId)
		_, found := keeper.GetHostChain(ctx, item.ChainId)
		require.False(t, found)
	}
}

func TestHostChainGetAll(t *testing.T) {
	keeper, ctx := keepertest.YieldKeeper(t)
	items := createNHostChain(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllHostChain(ctx)),
	)
}

// func (s *KeeperTestSuite) TestGetHostChainFromTransferChannelID() {
// 	// Store 5 host chains
// 	expectedHostChains := map[string]types.HostChain{}
// 	for i := 0; i < 5; i++ {
// 		chainId := fmt.Sprintf("chain-%d", i)
// 		channelId := fmt.Sprintf("channel-%d", i)

// 		HostChain := types.HostChain{
// 			ChainId:           chainId,
// 			TransferChannelId: channelId,
// 		}
// 		s.App.StakeibcKeeper.SetHostChain(s.Ctx, HostChain)
// 		expectedHostChains[channelId] = HostChain
// 	}

// 	// Look up each host chain by the channel ID
// 	for i := 0; i < 5; i++ {
// 		channelId := fmt.Sprintf("channel-%d", i)

// 		expectedHostChain := expectedHostChains[channelId]
// 		actualHostChain, found := s.App.StakeibcKeeper.GetHostChainFromTransferChannelID(s.Ctx, channelId)

// 		s.Require().True(found, "found host chain %d", i)
// 		s.Require().Equal(expectedHostChain.ChainId, actualHostChain.ChainId, "host chain %d chain-id", i)
// 	}

// 	// Lookup a non-existent host chain - should not be found
// 	_, found := s.App.StakeibcKeeper.GetHostChainFromTransferChannelID(s.Ctx, "fake_channel")
// 	s.Require().False(found, "fake channel should not be found")
// }
