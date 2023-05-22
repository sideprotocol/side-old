package keeper_test

import (
	"fmt"
	"sidechain/x/devearn/types"

	"github.com/ethereum/go-ethereum/common"
)

func (suite *KeeperTestSuite) TestGetAllDevInfos() {
	var expRes []types.DevEarnInfo

	testCases := []struct {
		name     string
		malleate func()
	}{
		{
			"no incentive registered",
			func() { expRes = []types.DevEarnInfo{} },
		},
		{
			"1 pair registered",
			func() {
				in := types.NewDevEarn(contract, 0, epochs, suite.priv.PubKey().Address().String())
				suite.app.DevearnKeeper.SetDevEarnInfo(suite.ctx, in)
				suite.Commit()

				expRes = []types.DevEarnInfo{in}
			},
		},
		{
			"2 pairs registered",
			func() {
				in := types.NewDevEarn(contract, 0, epochs, suite.priv.PubKey().Address().String())
				in2 := types.NewDevEarn(contract2, 0, epochs, suite.priv.PubKey().Address().String())
				suite.app.DevearnKeeper.SetDevEarnInfo(suite.ctx, in)
				suite.app.DevearnKeeper.SetDevEarnInfo(suite.ctx, in2)
				suite.Commit()

				expRes = []types.DevEarnInfo{in, in2}
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			suite.deployContracts()
			tc.malleate()
			res := suite.app.DevearnKeeper.GetAllDevEarnInfos(suite.ctx)
			suite.Require().ElementsMatch(expRes, res, tc.name)
		})
	}
}

func (suite *KeeperTestSuite) TestGetDevInfo() {
	suite.deployContracts()
	expIn := types.NewDevEarn(contract, 0, epochs, suite.priv.PubKey().Address().String())
	suite.app.DevearnKeeper.SetDevEarnInfo(suite.ctx, expIn)
	suite.Commit()

	testCases := []struct {
		name     string
		contract common.Address
		ok       bool
	}{
		{"nil address", common.Address{}, false},
		{"valid id", common.HexToAddress(expIn.Contract), true},
	}
	for _, tc := range testCases {
		in, found := suite.app.DevearnKeeper.GetDevEarnInfo(suite.ctx, tc.contract)
		if tc.ok {
			suite.Require().True(found, tc.name)
			suite.Require().Equal(expIn, in, tc.name)
		} else {
			suite.Require().False(found, tc.name)
		}
	}
}

func (suite *KeeperTestSuite) TestDeleteDevEarnInfo() {
	suite.deployContracts()
	// Register Incentive
	_, err := suite.app.DevearnKeeper.RegisterDevEarnInfo(
		suite.ctx,
		contract,
		epochs,
		suite.priv.PubKey().Address().String(),
	)
	suite.Require().NoError(err)

	regIn, found := suite.app.DevearnKeeper.GetDevEarnInfo(suite.ctx, contract)
	suite.Require().True(found)

	testCases := []struct {
		name     string
		malleate func()
		ok       bool
	}{
		{"valid incentive", func() {}, true},
		{
			"deleted incentive",
			func() {
				suite.app.DevearnKeeper.DeleteDevEarnInfo(suite.ctx, regIn)
			},
			false,
		},
	}
	for _, tc := range testCases {
		tc.malleate()
		in, found := suite.app.DevearnKeeper.GetDevEarnInfo(
			suite.ctx,
			contract,
		)
		if tc.ok {
			suite.Require().True(found, tc.name)
			suite.Require().Equal(regIn, in, tc.name)
		} else {
			suite.Require().False(found, tc.name)
		}
	}
}

func (suite *KeeperTestSuite) TestIsDevEarnRegistered() {
	suite.SetupTest() // reset
	suite.deployContracts()
	regIn := types.NewDevEarn(contract, 0, epochs, suite.priv.PubKey().Address().String())
	suite.app.DevearnKeeper.SetDevEarnInfo(suite.ctx, regIn)
	suite.Commit()

	testCases := []struct {
		name     string
		contract common.Address
		ok       bool
	}{
		{"valid id", common.HexToAddress(regIn.Contract), true},
		{"pair not found", common.Address{}, false},
	}
	for _, tc := range testCases {
		found := suite.app.DevearnKeeper.IsDevEarnInfoRegistered(
			suite.ctx,
			tc.contract,
		)
		if tc.ok {
			suite.Require().True(found, tc.name)
		} else {
			suite.Require().False(found, tc.name)
		}
	}
}
