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
import type { ShippingRate } from './ShippingRate';
import {
    ShippingRateFromJSON,
    ShippingRateFromJSONTyped,
    ShippingRateToJSON,
} from './ShippingRate';

/**
 * 
 * @export
 * @interface ShippingsResponseShippingsInner
 */
export interface ShippingsResponseShippingsInner {
    /**
     * 配送設定ID
     * @type {string}
     * @memberof ShippingsResponseShippingsInner
     */
    id: string;
    /**
     * 配送設定名
     * @type {string}
     * @memberof ShippingsResponseShippingsInner
     */
    name: string;
    /**
     * 箱サイズ60の通常配送料一覧
     * @type {Array<ShippingRate>}
     * @memberof ShippingsResponseShippingsInner
     */
    box60Rates: Array<ShippingRate>;
    /**
     * 箱サイズ60の冷蔵便追加配送料
     * @type {number}
     * @memberof ShippingsResponseShippingsInner
     */
    box60Refrigerated: number;
    /**
     * 箱サイズ60の冷凍便追加配送料
     * @type {number}
     * @memberof ShippingsResponseShippingsInner
     */
    box60Frozen: number;
    /**
     * 箱サイズ80の通常配送料一覧
     * @type {Array<ShippingRate>}
     * @memberof ShippingsResponseShippingsInner
     */
    box80Rates: Array<ShippingRate>;
    /**
     * 箱サイズ80の冷蔵便追加配送料
     * @type {number}
     * @memberof ShippingsResponseShippingsInner
     */
    box80Refrigerated: number;
    /**
     * 箱サイズ80の冷凍便追加配送料
     * @type {number}
     * @memberof ShippingsResponseShippingsInner
     */
    box80Frozen: number;
    /**
     * 箱サイズ100の通常配送料一覧
     * @type {Array<ShippingRate>}
     * @memberof ShippingsResponseShippingsInner
     */
    box100Rates: Array<ShippingRate>;
    /**
     * 箱サイズ100の冷蔵便追加配送料
     * @type {number}
     * @memberof ShippingsResponseShippingsInner
     */
    box100Refrigerated: number;
    /**
     * 箱サイズ100の冷凍便追加配送料
     * @type {number}
     * @memberof ShippingsResponseShippingsInner
     */
    box100Frozen: number;
    /**
     * 送料無料オプションの有無
     * @type {boolean}
     * @memberof ShippingsResponseShippingsInner
     */
    hasFreeShipping: boolean;
    /**
     * 送料無料になる金額
     * @type {number}
     * @memberof ShippingsResponseShippingsInner
     */
    freeShippingRates: number;
    /**
     * 登録日時 (unixtime)
     * @type {number}
     * @memberof ShippingsResponseShippingsInner
     */
    createdAt: number;
    /**
     * 更新日時 (unixtime)
     * @type {number}
     * @memberof ShippingsResponseShippingsInner
     */
    updatedAt: number;
}

/**
 * Check if a given object implements the ShippingsResponseShippingsInner interface.
 */
export function instanceOfShippingsResponseShippingsInner(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "name" in value;
    isInstance = isInstance && "box60Rates" in value;
    isInstance = isInstance && "box60Refrigerated" in value;
    isInstance = isInstance && "box60Frozen" in value;
    isInstance = isInstance && "box80Rates" in value;
    isInstance = isInstance && "box80Refrigerated" in value;
    isInstance = isInstance && "box80Frozen" in value;
    isInstance = isInstance && "box100Rates" in value;
    isInstance = isInstance && "box100Refrigerated" in value;
    isInstance = isInstance && "box100Frozen" in value;
    isInstance = isInstance && "hasFreeShipping" in value;
    isInstance = isInstance && "freeShippingRates" in value;
    isInstance = isInstance && "createdAt" in value;
    isInstance = isInstance && "updatedAt" in value;

    return isInstance;
}

export function ShippingsResponseShippingsInnerFromJSON(json: any): ShippingsResponseShippingsInner {
    return ShippingsResponseShippingsInnerFromJSONTyped(json, false);
}

export function ShippingsResponseShippingsInnerFromJSONTyped(json: any, ignoreDiscriminator: boolean): ShippingsResponseShippingsInner {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'name': json['name'],
        'box60Rates': ((json['box60Rates'] as Array<any>).map(ShippingRateFromJSON)),
        'box60Refrigerated': json['box60Refrigerated'],
        'box60Frozen': json['box60Frozen'],
        'box80Rates': ((json['box80Rates'] as Array<any>).map(ShippingRateFromJSON)),
        'box80Refrigerated': json['box80Refrigerated'],
        'box80Frozen': json['box80Frozen'],
        'box100Rates': ((json['box100Rates'] as Array<any>).map(ShippingRateFromJSON)),
        'box100Refrigerated': json['box100Refrigerated'],
        'box100Frozen': json['box100Frozen'],
        'hasFreeShipping': json['hasFreeShipping'],
        'freeShippingRates': json['freeShippingRates'],
        'createdAt': json['createdAt'],
        'updatedAt': json['updatedAt'],
    };
}

export function ShippingsResponseShippingsInnerToJSON(value?: ShippingsResponseShippingsInner | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'name': value.name,
        'box60Rates': ((value.box60Rates as Array<any>).map(ShippingRateToJSON)),
        'box60Refrigerated': value.box60Refrigerated,
        'box60Frozen': value.box60Frozen,
        'box80Rates': ((value.box80Rates as Array<any>).map(ShippingRateToJSON)),
        'box80Refrigerated': value.box80Refrigerated,
        'box80Frozen': value.box80Frozen,
        'box100Rates': ((value.box100Rates as Array<any>).map(ShippingRateToJSON)),
        'box100Refrigerated': value.box100Refrigerated,
        'box100Frozen': value.box100Frozen,
        'hasFreeShipping': value.hasFreeShipping,
        'freeShippingRates': value.freeShippingRates,
        'createdAt': value.createdAt,
        'updatedAt': value.updatedAt,
    };
}

