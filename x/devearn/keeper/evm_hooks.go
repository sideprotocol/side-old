package keeper

import (
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"sidechain/x/devearn/types"
)

var _ evmtypes.EvmHooks = Hooks{}

// PostTxProcessing is a wrapper for calling the EVM PostTxProcessing hook on
// the module keeper
func (h Hooks) PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error {
	return h.k.PostTxProcessing(ctx, msg, receipt)
}

// PostTxProcessing implements EvmHooks.PostTxProcessing. After each successful
// interaction with an incentivized contract, the participants's GasUsed is
// added to its gasMeter.
func (k Keeper) PostTxProcessing(ctx sdk.Context, msg core.Message, receipt *ethtypes.Receipt) error {
	// check if the Incentives are globally enabled
	params := k.GetParams(ctx)
	if !params.EnableDevEarn {
		return nil
	}

	contract := msg.To()
	// If theres no dev earn registered for the contract, do nothing
	if contract == nil || !k.IsDevEarnInfoRegistered(ctx, *contract) {
		return nil
	}

	k.addGasToDevEarn(ctx, *contract, receipt.GasUsed)

	defer func() {
		telemetry.IncrCounter(
			1,
			"tx", "msg", "ethereum_tx", types.ModuleName, "total",
		)

		if receipt.GasUsed != 0 {
			telemetry.IncrCounter(
				float32(receipt.GasUsed),
				"tx", "msg", "ethereum_tx", types.ModuleName, "gas_used", "total",
			)
		}
	}()

	return nil
}

// addGasToIncentive adds gasUsed to an incentive's cumulated totalGas
func (k Keeper) addGasToDevEarn(
	ctx sdk.Context,
	contract common.Address,
	gasUsed uint64,
) {
	// NOTE: existence of contract incentive is already checked
	incentive, _ := k.GetDevEarnInfo(ctx, contract)
	incentive.GasMeter += gasUsed
	k.SetDevEarnInfo(ctx, incentive)
}
