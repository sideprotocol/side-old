package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/btcbridge/types"
)

func (k Keeper) EmitEvent(ctx sdk.Context, sender string, attr ...sdk.Attribute) {
	headerAttr := []sdk.Attribute{
		{
			Key:   "sender",
			Value: sender,
		},
	}

	headerAttr = append(headerAttr, attr...)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.ModuleName,
			// attr...,
			headerAttr...,
		),
	)
}
