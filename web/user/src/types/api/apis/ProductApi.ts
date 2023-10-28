/* tslint:disable */
/* eslint-disable */
/**
 * Marche Online
 * マルシェ購入者用API
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
  ProductResponse,
  ProductsResponse,
} from '../models/index';
import {
    ErrorResponseFromJSON,
    ErrorResponseToJSON,
    ProductResponseFromJSON,
    ProductResponseToJSON,
    ProductsResponseFromJSON,
    ProductsResponseToJSON,
} from '../models/index';

export interface V1GetProductRequest {
    productId: string;
}

export interface V1ListProductsRequest {
    limit?: number;
    offset?: number;
}

/**
 * 
 */
export class ProductApi extends runtime.BaseAPI {

    /**
     * 商品取得
     */
    async v1GetProductRaw(requestParameters: V1GetProductRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<ProductResponse>> {
        if (requestParameters.productId === null || requestParameters.productId === undefined) {
            throw new runtime.RequiredError('productId','Required parameter requestParameters.productId was null or undefined when calling v1GetProduct.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("bearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/v1/products/{productId}`.replace(`{${"productId"}}`, encodeURIComponent(String(requestParameters.productId))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => ProductResponseFromJSON(jsonValue));
    }

    /**
     * 商品取得
     */
    async v1GetProduct(requestParameters: V1GetProductRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<ProductResponse> {
        const response = await this.v1GetProductRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * 商品一覧取得
     */
    async v1ListProductsRaw(requestParameters: V1ListProductsRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<ProductsResponse>> {
        const queryParameters: any = {};

        if (requestParameters.limit !== undefined) {
            queryParameters['limit'] = requestParameters.limit;
        }

        if (requestParameters.offset !== undefined) {
            queryParameters['offset'] = requestParameters.offset;
        }

        const headerParameters: runtime.HTTPHeaders = {};

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("bearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/v1/products`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => ProductsResponseFromJSON(jsonValue));
    }

    /**
     * 商品一覧取得
     */
    async v1ListProducts(requestParameters: V1ListProductsRequest = {}, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<ProductsResponse> {
        const response = await this.v1ListProductsRaw(requestParameters, initOverrides);
        return await response.value();
    }

}
