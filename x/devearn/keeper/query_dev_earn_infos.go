package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/sideprotocol/sidechain/x/devearn/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) DevEarnInfos(goCtx context.Context, req *types.QueryDevEarnInfosRequest) (*types.QueryDevEarnInfosResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixDevEarn)
	var devEarnInfos []types.DevEarnInfo
	pageRes, err := query.Paginate(
		store,
		req.Pagination,
		func(key, value []byte) error {
			var devEarnInfo types.DevEarnInfo
			if err := devEarnInfo.Unmarshal(value); err != nil {
				return err
			}

			devEarnInfos = append(devEarnInfos, devEarnInfo)

			return nil
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryDevEarnInfosResponse{DevEarnInfos: devEarnInfos, Pagination: pageRes}, nil
}
