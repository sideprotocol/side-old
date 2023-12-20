---
title: "Interchainquery"
excerpt: ""
category: 6392913957c533007128548e
---

# Interchain Query

## Abstract

Side uses interchain queries and interchain accounts to perform multichain liquid staking. The `interchainquery` module creates a framework that allows other modules to query other appchains using IBC. The `interchainquery` module is used to make redemption rate ICQ queries on stride. The callback triggers update of redemption rate on side chain.

## Contents

1. **[Concepts](#concepts)**
2. **[State](#state)**
3. **[Events](#events)**
4. **[Keeper](#keeper)**
5. **[Msgs](#msgs)**  
6. **[Queries](#queries)**

## State

The `interchainquery` module keeps `Query` objects and modifies the information from query to query, as defined in `proto/interchainquery/v1/genesis.proto`

### InterchainQuery information type

`Query` has information types that pertain to the query itself. `Query` keeps the following:

1. `id`: query identification string.
2. `connection_id`: id of the connection between the controller and host chain.
3. `chain_id`: id of the queried chain.
4. `query_type`: type of interchain query (e.g. bank store query, redemption rate)
5. `request_data`: serialized request information (e.g. the address with which to query)
6. `callback_module`: name of the module that will handle the callback
7. `callback_id`: ID for the function that will be called after the response is returned
8. `callback_data`: optional serialized data associated with the callback
9. `timeout_policy`: specifies how to handle a timeout (fail the query, retry the query, or execute the callback with a timeout)
10. `timeout_duration`: the relative time from the current block with which the query should timeout
11. `timeout_timestamp`: the absolute time at which the query times out
12. `request_sent`: boolean indicating whether the query event has been emitted (and can be identified by a relayer)
13. `submission_height`: the light client hight of the queried chain at the time of query submission


`DataPoint` has information types that pertain to the data that is queried. `DataPoint` keeps the following:

1. `id` keeps the identification string of the datapoint
2. `remote_height` keeps the block height of the queried chain
3. `local_height` keeps the block height of the querying chain
4. `value` keeps the bytecode value of the data retrieved by the Query
