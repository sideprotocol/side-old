package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	"sidechain/x/devearn/types"
)

func (k Keeper) RegisterDevEarn(
	ctx sdk.Context,
	contract common.Address,
	epochs uint32,
	ownerAddr string,
) (*types.DevEarnInfo, error) {
	// Check if the Incentives are globally enabled
	params := k.GetParams(ctx)
	if !params.EnableDevEarn {
		return nil, errorsmod.Wrap(
			types.ErrInternalDevEarn,
			"dev earn are currently disabled by governance",
		)
	}

	// Check if contract exists
	acc := k.evmKeeper.GetAccountWithoutBalance(ctx, contract)
	if acc == nil || !acc.IsContract() {
		return nil, errorsmod.Wrapf(
			types.ErrContractNotFound,
			"contract doesn't exist: %s", contract,
		)
	}

	// Check if the incentive is already registered
	if k.IsDevEarnInfoRegistered(ctx, contract) {
		return nil, errorsmod.Wrapf(
			types.ErrInternalDevEarn,
			"incentive already registered: %s", contract,
		)
	}

	// create incentive and set to store
	devEarnInfo := types.NewDevEarn(contract, 0, epochs, ownerAddr)
	devEarnInfo.StartTime = ctx.BlockTime()
	k.SetDevEarnInfo(ctx, devEarnInfo)

	return &devEarnInfo, nil
}

// RegisterIncentive deletes the incentive for a contract
func (k Keeper) CancelDevEarn(
	ctx sdk.Context,
	contract common.Address,
) error {
	// Check if the Incentives are globally enabled
	params := k.GetParams(ctx)
	if !params.EnableDevEarn {
		return errorsmod.Wrap(
			types.ErrInternalDevEarn,
			"incentives are currently disabled by governance",
		)
	}

	devEarnInfo, found := k.GetDevEarnInfo(ctx, contract)
	if !found {
		return errorsmod.Wrapf(
			errortypes.ErrInvalidAddress,
			"unmatching contract '%s' ", contract,
		)
	}

	k.DeleteDevEarnInfo(ctx, devEarnInfo)
	return nil
}
