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
/**
 * 生産者情報
 * @export
 * @interface Producer
 */
export interface Producer {
    /**
     * 生産者ID
     * @type {string}
     * @memberof Producer
     */
    id: string;
    /**
     * 担当コーディネータID
     * @type {string}
     * @memberof Producer
     */
    coordinatorId: string;
    /**
     * 生産者名
     * @type {string}
     * @memberof Producer
     */
    username: string;
    /**
     * プロフィール
     * @type {string}
     * @memberof Producer
     */
    profile: string;
    /**
     * サムネイルURL
     * @type {string}
     * @memberof Producer
     */
    thumbnailUrl: string;
    /**
     * ヘッダー画像URL
     * @type {string}
     * @memberof Producer
     */
    headerUrl: string;
    /**
     * 紹介動画URL
     * @type {string}
     * @memberof Producer
     */
    promotionVideoUrl: string;
    /**
     * Instagramアカウント
     * @type {string}
     * @memberof Producer
     */
    instagramId: string;
    /**
     * Facebookアカウント
     * @type {string}
     * @memberof Producer
     */
    facebookId: string;
    /**
     * 都道府県
     * @type {string}
     * @memberof Producer
     */
    prefecture: string;
    /**
     * 市区町村
     * @type {string}
     * @memberof Producer
     */
    city: string;
}

/**
 * Check if a given object implements the Producer interface.
 */
export function instanceOfProducer(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "coordinatorId" in value;
    isInstance = isInstance && "username" in value;
    isInstance = isInstance && "profile" in value;
    isInstance = isInstance && "thumbnailUrl" in value;
    isInstance = isInstance && "headerUrl" in value;
    isInstance = isInstance && "promotionVideoUrl" in value;
    isInstance = isInstance && "instagramId" in value;
    isInstance = isInstance && "facebookId" in value;
    isInstance = isInstance && "prefecture" in value;
    isInstance = isInstance && "city" in value;

    return isInstance;
}

export function ProducerFromJSON(json: any): Producer {
    return ProducerFromJSONTyped(json, false);
}

export function ProducerFromJSONTyped(json: any, ignoreDiscriminator: boolean): Producer {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'coordinatorId': json['coordinatorId'],
        'username': json['username'],
        'profile': json['profile'],
        'thumbnailUrl': json['thumbnailUrl'],
        'headerUrl': json['headerUrl'],
        'promotionVideoUrl': json['promotionVideoUrl'],
        'instagramId': json['instagramId'],
        'facebookId': json['facebookId'],
        'prefecture': json['prefecture'],
        'city': json['city'],
    };
}

export function ProducerToJSON(value?: Producer | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'coordinatorId': value.coordinatorId,
        'username': value.username,
        'profile': value.profile,
        'thumbnailUrl': value.thumbnailUrl,
        'headerUrl': value.headerUrl,
        'promotionVideoUrl': value.promotionVideoUrl,
        'instagramId': value.instagramId,
        'facebookId': value.facebookId,
        'prefecture': value.prefecture,
        'city': value.city,
    };
}

