package keepers

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	consensusparamkeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/group"
	groupkeeper "github.com/cosmos/cosmos-sdk/x/group/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"

	swasm "github.com/sideprotocol/side/wasmbinding"

	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icahost "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	"github.com/cosmos/ibc-go/v7/modules/apps/transfer"

	ibcfeekeeper "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/keeper"
	ibcfeetypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"

	icacallbacksmodulekeeper "github.com/Stride-Labs/stride/v16/x/icacallbacks/keeper"
	icacallbacksmoduletypes "github.com/Stride-Labs/stride/v16/x/icacallbacks/types"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibcclient "github.com/cosmos/ibc-go/v7/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	ibcporttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"

	interchainquerykeeper "github.com/sideprotocol/side/x/interchainquery/keeper"

	gmmmodulekeeper "github.com/sideprotocol/side/x/gmm/keeper"
	gmmmoduletypes "github.com/sideprotocol/side/x/gmm/types"

	yieldmodulekeeper "github.com/sideprotocol/side/x/yield/keeper"
	yieldmoduletypes "github.com/sideprotocol/side/x/yield/types"

	// this line is used by starport scaffolding # stargate/app/moduleImport

	appparams "github.com/sideprotocol/side/app/params"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

)

const (
	AccountAddressPrefix = "side"
)

// var (
// 	govAuthor = "bc1qjs3uxh3w8n6qhlegzjymyq9xn4jp8jam2qgy7x" //govAuthor// authtypes.NewModuleAddress(govtypes.ModuleName).String()
// )

type AppKeepers struct {
	// keepers
	AccountKeeper         authkeeper.AccountKeeper
	AuthzKeeper           authzkeeper.Keeper
	BankKeeper            bankkeeper.Keeper
	CapabilityKeeper      *capabilitykeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	SlashingKeeper        slashingkeeper.Keeper
	MintKeeper            mintkeeper.Keeper
	DistrKeeper           distrkeeper.Keeper
	GovKeeper             *govkeeper.Keeper
	CrisisKeeper          *crisiskeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	ParamsKeeper          paramskeeper.Keeper
	IBCKeeper             *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	IBCFeeKeeper          ibcfeekeeper.Keeper
	InterchainQueryKeeper interchainquerykeeper.Keeper
	EvidenceKeeper        evidencekeeper.Keeper
	TransferKeeper        ibctransferkeeper.Keeper
	ICAHostKeeper         icahostkeeper.Keeper
	IcacallbacksKeeper    icacallbacksmodulekeeper.Keeper
	IcaControllerKeeper   icacontrollerkeeper.Keeper

	FeeGrantKeeper        feegrantkeeper.Keeper
	GroupKeeper           groupkeeper.Keeper
	ConsensusParamsKeeper consensusparamkeeper.Keeper
	WasmKeeper            wasmkeeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper   capabilitykeeper.ScopedKeeper
	ScopeIBCFeeKeeper capabilitykeeper.ScopedKeeper

	ScopedTransferKeeper capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper  capabilitykeeper.ScopedKeeper
	scopedWasmKeeper     capabilitykeeper.ScopedKeeper

	GmmKeeper gmmmodulekeeper.Keeper

	YieldKeeper yieldmodulekeeper.Keeper

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey
}

// InitSpecialKeepers initiates special keepers (crisis appkeeper, upgradekeeper, params keeper)
func (appKeepers *AppKeepers) InitSpecialKeepers(
	appCodec codec.Codec,
	bApp *baseapp.BaseApp,
	cdc *codec.LegacyAmino,
	invCheckPeriod uint,
	skipUpgradeHeights map[int64]bool,
	homePath string,
) {
	appKeepers.GenerateKeys()
	paramsKeeper := appKeepers.initParamsKeeper(appCodec, cdc, appKeepers.keys[paramstypes.StoreKey], appKeepers.tkeys[paramstypes.TStoreKey])
	appKeepers.ParamsKeeper = paramsKeeper

	govAuthor := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	// set the BaseApp's parameter store
	consensusParamsKeeper := consensusparamkeeper.NewKeeper(
		appCodec, appKeepers.keys[consensusparamtypes.StoreKey], govAuthor)
	appKeepers.ConsensusParamsKeeper = consensusParamsKeeper
	bApp.SetParamStore(&appKeepers.ConsensusParamsKeeper)

	upgradeKeeper := upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		appKeepers.keys[upgradetypes.StoreKey],
		appCodec,
		homePath,
		bApp,
		govAuthor,
	)
	appKeepers.UpgradeKeeper = upgradeKeeper

	// add capability keeper and ScopeToModule for ibc module
	appKeepers.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, appKeepers.keys[capabilitytypes.StoreKey], appKeepers.memKeys[capabilitytypes.MemStoreKey])
	appKeepers.ScopedIBCKeeper = appKeepers.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	appKeepers.ScopedICAHostKeeper = appKeepers.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
	appKeepers.ScopedTransferKeeper = appKeepers.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	appKeepers.scopedWasmKeeper = appKeepers.CapabilityKeeper.ScopeToModule(wasmtypes.ModuleName)

	// TODO: Make a SetInvCheckPeriod fn on CrisisKeeper.
	// IMO, its bad design atm that it requires this in state machine initialization
	crisisKeeper := crisiskeeper.NewKeeper(
		appCodec, appKeepers.keys[crisistypes.StoreKey], invCheckPeriod, appKeepers.BankKeeper, authtypes.FeeCollectorName, govAuthor,
	)
	appKeepers.CrisisKeeper = crisisKeeper
}

// InitNormalKeepers initializes all 'normal' keepers (account, app, bank, auth, staking, distribution, slashing, transfer, gamm, IBC router, pool incentives, governance, mint, txfees keepers).
func (appKeepers *AppKeepers) InitNormalKeepers(
	appCodec codec.Codec,
	encodingConfig appparams.EncodingConfig,
	bApp *baseapp.BaseApp,
	maccPerms map[string][]string,
	wasmDir string,
	wasmConfig wasmtypes.WasmConfig,
	wasmOpts []wasmkeeper.Option,
	blockedAddress map[string]bool,
) {
	legacyAmino := encodingConfig.Amino
	// add keepers
	govAuthor := authtypes.NewModuleAddress(govtypes.ModuleName).String()
	appKeepers.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		appKeepers.keys[authtypes.StoreKey],
		authtypes.ProtoBaseAccount,
		maccPerms,
		sdk.Bech32PrefixAccAddr,
		govAuthor,
	)

	appKeepers.AuthzKeeper = authzkeeper.NewKeeper(
		appKeepers.keys[authz.ModuleName],
		appCodec,
		bApp.MsgServiceRouter(),
		appKeepers.AccountKeeper,
	)

	appKeepers.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		appKeepers.keys[banktypes.StoreKey],
		appKeepers.AccountKeeper,
		blockedAddress,
		govAuthor,
	)

	appKeepers.StakingKeeper = stakingkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[stakingtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		govAuthor,
	)

	appKeepers.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[feegrant.StoreKey],
		appKeepers.AccountKeeper,
	)

	appKeepers.MintKeeper = mintkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[minttypes.StoreKey],
		appKeepers.StakingKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		authtypes.FeeCollectorName,
		govAuthor,
	)

	appKeepers.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[distrtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		authtypes.FeeCollectorName,
		govAuthor,
	)

	slashingKeeper := slashingkeeper.NewKeeper(
		appCodec,
		legacyAmino,
		appKeepers.keys[slashingtypes.StoreKey],
		appKeepers.StakingKeeper,
		govAuthor,
	)

	appKeepers.SlashingKeeper = slashingKeeper

	groupConfig := group.DefaultConfig()
	/*
		Example of setting group params:
		groupConfig.MaxMetadataLen = 1000
	*/
	appKeepers.GroupKeeper = groupkeeper.NewKeeper(
		appKeepers.keys[group.StoreKey],
		appCodec,
		bApp.MsgServiceRouter(),
		appKeepers.AccountKeeper,
		groupConfig,
	)

	// Create IBC Keeper
	appKeepers.IBCKeeper = ibckeeper.NewKeeper(
		appCodec, appKeepers.keys[ibcexported.StoreKey],
		appKeepers.GetSubspace(ibcexported.ModuleName),
		appKeepers.StakingKeeper,
		appKeepers.UpgradeKeeper,
		appKeepers.ScopedIBCKeeper,
	)
	// IBC Fee Module keeper
	appKeepers.IBCFeeKeeper = ibcfeekeeper.NewKeeper(
		appCodec, appKeepers.keys[ibcfeetypes.StoreKey],
		appKeepers.IBCKeeper.ChannelKeeper, // may be replaced with IBC middleware
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper, appKeepers.AccountKeeper, appKeepers.BankKeeper,
	)

	appKeepers.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec, appKeepers.keys[icahosttypes.StoreKey],
		appKeepers.GetSubspace(icahosttypes.SubModuleName),
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.ScopedICAHostKeeper,
		bApp.MsgServiceRouter(),
	)

	appKeepers.IcacallbacksKeeper = *icacallbacksmodulekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[icacallbacksmoduletypes.StoreKey],
		appKeepers.keys[icacallbacksmoduletypes.MemStoreKey],
		appKeepers.GetSubspace(icacallbacksmoduletypes.ModuleName),
		*appKeepers.IBCKeeper,
	)

	scopedICAControllerKeeper := appKeepers.CapabilityKeeper.ScopeToModule(icacontrollertypes.SubModuleName)

	appKeepers.IcaControllerKeeper = icacontrollerkeeper.NewKeeper(
		appCodec, appKeepers.keys[icacontrollertypes.StoreKey],
		appKeepers.GetSubspace(icacontrollertypes.SubModuleName),
		appKeepers.IBCKeeper.ChannelKeeper, // may be replaced with middleware such as ics29 fee
		appKeepers.IBCKeeper.ChannelKeeper, &appKeepers.IBCKeeper.PortKeeper,
		scopedICAControllerKeeper, bApp.MsgServiceRouter(),
	)

	// Create Transfer Keepers
	appKeepers.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibctransfertypes.StoreKey],
		appKeepers.GetSubspace(ibctransfertypes.ModuleName),
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.ScopedTransferKeeper,
	)

	// Create evidence Keeper for to register the IBC light client misbehavior evidence route
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[evidencetypes.StoreKey],
		appKeepers.StakingKeeper,
		appKeepers.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	appKeepers.EvidenceKeeper = *evidenceKeeper

	govConfig := govtypes.DefaultConfig()

	govKeeper := govkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[govtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		bApp.MsgServiceRouter(),
		govConfig,
		govAuthor,
	)
	appKeepers.GovKeeper = govKeeper

	govRouter := govv1beta1.NewRouter()
	govRouter.
		AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(appKeepers.ParamsKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(appKeepers.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibcclient.NewClientProposalHandler(appKeepers.IBCKeeper.ClientKeeper))
	govKeeper.SetLegacyRouter(govRouter)

	appKeepers.GmmKeeper = *gmmmodulekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[gmmmoduletypes.StoreKey],
		appKeepers.keys[gmmmoduletypes.MemStoreKey],
		appKeepers.GetSubspace(gmmmoduletypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
	)

	appKeepers.YieldKeeper = *yieldmodulekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[yieldmoduletypes.StoreKey],
		appKeepers.keys[yieldmoduletypes.MemStoreKey],
		appKeepers.GetSubspace(yieldmoduletypes.ModuleName),

		appKeepers.BankKeeper,
		appKeepers.InterchainQueryKeeper,
		appKeepers.IBCKeeper,
		appKeepers.TransferKeeper,
		appKeepers.IcacallbacksKeeper,
	)

	// The last arguments can contain custom message handlers, and custom query handlers,
	// if we want to allow any custom callbacks
	supportedFeatures := "iterator,staking,stargate,side,cosmwasm_1_1,cosmwasm_1_2,cosmwasm_1_4"

	wasmOpts = append(swasm.RegisterCustomPlugins(&appKeepers.BankKeeper, &appKeepers.GmmKeeper), wasmOpts...)

	// Create Wasmd Keepers
	// this line is used by starport scaffolding # stargate/app/scopedKeeper
	// availableCapabilities := strings.Join(AllCapabilities(), ",")
	appKeepers.WasmKeeper = wasmkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[wasmtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		distrkeeper.NewQuerier(appKeepers.DistrKeeper),
		appKeepers.IBCFeeKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.scopedWasmKeeper,
		appKeepers.TransferKeeper,
		bApp.MsgServiceRouter(),
		bApp.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		supportedFeatures,
		govAuthor,
		wasmOpts...,
	)

	transferIBCModule := transfer.NewIBCModule(appKeepers.TransferKeeper)
	icaHostIBCModule := icahost.NewIBCModule(appKeepers.ICAHostKeeper)

	appKeepers.CapabilityKeeper.Seal()
	// wire up x/wasm to IBC
	// Create static IBC router, add transfer route, then set and seal it

<<<<<<< HEAD
	//var ics101WasmStack ibcporttypes.IBCModule
=======
>>>>>>> main
	ics101WasmStack := wasm.NewIBCHandler(appKeepers.WasmKeeper, appKeepers.IBCKeeper.ChannelKeeper, appKeepers.IBCKeeper.ChannelKeeper)
	// ics101WasmStack = packetforward.NewIBCMiddleware(
	// 	ics101WasmStack,
	// 	appKeepers.PacketForwardKeeper,
	// 	0, // retries on timeout
	// 	packetforwardkeeper.DefaultForwardTransferPacketTimeoutTimestamp, // forward timeout
	// 	packetforwardkeeper.DefaultRefundTransferPacketTimeoutTimestamp,  // refund timeout
	// )

	ibcRouter := ibcporttypes.NewRouter()
	ibcRouter.AddRoute(icahosttypes.SubModuleName, icaHostIBCModule).
		AddRoute(ibctransfertypes.ModuleName, transferIBCModule)
	ibcRouter.AddRoute(wasmtypes.ModuleName, ics101WasmStack)
	// this line is used by starport scaffolding # ibc/app/router
	appKeepers.IBCKeeper.SetRouter(ibcRouter)
}

// initParamsKeeper init params keeper and its subspaces.
func (appKeepers *AppKeepers) initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govv1.ParamKeyTable()) //nolint:staticcheck
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibcexported.ModuleName)
	paramsKeeper.Subspace(icacontrollertypes.SubModuleName)
	paramsKeeper.Subspace(icahosttypes.SubModuleName)
	paramsKeeper.Subspace(gmmmoduletypes.ModuleName)
	paramsKeeper.Subspace(yieldmoduletypes.ModuleName)

	return paramsKeeper
}

// SetupHooks sets up hooks for modules.
func (appKeepers *AppKeepers) SetupHooks() {
	appKeepers.GovKeeper.SetHooks(
		govtypes.NewMultiGovHooks(
		// insert governance hooks receivers here
		),
	)

	appKeepers.StakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(
			// insert staking hooks receivers here
			appKeepers.DistrKeeper.Hooks(),
			appKeepers.SlashingKeeper.Hooks(),
		),
	)
}

// TODO: We need to automate this, by bundling with a module struct...
func KVStoreKeys() []string {
	return []string{
		authtypes.StoreKey, authz.ModuleName, banktypes.StoreKey, stakingtypes.StoreKey,
		crisistypes.StoreKey, minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, ibcexported.StoreKey, upgradetypes.StoreKey,
		feegrant.StoreKey, evidencetypes.StoreKey, ibctransfertypes.StoreKey, icahosttypes.StoreKey,
		capabilitytypes.StoreKey, group.StoreKey, icacontrollertypes.StoreKey, consensusparamtypes.StoreKey, wasmtypes.StoreKey,
		ibcfeetypes.StoreKey,
		gmmmoduletypes.StoreKey,
		yieldmoduletypes.StoreKey,
	}
}

func AllCapabilities() []string {
	return []string{
		"iterator",
		"staking",
		"stargate",
		"cosmwasm_1_1",
		"cosmwasm_1_2",
		"cosmwasm_1_3",
	}
}
