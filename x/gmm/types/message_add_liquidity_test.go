package types

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
				Sender: "invalid_address",
			},
			err: ErrInvalidAddress,
		},
		{
			name: "invalid poolID",
			msg: MsgAddLiquidity{
				Sender: sample.AccAddress(),
				PoolId: "",
			},
			err: ErrInvalidPoolID,
		},
		{
			name: "invalid liquidity (length = 0)",
			msg: MsgAddLiquidity{
				Sender:    sample.AccAddress(),
				PoolId:    "test",
				Liquidity: []sdk.Coin{},
			},
			err: ErrInvalidLiquidityInLength,
		},
		{
			name: "invalid liquidity (length > 2)",
			msg: MsgAddLiquidity{
				Sender: sample.AccAddress(),
				PoolId: "test",
				Liquidity: []sdk.Coin{
					{Denom: "test1", Amount: sdkmath.NewInt(1000000000)},
					{Denom: "test2", Amount: sdkmath.NewInt(1000000000)},
					{Denom: "test3", Amount: sdkmath.NewInt(1000000000)},
				},
			},
			err: ErrInvalidLiquidityInLength,
		},
		{
			name: "invalid liquidity (amounts)",
			msg: MsgAddLiquidity{
				Sender: sample.AccAddress(),
				PoolId: "test",
				Liquidity: []sdk.Coin{
					{Denom: "test", Amount: sdkmath.NewInt(0)},
					{Denom: "test2", Amount: sdkmath.NewInt(1000000000)},
				},
			},
			err: ErrInvalidLiquidityAmount,
		},
		{
			name: "valid addLiquidity message",
			msg: MsgAddLiquidity{
				Sender: sample.AccAddress(),
				PoolId: "test",
				Liquidity: []sdk.Coin{
					{Denom: "test", Amount: sdkmath.NewInt(1000000000)},
					{Denom: "test2", Amount: sdkmath.NewInt(1000000000)},
				},
			},
			err: nil,
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
