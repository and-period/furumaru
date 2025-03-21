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
  CreateSpotRequest,
  ErrorResponse,
  SpotResponse,
  SpotsResponse,
  UpdateSpotRequest,
} from '../models/index';
import {
    CreateSpotRequestFromJSON,
    CreateSpotRequestToJSON,
    ErrorResponseFromJSON,
    ErrorResponseToJSON,
    SpotResponseFromJSON,
    SpotResponseToJSON,
    SpotsResponseFromJSON,
    SpotsResponseToJSON,
    UpdateSpotRequestFromJSON,
    UpdateSpotRequestToJSON,
} from '../models/index';

export interface V1CreateSpotRequest {
    body: CreateSpotRequest;
}

export interface V1DeleteSpotRequest {
    spotId: string;
}

export interface V1GetSpotRequest {
    spotId: string;
}

export interface V1ListSpotsRequest {
    longitude: number;
    latitude: number;
    radius?: number;
}

export interface V1UpdateSpotRequest {
    spotId: string;
    body: UpdateSpotRequest;
}

/**
 * 
 */
export class SpotApi extends runtime.BaseAPI {

    /**
     * スポット登録
     */
    async v1CreateSpotRaw(requestParameters: V1CreateSpotRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<SpotResponse>> {
        if (requestParameters['body'] == null) {
            throw new runtime.RequiredError(
                'body',
                'Required parameter "body" was null or undefined when calling v1CreateSpot().'
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
            path: `/v1/spots`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: requestParameters['body'] as any,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => SpotResponseFromJSON(jsonValue));
    }

    /**
     * スポット登録
     */
    async v1CreateSpot(requestParameters: V1CreateSpotRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<SpotResponse> {
        const response = await this.v1CreateSpotRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * スポット削除
     */
    async v1DeleteSpotRaw(requestParameters: V1DeleteSpotRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<object>> {
        if (requestParameters['spotId'] == null) {
            throw new runtime.RequiredError(
                'spotId',
                'Required parameter "spotId" was null or undefined when calling v1DeleteSpot().'
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
            path: `/v1/spots/{spotId}`.replace(`{${"spotId"}}`, encodeURIComponent(String(requestParameters['spotId']))),
            method: 'DELETE',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse<any>(response);
    }

    /**
     * スポット削除
     */
    async v1DeleteSpot(requestParameters: V1DeleteSpotRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<object> {
        const response = await this.v1DeleteSpotRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * スポット取得
     */
    async v1GetSpotRaw(requestParameters: V1GetSpotRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<SpotResponse>> {
        if (requestParameters['spotId'] == null) {
            throw new runtime.RequiredError(
                'spotId',
                'Required parameter "spotId" was null or undefined when calling v1GetSpot().'
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
            path: `/v1/spots/{spotId}`.replace(`{${"spotId"}}`, encodeURIComponent(String(requestParameters['spotId']))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => SpotResponseFromJSON(jsonValue));
    }

    /**
     * スポット取得
     */
    async v1GetSpot(requestParameters: V1GetSpotRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<SpotResponse> {
        const response = await this.v1GetSpotRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * スポット一覧取得
     */
    async v1ListSpotsRaw(requestParameters: V1ListSpotsRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<SpotsResponse>> {
        if (requestParameters['longitude'] == null) {
            throw new runtime.RequiredError(
                'longitude',
                'Required parameter "longitude" was null or undefined when calling v1ListSpots().'
            );
        }

        if (requestParameters['latitude'] == null) {
            throw new runtime.RequiredError(
                'latitude',
                'Required parameter "latitude" was null or undefined when calling v1ListSpots().'
            );
        }

        const queryParameters: any = {};

        if (requestParameters['longitude'] != null) {
            queryParameters['longitude'] = requestParameters['longitude'];
        }

        if (requestParameters['latitude'] != null) {
            queryParameters['latitude'] = requestParameters['latitude'];
        }

        if (requestParameters['radius'] != null) {
            queryParameters['radius'] = requestParameters['radius'];
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
            path: `/v1/spots`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => SpotsResponseFromJSON(jsonValue));
    }

    /**
     * スポット一覧取得
     */
    async v1ListSpots(requestParameters: V1ListSpotsRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<SpotsResponse> {
        const response = await this.v1ListSpotsRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * スポット更新
     */
    async v1UpdateSpotRaw(requestParameters: V1UpdateSpotRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<object>> {
        if (requestParameters['spotId'] == null) {
            throw new runtime.RequiredError(
                'spotId',
                'Required parameter "spotId" was null or undefined when calling v1UpdateSpot().'
            );
        }

        if (requestParameters['body'] == null) {
            throw new runtime.RequiredError(
                'body',
                'Required parameter "body" was null or undefined when calling v1UpdateSpot().'
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
            path: `/v1/spots/{spotId}`.replace(`{${"spotId"}}`, encodeURIComponent(String(requestParameters['spotId']))),
            method: 'PATCH',
            headers: headerParameters,
            query: queryParameters,
            body: requestParameters['body'] as any,
        }, initOverrides);

        return new runtime.JSONApiResponse<any>(response);
    }

    /**
     * スポット更新
     */
    async v1UpdateSpot(requestParameters: V1UpdateSpotRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<object> {
        const response = await this.v1UpdateSpotRaw(requestParameters, initOverrides);
        return await response.value();
    }

}
