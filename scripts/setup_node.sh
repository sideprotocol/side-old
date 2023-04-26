#!/bin/bash
# Default configuration variables
BINARY="sidechaind --home $HOME/.sidechaind"
DEFAULT_CHAIN_ID="sidechain_7070-1"
DEFAULT_DENOM="aside"
BOB_MNEMONIC="actress letter whip youth flip sort announce chief traffic side destroy seek parade warrior awake scan panther nominee harsh spawn differ enroll glue become"
DEPOSIT="10aside"
PROPOSAL_FILE="proposal.json"

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
$BINARY init sidenode --chain-id $CHAIN_ID --keyring-backend test

# Set custom denom
GENESIS_FILE="$HOME/.sidechaind/config/genesis.json"
TMP_GENESIS_FILE="$HOME/.sidechaind/config/genesis_modified.json"
jq --arg denom "$DENOM" '.app_state.crisis.constant_fee.denom = $denom | .app_state.gov.deposit_params.min_deposit[0].denom = $denom | .app_state.mint.params.mint_denom = $denom | .app_state.staking.params.bond_denom = $denom' $GENESIS_FILE > "$TMP_GENESIS_FILE"
mv "$TMP_GENESIS_FILE" $GENESIS_FILE

# Create Alice's account
echo "Creating Alice's account..."
$BINARY keys add alice --keyring-backend test
$BINARY add-genesis-account alice $ALICE_COINS

# Create Bob's account
echo "Creating Bob's account..."
$BINARY keys add bob  --keyring-backend test
$BINARY add-genesis-account bob $ALICE_COINS

# Bond token as a validator
echo "Bond token as a validator"
INIT_COINS="10000000000000000000000${DENOM}"
$BINARY gentx alice $INIT_COINS --chain-id $CHAIN_ID
$BINARY collect-gentxs

# Validate genesis file
echo "Validating genesis file..."
$BINARY validate-genesis

echo "Node setup completed."
