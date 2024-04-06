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
 * @interface Thumbnail
 */
export interface Thumbnail {
    /**
     * リサイズ済みサムネイルURL
     * @type {string}
     * @memberof Thumbnail
     * @deprecated
     */
    url: string;
    /**
     * 
     * @type {ImageSize}
     * @memberof Thumbnail
     */
    size: ImageSize;
}

/**
 * Check if a given object implements the Thumbnail interface.
 */
export function instanceOfThumbnail(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "url" in value;
    isInstance = isInstance && "size" in value;

    return isInstance;
}

export function ThumbnailFromJSON(json: any): Thumbnail {
    return ThumbnailFromJSONTyped(json, false);
}

export function ThumbnailFromJSONTyped(json: any, ignoreDiscriminator: boolean): Thumbnail {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'url': json['url'],
        'size': ImageSizeFromJSON(json['size']),
    };
}

export function ThumbnailToJSON(value?: Thumbnail | null): any {
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

