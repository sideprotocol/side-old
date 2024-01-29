package wasmbinding

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	"github.com/sideprotocol/side/wasmbinding/bindings"

	gmmkeeper "github.com/sideprotocol/side/x/gmm/keeper"
	gmmtypes "github.com/sideprotocol/side/x/gmm/types"
)

// CustomMessageDecorator returns decorator for custom CosmWasm bindings messages
func CustomMessageDecorator(bank *bankkeeper.BaseKeeper, gmm *gmmkeeper.Keeper) func(wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			wrapped: old,
			bank:    bank,
			gmm:     gmm,
		}
	}
}

type CustomMessenger struct {
	wrapped wasmkeeper.Messenger
	bank    *bankkeeper.BaseKeeper
	gmm     *gmmkeeper.Keeper
}

var _ wasmkeeper.Messenger = (*CustomMessenger)(nil)

// DispatchMsg executes on the contractMsg.
func (m *CustomMessenger) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) ([]sdk.Event, [][]byte, error) {
	if msg.Custom != nil {
		// only handle the happy path where this is really creating / minting / swapping ...
		// leave everything else for the wrapped version
		var contractMsg bindings.OsmosisMsg
		if err := json.Unmarshal(msg.Custom, &contractMsg); err != nil {
			return nil, nil, errorsmod.Wrap(err, "osmosis msg")
		}
		if contractMsg.Swap != nil {
			return m.swap(ctx, contractAddr, contractMsg.Swap)
		}
	}
	return m.wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}

// swap using gmm pool
func (m *CustomMessenger) swap(ctx sdk.Context, contractAddr sdk.AccAddress, createDenom *bindings.Swap) ([]sdk.Event, [][]byte, error) {
	err := PerformSwap(m.gmm, m.bank, ctx, contractAddr, createDenom)
	if err != nil {
		return nil, nil, errorsmod.Wrap(err, "perform swap token")
	}
	return nil, nil, nil
}

// PerformSwap is used with swap to swap a token denom; validates the msgSwap.
func PerformSwap(f *gmmkeeper.Keeper, b *bankkeeper.BaseKeeper, ctx sdk.Context, contractAddr sdk.AccAddress, swap *bindings.Swap) error {
	if swap == nil {
		return wasmvmtypes.InvalidRequest{Err: "swap token null swap token"}
	}

	msgServer := gmmkeeper.NewMsgServerImpl(*f)

	msgSwap := gmmtypes.NewMsgSwap(contractAddr.String(), swap.PoolId, swap.TokenIn, swap.TokenOut)

	if err := msgSwap.ValidateBasic(); err != nil {
		return errorsmod.Wrap(err, "failed validating MsgSwap")
	}

	// Swap token
	_, err := msgServer.Swap(
		sdk.WrapSDKContext(ctx),
		msgSwap,
	)
	if err != nil {
		return errorsmod.Wrap(err, "creating denom")
	}
	return nil
}
