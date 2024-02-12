package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgWithdraw_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgWithdraw
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgWithdraw{
				Sender: "invalid_address",
			},
			err: ErrInvalidAddress,
		},
		{
			name: "invalid receiver address",
			msg: MsgWithdraw{
				Sender:   sample.AccAddress(),
				Receiver: "invalid_address",
			},
			err: ErrInvalidAddress,
		},
		{
			name: "invalid poolID",
			msg: MsgWithdraw{
				Sender:   sample.AccAddress(),
				Receiver: sample.AccAddress(),
				PoolId:   "",
			},
			err: ErrInvalidPoolID,
		},
		{
			name: "invalid share amount",
			msg: MsgWithdraw{
				Sender:   sample.AccAddress(),
				Receiver: sample.AccAddress(),
				PoolId:   "test1",
				Share: sdk.NewCoin(
					"test",
					sdk.NewInt(0),
				),
			},
			err: ErrInvalidTokenAmount,
		},
		{
			name: "mismatched share denom",
			msg: MsgWithdraw{
				Sender:   sample.AccAddress(),
				Receiver: sample.AccAddress(),
				PoolId:   "test1",
				Share: sdk.NewCoin(
					"test",
					sdk.NewInt(1),
				),
			},
			err: ErrMismatchedShareDenom,
		},
		{
			name: "valid message",
			msg: MsgWithdraw{
				Sender:   sample.AccAddress(),
				Receiver: sample.AccAddress(),
				PoolId:   "test1",
				Share: sdk.NewCoin(
					"side/gmm/test1",
					sdk.NewInt(10),
				),
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
