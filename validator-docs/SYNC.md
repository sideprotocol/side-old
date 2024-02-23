# Sync with state-sync

state-sync is a feature that allows nodes to quickly sync their state by fetching a snapshot of the application state at a specific block height. This greatly reduces the time required for node to sync with the network, compared to the default method of replaying all blocks from the genesis block.
An advantage of state-sync is that the database is very small in comparison to a fully synced node, therefore using state-sync to resync your node to the network can help keep running costs lower by minimizing storage usage.

- Note: A snapshot-enabled RPC and from a trusted block height is required for state-sync.

1. Set SNAP_RPC variable to the snapshot RPC
    - `SNAP_RPC="https://devnet-rpc.side.one:443"`

2. Fetch the `LATEST_HEIGHT` from the snapshot RPC, set the state-sync `BLOCK_HEIGHT` and fetch the `TRUST_HASH` from the snapshot RPC. The `BLOCK_HEIGHT` to sync is determined by subtracting the snapshot-interval from the `LATEST_HEIGHT`.
    ```
    LATEST_HEIGHT=$(curl -s $SNAP_RPC/block | jq -r .result.block.header.height); \
    BLOCK_HEIGHT=$((LATEST_HEIGHT - 2000)); \
    TRUST_HASH=$(curl -s "$SNAP_RPC/block?height=$BLOCK_HEIGHT" | jq -r .result.block_id.hash)
    ```

3. Check variables to ensure they have been set
    `echo $LATEST_HEIGHT $BLOCK_HEIGHT $TRUST_HASH`

4. Set the required variables in ~/.side/config/config.toml
    ```
    sed -i.bak -E "s|^(enable[[:space:]]+=[[:space:]]+).*$|\1true| ; \
    s|^(rpc_servers[[:space:]]+=[[:space:]]+).*$|\1\"$SNAP_RPC,$SNAP_RPC\"| ; \
    s|^(trust_height[[:space:]]+=[[:space:]]+).*$|\1$BLOCK_HEIGHT| ; \
    s|^(trust_hash[[:space:]]+=[[:space:]]+).*$|\1\"$TRUST_HASH\"|" $HOME/.side/config/config.toml
    ```

5. Stop the node and reset the node database
    ```
    WARNING: This will erase your node database. If you are already running validator, be sure you backed up your config/priv_validator_key.json and config/node_key.json prior to running unsafe-reset-all.
    It is recommended to copy data/priv_validator_state.json to a backup and restore it after unsafe-reset-all to avoid potential double signing.
    ```

6. Restart node and check logs
    ```
    sided tendermint unsafe-reset-all
    sided start
    ```