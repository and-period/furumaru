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
 * @interface CreateAuthRequest
 */
export interface CreateAuthRequest {
    /**
     * メールアドレス
     * @type {string}
     * @memberof CreateAuthRequest
     */
    email: string;
    /**
     * 電話番号(国際番号 + 電話番号)
     * @type {string}
     * @memberof CreateAuthRequest
     */
    phoneNumber: string;
    /**
     * パスワード(8~32文字, 英小文字,数字を少なくとも1文字ずつは含む)
     * @type {string}
     * @memberof CreateAuthRequest
     */
    password: string;
    /**
     * パスワード(確認用)
     * @type {string}
     * @memberof CreateAuthRequest
     */
    passwordConfirmation: string;
}

/**
 * Check if a given object implements the CreateAuthRequest interface.
 */
export function instanceOfCreateAuthRequest(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "email" in value;
    isInstance = isInstance && "phoneNumber" in value;
    isInstance = isInstance && "password" in value;
    isInstance = isInstance && "passwordConfirmation" in value;

    return isInstance;
}

export function CreateAuthRequestFromJSON(json: any): CreateAuthRequest {
    return CreateAuthRequestFromJSONTyped(json, false);
}

export function CreateAuthRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): CreateAuthRequest {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'email': json['email'],
        'phoneNumber': json['phoneNumber'],
        'password': json['password'],
        'passwordConfirmation': json['passwordConfirmation'],
    };
}

export function CreateAuthRequestToJSON(value?: CreateAuthRequest | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'email': value.email,
        'phoneNumber': value.phoneNumber,
        'password': value.password,
        'passwordConfirmation': value.passwordConfirmation,
    };
}

