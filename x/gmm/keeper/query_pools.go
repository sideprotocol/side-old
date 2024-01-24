package keeper

import (
	"context"
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sideprotocol/side/x/gmm/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Pool(goCtx context.Context, req *types.QueryPoolRequest) (*types.QueryPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	pool, found := k.GetPool(ctx, req.PoolId)
	if !found {
		return &types.QueryPoolResponse{}, types.ErrPoolNotFound
	}

	poolI := convertPoolForWasm(pool)
	return &types.QueryPoolResponse{Pool: &poolI}, nil
}

func (k Keeper) Pools(goCtx context.Context, req *types.QueryAllPoolsRequest) (*types.QueryPoolsResponse, error) {
	// if req == nil {
	// 	return nil, status.Error(codes.InvalidArgument, "invalid request")
	// }

	pools := []types.PoolI{}
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Use the store with the mapping of poolId to its count
	store := ctx.KVStore(k.storeKey)
	poolIDToCountStore := prefix.NewStore(store, types.KeyPoolIDToCountPrefix)

	// A counter for pagination
	var counter int

	pageRes, err := query.Paginate(poolIDToCountStore, req.Pagination, func(key []byte, value []byte) error {
		counter++

		// Decode the count for the given poolId
		count := binary.BigEndian.Uint64(value)

		// Get the actual pool using the count
		poolStore := prefix.NewStore(store, types.KeyPrefix(string(types.KeyPoolsPrefix)))
		poolBytes := poolStore.Get(GetPoolKey(count))
		if poolBytes == nil {
			return nil
		}

		var pool types.Pool
		if err := k.cdc.Unmarshal(poolBytes, &pool); err != nil {
			return nil
		}
		poolI := convertPoolForWasm(pool)
		pools = append(pools, poolI)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPoolsResponse{Pools: pools, Pagination: pageRes}, nil
}

func (k Keeper) MyPools(goCtx context.Context, req *types.QueryPoolsRequest) (*types.QueryPoolsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var pools []types.PoolI
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Use the store with the mapping of poolId to its count
	store := ctx.KVStore(k.storeKey)
	poolIDToCountStore := prefix.NewStore(store, types.KeyPoolIDToCountPrefix)

	// A counter for pagination
	var counter int

	pageRes, err := query.Paginate(poolIDToCountStore, req.Pagination, func(key []byte, value []byte) error {
		counter++

		// Decode the count for the given poolId
		count := binary.BigEndian.Uint64(value)

		// Get the actual pool using the count
		poolStore := prefix.NewStore(store, types.KeyPrefix(string(types.KeyPoolsPrefix)))
		poolBytes := poolStore.Get(GetPoolKey(count))
		if poolBytes == nil {
			return nil
		}

		var pool types.Pool
		if err := k.cdc.Unmarshal(poolBytes, &pool); err != nil {
			return err
		}
		if pool.Sender == req.Creator {
			poolI := convertPoolForWasm(pool)
			pools = append(pools, poolI)
		}
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPoolsResponse{Pools: pools, Pagination: pageRes}, nil
}

func convertPoolForWasm(pool types.Pool) types.PoolI {
	assets := []*types.PoolWasmAsset{}
	for _, asset := range pool.Assets {
		assetCopy := asset
		assets = append(assets, &types.PoolWasmAsset{
			Balance: &assetCopy.Token,
			Decimal: uint32(asset.Decimal.Int64()),
			Weight:  uint32(asset.Weight.Int64()),
		})
	}
	poolI := types.PoolI{
		Id:            pool.PoolId,
		SourceCreator: pool.Sender,
		Assets:        assets,
		SwapFee:       uint32(pool.PoolParams.SwapFee.RoundInt().Int64()),
		Supply:        &pool.TotalShares,
	}
	return poolI
}
