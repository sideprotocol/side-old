package cli

import (
	"fmt"
	"strconv"
	"strings"

	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/gmm/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdCreatePool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pool [tokens] [decimals] [weights] [swap-fee] ",
		Short: "Broadcast message create-pool",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			tokens, err := sdk.ParseCoinsNormalized(args[0])
			if err != nil {
				return err
			}

			weights, err := parseWeights(args[2])
			if err != nil {
				return err
			}

			decimals, err := parseDecimals(args[1])
			if err != nil {
				return err
			}

			swapFee, err := strconv.Atoi(args[3])
			if err != nil {
				return err
			}

			// swapFee is ranged from 0 to 10000.
			if swapFee > 10000 {
				return fmt.Errorf("swap fee must be less than 10000")
			}

			liquidity := []types.PoolAsset{}
			for i := 0; i < len(weights); i++ {
				if len(weights) != len(decimals) {
					return fmt.Errorf("weights and decimals must have the same length")
				}
				if len(tokens) != len(decimals) {
					return fmt.Errorf("liquidity and weights must have the same length")
				}
				weight := sdkmath.NewInt(int64(weights[i]))
				liquidity = append(liquidity, types.PoolAsset{
					Token:   tokens[i],
					Weight:  &weight,
					Decimal: sdkmath.NewInt(int64(decimals[i])),
				})
			}

			amp := math.NewInt(0)
			msg := types.NewMsgCreatePool(
				clientCtx.GetFromAddress().String(),
				types.PoolParams{
					SwapFee: sdkmath.LegacyNewDec(int64(swapFee)),
					Amp:     &amp,
				},
				liquidity,
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

func parseDecimals(decimalsStr string) ([]uint32, error) {
	decimalList := strings.Split(decimalsStr, ",")
	decimals := make([]uint32, 0, len(decimalList))

	for _, decimalStr := range decimalList {
		decimal, err := strconv.Atoi(decimalStr)
		if err != nil {
			return nil, fmt.Errorf("invalid decimal %s", decimalStr)
		}
		decimals = append(decimals, uint32(decimal))
	}

	if len(decimals) != 2 {
		return nil, fmt.Errorf("invalid decimals length %v", decimals)
	}

	return decimals, nil
}

func parseWeights(weightsStr string) ([]uint32, error) {
	weights := strings.Split(weightsStr, ",")
	if len(weights) != 2 {
		return nil, fmt.Errorf("invalid weights length %v", weights)
	}

	totalWeight := 0
	weightsAsInt := []uint32{}
	for _, weight := range weights {
		weightAsInt, err := strconv.Atoi(weight)
		if err != nil || weightAsInt <= 0 {
			return nil, fmt.Errorf("can't parse weight value %v", err)
		}
		totalWeight += weightAsInt
		weightsAsInt = append(weightsAsInt, uint32(weightAsInt))
	}

	if totalWeight != 100 {
		return nil, fmt.Errorf("weight sum has to be equal to 100 %v", totalWeight)
	}
	return weightsAsInt, nil
}
