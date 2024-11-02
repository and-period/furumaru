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
  ArchiveSchedulesResponse,
  CreateGuestLiveCommentRequest,
  CreateLiveCommentRequest,
  ErrorResponse,
  LiveCommentsResponse,
  LiveSchedulesResponse,
  ScheduleResponse,
} from '../models/index';
import {
    ArchiveSchedulesResponseFromJSON,
    ArchiveSchedulesResponseToJSON,
    CreateGuestLiveCommentRequestFromJSON,
    CreateGuestLiveCommentRequestToJSON,
    CreateLiveCommentRequestFromJSON,
    CreateLiveCommentRequestToJSON,
    ErrorResponseFromJSON,
    ErrorResponseToJSON,
    LiveCommentsResponseFromJSON,
    LiveCommentsResponseToJSON,
    LiveSchedulesResponseFromJSON,
    LiveSchedulesResponseToJSON,
    ScheduleResponseFromJSON,
    ScheduleResponseToJSON,
} from '../models/index';

export interface V1ArchiveSchedulesRequest {
    limit?: number;
    offset?: number;
    coordinator?: string;
    producer?: string;
}

export interface V1CreateGuestLiveCommentRequest {
    scheduleId: string;
    body: CreateGuestLiveCommentRequest;
}

export interface V1CreateLiveCommentRequest {
    scheduleId: string;
    body: CreateLiveCommentRequest;
}

export interface V1GetScheduleRequest {
    scheduleId: string;
}

export interface V1ListLiveCommentsRequest {
    scheduleId: string;
    limit?: number;
    next?: string;
    start?: number;
    end?: number;
}

export interface V1LiveSchedulesRequest {
    limit?: number;
    offset?: number;
    coordinator?: string;
    producer?: string;
}

/**
 * 
 */
export class ScheduleApi extends runtime.BaseAPI {

    /**
     * 過去のマルシェ一覧取得
     */
    async v1ArchiveSchedulesRaw(requestParameters: V1ArchiveSchedulesRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<ArchiveSchedulesResponse>> {
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

        if (requestParameters['producer'] != null) {
            queryParameters['producer'] = requestParameters['producer'];
        }

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/v1/schedules/archives`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => ArchiveSchedulesResponseFromJSON(jsonValue));
    }

    /**
     * 過去のマルシェ一覧取得
     */
    async v1ArchiveSchedules(requestParameters: V1ArchiveSchedulesRequest = {}, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<ArchiveSchedulesResponse> {
        const response = await this.v1ArchiveSchedulesRaw(requestParameters, initOverrides);
        return await response.value();
    }

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
     * ライブ配信コメント投稿
     */
    async v1CreateLiveCommentRaw(requestParameters: V1CreateLiveCommentRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<void>> {
        if (requestParameters['scheduleId'] == null) {
            throw new runtime.RequiredError(
                'scheduleId',
                'Required parameter "scheduleId" was null or undefined when calling v1CreateLiveComment().'
            );
        }

        if (requestParameters['body'] == null) {
            throw new runtime.RequiredError(
                'body',
                'Required parameter "body" was null or undefined when calling v1CreateLiveComment().'
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
            path: `/v1/schedules/{scheduleId}/comments`.replace(`{${"scheduleId"}}`, encodeURIComponent(String(requestParameters['scheduleId']))),
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: requestParameters['body'] as any,
        }, initOverrides);

        return new runtime.VoidApiResponse(response);
    }

    /**
     * ライブ配信コメント投稿
     */
    async v1CreateLiveComment(requestParameters: V1CreateLiveCommentRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<void> {
        await this.v1CreateLiveCommentRaw(requestParameters, initOverrides);
    }

    /**
     * マルシェ開催スケジュール取得
     */
    async v1GetScheduleRaw(requestParameters: V1GetScheduleRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<ScheduleResponse>> {
        if (requestParameters['scheduleId'] == null) {
            throw new runtime.RequiredError(
                'scheduleId',
                'Required parameter "scheduleId" was null or undefined when calling v1GetSchedule().'
            );
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/v1/schedules/{scheduleId}`.replace(`{${"scheduleId"}}`, encodeURIComponent(String(requestParameters['scheduleId']))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => ScheduleResponseFromJSON(jsonValue));
    }

    /**
     * マルシェ開催スケジュール取得
     */
    async v1GetSchedule(requestParameters: V1GetScheduleRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<ScheduleResponse> {
        const response = await this.v1GetScheduleRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * ライブ配信コメント取得
     */
    async v1ListLiveCommentsRaw(requestParameters: V1ListLiveCommentsRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<LiveCommentsResponse>> {
        if (requestParameters['scheduleId'] == null) {
            throw new runtime.RequiredError(
                'scheduleId',
                'Required parameter "scheduleId" was null or undefined when calling v1ListLiveComments().'
            );
        }

        const queryParameters: any = {};

        if (requestParameters['limit'] != null) {
            queryParameters['limit'] = requestParameters['limit'];
        }

        if (requestParameters['next'] != null) {
            queryParameters['next'] = requestParameters['next'];
        }

        if (requestParameters['start'] != null) {
            queryParameters['start'] = requestParameters['start'];
        }

        if (requestParameters['end'] != null) {
            queryParameters['end'] = requestParameters['end'];
        }

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/v1/schedules/{scheduleId}/comments`.replace(`{${"scheduleId"}}`, encodeURIComponent(String(requestParameters['scheduleId']))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => LiveCommentsResponseFromJSON(jsonValue));
    }

    /**
     * ライブ配信コメント取得
     */
    async v1ListLiveComments(requestParameters: V1ListLiveCommentsRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<LiveCommentsResponse> {
        const response = await this.v1ListLiveCommentsRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * 配信中・配信予定のマルシェ一覧取得
     */
    async v1LiveSchedulesRaw(requestParameters: V1LiveSchedulesRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<LiveSchedulesResponse>> {
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

        if (requestParameters['producer'] != null) {
            queryParameters['producer'] = requestParameters['producer'];
        }

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/v1/schedules/lives`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => LiveSchedulesResponseFromJSON(jsonValue));
    }

    /**
     * 配信中・配信予定のマルシェ一覧取得
     */
    async v1LiveSchedules(requestParameters: V1LiveSchedulesRequest = {}, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<LiveSchedulesResponse> {
        const response = await this.v1LiveSchedulesRaw(requestParameters, initOverrides);
        return await response.value();
    }

}
