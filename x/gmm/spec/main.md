---
title: In-chain Automatic Market Maker.
stage: draft
category: Side/AMM
kind: instantiation
author: Marian <marian@side.one>
created: 2023-11-08
modified: 2023-11-08
---

## Synopsis

This specification delineates the core mathematical models and algorithmic principles for an in-chain automated market maker (AMM) module. The AMM is designed to facilitate a decentralized, permissionless exchange mechanism within a single blockchain ecosystem. The module supports a variety of liquidity pool configurations, including weighted and stable pools, harnessing the robustness of AMM principles to provide a resilient and flexible trading platform.

### Motivation

The Side AMM enables permissionless token exchanges, allowing anyone to create token pair pools, contribute as liquidity providers, and earn profits from swap fees and related products. It aims to support different pool styles and offers features like single-asset deposits for user convenience and arbitrage opportunities.

### Definitions

`Automated market makers (AMM)`: Decentralized protocols within a single blockchain that facilitate token trading by using liquidity pools and predetermined formulas to calculate prices.

`Weighted pool`: Liquidity pool where each token's denomination has a specific weight, influencing the [pool's pricing mechanism](https://docs.balancer.fi/reference/math/weighted-math.html).

`Stable pool`: Liquidity pool designed for assets that are expected to maintain a stable price relative to each other, typically stablecoins based on [stable math](https://docs.balancer.fi/reference/math/stable-math.html).

`Single-asset deposit`: a deposit into a liquidity pool that does not require users to deposit both token denominations -- one is enough. While this deposit method will impact the exchange price of the pool, it also provides an opportunity for arbitrage.

`Multi-asset deposit`: a deposit into a liquidity pool that require users to deposit both token denominations. This deposit method will not impact the exchange price of the pool.

`Left-side swap`: a token exchange that specifies the desired quantity to be sold.

`Right-side swap`: a token exchange that specifies the desired quantity to be purchased.

`Pool state`: the entire state of a liquidity pool including its invariant value which is derived from its token balances and weights inside.

### Desired Properties

- `Permissionless`: no need to whitelist connections, modules, or denominations. Individual implementations may have their own permissioning scheme, however the protocol must not require permissioning from a trusted party to be secure.
- `Decentralized`: all parameters are managed on chain via governance. Does not require any central authority or entity to function. Also does not require a single blockchain, acting as a hub, to function.
- `Guarantee of Exchange`: no occurrence of a user receiving tokens without the equivalent promised exchange.
- `Liquidity Incentives`: supports the collection of fees which are distributed to liquidity providers and acts as incentive for liquidity participation.

## Technical Specification

### Liquidity Provision and Pool Creation

Liquidity providers (LPs) are able to create pools within the AMM module to earn passive income through a permissionless swap service. Upon creating a pool, LPs have the option to select the type of pool they wish to establish. The Side/AMM module currently supports two primary pool types:

`Stable Pools`: Optimized for assets that maintain a stable value relative to one another, commonly used for stablecoins.
`Weighted Pools`: Allow for variable asset ratios within the pool, with prices determined by the relative weights assigned to each asset.

In addition to selecting the pool type, LPs can set the fee structure for the pool, determining the swap fees that will be incurred by users when executing trades.

### Joining Existing Pools

New liquidity providers can join an existing pool at any time. By contributing liquidity to the pool, they receive liquidity tokens which represent their share of the pool. The value of these tokens can increase over time as swap fees accumulate within the pool.

### Swap Fee Accumulation and LP Token Valuation

Swap transactions within the pool accrue fees, which are added to the pool's total value. As a result, the price of the liquidity tokens associated with the pool increases correspondingly. The value of LP tokens is calculated using a precise mathematical model, which may vary depending on whether the pool operates on weighted math or stable math principles.

### Swap Transactions and Pool Elimination

Users holding tokens that they wish to exchange can initiate a swap transaction (referred to as a left swap) to obtain their desired target token. After deducting the swap fee, the remaining tokens are immediately transferred to the user once the transaction is confirmed.

Should a situation arise where all liquidity tokens are returned to the module's account, signifying the withdrawal of all liquidity providers, the pool is automatically terminated and removed from the system.

Here is simple diagram.

```
+----------------------+-------------------+---------------------+-------------------+
| Liquidity Provider   | AMM Module        | Liquidity Pool      | User (Trader)     |
+----------------------+-------------------+---------------------+-------------------+
| Create pool          |                   |                     |                   |
| (select type & fee)  |                   |                     |                   |
| -------------------->| Initialize pool   |                     |                   |
|                      |------------------>|                     |                   |
| Provide initial      |                   |                     |                   |
| liquidity            |                   |                     |                   |
| -------------------->|                   | Issue liquidity    |                   |
|                      |                   | tokens              |                   |
|                      |                   |<--------------------|                   |
|                      |                   |                     |                   |
|                      |                   | Other LPs can join  |                   |
|                      |                   | the pool anytime    |                   |
|                      |                   |                     |                   |
|                      |                   |                     | Request swap      |
|                      |                   |                     | (specify tokens)  |
|                      |                   |                     |------------------>|
|                      | Calculate swap    |                     |                   |
|                      | & update state    |                     |                   |
|                      |------------------>|                     |                   |
|                      |                   | Transfer target     |                   |
|                      |                   | tokens (minus fee)  |                   |
|                      |                   |-------------------->|                   |
|                      |                   |                     |                   |
|                      |                   | Swap fees accrue    |                   |
|                      |                   | in pool             |                   |
|                      |                   |                     |                   |
| Request withdrawal   |                   |                     |                   |
| -------------------->|                   | Release             |                   |
|                      |                   | proportional assets |                   |
|                      |                   |<--------------------|                   |
|                      |                   |                     |                   |
|                      | If all liquidity  |                     |                   |
|                      | is withdrawn      |                     |                   |
|                      |------------------>| Eliminate pool      |                   |
|                      |                   |                     |                   |
+----------------------+-------------------+---------------------+-------------------+

```

## History

Oct 9, 2022 - Draft written

Oct 11, 2022 - Draft revised

## References

https://dev.balancer.fi/resources/pool-math

## Copyright

All content herein is licensed under [Apache 2.0](https://www.apache.org/licenses/LICENSE-2.0).
