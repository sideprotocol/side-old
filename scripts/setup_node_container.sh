#!/bin/bash
set -e

# Default configuration variables
BINARY="sidechaind --home /root/.sidechaind"
DEFAULT_CHAIN_ID="sidehub-1818-1"
DEFAULT_DENOM="aside"

# Parse command-line options
while getopts ":c:d:" opt; do
  case ${opt} in
    c )
      CHAIN_ID=$OPTARG
      ;;
    d )
      DENOM=$OPTARG
      ;;
    \? )
      echo "Usage: ${0} [-c chain_id] [-d denom]"
      exit 1
      ;;
    : )
      echo "Option -$OPTARG requires an argument."
      exit 1
      ;;
  esac
done

# Set default values if not provided as options
CHAIN_ID=${CHAIN_ID:-$DEFAULT_CHAIN_ID}
DENOM=${DENOM:-$DEFAULT_DENOM}
ALICE_COINS="990000000000000000000000${DENOM}"
BOB_COINS="10000000000000000000000${DENOM}"

# Initialize the node
echo "Initialize the node"
$BINARY init sidenode --chain-id $CHAIN_ID --keyring-backend test


# Set custom denom
echo "Set custom denom"
GENESIS_FILE="/root/.sidechaind/config/genesis.json"
TMP_GENESIS_FILE="/root/.sidechaind/config/genesis_modified.json"
jq --arg denom "$DENOM" '.app_state.crisis.constant_fee.denom = $denom | .app_state.gov.deposit_params.min_deposit[0].denom = $denom | .app_state.mint.params.mint_denom = $denom | .app_state.staking.params.bond_denom = $denom' $GENESIS_FILE > "$TMP_GENESIS_FILE"
mv "$TMP_GENESIS_FILE" $GENESIS_FILE

# Set custom chain_id
echo "Set custom chain_id"
jq --arg chain_id "$CHAIN_ID" '.chain_id = $chain_id' $GENESIS_FILE > "$TMP_GENESIS_FILE"
mv "$TMP_GENESIS_FILE" $GENESIS_FILE


# Create Alice's account
echo "Creating Alice's account..."
ALICE_OUTPUT=$($BINARY keys add alice --keyring-backend test)
echo "Alice output: $ALICE_OUTPUT"
ALICE_ADDRESS=$(echo "$ALICE_OUTPUT" | grep -oPm1 "(?<=address: )[^\s]+")
echo "Alice address: $ALICE_ADDRESS"
$BINARY add-genesis-account $ALICE_ADDRESS $ALICE_COINS


# Create Bob's account
echo "Creating Bob's account..."
BOB_OUTPUT=$($BINARY keys add bob --keyring-backend test)
echo "Alice output: $BOB_OUTPUT"
BOB_ADDRESS=$(echo "$BOB_OUTPUT" | grep -oPm1 "(?<=address: )[^\s]+")
echo "Bob address: $BOB_ADDRESS"
$BINARY add-genesis-account $BOB_ADDRESS $BOB_COINS

# Bond token as a validator
echo "Bond token as a validator"
INIT_COINS="10000000000000000000000${DENOM}"
$BINARY gentx alice $INIT_COINS --chain-id $CHAIN_ID --keyring-backend test
$BINARY collect-gentxs

# Validate genesis file
echo "Validating genesis file..."
$BINARY validate-genesis

echo "Node setup completed."
