package keeper

import (
	"encoding/json"
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"

	sdkmath "cosmossdk.io/math"
	//wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	gmmtypes "github.com/sideprotocol/side/x/gmm/types"
	"github.com/sideprotocol/side/x/router/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace

		channelKeeper types.ChannelKeeper
		portKeeper    types.PortKeeper
		scopedKeeper  exported.ScopedKeeper
		ics4Wrapper   porttypes.ICS4Wrapper
		//ibcContractKeeper wasmtypes.IBCContractKeeper

		wasmKeeper types.WasmKeeper
		gmmKeeper  types.GmmKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	channelKeeper types.ChannelKeeper,
	portKeeper types.PortKeeper,
	scopedKeeper exported.ScopedKeeper,
	wasmKeeper types.WasmKeeper,
	gmmKeeper types.GmmKeeper,

) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,

		channelKeeper: channelKeeper,
		portKeeper:    portKeeper,
		scopedKeeper:  scopedKeeper,
		wasmKeeper:    wasmKeeper,
		gmmKeeper:     gmmKeeper,
	}
}

// ----------------------------------------------------------------------------
// IBC Keeper Logic
// ----------------------------------------------------------------------------

// ChanCloseInit defines a wrapper function for the channel Keeper's function.
func (k Keeper) ChanCloseInit(ctx sdk.Context, portID, channelID string) error {
	capName := host.ChannelCapabilityPath(portID, channelID)
	chanCap, ok := k.scopedKeeper.GetCapability(ctx, capName)
	if !ok {
		return sdkerrors.Wrapf(channeltypes.ErrChannelCapabilityNotFound, "could not retrieve channel capability at: %s", capName)
	}
	return k.channelKeeper.ChanCloseInit(ctx, portID, channelID, chanCap)
}

// IsBound checks if the IBC app module is already bound to the desired port
func (k Keeper) IsBound(ctx sdk.Context, portID string) bool {
	_, ok := k.scopedKeeper.GetCapability(ctx, host.PortPath(portID))
	return ok
}

// BindPort defines a wrapper function for the port Keeper's function in
// order to expose it to module's InitGenesis function
func (k Keeper) BindPort(ctx sdk.Context, portID string) error {
	cap := k.portKeeper.BindPort(ctx, portID)
	return k.ClaimCapability(ctx, cap, host.PortPath(portID))
}

// GetPort returns the portID for the IBC app module. Used in ExportGenesis
func (k Keeper) GetPort(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get(types.PortKey))
}

// SetPort sets the portID for the IBC app module. Used in InitGenesis
func (k Keeper) SetPort(ctx sdk.Context, portID string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PortKey, []byte(portID))
}

// AuthenticateCapability wraps the scopedKeeper's AuthenticateCapability function
func (k Keeper) AuthenticateCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) bool {
	return k.scopedKeeper.AuthenticateCapability(ctx, cap, name)
}

// ClaimCapability allows the IBC app module to claim a capability that core IBC
// passes to it
func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) SwapWithRouters(ctx sdk.Context, tokenIn sdk.Coin, routes types.SwapAmountInRoutes, contract sdk.AccAddress, slippage sdkmath.Int) error {
	in := tokenIn
	for _, router := range routes {
		pool, _, err := k.getPool(ctx, router.PoolId, contract)
		if err != nil {
			return err
		}
		out, err := pool.EstimateSwap(in, router.OutDenom)
		if err != nil {
			return err
		}
		in = out
	}
	return k.performSwap(ctx, in, routes, contract)
}

func (k Keeper) getPool(ctx sdk.Context, poolID string, contract sdk.AccAddress) (gmmtypes.Pool, bool, error) {
	poolRes, err := k.gmmKeeper.Pool(ctx, &gmmtypes.QueryPoolRequest{PoolId: poolID})
	if err == nil {
		return *poolRes.Pool, false, nil
	}

	poolQuery := types.PoolQuery{PoolID: poolID}
	rawPoolQueryData, err := json.Marshal(poolQuery)
	if err != nil {
		return gmmtypes.Pool{}, true, err
	}
	rawPoolRes, err := k.wasmKeeper.QuerySmart(ctx, contract, rawPoolQueryData)
	if err != nil {
		return gmmtypes.Pool{}, true, err
	}
	var iPoolRes types.InterchainPoolResponse
	if err = json.Unmarshal(rawPoolRes, &iPoolRes); err != nil {
		return gmmtypes.Pool{}, true, err
	}
	return iPoolRes.ToGmmPool(), true, nil
}

func (k Keeper) performSwap(ctx sdk.Context, in sdk.Coin, routes types.SwapAmountInRoutes, contract sdk.AccAddress) error {
	for _, router := range routes {
		pool, isWasmPool, err := k.getPool(ctx, router.PoolId, contract)
		if err != nil {
			return err
		}

		out, err := pool.EstimateSwap(in, router.OutDenom)
		if err != nil {
			return err
		}

		msgSwap := &gmmtypes.MsgSwap{
			Sender:   types.ModuleName,
			PoolId:   pool.PoolId,
			TokenIn:  in,
			TokenOut: out,
			Slippage: sdkmath.NewInt(0),
		}

		if _, err = k.executeSwap(ctx, pool, contract, msgSwap, isWasmPool); err != nil {
			return err
		}
		in = out
	}
	return nil
}

func (k Keeper) executeSwap(ctx sdk.Context, pool gmmtypes.Pool, contract sdk.AccAddress, msgSwap *gmmtypes.MsgSwap, isWasmPool bool) (*gmmtypes.MsgSwapResponse, error) {
	if !isWasmPool {
		return k.gmmKeeper.Swap(ctx, msgSwap)
	}
	_, err := k.wasmKeeper.Execute(ctx, contract, sdk.AccAddress(types.ModuleName), []byte{}, sdk.NewCoins(msgSwap.TokenIn))
	return &gmmtypes.MsgSwapResponse{}, err
}

// func (k Keeper) SendPacket(
// 	ctx sdk.Context,
// 	sourcePort,
// 	sourceChannel string,
// 	timeoutHeight clienttypes.Height,
// 	timeoutTimestamp uint64,
// 	data []byte,
// ) (*uint64, error) {
// 	_, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
// 	if !found {
// 		return nil, errorsmod.Wrapf(channeltypes.ErrChannelNotFound, "port ID (%s) channel ID (%s)", sourcePort, sourceChannel)
// 	}

// 	// get the next sequence
// 	_, found = k.channelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
// 	if !found {
// 		return nil, errorsmod.Wrapf(
// 			channeltypes.ErrSequenceSendNotFound,
// 			"source port: %s, source channel: %s", sourcePort, sourceChannel,
// 		)
// 	}

// 	// begin createOutgoingPacket logic
// 	// See spec for this logic: https://github.com/cosmos/ibc/tree/master/spec/app/ics-020-fungible-token-transfer#packet-relay
// 	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
// 	if !ok {
// 		return nil, errorsmod.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
// 	}

// 	sequence, err := k.ics4Wrapper.SendPacket(ctx, channelCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
// 	return &sequence, err
// }

func (k *Keeper) TimeoutShouldRetry(ctx sdk.Context, packet channeltypes.Packet) {

}
