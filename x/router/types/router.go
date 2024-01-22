package types

type SwapAmountInRoute struct {
	PoolId   string `protobuf:"varint,1,opt,name=pool_id,json=poolId,proto3" json:"pool_id,omitempty" yaml:"pool_id"`
	OutDenom string `protobuf:"bytes,2,opt,name=token_out_denom,json=tokenOutDenom,proto3" json:"token_out_denom,omitempty" yaml:"token_out_denom"`
}

type SwapAmountInRoutes []SwapAmountInRoute
