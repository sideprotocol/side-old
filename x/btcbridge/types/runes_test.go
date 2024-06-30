package types_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/btcsuite/btcd/wire"

	"github.com/sideprotocol/side/x/btcbridge/types"
)

func TestParseRunes(t *testing.T) {
	testCases := []struct {
		name        string
		pkScriptHex string
		edicts      []*types.Edict
		expectPass  bool
	}{
		{
			name:        "valid runes edict",
			pkScriptHex: "6a5d0b00c0a2330380cab5ee0101",
			edicts: []*types.Edict{
				{
					Id:     &types.RuneId{Block: 840000, Tx: 3},
					Amount: "500000000",
					Output: 1,
				},
			},
			expectPass: true,
		},
		{
			name:        "output index is out of range",
			pkScriptHex: "6a5d0b00c0a2330380cab5ee0102",
			expectPass:  false,
		},
		{
			name:        "no OP_RETURN",
			pkScriptHex: "615d0b00c0a2330380cab5ee0102",
			expectPass:  true,
			edicts:      nil,
		},
		{
			name:        "no runes magic number",
			pkScriptHex: "6a5c0b00c0a2330380cab5ee0102",
			expectPass:  true,
			edicts:      nil,
		},
		{
			name:        "non data push op",
			pkScriptHex: "6a5d4f00c0a2330380cab5ee0102",
			expectPass:  false,
		},
		{
			name:        "no tag body for edicts",
			pkScriptHex: "6a5d0b01c0a2330380cab5ee0102",
			expectPass:  false,
		},
		{
			name:        "invalid edict",
			pkScriptHex: "6a5d0b00c0a2330380cab5ee01",
			expectPass:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pkScript, err := hex.DecodeString(tc.pkScriptHex)
			require.NoError(t, err)

			tx := wire.NewMsgTx(types.TxVersion)
			tx.AddTxOut(wire.NewTxOut(0, pkScript))

			edicts, err := types.ParseRunes(tx)
			if tc.expectPass {
				require.NoError(t, err)
				require.EqualValues(t, tc.edicts, edicts)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}
