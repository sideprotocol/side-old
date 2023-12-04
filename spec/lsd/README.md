# Liquid Staking Derivative Specification

## Introduction

The Liquid Staking Derivative aims to provide users with the ability to stake assets while maintaining liquidity through the issuance of derivative tokens. This specification outlines the technical details and features implemented using the Cosmos SDK.

## Architecture

The architecture comprises ICA, ICQ modules, minting module, allowing users to stake assets and receive derivative tokens representing their staked holdings.

![Alt text](image.png)

## Core features

### Liquid Staking
- Users can stake assets to receive derivative tokens.
- Minting and burning of derivative tokens based on staking and unstaking actions.
- Users can then use derivatives tokens in DeFi

### Staking Rewards
- Calculation and distribution of staking rewards.
- Users rewards are accumulated in derivative tokens i.e value of derivative token increases over time.

## Techincal Details

The technical implementation is divided into 4 flows which are as follows:

1. Deposit and mintLSD flow:
    - Users initiate the deposit process by sending a specified amount of stakable assets to the liquid staking module.
    - The module verifies the deposit and mints an equivalent amount of Liquid Staking Derivative (LSD) tokens, representing the staked assets.
    - More details can be found [here](./DEPOSIT.md)

2. Claim and Restake flow:
    - Module claims the rewards, stake them and increase the exchange rate of LSD.
    - More details can be found [here](./CLAIMANDRESTAKE.md)

3. Unbond requests and claim unbonded amount:
    - Users can initiate unbonding requests to withdraw their staked assets, subject to a predefined unbonding period.
    - After unbonding period is over, users can claim their tokens
    - More details can be found [here](./WITHDRAW.md)

4. Epoch module:
    - Introduction of an epoch module to manage time-based operations and ensure the consistency of staking rewards and unbonding periods.
    - Triggers epoch-based events, such as reward claim, unbonding requests and claimable requests.

5. LSM:
    - This allows users who have already staked their assets and don't want to unstake to move to LSD, they can do so directly.