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
  CreateProductReviewRequest,
  ErrorResponse,
  ProductResponse,
  ProductReviewResponse,
  ProductReviewsResponse,
  ProductsResponse,
  UpdateProductReviewRequest,
} from '../models/index';
import {
    CreateProductReviewRequestFromJSON,
    CreateProductReviewRequestToJSON,
    ErrorResponseFromJSON,
    ErrorResponseToJSON,
    ProductResponseFromJSON,
    ProductResponseToJSON,
    ProductReviewResponseFromJSON,
    ProductReviewResponseToJSON,
    ProductReviewsResponseFromJSON,
    ProductReviewsResponseToJSON,
    ProductsResponseFromJSON,
    ProductsResponseToJSON,
    UpdateProductReviewRequestFromJSON,
    UpdateProductReviewRequestToJSON,
} from '../models/index';

export interface V1CreateProductReviewRequest {
    productId: string;
    body: CreateProductReviewRequest;
}

export interface V1DeleteProductReviewRequest {
    productId: string;
    reviewId: string;
}

export interface V1GetProductRequest {
    productId: string;
}

export interface V1GetProductReviewRequest {
    productId: string;
    reviewId: string;
}

export interface V1ListProductReviewsRequest {
    productId: string;
    userId?: string;
    limit?: number;
    nextToken?: string;
    rates?: string;
}

export interface V1ListProductsRequest {
    limit?: number;
    offset?: number;
    coordinatorId?: string;
}

export interface V1UpdateProductReviewRequest {
    productId: string;
    reviewId: string;
    body: UpdateProductReviewRequest;
}

/**
 * 
 */
export class ProductApi extends runtime.BaseAPI {

    /**
     * 商品レビュー投稿
     */
    async v1CreateProductReviewRaw(requestParameters: V1CreateProductReviewRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<void>> {
        if (requestParameters['productId'] == null) {
            throw new runtime.RequiredError(
                'productId',
                'Required parameter "productId" was null or undefined when calling v1CreateProductReview().'
            );
        }

        if (requestParameters['body'] == null) {
            throw new runtime.RequiredError(
                'body',
                'Required parameter "body" was null or undefined when calling v1CreateProductReview().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("bearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/v1/products/{productId}/reviews`.replace(`{${"productId"}}`, encodeURIComponent(String(requestParameters['productId']))),
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: requestParameters['body'] as any,
        }, initOverrides);

        return new runtime.VoidApiResponse(response);
    }

    /**
     * 商品レビュー投稿
     */
    async v1CreateProductReview(requestParameters: V1CreateProductReviewRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<void> {
        await this.v1CreateProductReviewRaw(requestParameters, initOverrides);
    }

    /**
     * 商品レビュー削除
     */
    async v1DeleteProductReviewRaw(requestParameters: V1DeleteProductReviewRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<void>> {
        if (requestParameters['productId'] == null) {
            throw new runtime.RequiredError(
                'productId',
                'Required parameter "productId" was null or undefined when calling v1DeleteProductReview().'
            );
        }

        if (requestParameters['reviewId'] == null) {
            throw new runtime.RequiredError(
                'reviewId',
                'Required parameter "reviewId" was null or undefined when calling v1DeleteProductReview().'
            );
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
            path: `/v1/products/{productId}/reviews/{reviewId}`.replace(`{${"productId"}}`, encodeURIComponent(String(requestParameters['productId']))).replace(`{${"reviewId"}}`, encodeURIComponent(String(requestParameters['reviewId']))),
            method: 'DELETE',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.VoidApiResponse(response);
    }

    /**
     * 商品レビュー削除
     */
    async v1DeleteProductReview(requestParameters: V1DeleteProductReviewRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<void> {
        await this.v1DeleteProductReviewRaw(requestParameters, initOverrides);
    }

    /**
     * 商品取得
     */
    async v1GetProductRaw(requestParameters: V1GetProductRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<ProductResponse>> {
        if (requestParameters['productId'] == null) {
            throw new runtime.RequiredError(
                'productId',
                'Required parameter "productId" was null or undefined when calling v1GetProduct().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/v1/products/{productId}`.replace(`{${"productId"}}`, encodeURIComponent(String(requestParameters['productId']))),
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
     * 商品レビュー取得
     */
    async v1GetProductReviewRaw(requestParameters: V1GetProductReviewRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<ProductReviewResponse>> {
        if (requestParameters['productId'] == null) {
            throw new runtime.RequiredError(
                'productId',
                'Required parameter "productId" was null or undefined when calling v1GetProductReview().'
            );
        }

        if (requestParameters['reviewId'] == null) {
            throw new runtime.RequiredError(
                'reviewId',
                'Required parameter "reviewId" was null or undefined when calling v1GetProductReview().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/v1/products/{productId}/reviews/{reviewId}`.replace(`{${"productId"}}`, encodeURIComponent(String(requestParameters['productId']))).replace(`{${"reviewId"}}`, encodeURIComponent(String(requestParameters['reviewId']))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => ProductReviewResponseFromJSON(jsonValue));
    }

    /**
     * 商品レビュー取得
     */
    async v1GetProductReview(requestParameters: V1GetProductReviewRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<ProductReviewResponse> {
        const response = await this.v1GetProductReviewRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * 商品レビュー一覧取得
     */
    async v1ListProductReviewsRaw(requestParameters: V1ListProductReviewsRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<ProductReviewsResponse>> {
        if (requestParameters['productId'] == null) {
            throw new runtime.RequiredError(
                'productId',
                'Required parameter "productId" was null or undefined when calling v1ListProductReviews().'
            );
        }

        const queryParameters: any = {};

        if (requestParameters['userId'] != null) {
            queryParameters['userId'] = requestParameters['userId'];
        }

        if (requestParameters['limit'] != null) {
            queryParameters['limit'] = requestParameters['limit'];
        }

        if (requestParameters['nextToken'] != null) {
            queryParameters['nextToken'] = requestParameters['nextToken'];
        }

        if (requestParameters['rates'] != null) {
            queryParameters['rates'] = requestParameters['rates'];
        }

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/v1/products/{productId}/reviews`.replace(`{${"productId"}}`, encodeURIComponent(String(requestParameters['productId']))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => ProductReviewsResponseFromJSON(jsonValue));
    }

    /**
     * 商品レビュー一覧取得
     */
    async v1ListProductReviews(requestParameters: V1ListProductReviewsRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<ProductReviewsResponse> {
        const response = await this.v1ListProductReviewsRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * 商品一覧取得
     */
    async v1ListProductsRaw(requestParameters: V1ListProductsRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<ProductsResponse>> {
        const queryParameters: any = {};

        if (requestParameters['limit'] != null) {
            queryParameters['limit'] = requestParameters['limit'];
        }

        if (requestParameters['offset'] != null) {
            queryParameters['offset'] = requestParameters['offset'];
        }

        if (requestParameters['coordinatorId'] != null) {
            queryParameters['coordinatorId'] = requestParameters['coordinatorId'];
        }

        const headerParameters: runtime.HTTPHeaders = {};

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

    /**
     * 商品レビュー更新
     */
    async v1UpdateProductReviewRaw(requestParameters: V1UpdateProductReviewRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<void>> {
        if (requestParameters['productId'] == null) {
            throw new runtime.RequiredError(
                'productId',
                'Required parameter "productId" was null or undefined when calling v1UpdateProductReview().'
            );
        }

        if (requestParameters['reviewId'] == null) {
            throw new runtime.RequiredError(
                'reviewId',
                'Required parameter "reviewId" was null or undefined when calling v1UpdateProductReview().'
            );
        }

        if (requestParameters['body'] == null) {
            throw new runtime.RequiredError(
                'body',
                'Required parameter "body" was null or undefined when calling v1UpdateProductReview().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("bearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/v1/products/{productId}/reviews/{reviewId}`.replace(`{${"productId"}}`, encodeURIComponent(String(requestParameters['productId']))).replace(`{${"reviewId"}}`, encodeURIComponent(String(requestParameters['reviewId']))),
            method: 'PATCH',
            headers: headerParameters,
            query: queryParameters,
            body: requestParameters['body'] as any,
        }, initOverrides);

        return new runtime.VoidApiResponse(response);
    }

    /**
     * 商品レビュー更新
     */
    async v1UpdateProductReview(requestParameters: V1UpdateProductReviewRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<void> {
        await this.v1UpdateProductReviewRaw(requestParameters, initOverrides);
    }

}
