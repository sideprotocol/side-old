package keeper

import (
	"sidechain/x/devearn/types"
	erc20types "sidechain/x/erc20/types"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
)

func (k Keeper) RegisterDevEarnInfo(
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
func (k Keeper) CancelDevEarnInfo(
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

// Add asset to whitelist storage
func (k Keeper) AddAssetToWhitelist(
	ctx sdk.Context,
	denom string,
) (*types.Assets, error) {
	// Check if the Incentives are globally enabled
	params := k.GetParams(ctx)
	if !params.EnableDevEarn {
		return nil, errorsmod.Wrap(
			types.ErrInternalDevEarn,
			"dev earn are currently disabled by governance",
		)
	}

	// Check if the asset is already registered
	if k.IsAssetRegistered(ctx, denom) {
		return nil, errorsmod.Wrapf(
			types.ErrInternalDevEarn,
			"asset already registered: %s", denom,
		)
	}

	// Check if asset is registered in erc20 module
	_, tokenPairErr := k.erc20Keeper.TokenPair(
		ctx, &erc20types.QueryTokenPairRequest{Token: denom})
	if tokenPairErr != nil {
		return nil, tokenPairErr
	}

	assets := types.Assets{Denom: denom}
	k.SetAssets(ctx, assets)

	return &assets, nil
}

// Remove asset from whitelist
func (k Keeper) RemoveAssetFromWhitelist(
	ctx sdk.Context,
	denom string,
) error {
	// Check if the Incentives are globally enabled
	params := k.GetParams(ctx)
	if !params.EnableDevEarn {
		return errorsmod.Wrap(
			types.ErrInternalDevEarn,
			"incentives are currently disabled by governance",
		)
	}

	if !k.IsAssetRegistered(ctx, denom) {
		return errorsmod.Wrapf(
			types.ErrInternalDevEarn,
			"asset is not in whitelist: %s", denom,
		)
	}

	k.RemoveAssets(ctx, denom)
	return nil
}
