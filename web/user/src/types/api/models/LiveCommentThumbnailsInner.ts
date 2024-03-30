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
 * @interface LiveCommentThumbnailsInner
 */
export interface LiveCommentThumbnailsInner {
    /**
     * リサイズ済みサムネイルURL
     * @type {string}
     * @memberof LiveCommentThumbnailsInner
     */
    url: string;
    /**
     * 
     * @type {ImageSize}
     * @memberof LiveCommentThumbnailsInner
     */
    size: ImageSize;
}

/**
 * Check if a given object implements the LiveCommentThumbnailsInner interface.
 */
export function instanceOfLiveCommentThumbnailsInner(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "url" in value;
    isInstance = isInstance && "size" in value;

    return isInstance;
}

export function LiveCommentThumbnailsInnerFromJSON(json: any): LiveCommentThumbnailsInner {
    return LiveCommentThumbnailsInnerFromJSONTyped(json, false);
}

export function LiveCommentThumbnailsInnerFromJSONTyped(json: any, ignoreDiscriminator: boolean): LiveCommentThumbnailsInner {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'url': json['url'],
        'size': ImageSizeFromJSON(json['size']),
    };
}

export function LiveCommentThumbnailsInnerToJSON(value?: LiveCommentThumbnailsInner | null): any {
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
