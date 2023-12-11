# Deposit flow
User Deposits Stakable Assets to Module

## Description:

Users can deposit stakable assets to the module, initiating various actions on the Host Chain.

## Implementation Steps:

- Users initiate a deposit transaction to the module.
- The module creates an (Interchain account)[https://tutorials.cosmos.network/academy/3-ibc/8-ica.html] on the Host Chain as a prerequisite for further interactions which are as follows:

    - `MsgDelegate`
        ```go
        type MsgDelegate struct {
            DelegatorAddress string      `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
            ValidatorAddress string      `protobuf:"bytes,2,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
            Amount           types1.Coin `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount"`
        }
        ```
    - `MsgUndelegate`
        ```go
        type MsgUndelegate struct {
            DelegatorAddress string      `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
            ValidatorAddress string      `protobuf:"bytes,2,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
            Amount           types1.Coin `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount"`
        }
        ```
    - `MsgWithdrawDelegatorReward`
        ```go
        type MsgWithdrawDelegatorReward struct {
            DelegatorAddress string `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
            ValidatorAddress string `protobuf:"bytes,2,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
        }
        ```
    - `MsgSetWithdrawAddress`
        ```go
        type MsgSetWithdrawAddress struct {
            DelegatorAddress string `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
            WithdrawAddress  string `protobuf:"bytes,2,opt,name=withdraw_address,json=withdrawAddress,proto3" json:"withdraw_address,omitempty"`
        }
        ```
    - `Parameters`:
        DelegatorAddress: The address of the user delegating assets.
        ValidatorAddress: The address of the validator receiving the delegation.
        Amount: The amount of stakable assets being delegated.
        WithdrawAddress: The new withdrawal address for delegated rewards.

- After a certain time-period, a epoch event is trigerred,
    1. Controller chain sends funds to Host Chain Side-ICA for staking.
        ```go
        MsgTransfer{
            SourcePort:       sourcePort,
            SourceChannel:    sourceChannel,
            Token:            stakable_token,
            Sender:           module_account_address,
            Receiver:         side-ica,
            TimeoutHeight:    timeoutHeight,
            TimeoutTimestamp: timeoutTimestamp,
            Memo:             memo,
        }
        ```
    2. Controller chain send delegate message to Host Chain SIDE-ICA.

- `Exchange Rate Calculation`: The exchange rate (ex_rt) is calculated based on the total token staked through SIDE-ICA and the total supply of LSD tokens issued by the SIDE-Chain.
    - Formula
    ```go
        ex_rt = total_token_staked / total_supply_LSD_token
    ```
    - Parameters
        - `total_token_staked`: The total value staked through the SIDE-ICA account.
        - `total_supply_LSD_token`: The total supply of LSD tokens issued by the SIDE-Chain to users in exchange for stakable assets.

- The SIDE-Chain issues LSD tokens to users based on the calculated exchange rate.
    - Formula
    ```go
        issued_LSD_tokens = deposited_tokens / ex_rt
    ```
    - Parameters:
        - `deposited_tokens`: The number of tokens deposited by users.
        - `ex_rt`: The calculated exchange rate.
    - Example:
    If the exchange rate is 1.2 and a user deposits 100 tokens, they will receive 100 / 1.2 LSD tokens.

## Data Structure

- `User Deposits`: These deposits are stored with each epoch number,
    ```go
    DepositStore {
        Epoch              uint64
        Amount             types1.Coin 
        HostChain          string
        Status             Deposit_Status
    }
    ```
    - After each epoch end, module transfers Amount to HostChain and update status
    ```go
    Status {
        DEPOSIT_PENDING,
        DEPOSIT_SENT,
        DEPOSIT_RECEIVED
        DEPOSIT_STAKED,
    }
    ```
- `Validator Set`: Validator parameters are stored as array of:
    ```go
    Set {
        Validator_address string
        Weight uint64
        Status VALIDATOR_STATUS
        Amount_staked uint64
    }
    ```
