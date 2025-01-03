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
/**
 * 
 * @export
 * @interface UpdateSpotRequest
 */
export interface UpdateSpotRequest {
    /**
     * スポット種別ID
     * @type {string}
     * @memberof UpdateSpotRequest
     */
    spotTypeId: string;
    /**
     * スポット名（64文字まで）
     * @type {string}
     * @memberof UpdateSpotRequest
     */
    name: string;
    /**
     * 説明（2000文字まで）
     * @type {string}
     * @memberof UpdateSpotRequest
     */
    description: string;
    /**
     * サムネイルURL
     * @type {string}
     * @memberof UpdateSpotRequest
     */
    thumbnailUrl: string;
    /**
     * 緯度
     * @type {number}
     * @memberof UpdateSpotRequest
     */
    latitude: number;
    /**
     * 経度
     * @type {number}
     * @memberof UpdateSpotRequest
     */
    longitude: number;
}

/**
 * Check if a given object implements the UpdateSpotRequest interface.
 */
export function instanceOfUpdateSpotRequest(value: object): value is UpdateSpotRequest {
    if (!('spotTypeId' in value) || value['spotTypeId'] === undefined) return false;
    if (!('name' in value) || value['name'] === undefined) return false;
    if (!('description' in value) || value['description'] === undefined) return false;
    if (!('thumbnailUrl' in value) || value['thumbnailUrl'] === undefined) return false;
    if (!('latitude' in value) || value['latitude'] === undefined) return false;
    if (!('longitude' in value) || value['longitude'] === undefined) return false;
    return true;
}

export function UpdateSpotRequestFromJSON(json: any): UpdateSpotRequest {
    return UpdateSpotRequestFromJSONTyped(json, false);
}

export function UpdateSpotRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): UpdateSpotRequest {
    if (json == null) {
        return json;
    }
    return {
        
        'spotTypeId': json['spotTypeId'],
        'name': json['name'],
        'description': json['description'],
        'thumbnailUrl': json['thumbnailUrl'],
        'latitude': json['latitude'],
        'longitude': json['longitude'],
    };
}

export function UpdateSpotRequestToJSON(value?: UpdateSpotRequest | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'spotTypeId': value['spotTypeId'],
        'name': value['name'],
        'description': value['description'],
        'thumbnailUrl': value['thumbnailUrl'],
        'latitude': value['latitude'],
        'longitude': value['longitude'],
    };
}

