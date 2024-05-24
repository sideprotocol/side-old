package cmd

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sideprotocol/side/app"
)

func initSDKConfig() {
	// Set prefixes
	accountPubKeyPrefix := app.BitcoinChainCfg.Bech32HRPSegwit + "pub"
	validatorAddressPrefix := app.BitcoinChainCfg.Bech32HRPSegwit + "valoper"
	validatorPubKeyPrefix := app.BitcoinChainCfg.Bech32HRPSegwit + "valoperpub"
	consNodeAddressPrefix := app.BitcoinChainCfg.Bech32HRPSegwit + "valcons"
	consNodePubKeyPrefix := app.BitcoinChainCfg.Bech32HRPSegwit + "valconspub"

	// Set and seal config
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(app.BitcoinChainCfg.Bech32HRPSegwit, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)

	config.Seal()
}
