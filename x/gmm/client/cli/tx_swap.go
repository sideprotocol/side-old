package cli

import (
	"fmt"
	"strconv"

	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/gmm/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdSwap() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "swap [pool_id] [tokenIn] [tokenOut] [slippage]",
		Short: "Broadcast message swap",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			tokenIn, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			tokenOut, err := sdk.ParseCoinNormalized(args[2])
			if err != nil {
				return err
			}

			slippage := sdkmath.NewUintFromString(args[3])
			// Define the range limits
			minSlippage, maxSlippage := sdkmath.NewUint(0), sdkmath.NewUint(100)

			// Check if slippage is within the range [0, 10000]
			if slippage.LT(minSlippage) || slippage.GT(maxSlippage) {
				return fmt.Errorf("slippage %d is out of range [0-10000]", slippage)
			}

			msg := types.NewMsgSwap(
				clientCtx.GetFromAddress().String(),
				args[0],
				tokenIn,
				tokenOut,
				sdkmath.Int(slippage),
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
