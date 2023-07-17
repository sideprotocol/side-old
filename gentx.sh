KEYS[0]="dev0"
KEYS[1]="dev1"
KEYS[2]="dev2"
MNEMONIC[0]="outside remember acquire boost brand gate pigeon follow involve zoo convince gym exhibit toss left blame inhale black imitate orchard entire position hollow destroy"
MNEMONIC[1]="crumble token absorb subject tenant young december sick fiction length pride cigar diagram dynamic innocent expand board shoulder auto deny swim bench sketch palace"
MNEMONIC[2]="clutch also human garment possible will pool cube best lens smart latin wear such skate moon water trap camp top spy radio jar solid"

CHAINID="sidechain_1818-1"
MONIKER="localtestnet"
BINARY="sidechaind"
DENOMS=("aside" "aetc" "ausdc" "aeth")
INITIAL_SUPPLY="100000000000000000000000000"
BLOCK_GAS=10000000
MAX_GAS=10000000000
VALIDATOR=1

# Remember to change to other types of keyring like 'file' in-case exposing to outside world,
# otherwise your balance will be wiped quickly
# The keyring test does not require private key to steal tokens from you
KEYRING="test"
KEYALGO="eth_secp256k1"
LOGLEVEL="info"
# Set dedicated home directory for the $BINARY instance
HOMEDIR="$HOME/.$BINARY-$VALIDATOR"
# to trace evm
#TRACE="--trace"
TRACE=""

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
	for KEY in "${!KEYS[@]}"; do
		#$BINARY keys add "$KEY" --keyring-backend $KEYRING --algo $KEYALGO --recover --home "$HOMEDIR"
        echo "${MNEMONIC[$KEY]}" | $BINARY keys add "${KEYS[$KEY]}" --keyring-backend $KEYRING --algo $KEYALGO --recover --home "$HOMEDIR"
    done

	# Set moniker and chain-id for Cascadia (Moniker can be anything, chain-id must be an integer)
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

	sed -i.bak "s/^persistent_peers =.*/persistent_peers = \"$PEERS\"/" "$HOMEDIR/config/config.toml"

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
    $BINARY gentx "${KEYS[$VALIDATOR]}" ${INITIAL_SUPPLY}${DENOMS[0]} --keyring-backend $KEYRING --chain-id $CHAINID --home "$HOMEDIR"
fi


