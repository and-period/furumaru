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
import type { ImageSize } from './ImageSize';
import {
    ImageSizeFromJSON,
    ImageSizeFromJSONTyped,
    ImageSizeToJSON,
} from './ImageSize';

/**
 * 
 * @export
 * @interface ProductMediaInnerImagesInner
 */
export interface ProductMediaInnerImagesInner {
    /**
     * リサイズ済み画像URL
     * @type {string}
     * @memberof ProductMediaInnerImagesInner
     */
    url: string;
    /**
     * 
     * @type {ImageSize}
     * @memberof ProductMediaInnerImagesInner
     */
    size: ImageSize;
}

/**
 * Check if a given object implements the ProductMediaInnerImagesInner interface.
 */
export function instanceOfProductMediaInnerImagesInner(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "url" in value;
    isInstance = isInstance && "size" in value;

    return isInstance;
}

export function ProductMediaInnerImagesInnerFromJSON(json: any): ProductMediaInnerImagesInner {
    return ProductMediaInnerImagesInnerFromJSONTyped(json, false);
}

export function ProductMediaInnerImagesInnerFromJSONTyped(json: any, ignoreDiscriminator: boolean): ProductMediaInnerImagesInner {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'url': json['url'],
        'size': ImageSizeFromJSON(json['size']),
    };
}

export function ProductMediaInnerImagesInnerToJSON(value?: ProductMediaInnerImagesInner | null): any {
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

