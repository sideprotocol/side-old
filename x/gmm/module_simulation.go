package gmm

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/sideprotocol/side/testutil/sample"
	gmmsimulation "github.com/sideprotocol/side/x/gmm/simulation"
	"github.com/sideprotocol/side/x/gmm/types"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = gmmsimulation.FindAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
	_ = rand.Rand{}
)

const (
	opWeightMsgAddLiquidity = "op_weight_msg_add_liquidity" // #nosec
	// TODO: Determine the simulation weight value
	defaultWeightMsgAddLiquidity int = 100

	opWeightMsgCreatePool = "op_weight_msg_create_pool" // #nosec
	// TODO: Determine the simulation weight value
	defaultWeightMsgCreatePool int = 100

	opWeightMsgWithdraw = "op_weight_msg_withdraw" // #nosec
	// TODO: Determine the simulation weight value
	defaultWeightMsgWithdraw int = 100

	opWeightMsgSwap = "op_weight_msg_swap" // #nosec
	// TODO: Determine the simulation weight value
	defaultWeightMsgSwap int = 100

	opWeightMsgDeletePool = "op_weight_msg_delete_pool" // #nosec
	// TODO: Determine the simulation weight value
	defaultWeightMsgDeletePool int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	gmmGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&gmmGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// ProposalContents doesn't return any content functions for governance proposals.
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalMsg {
	return nil
}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	_ = simState
	var weightMsgAddLiquidity int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgAddLiquidity, &weightMsgAddLiquidity, nil,
		func(_ *rand.Rand) {
			weightMsgAddLiquidity = defaultWeightMsgAddLiquidity
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgAddLiquidity,
		gmmsimulation.SimulateMsgAddLiquidity(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCreatePool int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgCreatePool, &weightMsgCreatePool, nil,
		func(_ *rand.Rand) {
			weightMsgCreatePool = defaultWeightMsgCreatePool
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreatePool,
		gmmsimulation.SimulateMsgCreatePool(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgWithdraw int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgWithdraw, &weightMsgWithdraw, nil,
		func(_ *rand.Rand) {
			weightMsgWithdraw = defaultWeightMsgWithdraw
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgWithdraw,
		gmmsimulation.SimulateMsgWithdraw(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgSwap int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgSwap, &weightMsgSwap, nil,
		func(_ *rand.Rand) {
			weightMsgSwap = defaultWeightMsgSwap
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSwap,
		gmmsimulation.SimulateMsgSwap(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgDeletePool int
	simState.AppParams.GetOrGenerate(simState.Cdc, opWeightMsgDeletePool, &weightMsgDeletePool, nil,
		func(_ *rand.Rand) {
			weightMsgDeletePool = defaultWeightMsgDeletePool
		},
	)

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	_ = simState
	return []simtypes.WeightedProposalMsg{
		simulation.NewWeightedProposalMsg(
			opWeightMsgAddLiquidity,
			defaultWeightMsgAddLiquidity,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				gmmsimulation.SimulateMsgAddLiquidity(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgCreatePool,
			defaultWeightMsgCreatePool,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				gmmsimulation.SimulateMsgCreatePool(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgWithdraw,
			defaultWeightMsgWithdraw,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				gmmsimulation.SimulateMsgWithdraw(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),
		simulation.NewWeightedProposalMsg(
			opWeightMsgSwap,
			defaultWeightMsgSwap,
			func(r *rand.Rand, ctx sdk.Context, accs []simtypes.Account) sdk.Msg {
				gmmsimulation.SimulateMsgSwap(am.accountKeeper, am.bankKeeper, am.keeper)
				return nil
			},
		),

		// this line is used by starport scaffolding # simapp/module/OpMsg
	}
}
