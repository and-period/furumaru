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
  VideoResponse,
  VideosResponse,
} from '../models/index';
import {
    ErrorResponseFromJSON,
    ErrorResponseToJSON,
    VideoResponseFromJSON,
    VideoResponseToJSON,
    VideosResponseFromJSON,
    VideosResponseToJSON,
} from '../models/index';

export interface V1GetVideoRequest {
    videoId: string;
}

export interface V1VideosRequest {
    limit?: number;
    offset?: number;
    coordinator?: string;
    category?: string;
}

/**
 * 
 */
export class VideoApi extends runtime.BaseAPI {

    /**
     * オンデマンド配信取得
     */
    async v1GetVideoRaw(requestParameters: V1GetVideoRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<VideoResponse>> {
        if (requestParameters['videoId'] == null) {
            throw new runtime.RequiredError(
                'videoId',
                'Required parameter "videoId" was null or undefined when calling v1GetVideo().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/v1/videos/{videoId}`.replace(`{${"videoId"}}`, encodeURIComponent(String(requestParameters['videoId']))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => VideoResponseFromJSON(jsonValue));
    }

    /**
     * オンデマンド配信取得
     */
    async v1GetVideo(requestParameters: V1GetVideoRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<VideoResponse> {
        const response = await this.v1GetVideoRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * オンデマンド配信一覧取得
     */
    async v1VideosRaw(requestParameters: V1VideosRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<VideosResponse>> {
        const queryParameters: any = {};

        if (requestParameters['limit'] != null) {
            queryParameters['limit'] = requestParameters['limit'];
        }

        if (requestParameters['offset'] != null) {
            queryParameters['offset'] = requestParameters['offset'];
        }

        if (requestParameters['coordinator'] != null) {
            queryParameters['coordinator'] = requestParameters['coordinator'];
        }

        if (requestParameters['category'] != null) {
            queryParameters['category'] = requestParameters['category'];
        }

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/v1/videos`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => VideosResponseFromJSON(jsonValue));
    }

    /**
     * オンデマンド配信一覧取得
     */
    async v1Videos(requestParameters: V1VideosRequest = {}, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<VideosResponse> {
        const response = await this.v1VideosRaw(requestParameters, initOverrides);
        return await response.value();
    }

}