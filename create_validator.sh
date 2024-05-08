#!/bin/bash

KEYS=("dev0" "dev1")
CHAINID="S2-testnet-1"
MONIKER="Side Labs"
BINARY="$HOME/go/bin/sided"

# Remember to change to other types of keyring like 'file' in-case exposing to outside world,
# otherwise your balance will be wiped quickly
# The keyring test does not require private key to steal tokens from you
KEYRING="test"
#KEYALGO="secp256k1"
KEYALGO="segwit"
LOGLEVEL="info"
# Set dedicated home directory for the $BINARY instance
HOMEDIR="$HOME/.side2"

# Path variables
CONFIG=$HOMEDIR/config/config.toml
APP_TOML=$HOMEDIR/config/app.toml
GENESIS=$HOMEDIR/config/genesis.json
TMP_GENESIS=$HOMEDIR/config/tmp_genesis.json
PEER_ID=""

# validate dependencies are installed
command -v jq >/dev/null 2>&1 || {
	echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"
	exit 1
}

# used to exit on first error (any non-zero exit code)
set -e

# Set client config
$BINARY config keyring-backend $KEYRING --home "$HOMEDIR"
$BINARY config chain-id $CHAINID --home "$HOMEDIR"

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
$BINARY tx staking create-validator -y --moniker segwit --pubkey=$($BINARY tendermint show-validator --home $HOMEDIR) --amount="10000000uside" --commission-rate 0.1 --commission-max-rate 0.1 --commission-max-change-rate 0.1 --min-self-delegation 1 --chain-id $CHAINID --fees 2000uside --home "$HOMEDIR" --from segwit_key
sleep 6
$BINARY tx staking create-validator -y --moniker secp256 --pubkey=$($BINARY tendermint show-validator --home $HOMEDIR) --amount="15000000uside" --commission-rate 0.1 --commission-max-rate 0.1 --commission-max-change-rate 0.1 --min-self-delegation 1 --chain-id $CHAINID --fees 2000uside --home "$HOMEDIR" --from secp256k1_key

