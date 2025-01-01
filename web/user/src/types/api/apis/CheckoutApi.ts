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
  CheckoutExperienceRequest,
  CheckoutProductRequest,
  CheckoutResponse,
  CheckoutStateResponse,
  ErrorResponse,
  GuestCheckoutExperienceRequest,
  GuestCheckoutProductRequest,
  GuestCheckoutResponse,
  GuestCheckoutStateResponse,
  GuestPreCheckoutExperienceResponse,
  PreCheckoutExperienceResponse,
} from '../models/index';
import {
    CheckoutExperienceRequestFromJSON,
    CheckoutExperienceRequestToJSON,
    CheckoutProductRequestFromJSON,
    CheckoutProductRequestToJSON,
    CheckoutResponseFromJSON,
    CheckoutResponseToJSON,
    CheckoutStateResponseFromJSON,
    CheckoutStateResponseToJSON,
    ErrorResponseFromJSON,
    ErrorResponseToJSON,
    GuestCheckoutExperienceRequestFromJSON,
    GuestCheckoutExperienceRequestToJSON,
    GuestCheckoutProductRequestFromJSON,
    GuestCheckoutProductRequestToJSON,
    GuestCheckoutResponseFromJSON,
    GuestCheckoutResponseToJSON,
    GuestCheckoutStateResponseFromJSON,
    GuestCheckoutStateResponseToJSON,
    GuestPreCheckoutExperienceResponseFromJSON,
    GuestPreCheckoutExperienceResponseToJSON,
    PreCheckoutExperienceResponseFromJSON,
    PreCheckoutExperienceResponseToJSON,
} from '../models/index';

export interface V1CheckoutProductRequest {
    body: CheckoutProductRequest;
}

export interface V1CheckoutsExperiencesExperienceIdGetRequest {
    experienceId: string;
    promotion?: string;
    adult?: number;
    juniorHighSchool?: number;
    elementarySchool?: number;
    preschool?: number;
    senior?: number;
}

export interface V1CheckoutsExperiencesExperienceIdPostRequest {
    experienceId: string;
    body: CheckoutExperienceRequest;
}

export interface V1GetCheckoutStateRequest {
    transactionId: string;
}

export interface V1GetGuestCheckoutStateRequest {
    transactionId: string;
}

export interface V1GuestCheckoutExperienceRequest {
    experienceId: string;
    body: GuestCheckoutExperienceRequest;
}

export interface V1GuestCheckoutProductRequest {
    body: GuestCheckoutProductRequest;
}

export interface V1GuestsCheckoutsExperiencesExperienceIdGetRequest {
    experienceId: string;
    promotion?: string;
    adult?: number;
    juniorHighSchool?: number;
    elementarySchool?: number;
    preschool?: number;
    senior?: number;
}

/**
 * 
 */
export class CheckoutApi extends runtime.BaseAPI {

    /**
     * 商品購入
     */
    async v1CheckoutProductRaw(requestParameters: V1CheckoutProductRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<CheckoutResponse>> {
        if (requestParameters['body'] == null) {
            throw new runtime.RequiredError(
                'body',
                'Required parameter "body" was null or undefined when calling v1CheckoutProduct().'
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
            path: `/v1/checkouts/products`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: requestParameters['body'] as any,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => CheckoutResponseFromJSON(jsonValue));
    }

    /**
     * 商品購入
     */
    async v1CheckoutProduct(requestParameters: V1CheckoutProductRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<CheckoutResponse> {
        const response = await this.v1CheckoutProductRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * 体験購入前確認
     */
    async v1CheckoutsExperiencesExperienceIdGetRaw(requestParameters: V1CheckoutsExperiencesExperienceIdGetRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<PreCheckoutExperienceResponse>> {
        if (requestParameters['experienceId'] == null) {
            throw new runtime.RequiredError(
                'experienceId',
                'Required parameter "experienceId" was null or undefined when calling v1CheckoutsExperiencesExperienceIdGet().'
            );
        }

        const queryParameters: any = {};

        if (requestParameters['promotion'] != null) {
            queryParameters['promotion'] = requestParameters['promotion'];
        }

        if (requestParameters['adult'] != null) {
            queryParameters['adult'] = requestParameters['adult'];
        }

        if (requestParameters['juniorHighSchool'] != null) {
            queryParameters['juniorHighSchool'] = requestParameters['juniorHighSchool'];
        }

        if (requestParameters['elementarySchool'] != null) {
            queryParameters['elementarySchool'] = requestParameters['elementarySchool'];
        }

        if (requestParameters['preschool'] != null) {
            queryParameters['preschool'] = requestParameters['preschool'];
        }

        if (requestParameters['senior'] != null) {
            queryParameters['senior'] = requestParameters['senior'];
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
            path: `/v1/checkouts/experiences/{experienceId}`.replace(`{${"experienceId"}}`, encodeURIComponent(String(requestParameters['experienceId']))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => PreCheckoutExperienceResponseFromJSON(jsonValue));
    }

    /**
     * 体験購入前確認
     */
    async v1CheckoutsExperiencesExperienceIdGet(requestParameters: V1CheckoutsExperiencesExperienceIdGetRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<PreCheckoutExperienceResponse> {
        const response = await this.v1CheckoutsExperiencesExperienceIdGetRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * 体験購入
     */
    async v1CheckoutsExperiencesExperienceIdPostRaw(requestParameters: V1CheckoutsExperiencesExperienceIdPostRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<CheckoutResponse>> {
        if (requestParameters['experienceId'] == null) {
            throw new runtime.RequiredError(
                'experienceId',
                'Required parameter "experienceId" was null or undefined when calling v1CheckoutsExperiencesExperienceIdPost().'
            );
        }

        if (requestParameters['body'] == null) {
            throw new runtime.RequiredError(
                'body',
                'Required parameter "body" was null or undefined when calling v1CheckoutsExperiencesExperienceIdPost().'
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
            path: `/v1/checkouts/experiences/{experienceId}`.replace(`{${"experienceId"}}`, encodeURIComponent(String(requestParameters['experienceId']))),
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: requestParameters['body'] as any,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => CheckoutResponseFromJSON(jsonValue));
    }

    /**
     * 体験購入
     */
    async v1CheckoutsExperiencesExperienceIdPost(requestParameters: V1CheckoutsExperiencesExperienceIdPostRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<CheckoutResponse> {
        const response = await this.v1CheckoutsExperiencesExperienceIdPostRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * 注文情報の取得
     */
    async v1GetCheckoutStateRaw(requestParameters: V1GetCheckoutStateRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<CheckoutStateResponse>> {
        if (requestParameters['transactionId'] == null) {
            throw new runtime.RequiredError(
                'transactionId',
                'Required parameter "transactionId" was null or undefined when calling v1GetCheckoutState().'
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
            path: `/v1/checkouts/{transactionId}`.replace(`{${"transactionId"}}`, encodeURIComponent(String(requestParameters['transactionId']))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => CheckoutStateResponseFromJSON(jsonValue));
    }

    /**
     * 注文情報の取得
     */
    async v1GetCheckoutState(requestParameters: V1GetCheckoutStateRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<CheckoutStateResponse> {
        const response = await this.v1GetCheckoutStateRaw(requestParameters, initOverrides);
        return await response.value();
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

    /**
     * 体験購入前確認
     */
    async v1GuestsCheckoutsExperiencesExperienceIdGetRaw(requestParameters: V1GuestsCheckoutsExperiencesExperienceIdGetRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<GuestPreCheckoutExperienceResponse>> {
        if (requestParameters['experienceId'] == null) {
            throw new runtime.RequiredError(
                'experienceId',
                'Required parameter "experienceId" was null or undefined when calling v1GuestsCheckoutsExperiencesExperienceIdGet().'
            );
        }

        const queryParameters: any = {};

        if (requestParameters['promotion'] != null) {
            queryParameters['promotion'] = requestParameters['promotion'];
        }

        if (requestParameters['adult'] != null) {
            queryParameters['adult'] = requestParameters['adult'];
        }

        if (requestParameters['juniorHighSchool'] != null) {
            queryParameters['juniorHighSchool'] = requestParameters['juniorHighSchool'];
        }

        if (requestParameters['elementarySchool'] != null) {
            queryParameters['elementarySchool'] = requestParameters['elementarySchool'];
        }

        if (requestParameters['preschool'] != null) {
            queryParameters['preschool'] = requestParameters['preschool'];
        }

        if (requestParameters['senior'] != null) {
            queryParameters['senior'] = requestParameters['senior'];
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
            path: `/v1/guests/checkouts/experiences/{experienceId}`.replace(`{${"experienceId"}}`, encodeURIComponent(String(requestParameters['experienceId']))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => GuestPreCheckoutExperienceResponseFromJSON(jsonValue));
    }

    /**
     * 体験購入前確認
     */
    async v1GuestsCheckoutsExperiencesExperienceIdGet(requestParameters: V1GuestsCheckoutsExperiencesExperienceIdGetRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<GuestPreCheckoutExperienceResponse> {
        const response = await this.v1GuestsCheckoutsExperiencesExperienceIdGetRaw(requestParameters, initOverrides);
        return await response.value();
    }

}
