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

import { mapValues } from '../runtime';
/**
 * 
 * @export
 * @interface AuthGoogleAccountResponse
 */
export interface AuthGoogleAccountResponse {
    /**
     * Google認証用のURL
     * @type {string}
     * @memberof AuthGoogleAccountResponse
     */
    url: string;
}

/**
 * Check if a given object implements the AuthGoogleAccountResponse interface.
 */
export function instanceOfAuthGoogleAccountResponse(value: object): value is AuthGoogleAccountResponse {
    if (!('url' in value) || value['url'] === undefined) return false;
    return true;
}

export function AuthGoogleAccountResponseFromJSON(json: any): AuthGoogleAccountResponse {
    return AuthGoogleAccountResponseFromJSONTyped(json, false);
}

export function AuthGoogleAccountResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): AuthGoogleAccountResponse {
    if (json == null) {
        return json;
    }
    return {
        
        'url': json['url'],
    };
}

export function AuthGoogleAccountResponseToJSON(value?: AuthGoogleAccountResponse | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'url': value['url'],
    };
}

