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
import type { Weekday } from './Weekday';
import {
    WeekdayFromJSON,
    WeekdayFromJSONTyped,
    WeekdayToJSON,
} from './Weekday';

/**
 * コーディネータ情報
 * @export
 * @interface Coordinator
 */
export interface Coordinator {
    /**
     * コーディネータID
     * @type {string}
     * @memberof Coordinator
     */
    id: string;
    /**
     * マルシェ名
     * @type {string}
     * @memberof Coordinator
     */
    marcheName: string;
    /**
     * コーディネータ名
     * @type {string}
     * @memberof Coordinator
     */
    username: string;
    /**
     * プロフィール
     * @type {string}
     * @memberof Coordinator
     */
    profile: string;
    /**
     * 取り扱い品目ID一覧
     * @type {Array<string>}
     * @memberof Coordinator
     */
    productTypeIds: Array<string>;
    /**
     * 営業曜日
     * @type {Array<Weekday>}
     * @memberof Coordinator
     */
    businessDays: Array<Weekday>;
    /**
     * サムネイルURL
     * @type {string}
     * @memberof Coordinator
     */
    thumbnailUrl: string;
    /**
     * ヘッダー画像URL
     * @type {string}
     * @memberof Coordinator
     */
    headerUrl: string;
    /**
     * 紹介動画URL
     * @type {string}
     * @memberof Coordinator
     */
    promotionVideoUrl: string;
    /**
     * Instagramアカウント
     * @type {string}
     * @memberof Coordinator
     */
    instagramId: string;
    /**
     * Facebookアカウント
     * @type {string}
     * @memberof Coordinator
     */
    facebookId: string;
    /**
     * 都道府県
     * @type {string}
     * @memberof Coordinator
     */
    prefecture: string;
    /**
     * 市区町村
     * @type {string}
     * @memberof Coordinator
     */
    city: string;
}

/**
 * Check if a given object implements the Coordinator interface.
 */
export function instanceOfCoordinator(value: object): value is Coordinator {
    if (!('id' in value) || value['id'] === undefined) return false;
    if (!('marcheName' in value) || value['marcheName'] === undefined) return false;
    if (!('username' in value) || value['username'] === undefined) return false;
    if (!('profile' in value) || value['profile'] === undefined) return false;
    if (!('productTypeIds' in value) || value['productTypeIds'] === undefined) return false;
    if (!('businessDays' in value) || value['businessDays'] === undefined) return false;
    if (!('thumbnailUrl' in value) || value['thumbnailUrl'] === undefined) return false;
    if (!('headerUrl' in value) || value['headerUrl'] === undefined) return false;
    if (!('promotionVideoUrl' in value) || value['promotionVideoUrl'] === undefined) return false;
    if (!('instagramId' in value) || value['instagramId'] === undefined) return false;
    if (!('facebookId' in value) || value['facebookId'] === undefined) return false;
    if (!('prefecture' in value) || value['prefecture'] === undefined) return false;
    if (!('city' in value) || value['city'] === undefined) return false;
    return true;
}

export function CoordinatorFromJSON(json: any): Coordinator {
    return CoordinatorFromJSONTyped(json, false);
}

export function CoordinatorFromJSONTyped(json: any, ignoreDiscriminator: boolean): Coordinator {
    if (json == null) {
        return json;
    }
    return {
        
        'id': json['id'],
        'marcheName': json['marcheName'],
        'username': json['username'],
        'profile': json['profile'],
        'productTypeIds': json['productTypeIds'],
        'businessDays': ((json['businessDays'] as Array<any>).map(WeekdayFromJSON)),
        'thumbnailUrl': json['thumbnailUrl'],
        'headerUrl': json['headerUrl'],
        'promotionVideoUrl': json['promotionVideoUrl'],
        'instagramId': json['instagramId'],
        'facebookId': json['facebookId'],
        'prefecture': json['prefecture'],
        'city': json['city'],
    };
}

export function CoordinatorToJSON(value?: Coordinator | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'id': value['id'],
        'marcheName': value['marcheName'],
        'username': value['username'],
        'profile': value['profile'],
        'productTypeIds': value['productTypeIds'],
        'businessDays': ((value['businessDays'] as Array<any>).map(WeekdayToJSON)),
        'thumbnailUrl': value['thumbnailUrl'],
        'headerUrl': value['headerUrl'],
        'promotionVideoUrl': value['promotionVideoUrl'],
        'instagramId': value['instagramId'],
        'facebookId': value['facebookId'],
        'prefecture': value['prefecture'],
        'city': value['city'],
    };
}

