package types

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgCreatePool_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCreatePool
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgCreatePool{
				Creator: "invalid_address",
			},
			err: ErrInvalidAddress,
		},
		{
			name: "valid address, invalid PoolParams",
			msg: MsgCreatePool{
				Creator: sample.AccAddress(),
				Params:  nil,
			},
			err: ErrInvalidPoolParams, // Replace with the actual error
		},
		{
			name: "valid address, valid PoolParams, empty Liquidity",
			msg: MsgCreatePool{
				Creator: sample.AccAddress(),
				Params: &PoolParams{
					Type:    PoolType_WEIGHT,
					SwapFee: sdkmath.LegacyNewDec(int64(0)),
					Amp:     sdk.NewInt(0),
				}, // Replace with a sample valid PoolParams
			},
			err: ErrEmptyLiquidity, // Replace with the actual error
		},
		{
			name: "valid address, valid PoolParams, valid Liquidity",
			msg: MsgCreatePool{
				Creator: sample.AccAddress(),
				Params: &PoolParams{
					Type:    PoolType_WEIGHT,
					SwapFee: sdkmath.LegacyNewDec(int64(0)),
					Amp:     sdk.NewInt(0),
				}, // Replace with a sample valid PoolParams
				Liquidity: []PoolAsset{
					{
						Token:  sdk.NewCoin("test1", sdk.NewInt(100000000000000)),
						Weight: sdkmath.NewInt(int64(1)),
					},
				}, // Replace with a sample valid Liquidity
			},
			err: ErrInvalidLiquidityInLength,
		},
		{
			name: "valid address, valid PoolParams, valid Liquidity",
			msg: MsgCreatePool{
				Creator: sample.AccAddress(),
				Params: &PoolParams{
					Type:    PoolType_WEIGHT,
					SwapFee: sdkmath.LegacyNewDec(int64(0)),
					Amp:     sdk.NewInt(0),
				}, // Replace with a sample valid PoolParams
				Liquidity: []PoolAsset{
					{
						Token:  sdk.NewCoin("test1", sdk.NewInt(100000000000000)),
						Weight: sdkmath.NewInt(int64(1)),
					},
					{
						Token:  sdk.NewCoin("test2", sdk.NewInt(100000000000000)),
						Weight: sdkmath.NewInt(int64(1)),
					},
				}, // Replace with a sample valid Liquidity
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
