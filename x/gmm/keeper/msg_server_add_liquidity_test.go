package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/sideprotocol/side/app"
	"github.com/sideprotocol/side/testutil/sample"
	"github.com/sideprotocol/side/x/gmm/types"
)

func (suite *KeeperTestSuite) TestMsgAddLiquidity() {
	// Create a new pool
	poolID := suite.CreateNewPool()

	// Add liquidity to the pool
	msg := types.MsgAddLiquidity{
		Sender: types.Carol,
		PoolId: poolID,
		Liquidity: []sdk.Coin{
			sdk.NewCoin(
				simapp.DefaultBondDenom,
				sdk.NewInt(100),
			),
		},
	}
	ctx := sdk.WrapSDKContext(suite.ctx)
	res, err := suite.msgServer.AddLiquidity(ctx, &msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
}

func (suite *KeeperTestSuite) TestMsgAddLiquidityFail() {
	var msg *types.MsgAddLiquidity
	// Create a new pool
	poolID := suite.CreateNewPool()

	testCases := []struct {
		name   string
		mallet func()
	}{
		{
			"invalid sender",
			func() {},
		},
		{
			"invalid poolID",
			func() {
				msg = &types.MsgAddLiquidity{
					Sender: sample.AccAddress(),
					PoolId: "",
					Liquidity: []sdk.Coin{
						sdk.NewCoin(
							simapp.DefaultBondDenom,
							sdk.NewInt(100),
						),
					},
				}
			},
		},
		{
			"not enough funds",
			func() {
				msg = &types.MsgAddLiquidity{
					Sender: sample.AccAddress(),
					PoolId: poolID,
					Liquidity: []sdk.Coin{
						sdk.NewCoin(
							simapp.DefaultBondDenom,
							sdk.NewInt(100),
						),
					},
				}
			},
		},
	}

	for _, tc := range testCases {

		msg = types.NewMsgAddLiquidity(
			"",
			poolID,
			[]sdk.Coin{},
		)
		tc.mallet()

		res, err := suite.msgServer.AddLiquidity(sdk.WrapSDKContext(suite.ctx), msg)
		suite.Require().Error(err)
		suite.Require().Nil(res)
	}
}
