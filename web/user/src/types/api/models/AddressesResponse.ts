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
import type { Address } from './Address';
import {
    AddressFromJSON,
    AddressFromJSONTyped,
    AddressToJSON,
} from './Address';

/**
 * 
 * @export
 * @interface AddressesResponse
 */
export interface AddressesResponse {
    /**
     * アドレス一覧
     * @type {Array<Address>}
     * @memberof AddressesResponse
     */
    addresses: Array<Address>;
    /**
     * 合計数
     * @type {number}
     * @memberof AddressesResponse
     */
    total: number;
}

/**
 * Check if a given object implements the AddressesResponse interface.
 */
export function instanceOfAddressesResponse(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "addresses" in value;
    isInstance = isInstance && "total" in value;

    return isInstance;
}

export function AddressesResponseFromJSON(json: any): AddressesResponse {
    return AddressesResponseFromJSONTyped(json, false);
}

export function AddressesResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): AddressesResponse {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'addresses': ((json['addresses'] as Array<any>).map(AddressFromJSON)),
        'total': json['total'],
    };
}

export function AddressesResponseToJSON(value?: AddressesResponse | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'addresses': ((value.addresses as Array<any>).map(AddressToJSON)),
        'total': value.total,
    };
}

