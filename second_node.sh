#!/bin/bash

KEYS=("dev0" "dev1")
CHAINID="S2-testnet-1"
MONIKER="Side Labs"
BINARY="$HOME/go/bin/sided"
DENOM_STR="uside,ubtct,uusdc,uusdc.axl,uusdc.noble,uusdt,uusdt.kava,uusdt.axl,uwbtc.axl,uwbtc.osmo,uwbtc"

set -f
IFS=,
DENOMS=($DENOM_STR)

IFS=";"


INITIAL_SUPPLY="500000000000000"
BLOCK_GAS=10000000
MAX_GAS=10000000000

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
PEER_ID="peer"

# validate dependencies are installed
command -v jq >/dev/null 2>&1 || {
	echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"
	exit 1
}

# used to exit on first error (any non-zero exit code)
set -e

# Reinstall daemon
# make install

# User prompt if an existing local node configuration is found.
if [ -d "$HOMEDIR" ]; then
	printf "\nAn existing folder at '%s' was found. You can choose to delete this folder and start a new local node with new keys from genesis. When declined, the existing local node is started. \n" "$HOMEDIR"
	echo "Overwrite the existing configuration and start a new local node? [y/n]"
	read -r overwrite
else
	overwrite="Y"
fi


# Setup local node if overwrite is set to Yes, otherwise skip setup
if [[ $overwrite == "y" || $overwrite == "Y" ]]; then
	# Remove the previous folder
	rm -rf "$HOMEDIR"

	$BINARY init $MONIKER -o --chain-id $CHAINID --home "$HOMEDIR"

	# Set client config
	cp -r ~/.side/config/genesis.json $HOMEDIR/config/genesis.json
	$BINARY config keyring-backend $KEYRING --home "$HOMEDIR"
	$BINARY config chain-id $CHAINID --home "$HOMEDIR"

	sed -i.bak 's/127.0.0.1:26657/0.0.0.0:16657/g' "$CONFIG"
	sed -i.bak 's/127.0.0.1:26658/0.0.0.0:16658/g' "$CONFIG"
	sed -i.bak 's/0.0.0.0:26656/0.0.0.0:16656/g' "$CONFIG"
	sed -i.bak 's/persistent_peers = ""/persistent_peers = "dfd3e3c99414aa850f6e269cf4a674a66062cd49@127.0.0.1:26656"/g' "$CONFIG"
	#sed -i 's/persistent_peers = "$PEERID"/g' "$CONFIG"

	sed -i.bak 's/swagger = false/swagger = true/g' $APP_TOML
	sed -i.bak 's/localhost:9090/localhost:8090/g' $APP_TOML

	$BINARY keys add secp256k1_key --key-type segwit --home "$HOMEDIR"
	$BINARY tx bank send dev1 $($BINARY keys show secp256k1_key -a --home "$HOMEDIR") 2000000000uside --chain-id $CHAINID --fees 200uside --yes
	sleep 6
	$BINARY keys add segwit_key --key-type segwit --home "$HOMEDIR"
	$BINARY tx bank send dev1 $($BINARY keys show segwit_key -a --home "$HOMEDIR") 2000000000uside --chain-id $CHAINID --fees 200uside --yes


fi

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
$BINARY start --log_level info --minimum-gas-prices=0.0001${DENOMS[0]} --home "$HOMEDIR" --pruning="everything"
