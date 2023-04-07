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
  CreateProducerRequest,
  ErrorResponse,
  ProducerResponse,
  ProducersResponse,
  UpdateProducerEmailRequest,
  UpdateProducerRequest,
  UploadImageResponse,
} from '../models';
import {
    CreateProducerRequestFromJSON,
    CreateProducerRequestToJSON,
    ErrorResponseFromJSON,
    ErrorResponseToJSON,
    ProducerResponseFromJSON,
    ProducerResponseToJSON,
    ProducersResponseFromJSON,
    ProducersResponseToJSON,
    UpdateProducerEmailRequestFromJSON,
    UpdateProducerEmailRequestToJSON,
    UpdateProducerRequestFromJSON,
    UpdateProducerRequestToJSON,
    UploadImageResponseFromJSON,
    UploadImageResponseToJSON,
} from '../models';

export interface V1CreateProducerRequest {
    body: CreateProducerRequest;
}

export interface V1DeleteProducerRequest {
    producerId: string;
}

export interface V1GetProducerRequest {
    producerId: string;
}

export interface V1ListProducersRequest {
    limit?: number;
    offset?: number;
    filters?: string;
}

export interface V1UpdateProducerRequest {
    producerId: string;
    body: UpdateProducerRequest;
}

export interface V1UpdateProducerEmailRequest {
    producerId: string;
    body: UpdateProducerEmailRequest;
}

export interface V1UpdateProducerPasswordRequest {
    producerId: string;
    body?: object;
}

export interface V1UploadProducerHeaderRequest {
    image?: Blob;
}

export interface V1UploadProducerThumbnailRequest {
    thumbnail?: Blob;
}

/**
 * 
 */
export class ProducerApi extends runtime.BaseAPI {

    /**
     * 生産者登録
     */
    async v1CreateProducerRaw(requestParameters: V1CreateProducerRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<ProducerResponse>> {
        if (requestParameters.body === null || requestParameters.body === undefined) {
            throw new runtime.RequiredError('body','Required parameter requestParameters.body was null or undefined when calling v1CreateProducer.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("BearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/v1/producers`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: requestParameters.body as any,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => ProducerResponseFromJSON(jsonValue));
    }

    /**
     * 生産者登録
     */
    async v1CreateProducer(requestParameters: V1CreateProducerRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<ProducerResponse> {
        const response = await this.v1CreateProducerRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * 生産者退会
     */
    async v1DeleteProducerRaw(requestParameters: V1DeleteProducerRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<object>> {
        if (requestParameters.producerId === null || requestParameters.producerId === undefined) {
            throw new runtime.RequiredError('producerId','Required parameter requestParameters.producerId was null or undefined when calling v1DeleteProducer.');
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
            path: `/v1/producers/{producerId}`.replace(`{${"producerId"}}`, encodeURIComponent(String(requestParameters.producerId))),
            method: 'DELETE',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse<any>(response);
    }

    /**
     * 生産者退会
     */
    async v1DeleteProducer(requestParameters: V1DeleteProducerRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<object> {
        const response = await this.v1DeleteProducerRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * 生産者取得
     */
    async v1GetProducerRaw(requestParameters: V1GetProducerRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<ProducerResponse>> {
        if (requestParameters.producerId === null || requestParameters.producerId === undefined) {
            throw new runtime.RequiredError('producerId','Required parameter requestParameters.producerId was null or undefined when calling v1GetProducer.');
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
            path: `/v1/producers/{producerId}`.replace(`{${"producerId"}}`, encodeURIComponent(String(requestParameters.producerId))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => ProducerResponseFromJSON(jsonValue));
    }

    /**
     * 生産者取得
     */
    async v1GetProducer(requestParameters: V1GetProducerRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<ProducerResponse> {
        const response = await this.v1GetProducerRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * 生産者一覧取得
     */
    async v1ListProducersRaw(requestParameters: V1ListProducersRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<ProducersResponse>> {
        const queryParameters: any = {};

        if (requestParameters.limit !== undefined) {
            queryParameters['limit'] = requestParameters.limit;
        }

        if (requestParameters.offset !== undefined) {
            queryParameters['offset'] = requestParameters.offset;
        }

        if (requestParameters.filters !== undefined) {
            queryParameters['filters'] = requestParameters.filters;
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
            path: `/v1/producers`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => ProducersResponseFromJSON(jsonValue));
    }

    /**
     * 生産者一覧取得
     */
    async v1ListProducers(requestParameters: V1ListProducersRequest = {}, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<ProducersResponse> {
        const response = await this.v1ListProducersRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * 生産者更新
     */
    async v1UpdateProducerRaw(requestParameters: V1UpdateProducerRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<object>> {
        if (requestParameters.producerId === null || requestParameters.producerId === undefined) {
            throw new runtime.RequiredError('producerId','Required parameter requestParameters.producerId was null or undefined when calling v1UpdateProducer.');
        }

        if (requestParameters.body === null || requestParameters.body === undefined) {
            throw new runtime.RequiredError('body','Required parameter requestParameters.body was null or undefined when calling v1UpdateProducer.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("BearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/v1/producers/{producerId}`.replace(`{${"producerId"}}`, encodeURIComponent(String(requestParameters.producerId))),
            method: 'PATCH',
            headers: headerParameters,
            query: queryParameters,
            body: requestParameters.body as any,
        }, initOverrides);

        return new runtime.JSONApiResponse<any>(response);
    }

    /**
     * 生産者更新
     */
    async v1UpdateProducer(requestParameters: V1UpdateProducerRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<object> {
        const response = await this.v1UpdateProducerRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * 生産者メールアドレス更新
     */
    async v1UpdateProducerEmailRaw(requestParameters: V1UpdateProducerEmailRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<object>> {
        if (requestParameters.producerId === null || requestParameters.producerId === undefined) {
            throw new runtime.RequiredError('producerId','Required parameter requestParameters.producerId was null or undefined when calling v1UpdateProducerEmail.');
        }

        if (requestParameters.body === null || requestParameters.body === undefined) {
            throw new runtime.RequiredError('body','Required parameter requestParameters.body was null or undefined when calling v1UpdateProducerEmail.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("BearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/v1/producers/{producerId}/email`.replace(`{${"producerId"}}`, encodeURIComponent(String(requestParameters.producerId))),
            method: 'PATCH',
            headers: headerParameters,
            query: queryParameters,
            body: requestParameters.body as any,
        }, initOverrides);

        return new runtime.JSONApiResponse<any>(response);
    }

    /**
     * 生産者メールアドレス更新
     */
    async v1UpdateProducerEmail(requestParameters: V1UpdateProducerEmailRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<object> {
        const response = await this.v1UpdateProducerEmailRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * 生産者パスワード更新(ランダム生成)
     */
    async v1UpdateProducerPasswordRaw(requestParameters: V1UpdateProducerPasswordRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<object>> {
        if (requestParameters.producerId === null || requestParameters.producerId === undefined) {
            throw new runtime.RequiredError('producerId','Required parameter requestParameters.producerId was null or undefined when calling v1UpdateProducerPassword.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("BearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const response = await this.request({
            path: `/v1/producers/{producerId}/password`.replace(`{${"producerId"}}`, encodeURIComponent(String(requestParameters.producerId))),
            method: 'PATCH',
            headers: headerParameters,
            query: queryParameters,
            body: requestParameters.body as any,
        }, initOverrides);

        return new runtime.JSONApiResponse<any>(response);
    }

    /**
     * 生産者パスワード更新(ランダム生成)
     */
    async v1UpdateProducerPassword(requestParameters: V1UpdateProducerPasswordRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<object> {
        const response = await this.v1UpdateProducerPasswordRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * 生産者ヘッダー画像アップロード
     */
    async v1UploadProducerHeaderRaw(requestParameters: V1UploadProducerHeaderRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<UploadImageResponse>> {
        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("BearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const consumes: runtime.Consume[] = [
            { contentType: 'multipart/form-data' },
        ];
        // @ts-ignore: canConsumeForm may be unused
        const canConsumeForm = runtime.canConsumeForm(consumes);

        let formParams: { append(param: string, value: any): any };
        let useForm = false;
        // use FormData to transmit files using content-type "multipart/form-data"
        useForm = canConsumeForm;
        if (useForm) {
            formParams = new FormData();
        } else {
            formParams = new URLSearchParams();
        }

        if (requestParameters.image !== undefined) {
            formParams.append('image', requestParameters.image as any);
        }

        const response = await this.request({
            path: `/v1/upload/producers/header`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: formParams,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => UploadImageResponseFromJSON(jsonValue));
    }

    /**
     * 生産者ヘッダー画像アップロード
     */
    async v1UploadProducerHeader(requestParameters: V1UploadProducerHeaderRequest = {}, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<UploadImageResponse> {
        const response = await this.v1UploadProducerHeaderRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * 生産者サムネイルアップロード
     */
    async v1UploadProducerThumbnailRaw(requestParameters: V1UploadProducerThumbnailRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<UploadImageResponse>> {
        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        if (this.configuration && this.configuration.accessToken) {
            const token = this.configuration.accessToken;
            const tokenString = await token("BearerAuth", []);

            if (tokenString) {
                headerParameters["Authorization"] = `Bearer ${tokenString}`;
            }
        }
        const consumes: runtime.Consume[] = [
            { contentType: 'multipart/form-data' },
        ];
        // @ts-ignore: canConsumeForm may be unused
        const canConsumeForm = runtime.canConsumeForm(consumes);

        let formParams: { append(param: string, value: any): any };
        let useForm = false;
        // use FormData to transmit files using content-type "multipart/form-data"
        useForm = canConsumeForm;
        if (useForm) {
            formParams = new FormData();
        } else {
            formParams = new URLSearchParams();
        }

        if (requestParameters.thumbnail !== undefined) {
            formParams.append('thumbnail', requestParameters.thumbnail as any);
        }

        const response = await this.request({
            path: `/v1/upload/producers/thumbnail`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: formParams,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => UploadImageResponseFromJSON(jsonValue));
    }

    /**
     * 生産者サムネイルアップロード
     */
    async v1UploadProducerThumbnail(requestParameters: V1UploadProducerThumbnailRequest = {}, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<UploadImageResponse> {
        const response = await this.v1UploadProducerThumbnailRaw(requestParameters, initOverrides);
        return await response.value();
    }

}
