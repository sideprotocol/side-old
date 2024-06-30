package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/btcbridge/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) QueryParams(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

func (k Keeper) QueryChainTip(goCtx context.Context, req *types.QueryChainTipRequest) (*types.QueryChainTipResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	best := k.GetBestBlockHeader(ctx)

	return &types.QueryChainTipResponse{
		Hash:   best.Hash,
		Height: best.Height,
	}, nil
}

// BlockHeaderByHash queries the block header by hash.
func (k Keeper) QueryBlockHeaderByHash(goCtx context.Context, req *types.QueryBlockHeaderByHashRequest) (*types.QueryBlockHeaderByHashResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	header := k.GetBlockHeader(ctx, req.Hash)
	if header == nil {
		return nil, status.Error(codes.NotFound, "block header not found")
	}

	return &types.QueryBlockHeaderByHashResponse{BlockHeader: header}, nil
}

func (k Keeper) QueryBlockHeaderByHeight(goCtx context.Context, req *types.QueryBlockHeaderByHeightRequest) (*types.QueryBlockHeaderByHeightResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	header := k.GetBlockHeaderByHeight(ctx, req.Height)
	if header == nil {
		return nil, status.Error(codes.NotFound, "block header not found")
	}

	return &types.QueryBlockHeaderByHeightResponse{BlockHeader: header}, nil
}

func (k Keeper) QuerySigningRequest(goCtx context.Context, req *types.QuerySigningRequestRequest) (*types.QuerySigningRequestResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if req.Status == types.SigningStatus_SIGNING_STATUS_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "invalid status")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	requests := k.FilterSigningRequestsByStatus(ctx, req)

	return &types.QuerySigningRequestResponse{Requests: requests}, nil
}

func (k Keeper) QuerySigningRequestByAddress(goCtx context.Context, req *types.QuerySigningRequestByAddressRequest) (*types.QuerySigningRequestByAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	requests := k.FilterSigningRequestsByAddr(ctx, req)

	return &types.QuerySigningRequestByAddressResponse{Requests: requests}, nil
}

func (k Keeper) QueryUTXOs(goCtx context.Context, req *types.QueryUTXOsRequest) (*types.QueryUTXOsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	utxos := k.GetAllUTXOs(ctx)

	return &types.QueryUTXOsResponse{Utxos: utxos}, nil
}

func (k Keeper) QueryUTXOsByAddress(goCtx context.Context, req *types.QueryUTXOsByAddressRequest) (*types.QueryUTXOsByAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	_, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	utxos := k.GetUTXOsByAddr(ctx, req.Address)

	return &types.QueryUTXOsByAddressResponse{Utxos: utxos}, nil
}
