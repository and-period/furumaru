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

import { mapValues } from '../runtime';
import type { UploadStatus } from './UploadStatus';
import {
    UploadStatusFromJSON,
    UploadStatusFromJSONTyped,
    UploadStatusToJSON,
} from './UploadStatus';

/**
 * 
 * @export
 * @interface UploadStateResponse
 */
export interface UploadStateResponse {
    /**
     * 参照先ファイルURL
     * @type {string}
     * @memberof UploadStateResponse
     */
    url: string;
    /**
     * 
     * @type {UploadStatus}
     * @memberof UploadStateResponse
     */
    status: UploadStatus;
}



/**
 * Check if a given object implements the UploadStateResponse interface.
 */
export function instanceOfUploadStateResponse(value: object): value is UploadStateResponse {
    if (!('url' in value) || value['url'] === undefined) return false;
    if (!('status' in value) || value['status'] === undefined) return false;
    return true;
}

export function UploadStateResponseFromJSON(json: any): UploadStateResponse {
    return UploadStateResponseFromJSONTyped(json, false);
}

export function UploadStateResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): UploadStateResponse {
    if (json == null) {
        return json;
    }
    return {
        
        'url': json['url'],
        'status': UploadStatusFromJSON(json['status']),
    };
}

export function UploadStateResponseToJSON(value?: UploadStateResponse | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'url': value['url'],
        'status': UploadStatusToJSON(value['status']),
    };
}

