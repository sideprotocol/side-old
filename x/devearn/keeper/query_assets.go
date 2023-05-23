package keeper

import (
	"context"

	"sidechain/x/devearn/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AssetsAll(goCtx context.Context, req *types.QueryAllAssetsRequest) (*types.QueryAllAssetsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var assetss []types.Assets
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	assetsStore := prefix.NewStore(store, types.KeyPrefix(types.AssetsKey))

	pageRes, err := query.Paginate(assetsStore, req.Pagination, func(key []byte, value []byte) error {
		var assets types.Assets
		if err := k.cdc.Unmarshal(value, &assets); err != nil {
			return err
		}

		assetss = append(assetss, assets)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllAssetsResponse{Assets: assetss, Pagination: pageRes}, nil
}

func (k Keeper) Assets(goCtx context.Context, req *types.QueryGetAssetsRequest) (*types.QueryGetAssetsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	found := k.IsAssetRegistered(ctx, req.Denom)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryGetAssetsResponse{Assets: types.Assets{Denom: req.Denom}}, nil
}
