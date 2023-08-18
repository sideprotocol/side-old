#!/bin/bash
# Remove original setup
rm -rf testnet

DEFAULTDENOM="aside"
TESTNET_PATH=${PWD}/testnet
cd build/

# Initialize nodes
for i in {0..4}
do
    ./sidechaind init "node$i" --home ${TESTNET_PATH}/node$i --chain-id side_1818-1
done

# Add keys for each node
for i in {0..4}
do
    yes y | ./sidechaind keys add "wm-node$i" --home ${TESTNET_PATH}/node$i
done

# Add genesis accounts with balances for each node
for i in {0..4}
do
    ADDRESS=$(./sidechaind keys show "wm-node$i" --home ${TESTNET_PATH}/node$i/ | grep address | awk 'BEGIN { FS=": " } { print $2 }')
    if [ -z "$ADDRESS" ]; then
        echo "Error fetching address for wm-node$i"
        exit 1
    fi
    ./sidechaind add-genesis-account "$ADDRESS" 100000000000000000000000aside,100000000000000000000000usdc --home ${TESTNET_PATH}/node0/
done

# Copy genesis.json to all nodes
for i in {1..4}
do
    cp ${TESTNET_PATH}/node0/config/genesis.json ${TESTNET_PATH}/node$i/config/
done

mkdir -p ${TESTNET_PATH}/gentx/

# Generate gentx for each node
for i in {0..4}
do
    ./sidechaind gentx "wm-node$i" 7000000000000aside --chain-id side_1818-1 --home ${TESTNET_PATH}/node$i/ --ip "side_testnet_node$i" --output-document ${TESTNET_PATH}/gentx/node$i.json
    if [ $? -ne 0 ]; then
        echo "Error generating gentx for wm-node$i"
        exit 1
    fi
done

# Collect gentxs
./sidechaind collect-gentxs --gentx-dir ${TESTNET_PATH}/gentx/ --home ${TESTNET_PATH}/node0

if [ $? -ne 0 ]; then
    echo "Error collecting gentxs"
    exit 1
fi

echo "Script executed successfully!"

# Change the goverance parameters
jq '.app_state.gov.voting_params.voting_period = "600s"' ${TESTNET_PATH}/node0/config/genesis.json > temp.json && mv temp.json ${TESTNET_PATH}/node0/config/genesis.json
jq '.app_state.gov.deposit_params.max_deposit_period = "600s"' ${TESTNET_PATH}/node0/config/genesis.json > temp.json && mv temp.json ${TESTNET_PATH}/node0/config/genesis.json
jq '.app_state.gov.deposit_params.min_deposit[0].denom = "gov"' ${TESTNET_PATH}/node0/config/genesis.json > temp.json && mv temp.json ${TESTNET_PATH}/node0/config/genesis.json

# Change Maximum validators
jq '.app_state.staking.params.max_validators = 50' ${TESTNET_PATH}/node0/config/genesis.json > temp.json && mv temp.json ${TESTNET_PATH}/node0/config/genesis.json

# Change Default denom
for i in {0..4}
do
    sed -i 's/"bond_denom": "stake"/"bond_denom": "aside"/g' ${TESTNET_PATH}/node$i/config/genesis.json
    sed -i 's/"mint_denom": "stake"/"mint_denom": "aside"/g' ${TESTNET_PATH}/node$i/config/genesis.json
done

# Enable API
# Detect the OS
OS=$(uname)

# Define the sed command based on the OS
if [ "$OS" = "Darwin" ]; then
    # macOS
    SED_CMD="sed -i '' -E"
else
    # Linux
    SED_CMD="sed -i"
fi

# Apply the sed command
$SED_CMD '0,/enable = false/ s/enable = false/enable = true/' ${TESTNET_PATH}/node0/config/app.toml
$SED_CMD '0,/enable = false/ s/enable = false/enable = true/' ${TESTNET_PATH}/node1/config/app.toml
$SED_CMD '0,/enable = false/ s/enable = false/enable = true/' ${TESTNET_PATH}/node2/config/app.toml
$SED_CMD '0,/enable = false/ s/enable = false/enable = true/' ${TESTNET_PATH}/node3/config/app.toml
$SED_CMD '0,/enable = false/ s/enable = false/enable = true/' ${TESTNET_PATH}/node4/config/app.toml


# Replace geneiss to every node 
cp ${TESTNET_PATH}/node0/config/genesis.json ${TESTNET_PATH}/node1/config/
cp ${TESTNET_PATH}/node0/config/genesis.json ${TESTNET_PATH}/node2/config/
cp ${TESTNET_PATH}/node0/config/genesis.json ${TESTNET_PATH}/node3/config/
cp ${TESTNET_PATH}/node0/config/genesis.json ${TESTNET_PATH}/node4/config/
