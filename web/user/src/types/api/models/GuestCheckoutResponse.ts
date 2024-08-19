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
 * @interface GuestCheckoutResponse
 */
export interface GuestCheckoutResponse {
    /**
     * 支払いページへのリダイレクトURL
     * @type {string}
     * @memberof GuestCheckoutResponse
     */
    url: string;
}

/**
 * Check if a given object implements the GuestCheckoutResponse interface.
 */
export function instanceOfGuestCheckoutResponse(value: object): value is GuestCheckoutResponse {
    if (!('url' in value) || value['url'] === undefined) return false;
    return true;
}

export function GuestCheckoutResponseFromJSON(json: any): GuestCheckoutResponse {
    return GuestCheckoutResponseFromJSONTyped(json, false);
}

export function GuestCheckoutResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): GuestCheckoutResponse {
    if (json == null) {
        return json;
    }
    return {
        
        'url': json['url'],
    };
}

export function GuestCheckoutResponseToJSON(value?: GuestCheckoutResponse | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'url': value['url'],
    };
}

