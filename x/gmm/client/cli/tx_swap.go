package cli

import (
	"strconv"

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
		Use:   "swap [pool_id] [tokenIn] [tokenOut]",
		Short: "Broadcast message swap",
		Args:  cobra.ExactArgs(0),
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
			msg := types.NewMsgSwap(
				clientCtx.GetFromAddress().String(),
				args[0],
				tokenIn,
				tokenOut,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
