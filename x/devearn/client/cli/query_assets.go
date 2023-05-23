package cli

import (
	"context"

	"sidechain/x/devearn/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdListAssets() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-assets",
		Short: "list all assets",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllAssetsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.AssetsAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowAssets() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-assets [denom]",
		Short: "shows a assets",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			// id, err := strconv.ParseUint(args[0], 10, 64)
			// if err != nil {
			// 	return err
			// }

			params := &types.QueryGetAssetsRequest{
				Denom: args[0],
			}

			res, err := queryClient.Assets(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
