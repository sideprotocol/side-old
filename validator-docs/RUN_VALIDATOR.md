# Run a Validator

## What is a Validator?

Validators are responsible for committing new blocks to the blockchain through an automated voting process. A validator's stake is slashed if they become unavailable or sign blocks at the same height. Because there is a chance of slashing, we suggest you read about Sentry Node Architecture to protect your node from DDOS attacks and to ensure high-availability.

The following instructions assume you have already set up a full-node and are synchronized to the latest block height.

## Creating Your Validator

### Get some SIDE tokens from faucet
​​
- https://faucet.side.one

### Check your balances

```
$ sided query bank balances side162ttnnflhvuazp42vxu5jf3fgxy750lrn3zrld
balances:
- amount: "999999900000000000"
  denom: uside
  pagination:
  next_key: null
  total: "0"
```

### Create Your Validator

```
sided tx staking create-validator \
--from=alice \
--amount=400000000uside \
--pubkey=$(sided tendermint show-validator)  \
--moniker="side_node" \
--security-contact="contact@side.one" \
--chain-id="side-testnet-1" \
--commission-rate="0.1" \
--commission-max-rate="0.2" \
--commission-max-change-rate="0.05" \
--min-self-delegation="400000000" \
```

If you need further explanation for each of these command flags:
- the `from` flag is the KEY_NAME you created when initializing the key on your keyring
- the `amount` flag is the amount you will place in your own validator in uside
- the `pubkey` is the validator public key found earlier
- the `moniker` is a human-readable name you choose for your validator
- the `security-contact` is an email your delegates are able to contact you at
- the `chain-id` is whatever chain-id you are working with
- the `commission-rate` is the rate you will charge your delegates
- the `commission-max-rate` is the most you are allowed to charge your delegates
- the `commission-max-change-rate` is how much you can increase your commission rate in a 24 hour period
- the `min-self-delegation` is the lowest amount of personal funds the validator is required to have in their own validator to stay bonded
- the `gas-prices` is the amount of gas used to send this create-validator transaction

## Editing Validator Description

You can edit your validator's public description. This info is to identify your validator, and will be relied on by delegators to decide which validators to stake to. Make sure to provide input for every flag below. If a flag is not included in the command the field will default to empty (`--moniker` defaults to the machine name) if the field has never been set or remain the same if it has been set in the past.

The `<key_name>` specifies which validator you are editing. If you choose to not include some flags below, remember that the `--from` flag must be included to identify the validator to update.

The `--identity` can be used as to verify identity with systems like Keybase or UPort. When using Keybase, `--identity` should be populated with a 16-digit string that is generated with a keybase.io (opens new window)account. It's a cryptographically secure method of verifying your identity across multiple online networks. The Keybase API allows us to retrieve your Keybase avatar. This is how you can add a logo to your validator profile.

```
sided tx staking edit-validator
  --moniker="choose a moniker" \
  --website="https://www.side.one" \
  --identity=6A0D65E29A4CBC8E \
  --details="To infinity and beyond!" \
  --chain-id="side-testnet-1" \
  --gas="auto" \
  --gas-prices="0.0025uside" \
  --from=<key_name> \
  --commission-rate="0.10"
```

`Note`: The `commission-rate` value must adhere to the following rules:
- Must be between 0 and the validator's `commission-max-rate`
- Must not exceed the validator's `commission-max-change-rate` which is maximum % point change rate per day. In other words, a validator can only change its commission once per day and within `commission-max-change-rate` bounds.

## Viewing Validator Description

View the validator's information with this command:
`sided query slashing signing-info <validator-pubkey> --chain-id="side-testnet-1"`

## Unjail Validator

When a validator is "jailed" for downtime, you must submit an Unjail transaction from the operator account in order to be able to get block proposer rewards again (depends on the zone fee distribution).

```
sided tx slashing unjail \
 --from=<key_name> \
 --chain-id=<chain_id>
```

## Confirm Your Validator is Running

Your validator is active if the following command returns anything:
`sided query tendermint-validator-set | grep "$(sided tendermint show-address)"`
You should now see your validator in one of SIDE blockchain explorers. You are looking for the bech32 encoded address in the `~/.sided/config/priv_validator.json` file.

## Halting Your Validator

When attempting to perform routine maintenance or planning for an upcoming coordinated upgrade, it can be useful to have your validator systematically and gracefully halt. You can achieve this by either setting the halt-height to the height at which you want your node to shutdown or by passing the --halt-height flag to sided. The node will shutdown with a zero exit code at that given height after committing the block.

## Advanced configuration
You can find more advanced information about running a node or a validator on the https://docs.cometbft.com/v0.34/core/validators.