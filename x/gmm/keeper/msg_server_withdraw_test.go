package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/gmm/types"
)

func (suite *KeeperTestSuite) TestMsgWithdraw() {
	// Create a new pool
	poolID := suite.CreateNewPool()

	// Add liquidity to the pool
	msg := types.MsgWithdraw{
		Creator:  types.Alice,
		Receiver: types.Carol,
		PoolId:   poolID,
		Share:    sdk.NewCoin(poolID, sdk.NewInt(100)),
	}

	ctx := sdk.WrapSDKContext(suite.ctx)
	res, err := suite.msgServer.Withdraw(ctx, &msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
}
