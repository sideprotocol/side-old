package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
)

func GetDefaultTimeOut(ctx *sdk.Context) (clienttypes.Height, uint64) {
	// 100 block later than current block
	outBlockHeight := ctx.BlockHeight() + 200
	// 10 min later current block time.
	waitDuration, _ := time.ParseDuration("10m")
	timeoutStamp := ctx.BlockTime().Add(waitDuration)
	timeoutHeight := clienttypes.NewHeight(0, uint64(outBlockHeight))
	return timeoutHeight, uint64(timeoutStamp.UTC().UnixNano())
}