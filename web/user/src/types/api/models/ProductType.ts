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
import type { ProductTypeIconsInner } from './ProductTypeIconsInner';
import {
    ProductTypeIconsInnerFromJSON,
    ProductTypeIconsInnerFromJSONTyped,
    ProductTypeIconsInnerToJSON,
} from './ProductTypeIconsInner';

/**
 * 品目情報
 * @export
 * @interface ProductType
 */
export interface ProductType {
    /**
     * 品目ID
     * @type {string}
     * @memberof ProductType
     */
    id: string;
    /**
     * 品目名
     * @type {string}
     * @memberof ProductType
     */
    name: string;
    /**
     * アイコンURL
     * @type {string}
     * @memberof ProductType
     */
    iconUrl: string;
    /**
     * リサイズ済みアイコンURL一覧
     * @type {Array<ProductTypeIconsInner>}
     * @memberof ProductType
     */
    icons: Array<ProductTypeIconsInner>;
    /**
     * 商品種別ID
     * @type {string}
     * @memberof ProductType
     */
    categoryId: string;
}

/**
 * Check if a given object implements the ProductType interface.
 */
export function instanceOfProductType(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "name" in value;
    isInstance = isInstance && "iconUrl" in value;
    isInstance = isInstance && "icons" in value;
    isInstance = isInstance && "categoryId" in value;

    return isInstance;
}

export function ProductTypeFromJSON(json: any): ProductType {
    return ProductTypeFromJSONTyped(json, false);
}

export function ProductTypeFromJSONTyped(json: any, ignoreDiscriminator: boolean): ProductType {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'name': json['name'],
        'iconUrl': json['iconUrl'],
        'icons': ((json['icons'] as Array<any>).map(ProductTypeIconsInnerFromJSON)),
        'categoryId': json['categoryId'],
    };
}

export function ProductTypeToJSON(value?: ProductType | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'name': value.name,
        'iconUrl': value.iconUrl,
        'icons': ((value.icons as Array<any>).map(ProductTypeIconsInnerToJSON)),
        'categoryId': value.categoryId,
    };
}

