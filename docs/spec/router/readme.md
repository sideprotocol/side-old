---
x/router:
title: Router
stage: draft
category: SIDEHUB/X
kind: routing
author: Marian <marian@side.one>
created: 2022-05-23
modified: 2023-07-24
---

## Synopsis

In this system, asynchronous acknowledgements play a crucial role in facilitating atomic multi-hop packet flows. Only after the completion of the forward/multi-hop sequence, whether it results in success or failure, is the acknowledgement recorded on the chain where the user initiated the packet flow. This asynchronous approach implies that an IBC application user needs to only monitor the chain where the initial transfer was dispatched for tracking the response of the entire transaction process. This setup effectively simplifies user interaction and reduces monitoring complexity, thus enhancing user experience in a technically advanced and professionally designed system.

## Motivation

The route mechanism has been conceptualized to address specific challenges within decentralized finance (DeFi) and interoperability. Its design incorporates two main principles:

- Flexibility of Token Exchange: The mechanism allows for a versatile exchange of any tokens through the SideHub. This design ensures that users can freely convert their assets between different tokens without worrying about the compatibility or transferability across multiple blockchain networks.

- Optimized Liquidity Utilization: One of the significant constraints in token exchange is the availability of liquidity. The route mechanism aims to circumvent this problem by effectively using SideHub as an intermediary. It leverages two types of liquidity pools in SideHub, offering users an alternative path for conversion in situations where direct liquidity between $B and $A might be limited or non-existent.

By creating a flexible route for token conversion, users can make more efficient and effective trades, contributing to the overall liquidity and robustness of the decentralized finance ecosystem. The mechanism also helps alleviate the 'slippage' problem often encountered in liquidity pools, wherein large trades significantly shift the price of a token.

This route mechanism, with its capability of nesting IBC calls within a single user transaction, makes it a unique solution within the Cosmos SDK, offering optimized transaction paths while minimizing transaction costs and maximizing transaction success rates.

Through this route mechanism, SideHub aims to play a crucial role in enhancing the interoperability of the Cosmos blockchain network, opening up numerous possibilities for token swaps and liquidity provision. It will be a significant step towards a more integrated and fluid decentralized financial landscape.

### Definitions

```
        channel-0 channel-1         channel-2 channel-3
┌───────┐       ibc        ┌───────┐        ibc       ┌───────┐
│Chain A│◄────────────────►│SideHub│◄────────────────►│Chain B│
└───────┘                  └───────┘                  └───────┘
     1. transfer 2. recv_packet     3. forward
         ─────────────────► packet  ─────────────────►
                            forward   4. timeout
                            middleware◄───────────────
                                    5. forward retry
                                    ─────────────────►
         7. ack ERR                 6. timeout
         ◄─────────────────         ◄─────────────────
```

### Desired Properties

- Efficiency: The routing module should allow for swift and cost-effective token swaps.

- Reliability: The mechanism should provide a reliable pathway for token conversion, with minimized risk of transaction failures.

- Flexibility: The routing module should support a wide range of tokens and facilitate seamless exchange between them.

- Optimized Liquidity Utilization: The routing module should effectively utilize the liquidity pools in SideHub to ensure successful token conversions.

- Transparency: The mechanism should provide users with clear information about the status of their transactions.

## Technical Specification

### General Design

The Router module will receive data from the Inter-Blockchain Communication (IBC) module. If a direct pool for the asset pair requested by the user doesn't exist, the module will find an alternate path through an intermediary token.

For instance, if a user wishes to swap TokenX for TokenY, but no direct pool exists, the module will suggest an alternate path such as TokenX -> TokenZ -> TokenY. The final result is subject to a slippage check.

The Sidehub, in this context, acts as a middleware layer for broadcasting IBCSwap messages received from other chains.
This will be implemented similarly with this [`packet-forward-middleware`](https://github.com/strangelove-ventures/packet-forward-middleware) module provided by `Strangelove Ventures`.

This module is designed for `Transfer` ibc call for ics-20 tokens. so need to customize this to support `atomicswap` and `interchainswap`

They implemented Transfer middleware but we will implement `atomicswap`, `interchainswap` broad casting function

Upon successful integration of this middleware, it becomes feasible to establish connections from any source chain to a target chain linked to the SideHub. This interconnected chain architecture further facilitates the relay of various Inter-Blockchain Communication (IBC) packets through this path. Operations under the protocols ics100 and ics101, which include MakeOrder, TakeOrder, CancelOrder, MakePool,TakePool, CancelPool, MakeMultiAssetDeposit, TakeMultiAssetDeposit,SingleAssetDeposit, and MultiAssetWithdraw, can effectively be implemented with two paths (source chain -> side hub -> destination chain).

The primary challenge arises in identifying viable paths for token pair swaps, a process that will be implemented off-chain. An in-depth understanding of all available token pairs is crucial for this task. This requirement calls for the registration and off-chain storage of all created pools.

When a user selects a token pair for a swap, the system must sift through the token pairs array to identify potential paths. An optimal path is then determined from this array, with the optimization rule centered around satisfying the final result slippage as much as possible.

As this process operates entirely off-chain, paths can be pre-calculated based on the token pair and stored off-chain. This path pre-calculation significantly enhances system performance and ensures seamless user experience by avoiding real-time computation.

For further ease of use, the frontend will access these pre-calculated paths, presenting them to users for selection. To avoid potential timeouts and maintain system responsiveness, the current implementation will limit the path to 2-3 steps.

Off-chain database storing the token pairs and paths must be updated regularly as chain states and token pair availability can change over time. This comprehensive, performance-optimized approach effectively combines on-chain operations and off-chain computations, offering a streamlined user experience while efficiently handling complex operations.

### Data Structure

this packet will be saved as a map data structure in genesis status.
so we will save all packets based on relayer address as a map data struture so can trace multi-hop process.

```typescript
interface InFlightPacket {
  originalSenderAddress: string;
  refundChannelId: string;
  refundPortId: string;
  packetSrcChannelId: string;
  packetSrcPortId: string;
  packetTimeoutTimestamp: number;
  packetTimeoutHeight: string;
  packetData: Uint8Array;
  refundSequence: number;
  retriesRemaining: number;
  timeout: number;
  nonrefundable: boolean;
}

interface GenesisState {
  params: Params; // Params is assumed to be another defined interface

  // The map is represented as an object in TypeScript
  // The keys are strings and the values are InFlightPacket interfaces
  inFlightPackets: { [key: string]: InFlightPacket };
}

interface IBCMiddleware {
  app: IBCModule; // Assuming IBCModule is defined elsewhere
  keeper: Keeper; // Assuming Keeper is defined elsewhere

  retriesOnTimeout: number;
  forwardTimeout: number; // Assuming this is in milliseconds
  refundTimeout: number; // Assuming this is in milliseconds
}
```

### Channel lifecycle management

```typescript

class NewIBCMiddleware {
   app: IBCModule;
   keeper: Keeper;
   retriesOnTimeout: number;
   forwardTimeout: number;
   refundTimeout: number;

   function new(
    app: IBCModule,
    keeper: Keeper,
    retriesOnTimeout: number,
    forwardTimeout: number,
    refundTimeout: number,
  ): IBCMiddleware {
    return {
        app: app,
        keeper: keeper,
        retriesOnTimeout: retriesOnTimeout,
        forwardTimeout: forwardTimeout,
        refundTimeout: refundTimeout,
    };
 }

  OnChanOpenInit(

  order: ChannelOrder,
  connectionHops: [Identifier],
  portIdentifier: Identifier,
  channelIdentifier: Identifier,
  counterpartyPortIdentifier: Identifier,
  counterpartyChannelIdentifier: Identifier,
  version: string) => (version: string, err: Error) {

      return this.app.OnChanOpenInit(ctx, order, connectionHops, portID, channelID, chanCap, counterparty, version)
  }

   OnChanOpenTry(
        ctx: sdk.Context,
        order: channeltypes.Order,
        connectionHops: string[],
        portID: string, channelID: string,
        chanCap: capabilitytypes.Capability,
        counterparty: channeltypes.Counterparty,
        counterpartyVersion: string
    ): [string, Error] {
        return this.app.OnChanOpenTry(ctx, order, connectionHops, portID, channelID, chanCap, counterparty, counterpartyVersion);
    }

    OnChanOpenAck(
        ctx: sdk.Context,
        portID: string,
        channelID: string,
        counterpartyChannelID: string,
        counterpartyVersion: string,
    ): Error {
        return this.app.OnChanOpenAck(ctx, portID, channelID, counterpartyChannelID, counterpartyVersion);
    }

    OnChanOpenConfirm(ctx: sdk.Context, portID: string, channelID: string): Error {
        return this.app.OnChanOpenConfirm(ctx, portID, channelID);
    }

    OnChanCloseInit(ctx: sdk.Context, portID: string, channelID: string): Error {
        return this.app.OnChanCloseInit(ctx, portID, channelID);
    }

    OnChanCloseConfirm(ctx: sdk.Context, portID: string, channelID: string): Error {
        return this.app.OnChanCloseConfirm(ctx, portID, channelID);
    }

   onRecvPacket(
    ctx: Context,
    packet: Channel.Packet,
    relayer: string
    ): Acknowledgement| null {

        // ...
        let metadata = m.forward;

        // The function getBoolFromAny() must be defined or imported
        let processed = getBoolFromAny(new ProcessedKey());
        let nonrefundable = getBoolFromAny(new NonrefundableKey());
        let disableDenomComposition = getBoolFromAny(new DisableDenomCompositionKey());

        if (metadata.validate() !== null) {
            return new Channel.ErrorAcknowledgement(metadata.validate().message);
        }

        if (!processed) {
            let ack = this.app.onRecvPacket(ctx, packet, relayer);
            if (ack === null || !ack.success()) {
                return ack;
            }
        }

        // The function getDenomForThisChain() must be defined or imported
        let denomOnThisChain = !disableDenomComposition ?
        getDenomForThisChain(packet.destinationPort, packet.destinationChannel, packet.sourcePort, packet.sourceChannel, data.denom) :
        data.denom;

        // ...

        // Some lines are omitted here

        if( metadata.Validate()) {
	        return channeltypes.NewErrorAcknowledgement(err)
        }

	// if this packet has been handled by another middleware in the stack there may be no need to call into the
	// underlying app, otherwise the transfer module's OnRecvPacket callback could be invoked more than once
	// which would mint/burn vouchers more than once
        if (!processed) {
	        ack := im.app.OnRecvPacket(ctx, packet, relayer)
	        if ack == nil || !ack.Success() {
	        	return ack
	        }
        }

	// if this packet's token denom is already the base denom for some native token on this chain,
	// we do not need to do any further composition of the denom before forwarding the packet
        let  denomOnThisChain = data.Denom
	      if !disableDenomComposition {
	      	denomOnThisChain = getDenomForThisChain(
	      		packet.DestinationPort, packet.DestinationChannel,
	      		packet.SourcePort, packet.SourceChannel,
	      		data.Denom,
	      	)
	      }


        abortTransactionUnless(data.Amount !== 0)

	      const token = {denom: denomOnThisChain, amount:amountInt}
        abortTransactionUnless(metadata.Timeout)

        const timeout = metadata.Timeout
	      const retries =  this.app.retriesOnTimeout


	      const err = this.app.keeper.ForwardTransferPacket(nil, packet, data.Sender, data.Receiver, metadata, token, retries, timeout, []metrics.Label{}, nonrefundable)

        if err != null {
		      return channeltypes.NewErrorAcknowledgement(err)
	      }

	      // returning nil ack will prevent WriteAcknowledgement from occurring for forwarded packet.
	      // This is intentional so that the acknowledgement will be written later based on the ack/timeout of the forwarded packet.
        return ;
    }

    onAcknowledgementPacket(
        packet: Channel.Packet,
        acknowledgement: Buffer,
        relayer: string
    ): Promise<Error | null> {

      let data: FungibleTokenPacketData;
      try {
        data = protobuf.parse(packet.getData());
      } catch (err) {
        this.keeper.logger(ctx).error('packetForwardMiddleware error parsing packet data from ack packet', {
            sequence: packet.sequence,
            srcChannel: packet.sourceChannel,
            srcPort: packet.sourcePort,
            dstChannel: packet.destinationChannel,
            dstPort: packet.destinationPort,
            error: err.message
        });
        return this.app.onAcknowledgementPacket(ctx, packet, acknowledgement, relayer);
      }


      var ack channeltypes.Acknowledgement
      const ack = channeltypes.SubModuleCdc.UnmarshalJSON(acknowledgement)


	    inFlightPacket := im.keeper.GetAndClearInFlightPacket(ctx, packet.SourceChannel, packet.SourcePort, packet.Sequence)
	    if inFlightPacket != nil {
		  // this is a forwarded packet, so override handling to avoid refund from being processed.
		  return this.keeper.WriteAcknowledgementForForwardedPacket(ctx, packet, data, inFlightPacket, ack)
	  }

	  return this.app.OnAcknowledgementPacket(ctx, packet, acknowledgement, relayer)
  }


  // OnTimeoutPacket implements the IBCModule interface.
  function OnTimeoutPacket(packet channeltypes.Packet, relayer sdk.AccAddress)  {
	  var data transfertypes.FungibleTokenPacketData
    const data = protobuf.parse(packet.GetData())
	  if (data == null) {
	  	this.keeper.Logger().Error("packetForwardMiddleware error parsing packet data from timeout packet",
	  		"sequence", packet.Sequence,
	  		"src-channel", packet.SourceChannel, "src-port", packet.SourcePort,
	  		"dst-channel", packet.DestinationChannel, "dst-port", packet.DestinationPort,
	  		"error", err,
	  	)
	  	return this.app.OnTimeoutPacket(packet, relayer)
	  }

	  this.keeper.Logger().Debug("packetForwardMiddleware OnAcknowledgementPacket",
	  	"sequence", packet.Sequence,
	  	"src-channel", packet.SourceChannel, "src-port", packet.SourcePort,
	  	"dst-channel", packet.DestinationChannel, "dst-port", packet.DestinationPort,
	  	"amount", data.Amount, "denom", data.Denom,
	  )

	  const inFlightPacket = this.keeper.TimeoutShouldRetry(packet)
    abortTransactionUnless(inFlightPacket !== undefined)
	  return im.app.OnTimeoutPacket(ctx, packet, relayer)
  }

  // SendPacket implements the ICS4 Wrapper interface.
  function SendPacket(
  	chanCap capabilitytypes.Capability,
  	sourcePort string, sourceChannel string,
  	timeoutHeight clienttypes.Height,
  	timeoutTimestamp uint64,
  	data []byte,
  ) {
  	return this.keeper.SendPacket(chanCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, data)
  }

  // WriteAcknowledgement implements the ICS4 Wrapper interface.
  function WriteAcknowledgement(
  	chanCap capabilitytypes.Capability,
  	packet ibcexported.PacketI,
  	ack ibcexported.Acknowledgement,
  ) error {
  	return this.keeper.WriteAcknowledgement(chanCap, packet, ack)
  }

  function GetAppVersion(
  	portID,
  	channelID string,
  ) (string, bool) {
  	return this.keeper.GetAppVersion(ctx, portID, channelID)
  }
}
```

### Improvements

- Timeout Adjustment: The timeout settings are crucial for efficient packet handling. It's important to ensure that the timeout duration is sufficient to complete all hop-ibc calls, preventing premature termination of processes.

- Path Calculation Prior to Packet Transmission: Before initiating packet transfers, especially with swaps, it's essential to first determine the optimal path. This means calculating the complete path in advance to ensure efficient routing and improved performance."
