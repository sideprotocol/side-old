package keeper_test

import (
	"fmt"

	"sidechain/x/devearn/types"
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
				in := types.NewDevEarn(contract, 0, epochs)
				suite.app.DevearnKeeper.SetDevEarnInfo(suite.ctx, in)
				suite.Commit()

				expRes = []types.DevEarnInfo{in}
			},
		},
		{
			"2 pairs registered",
			func() {
				in := types.NewDevEarn(contract, 0, epochs)
				in2 := types.NewDevEarn(contract2, 0, epochs)
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

//func (suite *KeeperTestSuite) TestGetIncetive() {
//	suite.deployContracts()
//	expIn := types.NewIncentive(contract, allocations, epochs)
//	suite.app.IncentivesKeeper.SetIncentive(suite.ctx, expIn)
//	suite.Commit()
//
//	testCases := []struct {
//		name     string
//		contract common.Address
//		ok       bool
//	}{
//		{"nil address", common.Address{}, false},
//		{"valid id", common.HexToAddress(expIn.Contract), true},
//	}
//	for _, tc := range testCases {
//		in, found := suite.app.IncentivesKeeper.GetIncentive(suite.ctx, tc.contract)
//		if tc.ok {
//			suite.Require().True(found, tc.name)
//			suite.Require().Equal(expIn, in, tc.name)
//		} else {
//			suite.Require().False(found, tc.name)
//		}
//	}
//}
//
//func (suite *KeeperTestSuite) TestDeleteIncentiveAndUpdateAllocationMeters() {
//	suite.deployContracts()
//	// Register Incentive
//	_, err := suite.app.IncentivesKeeper.RegisterIncentive(
//		suite.ctx,
//		contract,
//		mintAllocations,
//		epochs,
//	)
//	suite.Require().NoError(err)
//
//	regIn, found := suite.app.IncentivesKeeper.GetIncentive(suite.ctx, contract)
//	suite.Require().True(found)
//
//	testCases := []struct {
//		name     string
//		malleate func()
//		ok       bool
//	}{
//		{"valid incentive", func() {}, true},
//		{
//			"deleted incentive",
//			func() {
//				suite.app.IncentivesKeeper.DeleteIncentiveAndUpdateAllocationMeters(suite.ctx, regIn)
//			},
//			false,
//		},
//	}
//	for _, tc := range testCases {
//		tc.malleate()
//		in, found := suite.app.IncentivesKeeper.GetIncentive(
//			suite.ctx,
//			common.HexToAddress(regIn.Contract),
//		)
//		if tc.ok {
//			suite.Require().True(found, tc.name)
//			suite.Require().Equal(regIn, in, tc.name)
//		} else {
//			suite.Require().False(found, tc.name)
//		}
//	}
//}
//
//func (suite *KeeperTestSuite) TestIsIncentiveRegistered() {
//	suite.deployContracts()
//	regIn := types.NewIncentive(contract, allocations, epochs)
//	suite.app.IncentivesKeeper.SetIncentive(suite.ctx, regIn)
//	suite.Commit()
//
//	testCases := []struct {
//		name     string
//		contract common.Address
//		ok       bool
//	}{
//		{"valid id", common.HexToAddress(regIn.Contract), true},
//		{"pair not found", common.Address{}, false},
//	}
//	for _, tc := range testCases {
//		found := suite.app.IncentivesKeeper.IsIncentiveRegistered(
//			suite.ctx,
//			tc.contract,
//		)
//		if tc.ok {
//			suite.Require().True(found, tc.name)
//		} else {
//			suite.Require().False(found, tc.name)
//		}
//	}
//}
