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
import type { CoordinatorHeadersInner } from './CoordinatorHeadersInner';
import {
    CoordinatorHeadersInnerFromJSON,
    CoordinatorHeadersInnerFromJSONTyped,
    CoordinatorHeadersInnerToJSON,
} from './CoordinatorHeadersInner';
import type { Thumbnail } from './Thumbnail';
import {
    ThumbnailFromJSON,
    ThumbnailFromJSONTyped,
    ThumbnailToJSON,
} from './Thumbnail';
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
     * リサイズ済みサムネイルURL一覧
     * @type {Array<Thumbnail>}
     * @memberof Coordinator
     */
    thumbnails: Array<Thumbnail>;
    /**
     * ヘッダー画像URL
     * @type {string}
     * @memberof Coordinator
     */
    headerUrl: string;
    /**
     * リサイズ済みヘッダー画像URL一覧
     * @type {Array<CoordinatorHeadersInner>}
     * @memberof Coordinator
     */
    headers: Array<CoordinatorHeadersInner>;
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
export function instanceOfCoordinator(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "marcheName" in value;
    isInstance = isInstance && "username" in value;
    isInstance = isInstance && "profile" in value;
    isInstance = isInstance && "productTypeIds" in value;
    isInstance = isInstance && "businessDays" in value;
    isInstance = isInstance && "thumbnailUrl" in value;
    isInstance = isInstance && "thumbnails" in value;
    isInstance = isInstance && "headerUrl" in value;
    isInstance = isInstance && "headers" in value;
    isInstance = isInstance && "promotionVideoUrl" in value;
    isInstance = isInstance && "instagramId" in value;
    isInstance = isInstance && "facebookId" in value;
    isInstance = isInstance && "prefecture" in value;
    isInstance = isInstance && "city" in value;

    return isInstance;
}

export function CoordinatorFromJSON(json: any): Coordinator {
    return CoordinatorFromJSONTyped(json, false);
}

export function CoordinatorFromJSONTyped(json: any, ignoreDiscriminator: boolean): Coordinator {
    if ((json === undefined) || (json === null)) {
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
        'thumbnails': ((json['thumbnails'] as Array<any>).map(ThumbnailFromJSON)),
        'headerUrl': json['headerUrl'],
        'headers': ((json['headers'] as Array<any>).map(CoordinatorHeadersInnerFromJSON)),
        'promotionVideoUrl': json['promotionVideoUrl'],
        'instagramId': json['instagramId'],
        'facebookId': json['facebookId'],
        'prefecture': json['prefecture'],
        'city': json['city'],
    };
}

export function CoordinatorToJSON(value?: Coordinator | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'marcheName': value.marcheName,
        'username': value.username,
        'profile': value.profile,
        'productTypeIds': value.productTypeIds,
        'businessDays': ((value.businessDays as Array<any>).map(WeekdayToJSON)),
        'thumbnailUrl': value.thumbnailUrl,
        'thumbnails': ((value.thumbnails as Array<any>).map(ThumbnailToJSON)),
        'headerUrl': value.headerUrl,
        'headers': ((value.headers as Array<any>).map(CoordinatorHeadersInnerToJSON)),
        'promotionVideoUrl': value.promotionVideoUrl,
        'instagramId': value.instagramId,
        'facebookId': value.facebookId,
        'prefecture': value.prefecture,
        'city': value.city,
    };
}

