package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"sidechain/x/devearn/types"
)

type (
	Keeper struct {
		cdc           codec.BinaryCodec
		storeKey      storetypes.StoreKey
		memKey        storetypes.StoreKey
		paramstore    paramtypes.Subspace
		authority     sdk.AccAddress
		bankKeeper    types.BankKeeper
		accountKeeper types.AccountKeeper
		evmKeeper     types.EvmKeeper
		oracleKeeper  types.OracleKeeper
		erc20Keeper   types.Erc20Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	authority sdk.AccAddress,
	bankKeeper types.BankKeeper,
	accountKeeper types.AccountKeeper,
	evmKeeper types.EvmKeeper,
	oracleKeeper types.OracleKeeper,
	erc20Keeper types.Erc20Keeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		paramstore:    ps,
		authority:     authority,
		bankKeeper:    bankKeeper,
		accountKeeper: accountKeeper,
		evmKeeper:     evmKeeper,
		oracleKeeper:  oracleKeeper,
		erc20Keeper:   erc20Keeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
