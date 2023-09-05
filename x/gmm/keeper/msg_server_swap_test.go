package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/sideprotocol/side/app"
	"github.com/sideprotocol/side/x/gmm/types"
)

func (suite *KeeperTestSuite) TestMsgSwap() {
	// Create a new pool
	poolID := suite.CreateNewPool()

	// Add liquidity to the pool
	msg := types.MsgSwap{
		Creator:  types.Alice,
		PoolId:   poolID,
		TokenIn:  sdk.NewCoin(simapp.DefaultBondDenom, sdk.NewInt(100)),
		DenomOut: simapp.AltDenom,
	}

	ctx := sdk.WrapSDKContext(suite.ctx)
	queryResBeforeSwap, err := suite.queryClient.Pool(ctx, &types.QueryPoolRequest{
		PoolId: poolID,
	})

	outAssetBeforeSwap := queryResBeforeSwap.Pool.Assets[msg.DenomOut]

	estimatedOut, err := queryResBeforeSwap.Pool.EstimateSwap(msg.TokenIn, simapp.AltDenom)
	suite.Require().NoError(err)

	res, err := suite.msgServer.Swap(ctx, &msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	queryResAfterSwap, err := suite.queryClient.Pool(ctx, &types.QueryPoolRequest{
		PoolId: poolID,
	})
	outAssetAfterSwap := queryResAfterSwap.Pool.Assets[msg.DenomOut]

	out := outAssetBeforeSwap.Token.Sub(outAssetAfterSwap.Token)
	suite.Require().Equal(out, estimatedOut)
}
