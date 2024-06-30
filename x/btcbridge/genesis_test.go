package btcbridge_test

// Path: x/btcbridge/genesis_test.go

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
	keepertest "github.com/sideprotocol/side/testutil/keeper"
	"github.com/sideprotocol/side/testutil/nullify"
	"github.com/sideprotocol/side/x/btcbridge"
	"github.com/sideprotocol/side/x/btcbridge/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	mnemonic := "sunny bamboo garlic fold reopen exile letter addict forest vessel square lunar shell number deliver cruise calm artist fire just kangaroo suit wheel extend"
	println(mnemonic)

	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.BtcBridgeKeeper(t)
	btcbridge.InitGenesis(ctx, *k, genesisState)
	got := btcbridge.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}

// TestSubmitTx tests the SubmitTx function
// func TestSubmitTx(t *testing.T) {

// 	// test tx: https://blockchain.info/tx/b657e22827039461a9493ede7bdf55b01579254c1630b0bfc9185ec564fc05ab?format=json

// 	k, ctx := keepertest.BtcbridgeClientKeeper(t)

// 	txHex := "02000000000101e5df234dbeff74d6da14754ebeea8ab0e2f60b1884566846cf6b36e8ceb5f5350100000000fdffffff02f0490200000000001600142ac13073e8d491428790481321a636696d00107681d7e205000000001600142bf3aa186cbdcbe88b70a67edcd5a32ce5e8e6d8024730440220081ee61d749ce8cedcf6eedde885579af2eb65ca67d29e6ae2c37109d81cbbb202203a1891ce45f91f159ccf04348ef37a3d1a12d89e5e01426e061326057e6c128d012103036bbdd77c9a932f37bd66175967c7fb7eb75ece06b87c1ad1716770cb3ca4ee79fc2a00"
// 	prevTxHex := "0200000000010183372652f2af9ab34b3a003efada6b054c75583185ac130de72599dfdf4e462b0100000000fdffffff02f0490200000000001600142ac13073e8d491428790481321a636696d001076a64ee50500000000160014a03614eef338681373de94a2dc2574de55da1980024730440220274250f6036bea0947daf4455ab4976f81721257d163fd952fb5b0c70470edc602202fba816be260219bbc40a8983c459cf05cf2209bf1e62e7ccbf78aec54db607f0121031cee21ef69fe68b240c3032616fa310c6a60a856c0a7e0c1298815c92fb2c61788fb2a00"

// 	msg := &types.MsgSubmitTransactionRequest{
// 		Sender:      "",
// 		Blockhash:   "000000000d73ecf25d3bf8e6ae65c35aa2a90e3271edff8bab90d87ed875f13b",
// 		TxBytes:     "0100000001b3f7",
// 		PrevTxBytes: "0100000001b3f7",
// 		Proof:       []string{"0100000001b3f7"},
// 	}
// 	// this line is used by starport scaffolding # handler/test/submit
// 	err := k.ProcessBitcoinDepositTransaction(ctx, msg)
// 	require.NoError(t, err)
// }

// Decode transaction
func TestDecodeTransaction(t *testing.T) {
	hexStr := "02000000000101e5df234dbeff74d6da14754ebeea8ab0e2f60b1884566846cf6b36e8ceb5f5350100000000fdffffff02f0490200000000001600142ac13073e8d491428790481321a636696d00107681d7e205000000001600142bf3aa186cbdcbe88b70a67edcd5a32ce5e8e6d8024730440220081ee61d749ce8cedcf6eedde885579af2eb65ca67d29e6ae2c37109d81cbbb202203a1891ce45f91f159ccf04348ef37a3d1a12d89e5e01426e061326057e6c128d012103036bbdd77c9a932f37bd66175967c7fb7eb75ece06b87c1ad1716770cb3ca4ee79fc2a00"

	// Decode the hex string to transaction
	txBytes, err := hex.DecodeString(hexStr)
	require.NoError(t, err)

	// Create a new transaction
	var tx wire.MsgTx
	err = tx.Deserialize(bytes.NewReader(txBytes))
	require.NoError(t, err)

	uTx := btcutil.NewTx(&tx)

	for _, input := range uTx.MsgTx().TxIn {
		t.Log(input.PreviousOutPoint.String())
	}

	require.GreaterOrEqual(t, len(uTx.MsgTx().TxIn), 1)
}
