package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/sideprotocol/side/x/gmm/types"
)

func CmdQueryVolume() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "volume [pool-Id] [type](day or total)",
		Short: "shows the volume of the pool",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			if args[1] == "day" {
				res, err := queryClient.Volume24(cmd.Context(), &types.QueryVolumeRequest{
					PoolId: args[0],
				})
				if err != nil {
					return err
				}
				return clientCtx.PrintProto(res)
			} else if args[1] == "total" {
				res, err := queryClient.TotalVolume(cmd.Context(), &types.QueryTotalVolumeRequest{
					PoolId: args[0],
				})
				if err != nil {
					return err
				}
				return clientCtx.PrintProto(res)
			} else {
				return fmt.Errorf("invalid volume type")
			}
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
