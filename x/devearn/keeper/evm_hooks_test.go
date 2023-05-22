package keeper_test

import (
	"fmt"
	"math/big"
	"sidechain/testutil"
	utiltx "sidechain/testutil/tx"
	"sidechain/x/devearn/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	ethermint "github.com/evmos/ethermint/types"
	evm "github.com/evmos/ethermint/x/evm/types"
	vestingtypes "github.com/evmos/evmos/v11/x/vesting/types"
)

func (suite *KeeperTestSuite) TestEvmHooksStoreTxGasUsed() {
	var expGasUsed uint64

	testCases := []struct {
		name     string
		malleate func(common.Address)

		expPass bool
	}{
		{
			" dev earn is disabled globally",
			func(_ common.Address) {
				params := types.DefaultParams()
				params.EnableDevEarn = false
				suite.app.DevearnKeeper.SetParams(suite.ctx, params) //nolint:errcheck
			},
			false,
		},
		{
			"correct execution - one tx",
			func(contractAddr common.Address) {
				acc := &ethermint.EthAccount{
					BaseAccount: authtypes.NewBaseAccount(sdk.AccAddress(suite.address.Bytes()), nil, 0, 0),
					CodeHash:    common.BytesToHash(crypto.Keccak256(nil)).String(),
				}
				suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

				res := suite.MintERC20Token(contractAddr, suite.address, suite.address, big.NewInt(1000))
				expGasUsed = res.AsTransaction().Gas()
			},
			true,
		},
		{
			"correct execution with Base account - one tx",
			func(contractAddr common.Address) {
				acc := authtypes.NewBaseAccount(sdk.AccAddress(suite.address.Bytes()), nil, 0, 0)
				suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

				res := suite.MintERC20Token(contractAddr, suite.address, suite.address, big.NewInt(1000))
				expGasUsed = res.AsTransaction().Gas()
			},
			true,
		},
		{
			"correct execution with Vesting account - one tx",
			func(contractAddr common.Address) {
				acc := vestingtypes.NewClawbackVestingAccount(
					authtypes.NewBaseAccount(sdk.AccAddress(suite.address.Bytes()), nil, 0, 0),
					suite.address.Bytes(), nil, suite.ctx.BlockTime(), nil, nil,
				)

				suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

				res := suite.MintERC20Token(contractAddr, suite.address, suite.address, big.NewInt(1000))
				expGasUsed = res.AsTransaction().Gas()
			},
			true,
		},
		{
			"correct execution - two tx",
			func(contractAddr common.Address) {
				acc := &ethermint.EthAccount{
					BaseAccount: authtypes.NewBaseAccount(sdk.AccAddress(suite.address.Bytes()), nil, 0, 0),
					CodeHash:    common.BytesToHash(crypto.Keccak256(nil)).String(),
				}
				suite.app.AccountKeeper.SetAccount(suite.ctx, acc)

				res := suite.MintERC20Token(contractAddr, suite.address, suite.address, big.NewInt(500))
				res2 := suite.MintERC20Token(contractAddr, suite.address, suite.address, big.NewInt(500))
				expGasUsed = res.AsTransaction().Gas() + res2.AsTransaction().Gas()
			},
			true,
		},
		{
			"tx with non-incentivized contract",
			func(_ common.Address) {
				_ = suite.MintERC20Token(utiltx.GenerateAddress(), suite.address, suite.address, big.NewInt(1000))
			},
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.mintFeeCollector = true
			suite.SetupTest()
			suite.ensureHooksSet()
			// Deploy Contract
			contractAddr, err := suite.DeployContract("testcoin", "COIN", erc20Decimals)
			suite.Require().NoError(err)
			suite.Commit()

			// Register devearn
			_, err = suite.app.DevearnKeeper.RegisterDevEarnInfo(
				suite.ctx,
				contractAddr,
				epochs,
				"",
			)
			suite.Require().NoError(err)

			// Mint coins to pay gas fee
			coins := sdk.NewCoins(sdk.NewCoin(evm.DefaultEVMDenom, sdk.NewInt(30000000)))
			err = testutil.FundAccount(suite.ctx, suite.app.BankKeeper, sdk.AccAddress(suite.address.Bytes()), coins)
			suite.Require().NoError(err)

			// Submit tx
			tc.malleate(contractAddr)
			DevEarnInfo, found := suite.app.DevearnKeeper.GetDevEarnInfo(suite.ctx, contractAddr)
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().True(found)
				suite.Require().NotZero(DevEarnInfo.GasMeter)
				suite.Require().Equal(expGasUsed, DevEarnInfo.GasMeter)
			} else {
				suite.Require().NoError(err)
				suite.Require().Zero(DevEarnInfo.GasMeter)
			}
		})
	}
	suite.mintFeeCollector = false
}
