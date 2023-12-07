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
 * @interface AuthUserResponse
 */
export interface AuthUserResponse {
    /**
     * 購入者ID
     * @type {string}
     * @memberof AuthUserResponse
     */
    id: string;
    /**
     * 表示名
     * @type {string}
     * @memberof AuthUserResponse
     */
    username: string;
    /**
     * 姓(16文字まで)
     * @type {string}
     * @memberof AuthUserResponse
     */
    lastname: string;
    /**
     * 名(16文字まで)
     * @type {string}
     * @memberof AuthUserResponse
     */
    firstname: string;
    /**
     * 姓(かな)(ひらがな,32文字まで)
     * @type {string}
     * @memberof AuthUserResponse
     */
    lastnameKana: string;
    /**
     * 名(かな)(ひらがな,32文字まで)
     * @type {string}
     * @memberof AuthUserResponse
     */
    firstnameKana: string;
    /**
     * メールアドレス
     * @type {string}
     * @memberof AuthUserResponse
     */
    email: string;
    /**
     * サムネイルURL
     * @type {string}
     * @memberof AuthUserResponse
     */
    thumbnailUrl: string;
}

/**
 * Check if a given object implements the AuthUserResponse interface.
 */
export function instanceOfAuthUserResponse(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "username" in value;
    isInstance = isInstance && "lastname" in value;
    isInstance = isInstance && "firstname" in value;
    isInstance = isInstance && "lastnameKana" in value;
    isInstance = isInstance && "firstnameKana" in value;
    isInstance = isInstance && "email" in value;
    isInstance = isInstance && "thumbnailUrl" in value;

    return isInstance;
}

export function AuthUserResponseFromJSON(json: any): AuthUserResponse {
    return AuthUserResponseFromJSONTyped(json, false);
}

export function AuthUserResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): AuthUserResponse {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'username': json['username'],
        'lastname': json['lastname'],
        'firstname': json['firstname'],
        'lastnameKana': json['lastnameKana'],
        'firstnameKana': json['firstnameKana'],
        'email': json['email'],
        'thumbnailUrl': json['thumbnailUrl'],
    };
}

export function AuthUserResponseToJSON(value?: AuthUserResponse | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'username': value.username,
        'lastname': value.lastname,
        'firstname': value.firstname,
        'lastnameKana': value.lastnameKana,
        'firstnameKana': value.firstnameKana,
        'email': value.email,
        'thumbnailUrl': value.thumbnailUrl,
    };
}

