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
import type { ImageSize } from './ImageSize';
import {
    ImageSizeFromJSON,
    ImageSizeFromJSONTyped,
    ImageSizeToJSON,
} from './ImageSize';

/**
 * 
 * @export
 * @interface ProductTypesResponseProductTypesInnerIconsInner
 */
export interface ProductTypesResponseProductTypesInnerIconsInner {
    /**
     * リサイズ済みアイコンURL
     * @type {string}
     * @memberof ProductTypesResponseProductTypesInnerIconsInner
     */
    url: string;
    /**
     * 
     * @type {ImageSize}
     * @memberof ProductTypesResponseProductTypesInnerIconsInner
     */
    size: ImageSize;
}

/**
 * Check if a given object implements the ProductTypesResponseProductTypesInnerIconsInner interface.
 */
export function instanceOfProductTypesResponseProductTypesInnerIconsInner(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "url" in value;
    isInstance = isInstance && "size" in value;

    return isInstance;
}

export function ProductTypesResponseProductTypesInnerIconsInnerFromJSON(json: any): ProductTypesResponseProductTypesInnerIconsInner {
    return ProductTypesResponseProductTypesInnerIconsInnerFromJSONTyped(json, false);
}

export function ProductTypesResponseProductTypesInnerIconsInnerFromJSONTyped(json: any, ignoreDiscriminator: boolean): ProductTypesResponseProductTypesInnerIconsInner {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'url': json['url'],
        'size': ImageSizeFromJSON(json['size']),
    };
}

export function ProductTypesResponseProductTypesInnerIconsInnerToJSON(value?: ProductTypesResponseProductTypesInnerIconsInner | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'url': value.url,
        'size': ImageSizeToJSON(value.size),
    };
}

