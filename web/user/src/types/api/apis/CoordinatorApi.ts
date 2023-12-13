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
  CoordinatorResponse,
  ErrorResponse,
} from '../models/index';
import {
    CoordinatorResponseFromJSON,
    CoordinatorResponseToJSON,
    ErrorResponseFromJSON,
    ErrorResponseToJSON,
} from '../models/index';

export interface V1GetCoordinatorRequest {
    coordinatorId: string;
}

/**
 * 
 */
export class CoordinatorApi extends runtime.BaseAPI {

    /**
     * コーディネータ情報取得
     */
    async v1GetCoordinatorRaw(requestParameters: V1GetCoordinatorRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<CoordinatorResponse>> {
        if (requestParameters.coordinatorId === null || requestParameters.coordinatorId === undefined) {
            throw new runtime.RequiredError('coordinatorId','Required parameter requestParameters.coordinatorId was null or undefined when calling v1GetCoordinator.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/v1/coordinators/{coordinatorId}`.replace(`{${"coordinatorId"}}`, encodeURIComponent(String(requestParameters.coordinatorId))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => CoordinatorResponseFromJSON(jsonValue));
    }

    /**
     * コーディネータ情報取得
     */
    async v1GetCoordinator(requestParameters: V1GetCoordinatorRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<CoordinatorResponse> {
        const response = await this.v1GetCoordinatorRaw(requestParameters, initOverrides);
        return await response.value();
    }

}
