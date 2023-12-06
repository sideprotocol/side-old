# Claim and Restake flow
This step is triggered by epoch module.

## Registration of ICA Accounts and Rewards Management
Initial Registration of ICA Accounts

- Description:
    - Initially, two ICA (Interchain Account) accounts are registered, and the withdraw rewards address is set to rewards-ica.
- Process:
    - Two ICA accounts are registered on the Host Chain.
    - The withdraw rewards address for both accounts is set to rewards-ica.

## Rewards Account for Accurate Calculation

- Description:
    - A specific rewards account (rewards-ica) is designated to accurately calculate rewards accrued from staked balances.
- Purpose:
    - The rewards account is crucial for precise tracking and calculation of rewards generated from staked balances.

## Setting Withdraw Address for Host chain

```go
type MsgSetWithdrawAddress struct {
    DelegatorAddress string `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
    WithdrawAddress  string `protobuf:"bytes,2,opt,name=withdraw_address,json=withdrawAddress,proto3" json:"withdraw_address,omitempty"`
}
```
- Parameters:
    - `DelegatorAddress`: The address of the delegator setting the withdrawal address.
    - `WithdrawAddress`: The designated address for receiving rewards.

## Claiming and Depositing Rewards

### Claim Rewards Trigger:
- After a specific time period, a trigger initiates the process of claiming rewards.

### Claiming Rewards from Validators:
- All rewards from various validators are claimed during this step.
```go
type MsgWithdrawDelegatorReward struct {
    DelegatorAddress string `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
    ValidatorAddress string `protobuf:"bytes,2,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
}
```
- Parameters:
    - `DelegatorAddress`: The address of the delegator setting the withdrawal address.(side-ICA)
    - `WithdrawAddress`: The designated address for receiving rewards.(rewards-ICA)

### Deposit Request Preparation:
- A deposit request is prepared to stake the claimed rewards back to validators.

### Deposit to Validators:
- The deposit request is executed, and the claimed rewards are deposited back to the respective validators.

### Effect on Exchange Rate:
- The exchange rate gradually increases as a result of depositing the rewards back to validators.
