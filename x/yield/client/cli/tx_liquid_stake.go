package cli

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sideprotocol/side/x/yield/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdLiquidStake() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquid-stake [denom] [amount]",
		Short: "Broadcast message liquid-stake",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argDenom := args[0]
			argAmount, found := sdk.NewIntFromString(args[1])
			if !found {
				return errorsmod.Wrap(sdkerrors.ErrInvalidType, "can not convert string to int")
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgLiquidStake(
				clientCtx.GetFromAddress().String(),
				argDenom,
				argAmount,
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
