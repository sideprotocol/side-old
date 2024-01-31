package types

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgSwap_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgSwap
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgSwap{
				Sender: "invalid_address",
			},
			err: ErrInvalidAddress,
		},
		{
			name: "invalid poolID",
			msg: MsgSwap{
				Sender: sample.AccAddress(),
				PoolId: "",
			},
			err: ErrInvalidPoolID,
		},
		{
			name: "invalid tokenIn",
			msg: MsgSwap{
				Sender:  sample.AccAddress(),
				PoolId:  "test1",
				TokenIn: sdk.NewCoin("test1", sdkmath.NewInt(0)),
			},
			err: ErrInvalidTokenAmount,
		},
		{
			name: "invalid slippage",
			msg: MsgSwap{
				Sender:   sample.AccAddress(),
				PoolId:   "test1",
				TokenIn:  sdk.NewCoin("test1", sdkmath.NewInt(100)),
				TokenOut: sdk.NewCoin("test2", sdkmath.NewInt(100)),
				Slippage: sdkmath.NewInt(10000),
			},
			err: ErrInvalidSlippage,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
