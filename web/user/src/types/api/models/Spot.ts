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
import type { SpotUserType } from './SpotUserType';
import {
    SpotUserTypeFromJSON,
    SpotUserTypeFromJSONTyped,
    SpotUserTypeToJSON,
} from './SpotUserType';

/**
 * スポット情報
 * @export
 * @interface Spot
 */
export interface Spot {
    /**
     * スポットID
     * @type {string}
     * @memberof Spot
     */
    id: string;
    /**
     * スポット種別ID
     * @type {string}
     * @memberof Spot
     */
    spotTypeId: string;
    /**
     * スポット名
     * @type {string}
     * @memberof Spot
     */
    name: string;
    /**
     * スポット説明
     * @type {string}
     * @memberof Spot
     */
    description: string;
    /**
     * スポットURL
     * @type {string}
     * @memberof Spot
     */
    thumbnailUrl: string;
    /**
     * 経度
     * @type {number}
     * @memberof Spot
     */
    longitude: number;
    /**
     * 緯度
     * @type {number}
     * @memberof Spot
     */
    latitude: number;
    /**
     * 
     * @type {SpotUserType}
     * @memberof Spot
     */
    userType: SpotUserType;
    /**
     * 投稿者ID
     * @type {string}
     * @memberof Spot
     */
    userId: string;
    /**
     * 投稿者名
     * @type {string}
     * @memberof Spot
     */
    userName: string;
    /**
     * 投稿者サムネイルURL
     * @type {string}
     * @memberof Spot
     */
    userThumbnailUrl: string;
    /**
     * 登録日時 (unixtime)
     * @type {number}
     * @memberof Spot
     */
    createdAt: number;
    /**
     * 更新日時 (unixtime)
     * @type {number}
     * @memberof Spot
     */
    updatedAt: number;
}



/**
 * Check if a given object implements the Spot interface.
 */
export function instanceOfSpot(value: object): value is Spot {
    if (!('id' in value) || value['id'] === undefined) return false;
    if (!('spotTypeId' in value) || value['spotTypeId'] === undefined) return false;
    if (!('name' in value) || value['name'] === undefined) return false;
    if (!('description' in value) || value['description'] === undefined) return false;
    if (!('thumbnailUrl' in value) || value['thumbnailUrl'] === undefined) return false;
    if (!('longitude' in value) || value['longitude'] === undefined) return false;
    if (!('latitude' in value) || value['latitude'] === undefined) return false;
    if (!('userType' in value) || value['userType'] === undefined) return false;
    if (!('userId' in value) || value['userId'] === undefined) return false;
    if (!('userName' in value) || value['userName'] === undefined) return false;
    if (!('userThumbnailUrl' in value) || value['userThumbnailUrl'] === undefined) return false;
    if (!('createdAt' in value) || value['createdAt'] === undefined) return false;
    if (!('updatedAt' in value) || value['updatedAt'] === undefined) return false;
    return true;
}

export function SpotFromJSON(json: any): Spot {
    return SpotFromJSONTyped(json, false);
}

export function SpotFromJSONTyped(json: any, ignoreDiscriminator: boolean): Spot {
    if (json == null) {
        return json;
    }
    return {
        
        'id': json['id'],
        'spotTypeId': json['spotTypeId'],
        'name': json['name'],
        'description': json['description'],
        'thumbnailUrl': json['thumbnailUrl'],
        'longitude': json['longitude'],
        'latitude': json['latitude'],
        'userType': SpotUserTypeFromJSON(json['userType']),
        'userId': json['userId'],
        'userName': json['userName'],
        'userThumbnailUrl': json['userThumbnailUrl'],
        'createdAt': json['createdAt'],
        'updatedAt': json['updatedAt'],
    };
}

export function SpotToJSON(value?: Spot | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'id': value['id'],
        'spotTypeId': value['spotTypeId'],
        'name': value['name'],
        'description': value['description'],
        'thumbnailUrl': value['thumbnailUrl'],
        'longitude': value['longitude'],
        'latitude': value['latitude'],
        'userType': SpotUserTypeToJSON(value['userType']),
        'userId': value['userId'],
        'userName': value['userName'],
        'userThumbnailUrl': value['userThumbnailUrl'],
        'createdAt': value['createdAt'],
        'updatedAt': value['updatedAt'],
    };
}

