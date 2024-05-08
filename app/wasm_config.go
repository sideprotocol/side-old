package app

import (
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

const (
	// DefaultInstanceCost is initially set the same as in wasmd
	DefaultInstanceCost uint64 = 60_000
	// DefaultCompileCost set to a large number for testing
	DefaultCompileCost uint64 = 100
)

// MunGasRegisterConfig is defaults plus a custom compile amount
func GasRegisterConfig() wasmtypes.WasmGasRegisterConfig {
	gasConfig := wasmtypes.DefaultGasRegisterConfig()
	gasConfig.InstanceCost = DefaultInstanceCost
	gasConfig.CompileCost = DefaultCompileCost

	return gasConfig
}

func NewSideWasmGasRegister() wasmtypes.WasmGasRegister {
	return wasmtypes.NewWasmGasRegister(GasRegisterConfig())
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
