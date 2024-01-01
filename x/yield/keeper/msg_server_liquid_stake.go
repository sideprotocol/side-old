package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/side/x/yield/types"
)

func (k msgServer) LiquidStake(goCtx context.Context, msg *types.MsgLiquidStake) (*types.MsgLiquidStakeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the host chain from the base denom in the message (e.g. uatom)
	hostChain, err := k.GetHostChainFromHostDenom(ctx, msg.Denom)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidToken, "no host zone found for denom (%s)", msg.Denom)
	}

	// Get user and module account addresses
	liquidStakerAddress, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "user's address is invalid")
	}
	// hostZoneDepositAddress, err := sdk.AccAddressFromBech32(hostChain.DepositAddress)
	// if err != nil {
	// 	return nil, errorsmod.Wrapf(err, "host zone address is invalid")
	// }

	// The tokens that are sent to the protocol are denominated in the ibc hash of the native token on side (e.g. ibc/xxx)
	nativeDenom := hostChain.IbcDenom
	nativeCoin := sdk.NewCoin(nativeDenom, msg.Amount)
	if !types.IsIBCToken(nativeDenom) {
		return nil, errorsmod.Wrapf(types.ErrInvalidToken, "denom is not an IBC token (%s)", nativeDenom)
	}

	// Confirm the user has a sufficient balance to execute the liquid stake
	balance := k.bankKeeper.GetBalance(ctx, liquidStakerAddress, nativeDenom)
	if balance.IsLT(nativeCoin) {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInsufficientFunds, "balance is lower than staking amount. staking amount: %v, balance: %v", msg.Amount, balance.Amount)
	}

	// TODO: Handling the message
	// - Check balance
	// - Transfer amount given in message
	// - Deduce host denom from ibc from (msg.denom)
	// - Check for host zone using hose denom
	// - Add deposit record with status, amount and id
	// - Transfer native tokens(1st step of process)
	_ = ctx

	return &types.MsgLiquidStakeResponse{}, nil
}
