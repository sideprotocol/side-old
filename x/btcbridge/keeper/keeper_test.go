package keeper_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"lukechampine.com/uint128"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/btcutil/bech32"
	"github.com/cosmos/cosmos-sdk/crypto/keys/segwit"
	sdk "github.com/cosmos/cosmos-sdk/types"

	simapp "github.com/sideprotocol/side/app"
	"github.com/sideprotocol/side/x/btcbridge/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *simapp.App

	btcVault   string
	runesVault string
	sender     string

	btcVaultPkScript   []byte
	runesVaultPkScript []byte
	senderPkScript     []byte
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(suite.T())
	ctx := app.BaseApp.NewContext(false, tmproto.Header{Time: time.Now().UTC()})

	suite.ctx = ctx
	suite.app = app

	chainCfg := sdk.GetConfig().GetBtcChainCfg()

	suite.btcVault, _ = bech32.Encode(chainCfg.Bech32HRPSegwit, segwit.GenPrivKey().PubKey().Address().Bytes())
	suite.runesVault, _ = bech32.Encode(chainCfg.Bech32HRPSegwit, segwit.GenPrivKey().PubKey().Address())
	suite.sender, _ = bech32.Encode(chainCfg.Bech32HRPSegwit, segwit.GenPrivKey().PubKey().Address())

	suite.btcVaultPkScript = MustPkScriptFromAddress(suite.btcVault, chainCfg)
	suite.runesVaultPkScript = MustPkScriptFromAddress(suite.runesVault, chainCfg)
	suite.senderPkScript = MustPkScriptFromAddress(suite.sender, chainCfg)

	suite.setupParams(suite.btcVault, suite.runesVault)
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) setupParams(btcVault string, runesVault string) {
	suite.app.BtcBridgeKeeper.SetParams(suite.ctx, types.Params{Vaults: []*types.Vault{
		{
			Address:   btcVault,
			AssetType: types.AssetType_ASSET_TYPE_BTC,
		},
		{
			Address:   runesVault,
			AssetType: types.AssetType_ASSET_TYPE_RUNE,
		},
	}})
}

func (suite *KeeperTestSuite) setupUTXOs(utxos []*types.UTXO) {
	for _, utxo := range utxos {
		suite.app.BtcBridgeKeeper.SetUTXO(suite.ctx, utxo)
		suite.app.BtcBridgeKeeper.SetOwnerUTXO(suite.ctx, utxo)

		for _, r := range utxo.Runes {
			suite.app.BtcBridgeKeeper.SetOwnerRunesUTXO(suite.ctx, utxo, r.Id, r.Amount)
		}
	}
}

func (suite *KeeperTestSuite) TestMintRunes() {
	runeId := "840000:3"
	runeAmount := 500000000
	runeOutputIndex := 2

	runesScript, err := types.BuildEdictScript(runeId, uint128.From64(uint64(runeAmount)), uint32(runeOutputIndex))
	suite.NoError(err)

	tx := wire.NewMsgTx(types.TxVersion)
	tx.AddTxOut(wire.NewTxOut(0, runesScript))
	tx.AddTxOut(wire.NewTxOut(types.RunesOutValue, suite.senderPkScript))
	tx.AddTxOut(wire.NewTxOut(types.RunesOutValue, suite.runesVaultPkScript))

	denom := fmt.Sprintf("%s/%s", types.RunesProtocolName, runeId)

	balanceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(suite.sender), denom)
	suite.True(balanceBefore.Amount.IsZero(), "%s balance before mint should be zero", denom)

	recipient, err := suite.app.BtcBridgeKeeper.Mint(suite.ctx, btcutil.NewTx(tx), btcutil.NewTx(tx), 0)
	suite.NoError(err)
	suite.Equal(suite.sender, recipient.EncodeAddress(), "incorrect recipient")

	balanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(suite.sender), denom)
	suite.Equal(uint64(runeAmount), balanceAfter.Amount.Uint64(), "%s balance after mint should be %d", denom, runeAmount)

	utxos := suite.app.BtcBridgeKeeper.GetAllUTXOs(suite.ctx)
	suite.Len(utxos, 1, "there should be 1 utxo(s)")

	expectedUTXO := &types.UTXO{
		Txid:         tx.TxHash().String(),
		Vout:         uint64(runeOutputIndex),
		Address:      suite.runesVault,
		Amount:       types.RunesOutValue,
		PubKeyScript: suite.runesVaultPkScript,
		IsLocked:     false,
		Runes: []*types.RuneBalance{
			{
				Id:     runeId,
				Amount: fmt.Sprintf("%d", runeAmount),
			},
		},
	}

	suite.Equal(expectedUTXO, utxos[0], "utxos do not match")
}

func (suite *KeeperTestSuite) TestWithdrawRunes() {
	runeId := "840000:3"
	runeAmount := 500000000

	runesUTXOs := []*types.UTXO{
		{
			Txid:         chainhash.HashH([]byte("runes")).String(),
			Vout:         1,
			Address:      suite.runesVault,
			Amount:       types.RunesOutValue,
			PubKeyScript: suite.runesVaultPkScript,
			IsLocked:     false,
			Runes: []*types.RuneBalance{
				{
					Id:     runeId,
					Amount: fmt.Sprintf("%d", runeAmount),
				},
			},
		},
	}
	suite.setupUTXOs(runesUTXOs)

	feeRate := 100
	amount := runeAmount + 1

	denom := fmt.Sprintf("%s/%s", types.RunesProtocolName, runeId)
	coin := sdk.NewCoin(denom, sdk.NewInt(int64(amount)))

	_, err := suite.app.BtcBridgeKeeper.NewSigningRequest(suite.ctx, suite.sender, coin, int64(feeRate))
	suite.ErrorIs(err, types.ErrInsufficientUTXOs, "should fail due to insufficient runes utxos")

	amount = 100000000
	coin = sdk.NewCoin(denom, sdk.NewInt(int64(amount)))

	_, err = suite.app.BtcBridgeKeeper.NewSigningRequest(suite.ctx, suite.sender, coin, int64(feeRate))
	suite.ErrorIs(err, types.ErrInsufficientUTXOs, "should fail due to insufficient payment utxos")

	paymentUTXOs := []*types.UTXO{
		{
			Txid:         chainhash.HashH([]byte("payment")).String(),
			Vout:         1,
			Address:      suite.btcVault,
			Amount:       100000,
			PubKeyScript: suite.btcVaultPkScript,
			IsLocked:     false,
		},
	}
	suite.setupUTXOs(paymentUTXOs)

	req, err := suite.app.BtcBridgeKeeper.NewSigningRequest(suite.ctx, suite.sender, coin, int64(feeRate))
	suite.NoError(err)

	suite.True(suite.app.BtcBridgeKeeper.IsUTXOLocked(suite.ctx, runesUTXOs[0].Txid, runesUTXOs[0].Vout), "runes utxo should be locked")
	suite.True(suite.app.BtcBridgeKeeper.IsUTXOLocked(suite.ctx, paymentUTXOs[0].Txid, paymentUTXOs[0].Vout), "payment utxo should be locked")

	runesUTXOs = suite.app.BtcBridgeKeeper.GetUnlockedUTXOsByAddr(suite.ctx, suite.runesVault)
	suite.Len(runesUTXOs, 1, "there should be 1 unlocked runes utxo(s)")

	suite.Len(runesUTXOs[0].Runes, 1, "there should be 1 rune in the runes utxo")
	suite.Equal(runeId, runesUTXOs[0].Runes[0].Id, "incorrect rune id")
	suite.Equal(uint64(runeAmount-amount), types.RuneAmountFromString(runesUTXOs[0].Runes[0].Amount).Big().Uint64(), "incorrect rune amount")

	p, err := psbt.NewFromRawBytes(bytes.NewReader([]byte(req.Psbt)), true)
	suite.NoError(err)

	suite.Len(p.Inputs, 2, "there should be 2 inputs")
	suite.Equal(suite.runesVaultPkScript, p.Inputs[0].WitnessUtxo.PkScript, "the first input should be runes vault")
	suite.Equal(suite.btcVaultPkScript, p.Inputs[1].WitnessUtxo.PkScript, "the second input should be btc vault")

	expectedRunesScript, err := types.BuildEdictScript(runeId, uint128.From64(uint64(amount)), 2)
	suite.NoError(err)

	suite.Len(p.UnsignedTx.TxOut, 4, "there should be 4 outputs")
	suite.Equal(expectedRunesScript, p.UnsignedTx.TxOut[0].PkScript, "incorrect runes script")
	suite.Equal(suite.runesVaultPkScript, p.UnsignedTx.TxOut[1].PkScript, "the second output should be runes change output")
	suite.Equal(suite.senderPkScript, p.UnsignedTx.TxOut[2].PkScript, "the third output should be sender output")
	suite.Equal(suite.btcVaultPkScript, p.UnsignedTx.TxOut[3].PkScript, "the fouth output should be btc change output")
}

func MustPkScriptFromAddress(addr string, chainCfg *chaincfg.Params) []byte {
	address, err := btcutil.DecodeAddress(addr, chainCfg)
	if err != nil {
		panic(err)
	}

	pkScript, err := txscript.PayToAddrScript(address)
	if err != nil {
		panic(err)
	}

	return pkScript
}
