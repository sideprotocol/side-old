# Withdrawal Flow

## User Initiated Withdrawal

- Description:
    - Users initiate withdrawal of tokens with a withdrawal request to the SIDE Chain.
- Process:
    - Users sends a withdrawal request along with LSD tokens to the SIDE Chain.
    - The SIDE Chain stores the unbonding request in a data structure with the status set to PENDING.
    ```go
    type Unbonding struct {
        // Epoch number
        Epoch int64
        // Host chain name
        HostChain string
        // Time after which token is unbonded
        UnbondingTime time.Time
        // Amount of LSD tokens to burn
        BurnAmount types.Coin
        // Amount of tokens that is being unbonded
        UnbondAmount types.Coin
        // status of unbonding
        Status Unbonding_Status
    }
    ```

    ```go
    Unbonding_Status {
        PENDING,
        SENT,
        UNBONDING,
        UNBONDED
    }
    ```
    - User claimable storage
    ```go
    type UserUnbonding struct {
        // Epoch number
        Epoch int64
        // Host chain name
        HostChain string
        // Amount of LSD tokens to burn
        UnstakeAmount types.Coin
        // Amount of tokens that is being unbonded
        UnbondAmount types.Coin
    }
    ```

## Batched Unbonding Requests

- Description:
    - Due to a limit on simultaneous unbonding requests per account, unbonding requests are batched to ensure efficient processing.
- Limitation:
    - A maximum of 7 simultaneous unbonding requests is allowed for a single account.

- Process:
    - Unbonding requests are accumulated until a batch is formed.
    - Batch time is calculated as (unbonding period) / 7 to comply with the simultaneous unbonding request limit.

## Batch Unbonding Request to Host Chain

- Description:
    - After a certain epoch, the batched unbonding request is sent to the Host Chain, and the status is set to PENDING.
- Process:
    - Once the batch is formed, it is sent to the Host Chain for processing.
    - The status of unbonding requests is updated to PENDING.
    - Message sent to host chain, status is set to SENT.
    ```go
    type MsgUndelegate struct {
        DelegatorAddress string      `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
        ValidatorAddress string      `protobuf:"bytes,2,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
        Amount           types1.Coin `protobuf:"bytes,3,opt,name=amount,proto3" json:"amount"`
    }
    ```
    - Once ack is received, status is set to UNBONDING.

## Claim Unbonded Tokens

- Description:
    - After the unbonding period expires, an epoch event triggers the action to claim unbonded tokens and send to users.

- Process:
    - After the unbonding period lapses, an epoch event is triggered.
    - Transfer unstaked tokens from Host to SIDE chain, and the status of unbonding requests is set to CLAIMABLE.
    ```go
        MsgTransfer{
            SourcePort:       sourcePort,
            SourceChannel:    sourceChannel,
            Token:            stakable_token,
            Sender:           side-ica,
            Receiver:         module_account_address,
            TimeoutHeight:    timeoutHeight,
            TimeoutTimestamp: timeoutTimestamp,
            Memo:             memo,
        }
        ```
    - Users receive their unstaked tokens.(can be triggered via hooks or begin/end block)