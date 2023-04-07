/* tslint:disable */
/* eslint-disable */
/**
 * Marche Online
 * マルシェ管理者用API
 *
 * The version of the OpenAPI document: 0.1.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


import * as runtime from '../runtime';
import type {
  ErrorResponse,
  OrderResponse,
  OrdersResponse,
} from '../models';
import {
    ErrorResponseFromJSON,
    ErrorResponseToJSON,
    OrderResponseFromJSON,
    OrderResponseToJSON,
    OrdersResponseFromJSON,
    OrdersResponseToJSON,
} from '../models';

export interface V1GetOrderRequest {
    orderId: string;
}

export interface V1ListOrdersRequest {
    limit?: number;
    offset?: number;
    orders?: string;
}

/**
 * 
 */
export class OrderApi extends runtime.BaseAPI {

    /**
     * 注文取得
     */
    async v1GetOrderRaw(requestParameters: V1GetOrderRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<OrderResponse>> {
        if (requestParameters.orderId === null || requestParameters.orderId === undefined) {
            throw new runtime.RequiredError('orderId','Required parameter requestParameters.orderId was null or undefined when calling v1GetOrder.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("BearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/v1/orders/{orderId}`.replace(`{${"orderId"}}`, encodeURIComponent(String(requestParameters.orderId))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => OrderResponseFromJSON(jsonValue));
    }

    /**
     * 注文取得
     */
    async v1GetOrder(requestParameters: V1GetOrderRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<OrderResponse> {
        const response = await this.v1GetOrderRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * 注文一覧取得
     */
    async v1ListOrdersRaw(requestParameters: V1ListOrdersRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<OrdersResponse>> {
        const queryParameters: any = {};

        if (requestParameters.limit !== undefined) {
            queryParameters['limit'] = requestParameters.limit;
        }

        if (requestParameters.offset !== undefined) {
            queryParameters['offset'] = requestParameters.offset;
        }

        if (requestParameters.orders !== undefined) {
            queryParameters['orders'] = requestParameters.orders;
        }

        const headerParameters: runtime.HTTPHeaders = {};

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("BearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/v1/orders`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => OrdersResponseFromJSON(jsonValue));
    }

    /**
     * 注文一覧取得
     */
    async v1ListOrders(requestParameters: V1ListOrdersRequest = {}, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<OrdersResponse> {
        const response = await this.v1ListOrdersRaw(requestParameters, initOverrides);
        return await response.value();
    }

}
