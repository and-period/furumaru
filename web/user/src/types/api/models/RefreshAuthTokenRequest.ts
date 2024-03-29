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

import { exists, mapValues } from '../runtime';
/**
 * 
 * @export
 * @interface RefreshAuthTokenRequest
 */
export interface RefreshAuthTokenRequest {
    /**
     * 更新トークン
     * @type {string}
     * @memberof RefreshAuthTokenRequest
     */
    refreshToken: string;
}

/**
 * Check if a given object implements the RefreshAuthTokenRequest interface.
 */
export function instanceOfRefreshAuthTokenRequest(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "refreshToken" in value;

    return isInstance;
}

export function RefreshAuthTokenRequestFromJSON(json: any): RefreshAuthTokenRequest {
    return RefreshAuthTokenRequestFromJSONTyped(json, false);
}

export function RefreshAuthTokenRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): RefreshAuthTokenRequest {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'refreshToken': json['refreshToken'],
    };
}

export function RefreshAuthTokenRequestToJSON(value?: RefreshAuthTokenRequest | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'refreshToken': value.refreshToken,
    };
}

