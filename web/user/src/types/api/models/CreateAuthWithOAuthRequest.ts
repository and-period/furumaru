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
 * @interface CreateAuthWithOAuthRequest
 */
export interface CreateAuthWithOAuthRequest {
    /**
     * ユーザー名(表示用)(32文字まで)
     * @type {string}
     * @memberof CreateAuthWithOAuthRequest
     */
    username: string;
    /**
     * ユーザーID(検索用)(32文字まで)
     * @type {string}
     * @memberof CreateAuthWithOAuthRequest
     */
    accountId: string;
    /**
     * 姓(16文字まで)
     * @type {string}
     * @memberof CreateAuthWithOAuthRequest
     */
    lastname: string;
    /**
     * 名(16文字まで)
     * @type {string}
     * @memberof CreateAuthWithOAuthRequest
     */
    firstname: string;
    /**
     * 姓(かな)(ひらがな,32文字まで)
     * @type {string}
     * @memberof CreateAuthWithOAuthRequest
     */
    lastnameKana: string;
    /**
     * 名(かな)(ひらがな,32文字まで)
     * @type {string}
     * @memberof CreateAuthWithOAuthRequest
     */
    firstnameKana: string;
    /**
     * 電話番号(国際番号 + 電話番号)
     * @type {string}
     * @memberof CreateAuthWithOAuthRequest
     */
    phoneNumber: string;
}

/**
 * Check if a given object implements the CreateAuthWithOAuthRequest interface.
 */
export function instanceOfCreateAuthWithOAuthRequest(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "username" in value;
    isInstance = isInstance && "accountId" in value;
    isInstance = isInstance && "lastname" in value;
    isInstance = isInstance && "firstname" in value;
    isInstance = isInstance && "lastnameKana" in value;
    isInstance = isInstance && "firstnameKana" in value;
    isInstance = isInstance && "phoneNumber" in value;

    return isInstance;
}

export function CreateAuthWithOAuthRequestFromJSON(json: any): CreateAuthWithOAuthRequest {
    return CreateAuthWithOAuthRequestFromJSONTyped(json, false);
}

export function CreateAuthWithOAuthRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): CreateAuthWithOAuthRequest {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'username': json['username'],
        'accountId': json['accountId'],
        'lastname': json['lastname'],
        'firstname': json['firstname'],
        'lastnameKana': json['lastnameKana'],
        'firstnameKana': json['firstnameKana'],
        'phoneNumber': json['phoneNumber'],
    };
}

export function CreateAuthWithOAuthRequestToJSON(value?: CreateAuthWithOAuthRequest | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'username': value.username,
        'accountId': value.accountId,
        'lastname': value.lastname,
        'firstname': value.firstname,
        'lastnameKana': value.lastnameKana,
        'firstnameKana': value.firstnameKana,
        'phoneNumber': value.phoneNumber,
    };
}

