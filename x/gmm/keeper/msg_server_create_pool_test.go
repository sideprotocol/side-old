package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/gmm/keeper"
	"github.com/sideprotocol/side/x/gmm/types"
)

func (suite *KeeperTestSuite) TestMsgMakePool() {
	var msg *types.MsgCreatePool

	testCases := []struct {
		name    string
		expPass bool
	}{
		{
			"invalid sender",
			false,
		},
		// {
		// 	"channel does not exist",
		// 	false,
		// },
	}

	for _, tc := range testCases {
		suite.SetupTest()

		msg = types.NewMsgCreatePool(
			"",
			types.PoolParams{},
			[]types.PoolAsset{},
		)
		msgSrv := keeper.NewMsgServerImpl(suite.app.GmmKeeper)
		res, err := msgSrv.CreatePool(sdk.WrapSDKContext(suite.ctx), msg)

		if tc.expPass {
			suite.Require().NoError(err)
			suite.Require().NotNil(res)
		} else {
			suite.Require().Error(err)
			suite.Require().Nil(res)
		}
	}
}
