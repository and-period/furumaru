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
 * @interface CoordinatorsResponseCoordinatorsInnerThumbnailsInner
 */
export interface CoordinatorsResponseCoordinatorsInnerThumbnailsInner {
    /**
     * リサイズ済みサムネイルURL
     * @type {string}
     * @memberof CoordinatorsResponseCoordinatorsInnerThumbnailsInner
     */
    url: string;
    /**
     * 
     * @type {ImageSize}
     * @memberof CoordinatorsResponseCoordinatorsInnerThumbnailsInner
     */
    size: ImageSize;
}

/**
 * Check if a given object implements the CoordinatorsResponseCoordinatorsInnerThumbnailsInner interface.
 */
export function instanceOfCoordinatorsResponseCoordinatorsInnerThumbnailsInner(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "url" in value;
    isInstance = isInstance && "size" in value;

    return isInstance;
}

export function CoordinatorsResponseCoordinatorsInnerThumbnailsInnerFromJSON(json: any): CoordinatorsResponseCoordinatorsInnerThumbnailsInner {
    return CoordinatorsResponseCoordinatorsInnerThumbnailsInnerFromJSONTyped(json, false);
}

export function CoordinatorsResponseCoordinatorsInnerThumbnailsInnerFromJSONTyped(json: any, ignoreDiscriminator: boolean): CoordinatorsResponseCoordinatorsInnerThumbnailsInner {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'url': json['url'],
        'size': ImageSizeFromJSON(json['size']),
    };
}

export function CoordinatorsResponseCoordinatorsInnerThumbnailsInnerToJSON(value?: CoordinatorsResponseCoordinatorsInnerThumbnailsInner | null): any {
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

