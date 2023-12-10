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
import type { CoordinatorThumbnailsInner } from './CoordinatorThumbnailsInner';
import {
    CoordinatorThumbnailsInnerFromJSON,
    CoordinatorThumbnailsInnerFromJSONTyped,
    CoordinatorThumbnailsInnerToJSON,
} from './CoordinatorThumbnailsInner';

/**
 * 
 * @export
 * @interface TopLiveProductsInner
 */
export interface TopLiveProductsInner {
    /**
     * 商品ID
     * @type {string}
     * @memberof TopLiveProductsInner
     */
    id: string;
    /**
     * 商品名
     * @type {string}
     * @memberof TopLiveProductsInner
     */
    name: string;
    /**
     * 販売価格
     * @type {number}
     * @memberof TopLiveProductsInner
     */
    price: number;
    /**
     * 在庫数
     * @type {number}
     * @memberof TopLiveProductsInner
     */
    inventory: number;
    /**
     * サムネイルURL
     * @type {string}
     * @memberof TopLiveProductsInner
     */
    thumbnailUrl: string;
    /**
     * リサイズ済みサムネイルURL一覧
     * @type {Array<CoordinatorThumbnailsInner>}
     * @memberof TopLiveProductsInner
     */
    thumbnails: Array<CoordinatorThumbnailsInner>;
}

/**
 * Check if a given object implements the TopLiveProductsInner interface.
 */
export function instanceOfTopLiveProductsInner(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "name" in value;
    isInstance = isInstance && "price" in value;
    isInstance = isInstance && "inventory" in value;
    isInstance = isInstance && "thumbnailUrl" in value;
    isInstance = isInstance && "thumbnails" in value;

    return isInstance;
}

export function TopLiveProductsInnerFromJSON(json: any): TopLiveProductsInner {
    return TopLiveProductsInnerFromJSONTyped(json, false);
}

export function TopLiveProductsInnerFromJSONTyped(json: any, ignoreDiscriminator: boolean): TopLiveProductsInner {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'name': json['name'],
        'price': json['price'],
        'inventory': json['inventory'],
        'thumbnailUrl': json['thumbnailUrl'],
        'thumbnails': ((json['thumbnails'] as Array<any>).map(CoordinatorThumbnailsInnerFromJSON)),
    };
}

export function TopLiveProductsInnerToJSON(value?: TopLiveProductsInner | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'name': value.name,
        'price': value.price,
        'inventory': value.inventory,
        'thumbnailUrl': value.thumbnailUrl,
        'thumbnails': ((value.thumbnails as Array<any>).map(CoordinatorThumbnailsInnerToJSON)),
    };
}

