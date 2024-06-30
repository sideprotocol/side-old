package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/sideprotocol/side/x/btcbridge/types"
)

var DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())

// const (
// 	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
// 	listSeparator              = ","
// )

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdSubmitBlocks())
	cmd.AddCommand(CmdUpdateSenders())
	cmd.AddCommand(CmdWithdrawBitcoin())
	cmd.AddCommand(CmdSubmitWithdrawSignatures())

	return cmd
}

func CmdSubmitBlocks() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-blocks [file-path-to-block-headers.json]",
		Short: "Submit Bitcoin block headers to the chain",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// read the block headers from the file
			blockHeaders, err := readBlockHeadersFromFile(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgSubmitBlockHeaderRequest(
				clientCtx.GetFromAddress().String(),
				blockHeaders,
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

// Update Authorized Senders
func CmdUpdateSenders() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-senders [senders]",
		Short: "Update authorized senders",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// split the senders from args[0]
			senders := strings.Split(args[0], ",")
			if len(senders) == 0 {
				return fmt.Errorf("senders can not be empty")
			}

			msg := types.NewMsgUpdateSendersRequest(
				clientCtx.GetFromAddress().String(),
				senders,
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

// Withdraw Bitcoin
func CmdWithdrawBitcoin() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw [amount] [fee-rate]",
		Short: "Withdraw bitcoin to the given sender",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			_, err = sdk.ParseCoinsNormalized(args[0])
			if err != nil {
				return fmt.Errorf("invalid amount")
			}

			msg := types.NewMsgWithdrawBitcoinRequest(
				clientCtx.GetFromAddress().String(),
				args[0],
				args[1],
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

func CmdSubmitWithdrawSignatures() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-signature [psbt]",
		Short: "Submit signed withdrawal psbt",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			p, err := psbt.NewFromRawBytes(strings.NewReader(args[0]), true)
			if err != nil {
				return fmt.Errorf("invalid psbt")
			}

			signedTx, err := psbt.Extract(p)
			if err != nil {
				return fmt.Errorf("failed to extract tx from psbt")
			}

			msg := types.NewMsgSubmitWithdrawSignaturesRequest(
				clientCtx.GetFromAddress().String(),
				signedTx.TxHash().String(),
				args[0],
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

// readBlockHeadersFromFile reads the block headers from the file
func readBlockHeadersFromFile(filePath string) ([]*types.BlockHeader, error) {
	// read the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// read the block headers from the file
	var blockHeaders []*types.BlockHeader
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&blockHeaders); err != nil {
		return nil, err
	}
	return blockHeaders, nil
}
