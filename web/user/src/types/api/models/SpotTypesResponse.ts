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
import type { SpotType } from './SpotType';
import {
    SpotTypeFromJSON,
    SpotTypeFromJSONTyped,
    SpotTypeToJSON,
} from './SpotType';

/**
 * 
 * @export
 * @interface SpotTypesResponse
 */
export interface SpotTypesResponse {
    /**
     * スポット種別一覧
     * @type {Array<SpotType>}
     * @memberof SpotTypesResponse
     */
    spotTypes: Array<SpotType>;
    /**
     * 合計数
     * @type {number}
     * @memberof SpotTypesResponse
     */
    total: number;
}

/**
 * Check if a given object implements the SpotTypesResponse interface.
 */
export function instanceOfSpotTypesResponse(value: object): value is SpotTypesResponse {
    if (!('spotTypes' in value) || value['spotTypes'] === undefined) return false;
    if (!('total' in value) || value['total'] === undefined) return false;
    return true;
}

export function SpotTypesResponseFromJSON(json: any): SpotTypesResponse {
    return SpotTypesResponseFromJSONTyped(json, false);
}

export function SpotTypesResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): SpotTypesResponse {
    if (json == null) {
        return json;
    }
    return {
        
        'spotTypes': ((json['spotTypes'] as Array<any>).map(SpotTypeFromJSON)),
        'total': json['total'],
    };
}

export function SpotTypesResponseToJSON(value?: SpotTypesResponse | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'spotTypes': ((value['spotTypes'] as Array<any>).map(SpotTypeToJSON)),
        'total': value['total'],
    };
}

