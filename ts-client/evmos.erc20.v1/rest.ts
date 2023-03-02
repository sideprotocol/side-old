/* eslint-disable */
/* tslint:disable */
/*
 * ---------------------------------------------------------------
 * ## THIS FILE WAS GENERATED VIA SWAGGER-TYPESCRIPT-API        ##
 * ##                                                           ##
 * ## AUTHOR: acacode                                           ##
 * ## SOURCE: https://github.com/acacode/swagger-typescript-api ##
 * ---------------------------------------------------------------
 */

export interface Erc20V1Params {
  /** enable_erc20 is the parameter to enable the conversion of Cosmos coins <--> ERC20 tokens. */
  enable_erc20?: boolean;

  /**
   * enable_evm_hook is the parameter to enable the EVM hook that converts an ERC20 token to a Cosmos
   * Coin by transferring the Tokens through a MsgEthereumTx to the ModuleAddress Ethereum address.
   */
  enable_evm_hook?: boolean;
}

export interface ProtobufAny {
  "@type"?: string;
}

export interface RpcStatus {
  /** @format int32 */
  code?: number;
  message?: string;
  details?: ProtobufAny[];
}

export type V1MsgConvertCoinResponse = object;

export type V1MsgConvertERC20Response = object;

/**
* Owner enumerates the ownership of a ERC20 contract.

 - OWNER_UNSPECIFIED: OWNER_UNSPECIFIED defines an invalid/undefined owner.
 - OWNER_MODULE: OWNER_MODULE - erc20 is owned by the erc20 module account.
 - OWNER_EXTERNAL: OWNER_EXTERNAL - erc20 is owned by an external account.
*/
export enum V1Owner {
  OWNER_UNSPECIFIED = "OWNER_UNSPECIFIED",
  OWNER_MODULE = "OWNER_MODULE",
  OWNER_EXTERNAL = "OWNER_EXTERNAL",
}

/**
* QueryParamsResponse is the response type for the Query/Params RPC
method.
*/
export interface V1QueryParamsResponse {
  /** params are the erc20 module parameters */
  params?: Erc20V1Params;
}

/**
* QueryTokenPairResponse is the response type for the Query/TokenPair RPC
method.
*/
export interface V1QueryTokenPairResponse {
  /**
   * token_pairs returns the info about a registered token pair for the erc20 module
   * TokenPair defines an instance that records a pairing consisting of a native
   *  Cosmos Coin and an ERC20 token address.
   */
  token_pair?: V1TokenPair;
}

/**
* QueryTokenPairsResponse is the response type for the Query/TokenPairs RPC
method.
*/
export interface V1QueryTokenPairsResponse {
  /** token_pairs is a slice of registered token pairs for the erc20 module */
  token_pairs?: V1TokenPair[];

  /** pagination defines the pagination in the response. */
  pagination?: V1Beta1PageResponse;
}

/**
* TokenPair defines an instance that records a pairing consisting of a native
 Cosmos Coin and an ERC20 token address.
*/
export interface V1TokenPair {
  /** erc20_address is the hex address of ERC20 contract token */
  erc20_address?: string;

  /** denom defines the cosmos base denomination to be mapped to */
  denom?: string;

  /** enabled defines the token mapping enable status */
  enabled?: boolean;

  /**
   * contract_owner is the an ENUM specifying the type of ERC20 owner (0 invalid, 1 ModuleAccount, 2 external address)
   * Owner enumerates the ownership of a ERC20 contract.
   *
   *  - OWNER_UNSPECIFIED: OWNER_UNSPECIFIED defines an invalid/undefined owner.
   *  - OWNER_MODULE: OWNER_MODULE - erc20 is owned by the erc20 module account.
   *  - OWNER_EXTERNAL: OWNER_EXTERNAL - erc20 is owned by an external account.
   */
  contract_owner?: V1Owner;
}

/**
* Coin defines a token with a denomination and an amount.

NOTE: The amount field is an Int which implements the custom method
signatures required by gogoproto.
*/
export interface V1Beta1Coin {
  denom?: string;
  amount?: string;
}

/**
* message SomeRequest {
         Foo some_parameter = 1;
         PageRequest pagination = 2;
 }
*/
export interface V1Beta1PageRequest {
  /**
   * key is a value returned in PageResponse.next_key to begin
   * querying the next page most efficiently. Only one of offset or key
   * should be set.
   * @format byte
   */
  key?: string;

  /**
   * offset is a numeric offset that can be used when key is unavailable.
   * It is less efficient than using key. Only one of offset or key should
   * be set.
   * @format uint64
   */
  offset?: string;

  /**
   * limit is the total number of results to be returned in the result page.
   * If left empty it will default to a value to be set by each app.
   * @format uint64
   */
  limit?: string;

  /**
   * count_total is set to true  to indicate that the result set should include
   * a count of the total number of items available for pagination in UIs.
   * count_total is only respected when offset is used. It is ignored when key
   * is set.
   */
  count_total?: boolean;

  /**
   * reverse is set to true if results are to be returned in the descending order.
   *
   * Since: cosmos-sdk 0.43
   */
  reverse?: boolean;
}

/**
* PageResponse is to be embedded in gRPC response messages where the
corresponding request message has used PageRequest.

 message SomeResponse {
         repeated Bar results = 1;
         PageResponse page = 2;
 }
*/
export interface V1Beta1PageResponse {
  /**
   * next_key is the key to be passed to PageRequest.key to
   * query the next page most efficiently. It will be empty if
   * there are no more results.
   * @format byte
   */
  next_key?: string;

  /**
   * total is total number of results available if PageRequest.count_total
   * was set, its value is undefined otherwise
   * @format uint64
   */
  total?: string;
}

import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse, ResponseType } from "axios";

export type QueryParamsType = Record<string | number, any>;

export interface FullRequestParams extends Omit<AxiosRequestConfig, "data" | "params" | "url" | "responseType"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseType;
  /** request body */
  body?: unknown;
}

export type RequestParams = Omit<FullRequestParams, "body" | "method" | "query" | "path">;

export interface ApiConfig<SecurityDataType = unknown> extends Omit<AxiosRequestConfig, "data" | "cancelToken"> {
  securityWorker?: (
    securityData: SecurityDataType | null,
  ) => Promise<AxiosRequestConfig | void> | AxiosRequestConfig | void;
  secure?: boolean;
  format?: ResponseType;
}

export enum ContentType {
  Json = "application/json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
}

export class HttpClient<SecurityDataType = unknown> {
  public instance: AxiosInstance;
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private secure?: boolean;
  private format?: ResponseType;

  constructor({ securityWorker, secure, format, ...axiosConfig }: ApiConfig<SecurityDataType> = {}) {
    this.instance = axios.create({ ...axiosConfig, baseURL: axiosConfig.baseURL || "" });
    this.secure = secure;
    this.format = format;
    this.securityWorker = securityWorker;
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  private mergeRequestParams(params1: AxiosRequestConfig, params2?: AxiosRequestConfig): AxiosRequestConfig {
    return {
      ...this.instance.defaults,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...(this.instance.defaults.headers || {}),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    };
  }

  private createFormData(input: Record<string, unknown>): FormData {
    return Object.keys(input || {}).reduce((formData, key) => {
      const property = input[key];
      formData.append(
        key,
        property instanceof Blob
          ? property
          : typeof property === "object" && property !== null
          ? JSON.stringify(property)
          : `${property}`,
      );
      return formData;
    }, new FormData());
  }

  public request = async <T = any, _E = any>({
    secure,
    path,
    type,
    query,
    format,
    body,
    ...params
  }: FullRequestParams): Promise<AxiosResponse<T>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const responseFormat = (format && this.format) || void 0;

    if (type === ContentType.FormData && body && body !== null && typeof body === "object") {
      requestParams.headers.common = { Accept: "*/*" };
      requestParams.headers.post = {};
      requestParams.headers.put = {};

      body = this.createFormData(body as Record<string, unknown>);
    }

    return this.instance.request({
      ...requestParams,
      headers: {
        ...(type && type !== ContentType.FormData ? { "Content-Type": type } : {}),
        ...(requestParams.headers || {}),
      },
      params: query,
      responseType: responseFormat,
      data: body,
      url: path,
    });
  };
}

/**
 * @title evmos/erc20/v1/erc20.proto
 * @version version not set
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  /**
   * No description
   *
   * @tags Query
   * @name QueryParams
   * @summary Params retrieves the erc20 module params
   * @request GET:/evmos/erc20/v1/params
   */
  queryParams = (params: RequestParams = {}) =>
    this.request<V1QueryParamsResponse, RpcStatus>({
      path: `/evmos/erc20/v1/params`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryTokenPairs
   * @summary TokenPairs retrieves registered token pairs
   * @request GET:/evmos/erc20/v1/token_pairs
   */
  queryTokenPairs = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<V1QueryTokenPairsResponse, RpcStatus>({
      path: `/evmos/erc20/v1/token_pairs`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryTokenPair
   * @summary TokenPair retrieves a registered token pair
   * @request GET:/evmos/erc20/v1/token_pairs/{token}
   */
  queryTokenPair = (token: string, params: RequestParams = {}) =>
    this.request<V1QueryTokenPairResponse, RpcStatus>({
      path: `/evmos/erc20/v1/token_pairs/${token}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
 * No description
 * 
 * @tags Msg
 * @name MsgConvertCoin
 * @summary ConvertCoin mints a ERC20 representation of the native Cosmos coin denom
that is registered on the token mapping.
 * @request GET:/evmos/erc20/v1/tx/convert_coin
 */
  msgConvertCoin = (
    query?: { "coin.denom"?: string; "coin.amount"?: string; receiver?: string; sender?: string },
    params: RequestParams = {},
  ) =>
    this.request<V1MsgConvertCoinResponse, RpcStatus>({
      path: `/evmos/erc20/v1/tx/convert_coin`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
 * No description
 * 
 * @tags Msg
 * @name MsgConvertErc20
 * @summary ConvertERC20 mints a native Cosmos coin representation of the ERC20 token
contract that is registered on the token mapping.
 * @request GET:/evmos/erc20/v1/tx/convert_erc20
 */
  msgConvertErc20 = (
    query?: { contract_address?: string; amount?: string; receiver?: string; sender?: string },
    params: RequestParams = {},
  ) =>
    this.request<V1MsgConvertERC20Response, RpcStatus>({
      path: `/evmos/erc20/v1/tx/convert_erc20`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });
}
