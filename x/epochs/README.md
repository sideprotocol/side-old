# `epochs`

## Abstract

This document specifies the internal `x/epochs` module of the Sidechain.

Often, when working with the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk),
we would like to run certain pieces of code every so often.

The purpose of the `epochs` module is to allow other modules to maintain
that they would like to be signaled once in a time period.
So, another module can specify it wants to execute certain code once a week, starting at UTC-time = x.
`epochs` creates a generalized epoch interface to other modules so they can be more easily signaled upon such events.

## Contents

1. **[Concept](#concepts)**
2. **[State](#state)**
3. **[Events](#events)**
4. **[Keeper](#keepers)**
5. **[Hooks](#hooks)**
6. **[Queries](#queries)**
7. **[Future improvements](#future-improvements)**

## Concepts

The `epochs` module defines on-chain timers that execute at fixed time intervals.
Other sidechain modules can then register logic to be executed at the timer ticks.
We refer to the period in between two timer ticks as an "epoch".

Every timer has a unique identifier, and every epoch will have a start time and an end time,
where `end time = start time + timer interval`.


## State

### State Objects

The `x/epochs` module keeps the following `objects in state`:

| State Object | Description         | Key                  | Value               | Store |
|--------------|---------------------|----------------------|---------------------|-------|
| `EpochInfo`  | Epoch info bytecode | `[]byte{identifier}` | `[]byte{epochInfo}` | KV    |

#### EpochInfo

An `EpochInfo` defines several variables:

1. `identifier` keeps an epoch identification string
2. `start_time` keeps the start time for epoch counting:
   if block height passes `start_time`, then `epoch_counting_started` is set
3. `duration` keeps the target epoch duration
4. `current_epoch` keeps the current active epoch number
5. `current_epoch_start_time` keeps the start time of the current epoch
6. `epoch_counting_started` is a flag set with `start_time`, at which point `epoch_number` will be counted
7. `current_epoch_start_height` keeps the start block height of the current epoch

```protobuf
message EpochInfo {
    string identifier = 1;
    google.protobuf.Timestamp start_time = 2 [
        (gogoproto.stdtime) = true,
        (gogoproto.nullable) = false,
        (gogoproto.moretags) = "yaml:\"start_time\""
    ];
    google.protobuf.Duration duration = 3 [
        (gogoproto.nullable) = false,
        (gogoproto.stdduration) = true,
        (gogoproto.jsontag) = "duration,omitempty",
        (gogoproto.moretags) = "yaml:\"duration\""
    ];
    int64 current_epoch = 4;
    google.protobuf.Timestamp current_epoch_start_time = 5 [
        (gogoproto.stdtime) = true,
        (gogoproto.nullable) = false,
        (gogoproto.moretags) = "yaml:\"current_epoch_start_time\""
    ];
    bool epoch_counting_started = 6;
    reserved 7;
    int64 current_epoch_start_height = 8;
}
```

The `epochs` module keeps these `EpochInfo` objects in state, which are initialized at genesis
and are modified on begin blockers or end blockers.

#### Genesis State

The `x/epochs` module's `GenesisState` defines the state necessary for initializing the chain
from a previously exported height.
It contains a slice containing all the `EpochInfo` objects kept in state:

```go
// Genesis State defines the epoch module's genesis state
type GenesisState struct {
    // list of EpochInfo structs corresponding to all epochs
	Epochs []EpochInfo `protobuf:"bytes,1,rep,name=epochs,proto3" json:"epochs"`
}
```

## Events

The `x/epochs` module emits the following events:

### BeginBlocker

| Type          | Attribute Key     | Attribute Value   |
| ------------- | ----------------- | ----------------- |
| `epoch_start` | `"epoch_number"`  | `{epoch_number}`  |
| `epoch_start` | `"start_time"`    | `{start_time}`    |

### EndBlocker

| Type           | Attribute Key    | Attribute Value   |
| ------------- | ----------------- | ----------------- |
| `epoch_end`   | `"epoch_number"`  | `{epoch_number}`  |

## Keepers

The `x/epochs` module only exposes one keeper, the epochs keeper, which can be used to manage epochs.

### Epochs Keeper

Presently only one fully-permissioned epochs keeper is exposed,
which has the ability to both read and write the `EpochInfo` for all epochs,
and to iterate over all stored epochs.

```go
// Keeper of epoch nodule maintains collections of epochs and hooks.
type Keeper struct {
	cdc      codec.Codec
	storeKey storetypes.StoreKey
	hooks    types.EpochHooks
}
```

```go
// Keeper is the interface for epoch module keeper
type Keeper interface {
  // GetEpochInfo returns epoch info by identifier
  GetEpochInfo(ctx sdk.Context, identifier string) types.EpochInfo

  // SetEpochInfo set epoch info
  SetEpochInfo(ctx sdk.Context, epoch types.EpochInfo)

  // DeleteEpochInfo delete epoch info
  DeleteEpochInfo(ctx sdk.Context, identifier string)

  // IterateEpochInfo iterate through epochs
  IterateEpochInfo(ctx sdk.Context, fn func(index int64, epochInfo types.EpochInfo) (stop bool))

  // Get all epoch infos
  AllEpochInfos(ctx sdk.Context) []types.EpochInfo
}
```

## Hooks

The `x/epochs` module implements hooks so that other modules can use epochs
to allow facets of the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) to run on specific schedules.

### Hooks Implementation

```go
// combine multiple epoch hooks, all hook functions are run in array sequence
type MultiEpochHooks []types.EpochHooks

// AfterEpochEnd is called when epoch is going to be ended, epochNumber is the
// number of epoch that is ending
func (mh MultiEpochHooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {...}

// BeforeEpochStart is called when epoch is going to be started, epochNumber is
// the number of epoch that is starting
func (mh MultiEpochHooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {...}

// AfterEpochEnd executes the indicated hook after epochs ends
func (k Keeper) AfterEpochEnd(ctx sdk.Context, identifier string, epochNumber int64) {...}

// BeforeEpochStart executes the indicated hook before the epochs
func (k Keeper) BeforeEpochStart(ctx sdk.Context, identifier string, epochNumber int64) {...}
```

### Recieving Hooks

When other modules (outside of `x/epochs`) recieve hooks,
they need to filter the value `epochIdentifier`, and only do executions for a specific `epochIdentifier`.

The filtered values from `epochIdentifier` could be stored in the `Params` of other modules,
so they can be modified by governance.

Governance can change epoch periods from `week` to `day` as needed.


## Queries

The `x/epochs` module provides the following queries to check the module's state.

```protobuf
service Query {
  // EpochInfos provide running epochInfos
  rpc EpochInfos(QueryEpochsInfoRequest) returns (QueryEpochsInfoResponse) {}
  // CurrentEpoch provide current epoch of specified identifier
  rpc CurrentEpoch(QueryCurrentEpochRequest) returns (QueryCurrentEpochResponse) {}
}
```


## Future Improvements

### Correct Usage

In the current design, each epoch should be at least two blocks, as the start block should be different from the endblock.
Because of this, the time allocated to each epoch will be `max(block_time x 2, epoch_duration)`.
For example: if the `epoch_duration` is set to `1s`, and `block_time` is `5s`, actual epoch time should be `10s`.

It is recommended to configure `epoch_duration` to be more than two times the `block_time`, to use this module correctly.
If there is a mismatch between the `epoch_duration` and the actual epoch time, as in the example above,
then module logic could become invalid.

### Block-Time Drifts

This implementation of the `x/epochs` module has block-time drifts based on the value of `block_time`.
For example: if we have an epoch of 100 units that ends at `t=100`,
and we have a block at `t=97` and a block at `t=104` and `t=110`, this epoch ends at `t=104`,
and the new epoch will start at `t=110`.

There are time drifts here, varying about 1-2 blocks time, which will slow down epochs.