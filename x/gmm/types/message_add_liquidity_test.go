package types

import (
	"testing"

	"github.com/sideprotocol/side/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgAddLiquidity_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgAddLiquidity
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgAddLiquidity{
				Creator: "invalid_address",
			},
			err: ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgAddLiquidity{
				Creator: sample.AccAddress(),
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
