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
import type { ArchiveSummary } from './ArchiveSummary';
import {
    ArchiveSummaryFromJSON,
    ArchiveSummaryFromJSONTyped,
    ArchiveSummaryToJSON,
} from './ArchiveSummary';
import type { LiveSummary } from './LiveSummary';
import {
    LiveSummaryFromJSON,
    LiveSummaryFromJSONTyped,
    LiveSummaryToJSON,
} from './LiveSummary';
import type { Producer } from './Producer';
import {
    ProducerFromJSON,
    ProducerFromJSONTyped,
    ProducerToJSON,
} from './Producer';
import type { Product } from './Product';
import {
    ProductFromJSON,
    ProductFromJSONTyped,
    ProductToJSON,
} from './Product';

/**
 * 
 * @export
 * @interface ProducerResponse
 */
export interface ProducerResponse {
    /**
     * 
     * @type {Producer}
     * @memberof ProducerResponse
     */
    producer: Producer;
    /**
     * 開催中・開催予定のマルシェ一覧
     * @type {Array<LiveSummary>}
     * @memberof ProducerResponse
     */
    lives: Array<LiveSummary>;
    /**
     * 過去のマルシェ一覧
     * @type {Array<ArchiveSummary>}
     * @memberof ProducerResponse
     */
    archives: Array<ArchiveSummary>;
    /**
     * 商品一覧
     * @type {Array<Product>}
     * @memberof ProducerResponse
     */
    products: Array<Product>;
}

/**
 * Check if a given object implements the ProducerResponse interface.
 */
export function instanceOfProducerResponse(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "producer" in value;
    isInstance = isInstance && "lives" in value;
    isInstance = isInstance && "archives" in value;
    isInstance = isInstance && "products" in value;

    return isInstance;
}

export function ProducerResponseFromJSON(json: any): ProducerResponse {
    return ProducerResponseFromJSONTyped(json, false);
}

export function ProducerResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): ProducerResponse {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'producer': ProducerFromJSON(json['producer']),
        'lives': ((json['lives'] as Array<any>).map(LiveSummaryFromJSON)),
        'archives': ((json['archives'] as Array<any>).map(ArchiveSummaryFromJSON)),
        'products': ((json['products'] as Array<any>).map(ProductFromJSON)),
    };
}

export function ProducerResponseToJSON(value?: ProducerResponse | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'producer': ProducerToJSON(value.producer),
        'lives': ((value.lives as Array<any>).map(LiveSummaryToJSON)),
        'archives': ((value.archives as Array<any>).map(ArchiveSummaryToJSON)),
        'products': ((value.products as Array<any>).map(ProductToJSON)),
    };
}

