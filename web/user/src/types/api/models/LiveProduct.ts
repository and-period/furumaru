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
import type { Thumbnail } from './Thumbnail';
import {
    ThumbnailFromJSON,
    ThumbnailFromJSONTyped,
    ThumbnailToJSON,
} from './Thumbnail';

/**
 * マルシェに関連づく商品情報
 * @export
 * @interface LiveProduct
 */
export interface LiveProduct {
    /**
     * 商品ID
     * @type {string}
     * @memberof LiveProduct
     */
    id: string;
    /**
     * 商品名
     * @type {string}
     * @memberof LiveProduct
     */
    name: string;
    /**
     * 販売価格
     * @type {number}
     * @memberof LiveProduct
     */
    price: number;
    /**
     * 在庫数
     * @type {number}
     * @memberof LiveProduct
     */
    inventory: number;
    /**
     * サムネイルURL
     * @type {string}
     * @memberof LiveProduct
     */
    thumbnailUrl: string;
    /**
     * リサイズ済みサムネイルURL一覧
     * @type {Array<Thumbnail>}
     * @memberof LiveProduct
     */
    thumbnails: Array<Thumbnail>;
}

/**
 * Check if a given object implements the LiveProduct interface.
 */
export function instanceOfLiveProduct(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "name" in value;
    isInstance = isInstance && "price" in value;
    isInstance = isInstance && "inventory" in value;
    isInstance = isInstance && "thumbnailUrl" in value;
    isInstance = isInstance && "thumbnails" in value;

    return isInstance;
}

export function LiveProductFromJSON(json: any): LiveProduct {
    return LiveProductFromJSONTyped(json, false);
}

export function LiveProductFromJSONTyped(json: any, ignoreDiscriminator: boolean): LiveProduct {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'name': json['name'],
        'price': json['price'],
        'inventory': json['inventory'],
        'thumbnailUrl': json['thumbnailUrl'],
        'thumbnails': ((json['thumbnails'] as Array<any>).map(ThumbnailFromJSON)),
    };
}

export function LiveProductToJSON(value?: LiveProduct | null): any {
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
        'thumbnails': ((value.thumbnails as Array<any>).map(ThumbnailToJSON)),
    };
}
