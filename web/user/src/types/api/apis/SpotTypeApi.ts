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
  SpotTypesResponse,
} from '../models/index';
import {
    SpotTypesResponseFromJSON,
    SpotTypesResponseToJSON,
} from '../models/index';

export interface V1ListSpotTypesRequest {
    limit?: number;
    offset?: number;
    name?: string;
}

/**
 * 
 */
export class SpotTypeApi extends runtime.BaseAPI {

    /**
     * スポット種別一覧取得
     */
    async v1ListSpotTypesRaw(requestParameters: V1ListSpotTypesRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<SpotTypesResponse>> {
        const queryParameters: any = {};

        if (requestParameters['limit'] != null) {
            queryParameters['limit'] = requestParameters['limit'];
        }

        if (requestParameters['offset'] != null) {
            queryParameters['offset'] = requestParameters['offset'];
        }

        if (requestParameters['name'] != null) {
            queryParameters['name'] = requestParameters['name'];
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
            path: `/v1/spot-types`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => SpotTypesResponseFromJSON(jsonValue));
    }

    /**
     * スポット種別一覧取得
     */
    async v1ListSpotTypes(requestParameters: V1ListSpotTypesRequest = {}, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<SpotTypesResponse> {
        const response = await this.v1ListSpotTypesRaw(requestParameters, initOverrides);
        return await response.value();
    }

}
