package keeper

import (
	"context"
	"strings"

	errorsmod "cosmossdk.io/errors"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"

	"github.com/sideprotocol/sidechain/x/devearn/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) DevEarnInfo(goCtx context.Context, req *types.QueryDevEarnInfoRequest) (*types.QueryDevEarnInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	if strings.TrimSpace(req.Contract) == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"contract address is empty",
		)
	}
	if !common.IsHexAddress(req.Contract) {
		return nil, errorsmod.Wrapf(
			errortypes.ErrInvalidAddress, "address '%s' is not a valid ethereum hex address",
			req.Contract,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	contract := req.Contract
	contractAddr := common.HexToAddress(contract)
	devEarninfo, ok := k.GetDevEarnInfo(ctx, contractAddr)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "didn't found this devEarnInfo!!")
	}
	resp := &types.QueryDevEarnInfoResponse{DevEarnInfo: devEarninfo}
	return resp, nil
}
