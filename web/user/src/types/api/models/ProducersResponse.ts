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
import type { Producer } from './Producer';
import {
    ProducerFromJSON,
    ProducerFromJSONTyped,
    ProducerToJSON,
} from './Producer';

/**
 * 
 * @export
 * @interface ProducersResponse
 */
export interface ProducersResponse {
    /**
     * 生産者一覧
     * @type {Array<Producer>}
     * @memberof ProducersResponse
     */
    producers: Array<Producer>;
    /**
     * 合計数
     * @type {number}
     * @memberof ProducersResponse
     */
    total: number;
}

/**
 * Check if a given object implements the ProducersResponse interface.
 */
export function instanceOfProducersResponse(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "producers" in value;
    isInstance = isInstance && "total" in value;

    return isInstance;
}

export function ProducersResponseFromJSON(json: any): ProducersResponse {
    return ProducersResponseFromJSONTyped(json, false);
}

export function ProducersResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): ProducersResponse {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'producers': ((json['producers'] as Array<any>).map(ProducerFromJSON)),
        'total': json['total'],
    };
}

export function ProducersResponseToJSON(value?: ProducersResponse | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'producers': ((value.producers as Array<any>).map(ProducerToJSON)),
        'total': value.total,
    };
}

