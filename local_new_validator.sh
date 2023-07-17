#!/bin/bash
KEYS[0]="dev0"
KEYS[1]="dev1"
KEYS[2]="dev2"
CHAINID="sidechain_1818-1"
MONIKER="localtestnet"
BINARY="sidechaind"
DENOMS=("aside" "aetc" "ausdc" "aeth")
INITIAL_SUPPLY="100000000000000000000000000"
BLOCK_GAS=10000000
MAX_GAS=10000000000

# Remember to change to other types of keyring like 'file' in-case exposing to outside world,
# otherwise your balance will be wiped quickly
# The keyring test does not require private key to steal tokens from you
KEYRING="test"
KEYALGO="eth_secp256k1"
LOGLEVEL="info"
# Set dedicated home directory for the $BINARY instance
HOMEDIR="$HOME/.$BINARY"
# to trace evm


# Step 4: Submit a Transaction to Become a Validator
# Step 2: Create a Validator Key
$BINARY tx staking create-validator \
  --amount=${INITIAL_SUPPLY}${DENOMS[0]} \
  --pubkey=$(sidechaind tendermint show-validator) \
  --moniker="dev2" \
  --chain-id=$CHAINID \
  --commission-rate="0.10" \
  --commission-max-rate="0.20" \
  --commission-max-change-rate="0.01" \
  --min-self-delegation="1" \
  --gas="auto" \
  --gas-prices="0.1aside" \
  --from="dev2"