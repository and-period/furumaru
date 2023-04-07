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

import { exists, mapValues } from '../runtime';
import type { UsersResponseUsersInner } from './UsersResponseUsersInner';
import {
    UsersResponseUsersInnerFromJSON,
    UsersResponseUsersInnerFromJSONTyped,
    UsersResponseUsersInnerToJSON,
} from './UsersResponseUsersInner';

/**
 * 
 * @export
 * @interface UsersResponse
 */
export interface UsersResponse {
    /**
     * 購入者一覧
     * @type {Array<UsersResponseUsersInner>}
     * @memberof UsersResponse
     */
    users: Array<UsersResponseUsersInner>;
    /**
     * 合計数
     * @type {number}
     * @memberof UsersResponse
     */
    total: number;
}

/**
 * Check if a given object implements the UsersResponse interface.
 */
export function instanceOfUsersResponse(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "users" in value;
    isInstance = isInstance && "total" in value;

    return isInstance;
}

export function UsersResponseFromJSON(json: any): UsersResponse {
    return UsersResponseFromJSONTyped(json, false);
}

export function UsersResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): UsersResponse {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'users': ((json['users'] as Array<any>).map(UsersResponseUsersInnerFromJSON)),
        'total': json['total'],
    };
}

export function UsersResponseToJSON(value?: UsersResponse | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'users': ((value.users as Array<any>).map(UsersResponseUsersInnerToJSON)),
        'total': value.total,
    };
}

