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
 * @interface CreateContactRequest
 */
export interface CreateContactRequest {
    /**
     * 件名(64文字まで)
     * @type {string}
     * @memberof CreateContactRequest
     */
    title: string;
    /**
     * 内容(2000文字まで)
     * @type {string}
     * @memberof CreateContactRequest
     */
    content: string;
    /**
     * 氏名(64文字)
     * @type {string}
     * @memberof CreateContactRequest
     */
    username: string;
    /**
     * メールアドレス
     * @type {string}
     * @memberof CreateContactRequest
     */
    email: string;
    /**
     * 電話番号 (国際番号 + 電話番号)
     * @type {string}
     * @memberof CreateContactRequest
     */
    phoneNumber: string;
}

/**
 * Check if a given object implements the CreateContactRequest interface.
 */
export function instanceOfCreateContactRequest(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "title" in value;
    isInstance = isInstance && "content" in value;
    isInstance = isInstance && "username" in value;
    isInstance = isInstance && "email" in value;
    isInstance = isInstance && "phoneNumber" in value;

    return isInstance;
}

export function CreateContactRequestFromJSON(json: any): CreateContactRequest {
    return CreateContactRequestFromJSONTyped(json, false);
}

export function CreateContactRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): CreateContactRequest {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'title': json['title'],
        'content': json['content'],
        'username': json['username'],
        'email': json['email'],
        'phoneNumber': json['phoneNumber'],
    };
}

export function CreateContactRequestToJSON(value?: CreateContactRequest | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'title': value.title,
        'content': value.content,
        'username': value.username,
        'email': value.email,
        'phoneNumber': value.phoneNumber,
    };
}

