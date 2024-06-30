package cli

import (
	"context"
	"fmt"
	"strconv"

	// "strings"

	"github.com/sideprotocol/side/x/btcbridge/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(_ string) *cobra.Command {
	// Group yield queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdBestBlock())
	cmd.AddCommand(CmdQueryBlock())
	cmd.AddCommand(CmdQueryUTXOs())
	cmd.AddCommand(CmdQuerySigningRequest())
	// this line is used by starport scaffolding # 1

	return cmd
}

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "shows the parameters of the module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.QueryParams(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdBestBlock() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "best-block",
		Short: "shows the best block header of the light client",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.QueryChainTip(cmd.Context(), &types.QueryChainTipRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdQueryBlock returns the command to query the heights of the light client
func CmdQueryBlock() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "block [hash or height]",
		Short: "Query block by hash or height",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			height, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				res, err := queryClient.QueryBlockHeaderByHash(cmd.Context(), &types.QueryBlockHeaderByHashRequest{Hash: args[0]})
				if err != nil {
					return err
				}

				return clientCtx.PrintProto(res)
			}

			res, err := queryClient.QueryBlockHeaderByHeight(cmd.Context(), &types.QueryBlockHeaderByHeightRequest{Height: height})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// CmdQuerySigningRequest returns the command to query signing request
func CmdQuerySigningRequest() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signing-request [status or address]",
		Short: "Query signing requests by status or address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			status, err := strconv.ParseInt(args[0], 10, 32)
			if err != nil {
				_, err = sdk.AccAddressFromBech32(args[0])
				if err != nil {
					return err
				}

				res, err := queryClient.QuerySigningRequestByAddress(cmd.Context(), &types.QuerySigningRequestByAddressRequest{Address: args[0]})
				if err != nil {
					return err
				}

				return clientCtx.PrintProto(res)
			}

			res, err := queryClient.QuerySigningRequest(cmd.Context(), &types.QuerySigningRequestRequest{Status: types.SigningStatus(status)})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryUTXOs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "utxos [address]",
		Short: "query utxos with an optional address",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			if len(args) == 0 {
				return queryUTXOs(cmd.Context(), &clientCtx)
			}

			return queryUTXOsByAddr(cmd.Context(), &clientCtx, args[0])
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func queryUTXOs(cmdCtx context.Context, clientCtx *client.Context) error {
	queryClient := types.NewQueryClient(clientCtx)

	res, err := queryClient.QueryUTXOs(cmdCtx, &types.QueryUTXOsRequest{})
	if err != nil {
		return err
	}

	return clientCtx.PrintProto(res)
}

func queryUTXOsByAddr(cmdCtx context.Context, clientCtx *client.Context, addr string) error {
	queryClient := types.NewQueryClient(clientCtx)

	_, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return err
	}

	res, err := queryClient.QueryUTXOsByAddress(cmdCtx, &types.QueryUTXOsByAddressRequest{
		Address: addr,
	})
	if err != nil {
		return err
	}

	return clientCtx.PrintProto(res)
}
