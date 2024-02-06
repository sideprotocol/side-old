package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/gmm/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Volume24(goCtx context.Context, req *types.QueryVolumeRequest) (*types.QueryVolumeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	volumes := k.GetVolume24(ctx, req.PoolId)
	return &types.QueryVolumeResponse{Volumes: volumes}, nil
}

func (k Keeper) TotalVolume(goCtx context.Context, req *types.QueryTotalVolumeRequest) (*types.QueryTotalVolumeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	volumes := k.GetVolume24(ctx, req.PoolId)
	return &types.QueryTotalVolumeResponse{Volumes: volumes}, nil
}
