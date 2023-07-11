#!/bin/bash
set -e

# Default configuration variables
HOMEDIR="/root/.sidechaind/"
BINARY="sidechaind --home /root/.sidechaind"
DEFAULT_CHAIN_ID="sidehub_1818-1"
DEFAULT_DENOMS=("aside")
KEYS[0]="dev0"
KEYS[1]="dev1"
KEYS[2]="dev2"
MONIKER="localtestnet"
KEYRING="test"
KEYALGO="eth_secp256k1"
INITIAL_SUPPLY="100000000000000000000000000"
BLOCK_GAS=10000000
MAX_GAS=10000000000

# Path variables
CONFIG=$HOMEDIR/config/config.toml
APP_TOML=$HOMEDIR/config/app.toml
GENESIS=$HOMEDIR/config/genesis.json
TMP_GENESIS=$HOMEDIR/config/tmp_genesis.json


# Parse command-line options
while getopts ":c:d:" opt; do
  case ${opt} in
    c )
      CHAINID=$OPTARG
      ;;
    d )
      IFS=',' read -ra DENOMS <<< "$OPTARG"
      ;;
    o )
      IFS=',' read -ra overwrite <<< "$OPTARG"
      ;;
    \? )
      echo "Usage: ${0} [-c chain_id] [-d denom1,denom2,...]"
      exit 1
      ;;
    : )
      echo "Option -$OPTARG requires an argument."
      exit 1
      ;;
  esac
done

# Check if denoms are provided
if [ ${#DENOMS[@]} -eq 0 ]; then
  DENOMS=${DENOMS:-$DEFAULT_DENOMS}
fi

# Set default values if not provided as options
CHAINID=${CHAIN_ID:-$DEFAULT_CHAIN_ID}
DENOMS=${DENOMS:-$DEFAULT_DENOMS}
overwrite=${overwrite:-'y'} 
# validate dependencies are installed
command -v jq >/dev/null 2>&1 || {
	echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"
	exit 1
}

apt-get install bc

# # User prompt if an existing local node configuration is found.
if [ -d "$HOMEDIR" ]; then
	# printf "\nAn existing folder at '%s' was found. You can choose to delete this folder and start a new local node with new keys from genesis. When declined, the existing local node is started. \n" "$HOMEDIR"
	# echo "Overwrite the existing configuration and start a new local node? [y/n]"
	# read -r overwrite
  overwrite=$overwrite
else
	overwrite="Y"
fi


# Setup local node if overwrite is set to Yes, otherwise skip setup
if [[ $overwrite == "y" || $overwrite == "Y" ]]; then
	# Remove the previous folder
	rm -rf "$HOMEDIR"*

	# Set client config
	$BINARY config keyring-backend $KEYRING --home "$HOMEDIR"
	$BINARY config chain-id $CHAINID --home "$HOMEDIR"

	# If keys exist they should be deleted
	for KEY in "${KEYS[@]}"; do
		$BINARY keys add "$KEY" --keyring-backend $KEYRING --algo $KEYALGO --home "$HOMEDIR"
	done

	# Set moniker and chain-id for $BINARY (Moniker can be anything, chain-id must be an integer)
	$BINARY init $MONIKER -o --chain-id $CHAINID --home "$HOMEDIR"

	jq --arg denom "${DENOMS[0]}" '.app_state["staking"]["params"]["bond_denom"]=$denom' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	jq --arg denom "${DENOMS[0]}" '.app_state["crisis"]["constant_fee"]["denom"]=$denom' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	jq --arg denom "${DENOMS[0]}" '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]=$denom' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	jq --arg gas "$BLOCK_GAS" '.app_state["feemarket"]["block_gas"]=$gas' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	# Set gas limit in genesis
	jq --arg max_gas "$MAX_GAS" '.consensus_params["block"]["max_gas"]=$max_gas' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	# set custom pruning settings
	sed -i.bak 's/pruning = "default"/pruning = "custom"/g' "$APP_TOML"
	sed -i.bak 's/pruning-keep-recent = "0"/pruning-keep-recent = "2"/g' "$APP_TOML"
	sed -i.bak 's/pruning-interval = "0"/pruning-interval = "10"/g' "$APP_TOML"

	sed -i.bak 's/127.0.0.1:26657/0.0.0.0:26657/g' "$CONFIG"
	sed -i.bak 's/cors_allowed_origins\s*=\s*\[\]/cors_allowed_origins = ["*",]/g' "$CONFIG"
	sed -i.bak 's/swagger = false/swagger = true/g' $APP_TOML

	# Allocate genesis accounts (cosmos formatted addresses)
	for KEY in "${KEYS[@]}"; do
	    BALANCES=""
	    for DENOM in "${DENOMS[@]}"; do
	        BALANCES+=",${INITIAL_SUPPLY}$DENOM"
	    done
	    $BINARY add-genesis-account "$KEY" ${BALANCES:1} --keyring-backend $KEYRING --home "$HOMEDIR"
	done

	# Adjust total supply
	for DENOM in "${DENOMS[@]}"; do
	    total_supply=$(echo "${#KEYS[@]} * $INITIAL_SUPPLY" | bc)
	    if ! jq -e --arg denom "$DENOM" '.app_state["bank"]["supply"] | any(.denom == $denom)' "$GENESIS" >/dev/null; then
	        jq -r --arg total_supply "$total_supply" --arg denom "$DENOM" '.app_state["bank"]["supply"] += [{"denom": $denom, "amount": $total_supply}]' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	    fi
	done
	
	# Sign genesis transaction
	$BINARY gentx "${KEYS[0]}" ${INITIAL_SUPPLY}${DENOMS[0]} --keyring-backend $KEYRING --chain-id $CHAINID --home "$HOMEDIR"

	## In case you want to create multiple validators at genesis
	## 1. Back to `$BINARY keys add` step, init more keys
	## 2. Back to `$BINARY add-genesis-account` step, add balance for those
	## 3. Clone this ~/.$BINARY home directory into some others, let's say `~/.clonedCascadiad`
	## 4. Run `gentx` in each of those folders
	## 5. Copy the `gentx-*` folders under `~/.clonedCascadiad/config/gentx/` folders into the original `~/.$BINARY/config/gentx`

	# Collect genesis tx
	$BINARY collect-gentxs --home "$HOMEDIR"

	# Run this to ensure everything worked and that the genesis file is setup correctly
	$BINARY validate-genesis --home "$HOMEDIR"

	if [[ $1 == "pending" ]]; then
		echo "pending mode is on, please wait for the first block committed."
	fi
fi



