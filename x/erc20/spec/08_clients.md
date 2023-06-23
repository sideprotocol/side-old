<!--
order: 8
-->

# Clients

## CLI

Find below a list of  `sidechaind` commands added with the `x/erc20` module. You can obtain the full list by using the `sidechaind -h` command. A CLI command can look like this:

```bash
sidechaind query erc20 params
```

### Queries

| Command         | Subcommand    | Description                    |
| --------------- | ------------- | ------------------------------ |
| `query` `erc20` | `params`      | Get erc20 params               |
| `query` `erc20` | `token-pair`  | Get registered token pair      |
| `query` `erc20` | `token-pairs` | Get all registered token pairs |

### Transactions

| Command      | Subcommand      | Description                    |
| ------------ | --------------- | ------------------------------ |
| `tx` `erc20` | `convert-coin`  | Convert a Cosmos Coin to ERC20 |
| `tx` `erc20` | `convert-erc20` | Convert a ERC20 to Cosmos Coin |

### Proposals

The `tx gov submit-legacy-proposal` commands allow users to query create a proposal using the governance module CLI:

**`register-coin`**

Allows users to submit a `RegisterCoinProposal`. Submit a proposal to register a Cosmos coin to the erc20 along with an initial deposit. Upon passing, the proposal details must be supplied via a JSON file.

```bash
sidechaind tx gov submit-legacy-proposal register-coin METADATA_FILE [flags]
```

Where METADATA_FILE contains (example):

```json
{
  "metadata": [
    {
      "description": "The native staking and governance token of the Osmosis chain",
      "denom_units": [
        {
          "denom": "ibc/<HASH>",
          "exponent": 0,
          "aliases": ["ibcuosmo"]
        },
        {
          "denom": "OSMO",
          "exponent": 6
        }
      ],
      "base": "ibc/<HASH>",
      "display": "OSMO",
      "name": "Osmo",
      "symbol": "OSMO"
    }
  ]
}
```

**`register-erc20`**

Allows users to submit a `RegisterERC20Proposal`. Submit a proposal to register ERC20 tokens along with an initial deposit. To register multiple tokens in one proposal pass them after each other e.g. `register-erc20 <contract-address1> <contract-address2>`.

```bash
sidechaind tx gov submit-legacy-proposal register-erc20 ERC20_ADDRESS... [flags]
```

**`toggle-token-conversion`**

Allows users to submit a `ToggleTokenConversionProposal`.

```bash
sidechaind tx gov submit-legacy-proposal toggle-token-conversion TOKEN [flags]
```

**`param-change`**

Allows users to submit a `ParameterChangeProposal``.

```bash
sidechaind tx gov submit-legacy-proposal param-change PROPOSAL_FILE [flags]
```

## gRPC

### Queries

| Verb   | Method                                                    | Description                    |
| ------ | --------------------------------------------------------- | ------------------------------ |
| `gRPC` | `sidechain.erc20.v1.Query/Params`                         | Get erc20 params               |
| `gRPC` | `sidechain.erc20.v1.Query/TokenPair`                      | Get registered token pair      |
| `gRPC` | `sidechain.erc20.v1.Query/TokenPairs`                     | Get all registered token pairs |
| `GET`  | `/github.com/sideprotocol/sidechain/erc20/v1/params`      | Get erc20 params               |
| `GET`  | `/github.com/sideprotocol/sidechain/erc20/v1/token_pair`  | Get registered token pair      |
| `GET`  | `/github.com/sideprotocol/sidechain/erc20/v1/token_pairs` | Get all registered token pairs |

### Transactions

| Verb   | Method                                                         | Description                    |
| ------ | -------------------------------------------------------------- | ------------------------------ |
| `gRPC` | `sidechain.erc20.v1.Msg/ConvertCoin`                           | Convert a Cosmos Coin to ERC20 |
| `gRPC` | `sidechain.erc20.v1.Msg/ConvertERC20`                          | Convert a ERC20 to Cosmos Coin |
| `GET`  | `/github.com/sideprotocol/sidechain/erc20/v1/tx/convert_coin`  | Convert a Cosmos Coin to ERC20 |
| `GET`  | `/github.com/sideprotocol/sidechain/erc20/v1/tx/convert_erc20` | Convert a ERC20 to Cosmos Coin |
