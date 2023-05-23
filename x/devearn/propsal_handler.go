package devearn

import (
	"sidechain/x/devearn/keeper"
	"sidechain/x/devearn/types"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/ethereum/go-ethereum/common"
)

func NewDevEarnProposalHandler(k *keeper.Keeper) govv1beta1.Handler {
	return func(ctx sdk.Context, content govv1beta1.Content) error {
		switch c := content.(type) {
		case *types.RegisterDevEarnInfoProposal:
			return handleRegisterDevEarnProposal(ctx, k, c)
		case *types.CancelDevEarnInfoProposal:
			return handleCancelDevEarnProposal(ctx, k, c)
		case *types.AddAssetToWhitelistProposal:
			return handleAddAssetToWhitelistProposal(ctx, k, c)
		case *types.RemoveAssetFromWhitelistProposal:
			return handleRemoveAssetFromWhitelistProposal(ctx, k, c)
		default:
			return errorsmod.Wrapf(
				errortypes.ErrUnknownRequest,
				"unrecognized %s proposal content type: %T", types.ModuleName, c,
			)
		}
	}
}

func handleRegisterDevEarnProposal(ctx sdk.Context, k *keeper.Keeper, p *types.RegisterDevEarnInfoProposal) error {
	in, err := k.RegisterDevEarnInfo(ctx, common.HexToAddress(p.Contract), p.Epochs, p.OwnerAddress)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRegisterDevEarn,
			sdk.NewAttribute(types.AttributeKeyContract, in.Contract),
			sdk.NewAttribute(
				types.AttributeKeyEpochs,
				strconv.FormatUint(uint64(in.Epochs), 10),
			),
		),
	)
	return nil
}

func handleCancelDevEarnProposal(ctx sdk.Context, k *keeper.Keeper, p *types.CancelDevEarnInfoProposal) error {
	err := k.CancelDevEarnInfo(ctx, common.HexToAddress(p.Contract))
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCancelDevEarn,
			sdk.NewAttribute(types.AttributeKeyContract, p.Contract),
		),
	)
	return nil
}

func handleAddAssetToWhitelistProposal(ctx sdk.Context, k *keeper.Keeper, p *types.AddAssetToWhitelistProposal) error {
	in, err := k.AddAssetToWhitelist(ctx, p.Denom)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeAddAssetToWhitelist,
			sdk.NewAttribute(types.AttributeKeyAsset, in.Denom),
		),
	)
	return nil
}

func handleRemoveAssetFromWhitelistProposal(ctx sdk.Context, k *keeper.Keeper, p *types.RemoveAssetFromWhitelistProposal) error {
	err := k.RemoveAssetFromWhitelist(ctx, p.Denom)
	if err != nil {
		return err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRemoveAssetFromWhitelist,
			sdk.NewAttribute(types.AttributeKeyAsset, p.Denom),
		),
	)
	return nil
}
