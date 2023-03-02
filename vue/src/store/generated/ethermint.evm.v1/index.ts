import { Client, registry, MissingWalletError } from 'sidechain-client-ts'

import { EventEthereumTx } from "sidechain-client-ts/ethermint.evm.v1/types"
import { EventTxLog } from "sidechain-client-ts/ethermint.evm.v1/types"
import { EventMessage } from "sidechain-client-ts/ethermint.evm.v1/types"
import { EventBlockBloom } from "sidechain-client-ts/ethermint.evm.v1/types"
import { Params } from "sidechain-client-ts/ethermint.evm.v1/types"
import { ChainConfig } from "sidechain-client-ts/ethermint.evm.v1/types"
import { State } from "sidechain-client-ts/ethermint.evm.v1/types"
import { TransactionLogs } from "sidechain-client-ts/ethermint.evm.v1/types"
import { Log } from "sidechain-client-ts/ethermint.evm.v1/types"
import { TxResult } from "sidechain-client-ts/ethermint.evm.v1/types"
import { AccessTuple } from "sidechain-client-ts/ethermint.evm.v1/types"
import { TraceConfig } from "sidechain-client-ts/ethermint.evm.v1/types"
import { GenesisAccount } from "sidechain-client-ts/ethermint.evm.v1/types"
import { QueryTxLogsRequest } from "sidechain-client-ts/ethermint.evm.v1/types"
import { QueryTxLogsResponse } from "sidechain-client-ts/ethermint.evm.v1/types"
import { LegacyTx } from "sidechain-client-ts/ethermint.evm.v1/types"
import { AccessListTx } from "sidechain-client-ts/ethermint.evm.v1/types"
import { DynamicFeeTx } from "sidechain-client-ts/ethermint.evm.v1/types"
import { ExtensionOptionsEthereumTx } from "sidechain-client-ts/ethermint.evm.v1/types"


export { EventEthereumTx, EventTxLog, EventMessage, EventBlockBloom, Params, ChainConfig, State, TransactionLogs, Log, TxResult, AccessTuple, TraceConfig, GenesisAccount, QueryTxLogsRequest, QueryTxLogsResponse, LegacyTx, AccessListTx, DynamicFeeTx, ExtensionOptionsEthereumTx };

function initClient(vuexGetters) {
	return new Client(vuexGetters['common/env/getEnv'], vuexGetters['common/wallet/signer'])
}

function mergeResults(value, next_values) {
	for (let prop of Object.keys(next_values)) {
		if (Array.isArray(next_values[prop])) {
			value[prop]=[...value[prop], ...next_values[prop]]
		}else{
			value[prop]=next_values[prop]
		}
	}
	return value
}

type Field = {
	name: string;
	type: unknown;
}
function getStructure(template) {
	let structure: {fields: Field[]} = { fields: [] }
	for (const [key, value] of Object.entries(template)) {
		let field = { name: key, type: typeof value }
		structure.fields.push(field)
	}
	return structure
}
const getDefaultState = () => {
	return {
				Account: {},
				CosmosAccount: {},
				ValidatorAccount: {},
				Balance: {},
				Storage: {},
				Code: {},
				Params: {},
				EthCall: {},
				EstimateGas: {},
				TraceTx: {},
				TraceBlock: {},
				BaseFee: {},
				EthereumTx: {},
				
				_Structure: {
						EventEthereumTx: getStructure(EventEthereumTx.fromPartial({})),
						EventTxLog: getStructure(EventTxLog.fromPartial({})),
						EventMessage: getStructure(EventMessage.fromPartial({})),
						EventBlockBloom: getStructure(EventBlockBloom.fromPartial({})),
						Params: getStructure(Params.fromPartial({})),
						ChainConfig: getStructure(ChainConfig.fromPartial({})),
						State: getStructure(State.fromPartial({})),
						TransactionLogs: getStructure(TransactionLogs.fromPartial({})),
						Log: getStructure(Log.fromPartial({})),
						TxResult: getStructure(TxResult.fromPartial({})),
						AccessTuple: getStructure(AccessTuple.fromPartial({})),
						TraceConfig: getStructure(TraceConfig.fromPartial({})),
						GenesisAccount: getStructure(GenesisAccount.fromPartial({})),
						QueryTxLogsRequest: getStructure(QueryTxLogsRequest.fromPartial({})),
						QueryTxLogsResponse: getStructure(QueryTxLogsResponse.fromPartial({})),
						LegacyTx: getStructure(LegacyTx.fromPartial({})),
						AccessListTx: getStructure(AccessListTx.fromPartial({})),
						DynamicFeeTx: getStructure(DynamicFeeTx.fromPartial({})),
						ExtensionOptionsEthereumTx: getStructure(ExtensionOptionsEthereumTx.fromPartial({})),
						
		},
		_Registry: registry,
		_Subscriptions: new Set(),
	}
}

// initial state
const state = getDefaultState()

export default {
	namespaced: true,
	state,
	mutations: {
		RESET_STATE(state) {
			Object.assign(state, getDefaultState())
		},
		QUERY(state, { query, key, value }) {
			state[query][JSON.stringify(key)] = value
		},
		SUBSCRIBE(state, subscription) {
			state._Subscriptions.add(JSON.stringify(subscription))
		},
		UNSUBSCRIBE(state, subscription) {
			state._Subscriptions.delete(JSON.stringify(subscription))
		}
	},
	getters: {
				getAccount: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.Account[JSON.stringify(params)] ?? {}
		},
				getCosmosAccount: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.CosmosAccount[JSON.stringify(params)] ?? {}
		},
				getValidatorAccount: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.ValidatorAccount[JSON.stringify(params)] ?? {}
		},
				getBalance: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.Balance[JSON.stringify(params)] ?? {}
		},
				getStorage: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.Storage[JSON.stringify(params)] ?? {}
		},
				getCode: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.Code[JSON.stringify(params)] ?? {}
		},
				getParams: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.Params[JSON.stringify(params)] ?? {}
		},
				getEthCall: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.EthCall[JSON.stringify(params)] ?? {}
		},
				getEstimateGas: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.EstimateGas[JSON.stringify(params)] ?? {}
		},
				getTraceTx: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.TraceTx[JSON.stringify(params)] ?? {}
		},
				getTraceBlock: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.TraceBlock[JSON.stringify(params)] ?? {}
		},
				getBaseFee: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.BaseFee[JSON.stringify(params)] ?? {}
		},
				getEthereumTx: (state) => (params = { params: {}}) => {
					if (!(<any> params).query) {
						(<any> params).query=null
					}
			return state.EthereumTx[JSON.stringify(params)] ?? {}
		},
				
		getTypeStructure: (state) => (type) => {
			return state._Structure[type].fields
		},
		getRegistry: (state) => {
			return state._Registry
		}
	},
	actions: {
		init({ dispatch, rootGetters }) {
			console.log('Vuex module: ethermint.evm.v1 initialized!')
			if (rootGetters['common/env/client']) {
				rootGetters['common/env/client'].on('newblock', () => {
					dispatch('StoreUpdate')
				})
			}
		},
		resetState({ commit }) {
			commit('RESET_STATE')
		},
		unsubscribe({ commit }, subscription) {
			commit('UNSUBSCRIBE', subscription)
		},
		async StoreUpdate({ state, dispatch }) {
			state._Subscriptions.forEach(async (subscription) => {
				try {
					const sub=JSON.parse(subscription)
					await dispatch(sub.action, sub.payload)
				}catch(e) {
					throw new Error('Subscriptions: ' + e.message)
				}
			})
		},
		
		
		
		 		
		
		
		async QueryAccount({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params, query=null }) {
			try {
				const key = params ?? {};
				const client = initClient(rootGetters);
				let value= (await client.EthermintEvmV1.query.queryAccount( key.address)).data
				
					
				commit('QUERY', { query: 'Account', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryAccount', payload: { options: { all }, params: {...key},query }})
				return getters['getAccount']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new Error('QueryClient:QueryAccount API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryCosmosAccount({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params, query=null }) {
			try {
				const key = params ?? {};
				const client = initClient(rootGetters);
				let value= (await client.EthermintEvmV1.query.queryCosmosAccount( key.address)).data
				
					
				commit('QUERY', { query: 'CosmosAccount', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryCosmosAccount', payload: { options: { all }, params: {...key},query }})
				return getters['getCosmosAccount']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new Error('QueryClient:QueryCosmosAccount API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryValidatorAccount({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params, query=null }) {
			try {
				const key = params ?? {};
				const client = initClient(rootGetters);
				let value= (await client.EthermintEvmV1.query.queryValidatorAccount( key.cons_address)).data
				
					
				commit('QUERY', { query: 'ValidatorAccount', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryValidatorAccount', payload: { options: { all }, params: {...key},query }})
				return getters['getValidatorAccount']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new Error('QueryClient:QueryValidatorAccount API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryBalance({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params, query=null }) {
			try {
				const key = params ?? {};
				const client = initClient(rootGetters);
				let value= (await client.EthermintEvmV1.query.queryBalance( key.address)).data
				
					
				commit('QUERY', { query: 'Balance', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryBalance', payload: { options: { all }, params: {...key},query }})
				return getters['getBalance']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new Error('QueryClient:QueryBalance API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryStorage({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params, query=null }) {
			try {
				const key = params ?? {};
				const client = initClient(rootGetters);
				let value= (await client.EthermintEvmV1.query.queryStorage( key.address,  key.key)).data
				
					
				commit('QUERY', { query: 'Storage', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryStorage', payload: { options: { all }, params: {...key},query }})
				return getters['getStorage']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new Error('QueryClient:QueryStorage API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryCode({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params, query=null }) {
			try {
				const key = params ?? {};
				const client = initClient(rootGetters);
				let value= (await client.EthermintEvmV1.query.queryCode( key.address)).data
				
					
				commit('QUERY', { query: 'Code', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryCode', payload: { options: { all }, params: {...key},query }})
				return getters['getCode']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new Error('QueryClient:QueryCode API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryParams({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params, query=null }) {
			try {
				const key = params ?? {};
				const client = initClient(rootGetters);
				let value= (await client.EthermintEvmV1.query.queryParams()).data
				
					
				commit('QUERY', { query: 'Params', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryParams', payload: { options: { all }, params: {...key},query }})
				return getters['getParams']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new Error('QueryClient:QueryParams API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryEthCall({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params, query=null }) {
			try {
				const key = params ?? {};
				const client = initClient(rootGetters);
				let value= (await client.EthermintEvmV1.query.queryEthCall(query ?? undefined)).data
				
					
				while (all && (<any> value).pagination && (<any> value).pagination.next_key!=null) {
					let next_values=(await client.EthermintEvmV1.query.queryEthCall({...query ?? {}, 'pagination.key':(<any> value).pagination.next_key} as any)).data
					value = mergeResults(value, next_values);
				}
				commit('QUERY', { query: 'EthCall', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryEthCall', payload: { options: { all }, params: {...key},query }})
				return getters['getEthCall']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new Error('QueryClient:QueryEthCall API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryEstimateGas({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params, query=null }) {
			try {
				const key = params ?? {};
				const client = initClient(rootGetters);
				let value= (await client.EthermintEvmV1.query.queryEstimateGas(query ?? undefined)).data
				
					
				while (all && (<any> value).pagination && (<any> value).pagination.next_key!=null) {
					let next_values=(await client.EthermintEvmV1.query.queryEstimateGas({...query ?? {}, 'pagination.key':(<any> value).pagination.next_key} as any)).data
					value = mergeResults(value, next_values);
				}
				commit('QUERY', { query: 'EstimateGas', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryEstimateGas', payload: { options: { all }, params: {...key},query }})
				return getters['getEstimateGas']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new Error('QueryClient:QueryEstimateGas API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryTraceTx({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params, query=null }) {
			try {
				const key = params ?? {};
				const client = initClient(rootGetters);
				let value= (await client.EthermintEvmV1.query.queryTraceTx(query ?? undefined)).data
				
					
				while (all && (<any> value).pagination && (<any> value).pagination.next_key!=null) {
					let next_values=(await client.EthermintEvmV1.query.queryTraceTx({...query ?? {}, 'pagination.key':(<any> value).pagination.next_key} as any)).data
					value = mergeResults(value, next_values);
				}
				commit('QUERY', { query: 'TraceTx', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryTraceTx', payload: { options: { all }, params: {...key},query }})
				return getters['getTraceTx']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new Error('QueryClient:QueryTraceTx API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryTraceBlock({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params, query=null }) {
			try {
				const key = params ?? {};
				const client = initClient(rootGetters);
				let value= (await client.EthermintEvmV1.query.queryTraceBlock(query ?? undefined)).data
				
					
				while (all && (<any> value).pagination && (<any> value).pagination.next_key!=null) {
					let next_values=(await client.EthermintEvmV1.query.queryTraceBlock({...query ?? {}, 'pagination.key':(<any> value).pagination.next_key} as any)).data
					value = mergeResults(value, next_values);
				}
				commit('QUERY', { query: 'TraceBlock', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryTraceBlock', payload: { options: { all }, params: {...key},query }})
				return getters['getTraceBlock']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new Error('QueryClient:QueryTraceBlock API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async QueryBaseFee({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params, query=null }) {
			try {
				const key = params ?? {};
				const client = initClient(rootGetters);
				let value= (await client.EthermintEvmV1.query.queryBaseFee()).data
				
					
				commit('QUERY', { query: 'BaseFee', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'QueryBaseFee', payload: { options: { all }, params: {...key},query }})
				return getters['getBaseFee']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new Error('QueryClient:QueryBaseFee API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		
		
		 		
		
		
		async MsgEthereumTx({ commit, rootGetters, getters }, { options: { subscribe, all} = { subscribe:false, all:false}, params, query=null }) {
			try {
				const key = params ?? {};
				const client = initClient(rootGetters);
				let value= (await client.EthermintEvmV1.query.msgEthereumTx(query ?? undefined)).data
				
					
				while (all && (<any> value).pagination && (<any> value).pagination.next_key!=null) {
					let next_values=(await client.EthermintEvmV1.query.msgEthereumTx({...query ?? {}, 'pagination.key':(<any> value).pagination.next_key} as any)).data
					value = mergeResults(value, next_values);
				}
				commit('QUERY', { query: 'EthereumTx', key: { params: {...key}, query}, value })
				if (subscribe) commit('SUBSCRIBE', { action: 'MsgEthereumTx', payload: { options: { all }, params: {...key},query }})
				return getters['getEthereumTx']( { params: {...key}, query}) ?? {}
			} catch (e) {
				throw new Error('QueryClient:MsgEthereumTx API Node Unavailable. Could not perform query: ' + e.message)
				
			}
		},
		
		
		async sendMsgEthereumTx({ rootGetters }, { value, fee = {amount: [], gas: "200000"}, memo = '' }) {
			try {
				const client=await initClient(rootGetters)
				const fullFee = Array.isArray(fee)  ? {amount: fee, gas: "200000"} :fee;
				const result = await client.EthermintEvmV1.tx.sendMsgEthereumTx({ value, fee: fullFee, memo })
				return result
			} catch (e) {
				if (e == MissingWalletError) {
					throw new Error('TxClient:MsgEthereumTx:Init Could not initialize signing client. Wallet is required.')
				}else{
					throw new Error('TxClient:MsgEthereumTx:Send Could not broadcast Tx: '+ e.message)
				}
			}
		},
		
		async MsgEthereumTx({ rootGetters }, { value }) {
			try {
				const client=initClient(rootGetters)
				const msg = await client.EthermintEvmV1.tx.msgEthereumTx({value})
				return msg
			} catch (e) {
				if (e == MissingWalletError) {
					throw new Error('TxClient:MsgEthereumTx:Init Could not initialize signing client. Wallet is required.')
				} else{
					throw new Error('TxClient:MsgEthereumTx:Create Could not create message: ' + e.message)
				}
			}
		},
		
	}
}