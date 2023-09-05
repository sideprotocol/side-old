package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/gmm/types"
)

func (k Keeper) EmitEvent(ctx sdk.Context,
	action, poolID, sender string, attr ...sdk.Attribute,
) {
	headerAttr := []sdk.Attribute{
		{
			Key:   types.AttributeKeyAction,
			Value: action,
		},
		{
			Key:   types.AttributeKeyPoolID,
			Value: poolID,
		},
		{
			Key:   types.AttributeKeyName,
			Value: types.ModuleName,
		},
		{
			Key:   types.AttributeKeyMsgSender,
			Value: sender,
		},
	}

	headerAttr = append(headerAttr, attr...)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.ModuleName,
			headerAttr...,
		),
	)
}
