package cmd

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func initSDKConfig() {
	// Set prefixes
	bech32Prefix := "side"
	accountPubKeyPrefix := bech32Prefix + "pub"
	validatorAddressPrefix := bech32Prefix + "valoper"
	validatorPubKeyPrefix := bech32Prefix + "valoperpub"
	consNodeAddressPrefix := bech32Prefix + "valcons"
	consNodePubKeyPrefix := bech32Prefix + "valconspub"

	// Set and seal config
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(bech32Prefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	// config.SetBtcChainCfg(&chaincfg.SigNetParams)

	config.Seal()
}
