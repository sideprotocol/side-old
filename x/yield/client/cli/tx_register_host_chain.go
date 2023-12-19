package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/sideprotocol/side/x/yield/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdRegisterHostChain() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-host-chain [connection-id] [host-denom] [bech32prefix] [ibc-denom] [channel-id]",
		Short: "Broadcast message register-host-chain",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			connectionID := args[0]
			hostDenom := args[1]
			bech32prefix := args[2]
			ibcDenom := args[3]
			channelID := args[4]

			msg := types.NewMsgRegisterHostChain(
				clientCtx.GetFromAddress().String(),
				connectionID,
				bech32prefix,
				hostDenom,
				ibcDenom,
				channelID,
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
