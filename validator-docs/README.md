# SIDE Node Installation & Setup

Instruction to install and configure the sided binary

## Hardware Specifications
1. Minimum Requirements
    1. CPU: 2 cores
    2. RAM: 4 GB
    3. Storage: 200 GB
    4. Network: 1 Gbps

2. Recommended Specifications
    1. CPU: 4 cores
    2. RAM: 8 GB
    3. Storage: 500 GB
    4. Network: 1 Gbps

## Operating System

The choice of operating system for your node is entirely based on your personal preference. You can compile the sided daemon on a wide range of modern Linux distributions and recent versions of macOS.
```
The tutorial assumes that you are utilizing an Ubuntu LTS release. If you have opted for a different operating system, you will need to adjust the commands accordingly to align with your chosen OS.
```

## Prerequisites

Golang v1.20 [go releases and instructions][https://go.dev/dl/].

## Build sided from source

1. Ensure that you have the necessary version of Golang installed.

    `go version`

The output must align with the Golang version specified in the Prerequisites section.

2. Clone the source code from the repository and navigate to the cloned directory using the cd command.

```
git clone -b dev https://github.com/sideprotocol/sidechain.git
cd sidechain
git checkout v0.0.3
```

3. Compile the sided binary.

    `make install`

The provided command will compile the sided binary and save it in your $GOBIN directory. If $GOBIN is included in your $PATH, you should be able to execute the sided binary.
```
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

```
sided version
0.0.3
```

If you encounter any issues related to PATH settings, please consult the provided Go releases and instructions link mentioned in the prerequisites section.

# Run a Full Node

How to run a full node of SIDE blockchain

## Setting up the keyring (Optional)

The keyring holds the private/public keypairs used to interact with a node. For instance, a validator key needs to be set up before running the blockchain node, so that blocks can be correctly signed. The private key can be stored in different locations, called "backends", such as a file or the operating system's own key storage.

```
Available backends for the keyring
1. The os backend
2. The file backend
3. The pass backend
4. The kwallet backend
5. The test backend
6. The memory backend
You can visit https://docs.cosmos.network/v0.47/user/run-node/keyring for more details
```

`export SIDED_KEYRING_BACKEND=os`

## Adding keys to the keyring

```
sided keys add alice
Enter keyring passphrase (attempt 1/3):
​
- address: side1raru385lzypr5jw2al4dl9ls5900t02f99z7dw
  name: alice
  pubkey: '{"@type":"/cosmos.crypto.secp256k1.PubKey","key":"Ak43PcCCEr8J0ljTUN+xs0nJiLlwrqHVsii8uRNzX7M5"}'
  type: local
​
​
**Important** write this mnemonic phrase in a safe place.
It is the only way to recover your account if you ever forget your password.
​
enroll index bind glad tonight rhythm barely negative south quarter main length venue funny dance short loud foil electric thumb anger similar like nice
```

This command generates a new 24-word mnemonic phrase, persists it to the relevant backend, and outputs information about the keypair. If this keypair will be used to hold value-bearing tokens, be sure to write down the mnemonic phrase somewhere safe!

By default, the keyring generates a secp256k1 keypair. The keyring also supports ed25519 keys, which may be created by passing the --algo ed25519 flag. A keyring can of course hold both types of keys simultaneously, and the Cosmos SDK's x/auth module supports natively these two public key algorithms.

### Query your local keys:

`sided keys list`

## Initialize the SIDE Node

Before actually running the node, we need to initialize the blockchain, and most importantly its genesis file. This is done with the init subcommand:
Use `sided` to initialize your node:

`sided init <MY_SIDE_VALIDATOR> --chain-id=side-testnet-1 --home ~/.sided`

The command above creates all the configuration files needed for your node to run, as well as a default genesis file, which defines the initial state of the network. All these configuration files are in ~/.sided by default, but you can overwrite the location of this folder by passing the --home flag.

The ~/.sided folder has the following structure:

```
.                                   # ~/.sided
  |- data                           # Contains the databases used by the node.
  |- config/
      |- app.toml                   # Application-related configuration file.
      |- config.toml                # CometBFT-related configuration file.
      |- genesis.json               # The genesis file.
      |- node_key.json              # Private key to use for node authentication in the p2p protocol.
      |- priv_validator_key.json    # Private key to use as a validator in the consensus protocol.
```

## Configuring the Node Using app.toml and config.toml
The Cosmos SDK automatically generates two configuration files inside `~/.sided/config`:
- `config.toml`: used to configure the CometBFT, learn more on CometBFT's documentation,
- `app.toml`: generated by the Cosmos SDK, and used to configure your app, such as state pruning strategies, telemetry, gRPC and REST servers configuration, state sync...

Both files are heavily commented, please refer to them directly to tweak your node.

### Downloading the genesis file

Download the genesis file and replace your local one.

`wget https://raw.githubusercontent.com/sideprotocol/testnet/main/shambhala/genesis.json -O $HOME/.sided/config/genesis.json`

### Adding seeds and persistent peers
Open the config.toml to edit the seeds and persistent peers:

```
cd $HOME/.sided/config
nano config.toml
```

Modify the seed node or persistence node configuration by referring to the seeds and persistence details provided in the  file

```
# seeds = "2eba9c8e6fb9d56bbdd10d007a598541c37f6493@13.212.61.41:26656"
seeds = ""
```
or
```
# persistent_peers = "2eba9c8e6fb9d56bbdd10d007a598541c37f6493@13.212.61.41:26656"
persistent_peers = ""
```

## Setting up minimum gas prices

Open app.toml, find minimum-gas-prices, which defines the minimum gas prices the validator node is willing to accept for processing a transaction. Make sure to edit the field with some value, for example 0.005uside,

```
 # The minimum gas prices a validator is willing to accept for processing a
 # transaction. A transaction's fees must meet the minimum of any denomination
 # specified in this config (e.g. 0.25token1;0.0001token2).
 minimum-gas-prices = "0.005uside"
```

## Start Node and Sync

```
sided tendermint unsafe-reset-all 
sided start
```