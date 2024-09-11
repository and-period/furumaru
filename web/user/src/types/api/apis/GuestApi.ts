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
  CreateGuestLiveCommentRequest,
  CreateGuestVideoCommentRequest,
  ErrorResponse,
  GuestCheckoutProductRequest,
  GuestCheckoutResponse,
  GuestCheckoutStateResponse,
} from '../models/index';
import {
    CreateGuestLiveCommentRequestFromJSON,
    CreateGuestLiveCommentRequestToJSON,
    CreateGuestVideoCommentRequestFromJSON,
    CreateGuestVideoCommentRequestToJSON,
    ErrorResponseFromJSON,
    ErrorResponseToJSON,
    GuestCheckoutProductRequestFromJSON,
    GuestCheckoutProductRequestToJSON,
    GuestCheckoutResponseFromJSON,
    GuestCheckoutResponseToJSON,
    GuestCheckoutStateResponseFromJSON,
    GuestCheckoutStateResponseToJSON,
} from '../models/index';

export interface V1CreateGuestLiveCommentRequest {
    scheduleId: string;
    body: CreateGuestLiveCommentRequest;
}

export interface V1CreateGuestVideoCommentRequest {
    videoId: string;
    body: CreateGuestVideoCommentRequest;
}

export interface V1GetGuestCheckoutStateRequest {
    transactionId: string;
}

export interface V1GuestCheckoutExperienceRequest {
    experienceId: string;
    body: GuestCheckoutProductRequest;
}

export interface V1GuestCheckoutProductRequest {
    body: GuestCheckoutProductRequest;
}

/**
 * 
 */
export class GuestApi extends runtime.BaseAPI {

    /**
     * ライブ配信ゲストコメント投稿
     */
    async v1CreateGuestLiveCommentRaw(requestParameters: V1CreateGuestLiveCommentRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<void>> {
        if (requestParameters['scheduleId'] == null) {
            throw new runtime.RequiredError(
                'scheduleId',
                'Required parameter "scheduleId" was null or undefined when calling v1CreateGuestLiveComment().'
            );
        }

        if (requestParameters['body'] == null) {
            throw new runtime.RequiredError(
                'body',
                'Required parameter "body" was null or undefined when calling v1CreateGuestLiveComment().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/v1/guests/schedules/{scheduleId}/comments`.replace(`{${"scheduleId"}}`, encodeURIComponent(String(requestParameters['scheduleId']))),
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: requestParameters['body'] as any,
        }, initOverrides);

        return new runtime.VoidApiResponse(response);
    }

    /**
     * ライブ配信ゲストコメント投稿
     */
    async v1CreateGuestLiveComment(requestParameters: V1CreateGuestLiveCommentRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<void> {
        await this.v1CreateGuestLiveCommentRaw(requestParameters, initOverrides);
    }

    /**
     * オンデマンド配信ゲストコメント投稿
     */
    async v1CreateGuestVideoCommentRaw(requestParameters: V1CreateGuestVideoCommentRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<void>> {
        if (requestParameters['videoId'] == null) {
            throw new runtime.RequiredError(
                'videoId',
                'Required parameter "videoId" was null or undefined when calling v1CreateGuestVideoComment().'
            );
        }

        if (requestParameters['body'] == null) {
            throw new runtime.RequiredError(
                'body',
                'Required parameter "body" was null or undefined when calling v1CreateGuestVideoComment().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/v1/guests/videos/{videoId}/comments`.replace(`{${"videoId"}}`, encodeURIComponent(String(requestParameters['videoId']))),
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: requestParameters['body'] as any,
        }, initOverrides);

        return new runtime.VoidApiResponse(response);
    }

    /**
     * オンデマンド配信ゲストコメント投稿
     */
    async v1CreateGuestVideoComment(requestParameters: V1CreateGuestVideoCommentRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<void> {
        await this.v1CreateGuestVideoCommentRaw(requestParameters, initOverrides);
    }

    /**
     * ゲスト注文情報の取得
     */
    async v1GetGuestCheckoutStateRaw(requestParameters: V1GetGuestCheckoutStateRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<GuestCheckoutStateResponse>> {
        if (requestParameters['transactionId'] == null) {
            throw new runtime.RequiredError(
                'transactionId',
                'Required parameter "transactionId" was null or undefined when calling v1GetGuestCheckoutState().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/v1/guests/checkouts/{transactionId}`.replace(`{${"transactionId"}}`, encodeURIComponent(String(requestParameters['transactionId']))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => GuestCheckoutStateResponseFromJSON(jsonValue));
    }

    /**
     * ゲスト注文情報の取得
     */
    async v1GetGuestCheckoutState(requestParameters: V1GetGuestCheckoutStateRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<GuestCheckoutStateResponse> {
        const response = await this.v1GetGuestCheckoutStateRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * ゲスト体験購入
     */
    async v1GuestCheckoutExperienceRaw(requestParameters: V1GuestCheckoutExperienceRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<GuestCheckoutResponse>> {
        if (requestParameters['experienceId'] == null) {
            throw new runtime.RequiredError(
                'experienceId',
                'Required parameter "experienceId" was null or undefined when calling v1GuestCheckoutExperience().'
            );
        }

        if (requestParameters['body'] == null) {
            throw new runtime.RequiredError(
                'body',
                'Required parameter "body" was null or undefined when calling v1GuestCheckoutExperience().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/v1/guests/checkouts/experiences/{experienceId}`.replace(`{${"experienceId"}}`, encodeURIComponent(String(requestParameters['experienceId']))),
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: requestParameters['body'] as any,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => GuestCheckoutResponseFromJSON(jsonValue));
    }

    /**
     * ゲスト体験購入
     */
    async v1GuestCheckoutExperience(requestParameters: V1GuestCheckoutExperienceRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<GuestCheckoutResponse> {
        const response = await this.v1GuestCheckoutExperienceRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * ゲスト商品購入
     */
    async v1GuestCheckoutProductRaw(requestParameters: V1GuestCheckoutProductRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<GuestCheckoutResponse>> {
        if (requestParameters['body'] == null) {
            throw new runtime.RequiredError(
                'body',
                'Required parameter "body" was null or undefined when calling v1GuestCheckoutProduct().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/v1/guests/checkouts/products`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: requestParameters['body'] as any,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => GuestCheckoutResponseFromJSON(jsonValue));
    }

    /**
     * ゲスト商品購入
     */
    async v1GuestCheckoutProduct(requestParameters: V1GuestCheckoutProductRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<GuestCheckoutResponse> {
        const response = await this.v1GuestCheckoutProductRaw(requestParameters, initOverrides);
        return await response.value();
    }

}
