package types

import (
	"testing"

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
				PoolId:  "",
			},
			err: ErrInvalidPoolID,
		},
		{
			name: "invalid tokenIn",
			msg: MsgSwap{
				Sender: sample.AccAddress(),
				PoolId:  "test1",
				TokenIn: sdk.NewCoin("test1", sdk.NewInt(0)),
			},
			err: ErrInvalidTokenAmount,
		},
		{
			name: "invalid denomOut",
			msg: MsgSwap{
				Sender:  sample.AccAddress(),
				PoolId:   "test1",
				TokenIn:  sdk.NewCoin("test1", sdk.NewInt(100)),
				DenomOut: "",
			},
			err: ErrEmptyDenom,
		},
		{
			name: "invalid denomOut",
			msg: MsgSwap{
				Sender:  sample.AccAddress(),
				PoolId:   "test1",
				TokenIn:  sdk.NewCoin("test1", sdk.NewInt(100)),
				DenomOut: "test2",
			},
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
