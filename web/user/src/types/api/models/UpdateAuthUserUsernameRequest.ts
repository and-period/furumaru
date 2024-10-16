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
 * @interface UpdateAuthUserUsernameRequest
 */
export interface UpdateAuthUserUsernameRequest {
    /**
     * ユーザー名(表示用)(32文字まで)
     * @type {string}
     * @memberof UpdateAuthUserUsernameRequest
     */
    username: string;
}

/**
 * Check if a given object implements the UpdateAuthUserUsernameRequest interface.
 */
export function instanceOfUpdateAuthUserUsernameRequest(value: object): value is UpdateAuthUserUsernameRequest {
    if (!('username' in value) || value['username'] === undefined) return false;
    return true;
}

export function UpdateAuthUserUsernameRequestFromJSON(json: any): UpdateAuthUserUsernameRequest {
    return UpdateAuthUserUsernameRequestFromJSONTyped(json, false);
}

export function UpdateAuthUserUsernameRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): UpdateAuthUserUsernameRequest {
    if (json == null) {
        return json;
    }
    return {
        
        'username': json['username'],
    };
}

export function UpdateAuthUserUsernameRequestToJSON(value?: UpdateAuthUserUsernameRequest | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'username': value['username'],
    };
}

