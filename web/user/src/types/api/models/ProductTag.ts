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
/**
 * 商品タグ情報
 * @export
 * @interface ProductTag
 */
export interface ProductTag {
    /**
     * 商品タグID
     * @type {string}
     * @memberof ProductTag
     */
    id: string;
    /**
     * 商品タグ名(32文字まで)
     * @type {string}
     * @memberof ProductTag
     */
    name: string;
}

/**
 * Check if a given object implements the ProductTag interface.
 */
export function instanceOfProductTag(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "name" in value;

    return isInstance;
}

export function ProductTagFromJSON(json: any): ProductTag {
    return ProductTagFromJSONTyped(json, false);
}

export function ProductTagFromJSONTyped(json: any, ignoreDiscriminator: boolean): ProductTag {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'name': json['name'],
    };
}

export function ProductTagToJSON(value?: ProductTag | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'name': value.name,
    };
}

