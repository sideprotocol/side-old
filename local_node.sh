#!/bin/bash

KEYS=("validator" "test" "relayer" "btc-vault" "runes-vault")
CHAINID="grimoria-testnet-1"
MONIKER="Side Labs"
BINARY="$HOME/go/bin/sided"
DENOM_STR="uside,uusdc,uusdt"
INITIAL_ACCOUNT_STR="tb1qcr8te4kr609gcawutmrza0j4xv80jy8zmfp6l0"
set -f
IFS=,
DENOMS=($DENOM_STR)
INITIAL_ACCOUNTS=($INITIAL_ACCOUNT_STR)

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
HOMEDIR="$HOME/.side"

# Path variables
CONFIG=$HOMEDIR/config/config.toml
APP_TOML=$HOMEDIR/config/app.toml
GENESIS=$HOMEDIR/config/genesis.json
TMP_GENESIS=$HOMEDIR/config/tmp_genesis.json

# validate dependencies are installed
command -v jq >/dev/null 2>&1 || {
	echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"
	exit 1
}

# used to exit on first error (any non-zero exit code)
set -e

# Reinstall daemon
make install

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

	# Set client config
	$BINARY config keyring-backend $KEYRING --home "$HOMEDIR"
	$BINARY config chain-id $CHAINID --home "$HOMEDIR"

	# If keys exist they should be deleted
	for KEY in "${KEYS[@]}"; do
		$BINARY keys add "$KEY" --keyring-backend $KEYRING --algo $KEYALGO --home "$HOMEDIR"
	done
	# for KEY in "${KEYS[@]}"; do
    # # Add the --recover flag to initiate recovery mode
    # 	$BINARY keys add "$KEY" --keyring-backend $KEYRING --algo $KEYALGO --recover --home "$HOMEDIR"
	# done

	echo ""
	echo "☝️ Copy the above mnemonic phrases and import them to relayer! Press [Enter] to continue..."
	read -r continue

	# Set moniker and chain-id for Cascadia (Moniker can be anything, chain-id must be an integer)
	$BINARY init $MONIKER -o --chain-id $CHAINID --home "$HOMEDIR"

	jq --arg denom "${DENOMS[0]}" '.app_state["staking"]["params"]["bond_denom"]=$denom' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	jq --arg denom "${DENOMS[0]}" '.app_state["mint"]["params"]["mint_denom"]=$denom' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	jq --arg denom "${DENOMS[0]}" '.app_state["crisis"]["constant_fee"]["denom"]=$denom' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	jq --arg denom "${DENOMS[0]}" '.app_state["gov"]["deposit_params"]["min_deposit"][0]["denom"]=$denom' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	jq --arg denom "${DENOMS[0]}" '.app_state["gov"]["params"]["min_deposit"][0]["denom"]=$denom' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	jq --arg gas "$BLOCK_GAS" '.app_state["feemarket"]["block_gas"]=$gas' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	# Set gas limit in genesis
	jq --arg max_gas "$MAX_GAS" '.consensus_params["block"]["max_gas"]=$max_gas' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	# setup relayers
	RELAYER=$($BINARY keys show "${KEYS[2]}" -a --keyring-backend $KEYRING --home "$HOMEDIR")
	jq --arg relayer "$RELAYER" '.app_state["btcbridge"]["params"]["authorized_relayers"][0]=$relayer' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	# setup vaults
	BTC_VAULT=$($BINARY keys show "${KEYS[3]}" -a --keyring-backend $KEYRING --home "$HOMEDIR")
	jq --arg btc_vault "$BTC_VAULT" '.app_state["btcbridge"]["params"]["vaults"][0]["address"]=$btc_vault' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	PUKEY=$($BINARY keys show "${KEYS[3]}" --pubkeyhex --keyring-backend $KEYRING --home "$HOMEDIR")
	jq --arg pubkey "$PUKEY" '.app_state["btcbridge"]["params"]["vaults"][0]["pub_key"]=$pubkey' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	RUNES_VAULT=$($BINARY keys show "${KEYS[4]}" -a --keyring-backend $KEYRING --home "$HOMEDIR")
	jq --arg runes_vault "$RUNES_VAULT" '.app_state["btcbridge"]["params"]["vaults"][1]["address"]=$runes_vault' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	PUKEY=$($BINARY keys show "${KEYS[4]}" --pubkeyhex --keyring-backend $KEYRING --home "$HOMEDIR")
	jq --arg pubkey "$PUKEY" '.app_state["btcbridge"]["params"]["vaults"][1]["pub_key"]=$pubkey' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

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
	    for key in "${!DENOMS[@]}"; do
	        BALANCES+=",${INITIAL_SUPPLY}${DENOMS[$key]}"
	    done
	    echo ${BALANCES:1}
	    $BINARY add-genesis-account "$KEY" ${BALANCES:1} --keyring-backend $KEYRING --home "$HOMEDIR"
	done

	echo "Genesis accounts allocated for local accounts"

	# Allocate genesis accounts (cosmos formatted addresses)
	for ADDR in "${INITIAL_ACCOUNTS[@]}"; do
	    BALANCES=""
	    for key in "${!DENOMS[@]}"; do
	        BALANCES+=",${INITIAL_SUPPLY}${DENOMS[$key]}"
	    done
	    echo ${BALANCES:1}
	    $BINARY add-genesis-account "$ADDR" ${BALANCES:1} --home "$HOMEDIR"
	done
	echo "Genesis accounts allocated for initial accounts"



	# Sign genesis transaction
	# echo $INITIAL_SUPPLY${DENOMS[0]}
	$BINARY gentx "${KEYS[0]}" $INITIAL_SUPPLY${DENOMS[0]} --keyring-backend $KEYRING --chain-id $CHAINID --identity "666AC57CC678BEC4" --website="https://side.one" --home "$HOMEDIR"
	echo "Genesis transaction signed"

	## In case you want to create multiple validators at genesis
	## 1. Back to `$BINARY keys add` step, init more keys
	## 2. Back to `$BINARY add-genesis-account` step, add balance for those
	## 3. Clone this ~/.$BINARY home directory into some others, let's say `~/.clonedCascadiad`
	## 4. Run `gentx` in each of those folders
	## 5. Copy the `gentx-*` folders under `~/.clonedCascadiad/config/gentx/` folders into the original `~/.$BINARY/config/gentx`

	# Collect genesis tx
	$BINARY collect-gentxs --home "$HOMEDIR"
	echo "Genesis transactions collected"

	# Run this to ensure everything worked and that the genesis file is setup correctly
	$BINARY validate-genesis --home "$HOMEDIR"
	echo "Genesis file validated"

	if [[ $1 == "pending" ]]; then
		echo "pending mode is on, please wait for the first block committed."
	fi
fi


# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
$BINARY start --log_level info --minimum-gas-prices=0.0001${DENOMS[0]} --home "$HOMEDIR"
