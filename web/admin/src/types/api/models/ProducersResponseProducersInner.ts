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
import type { CoordinatorsResponseCoordinatorsInnerHeadersInner } from './CoordinatorsResponseCoordinatorsInnerHeadersInner';
import {
    CoordinatorsResponseCoordinatorsInnerHeadersInnerFromJSON,
    CoordinatorsResponseCoordinatorsInnerHeadersInnerFromJSONTyped,
    CoordinatorsResponseCoordinatorsInnerHeadersInnerToJSON,
} from './CoordinatorsResponseCoordinatorsInnerHeadersInner';
import type { CoordinatorsResponseCoordinatorsInnerThumbnailsInner } from './CoordinatorsResponseCoordinatorsInnerThumbnailsInner';
import {
    CoordinatorsResponseCoordinatorsInnerThumbnailsInnerFromJSON,
    CoordinatorsResponseCoordinatorsInnerThumbnailsInnerFromJSONTyped,
    CoordinatorsResponseCoordinatorsInnerThumbnailsInnerToJSON,
} from './CoordinatorsResponseCoordinatorsInnerThumbnailsInner';

/**
 * 
 * @export
 * @interface ProducersResponseProducersInner
 */
export interface ProducersResponseProducersInner {
    /**
     * システム管理者ID
     * @type {string}
     * @memberof ProducersResponseProducersInner
     */
    id: string;
    /**
     * 担当コーディネータID
     * @type {string}
     * @memberof ProducersResponseProducersInner
     */
    coordinatorId: string;
    /**
     * 姓
     * @type {string}
     * @memberof ProducersResponseProducersInner
     */
    lastname: string;
    /**
     * 名
     * @type {string}
     * @memberof ProducersResponseProducersInner
     */
    firstname: string;
    /**
     * 姓(かな)
     * @type {string}
     * @memberof ProducersResponseProducersInner
     */
    lastnameKana: string;
    /**
     * 名(かな)
     * @type {string}
     * @memberof ProducersResponseProducersInner
     */
    firstnameKana: string;
    /**
     * 店舗名
     * @type {string}
     * @memberof ProducersResponseProducersInner
     */
    storeName: string;
    /**
     * サムネイルURL
     * @type {string}
     * @memberof ProducersResponseProducersInner
     */
    thumbnailUrl: string;
    /**
     * リサイズ済みサムネイルURL一覧
     * @type {Array<CoordinatorsResponseCoordinatorsInnerThumbnailsInner>}
     * @memberof ProducersResponseProducersInner
     */
    thumbnails: Array<CoordinatorsResponseCoordinatorsInnerThumbnailsInner>;
    /**
     * ヘッダー画像URL
     * @type {string}
     * @memberof ProducersResponseProducersInner
     */
    headerUrl: string;
    /**
     * リサイズ済みヘッダー画像URL一覧
     * @type {Array<CoordinatorsResponseCoordinatorsInnerHeadersInner>}
     * @memberof ProducersResponseProducersInner
     */
    headers: Array<CoordinatorsResponseCoordinatorsInnerHeadersInner>;
    /**
     * メールアドレス
     * @type {string}
     * @memberof ProducersResponseProducersInner
     */
    email: string;
    /**
     * 電話番号 (国際番号 + 電話番号)
     * @type {string}
     * @memberof ProducersResponseProducersInner
     */
    phoneNumber: string;
    /**
     * 郵便番号
     * @type {string}
     * @memberof ProducersResponseProducersInner
     */
    postalCode: string;
    /**
     * 都道府県
     * @type {string}
     * @memberof ProducersResponseProducersInner
     */
    prefecture: string;
    /**
     * 市区町村
     * @type {string}
     * @memberof ProducersResponseProducersInner
     */
    city: string;
    /**
     * 町名・番地
     * @type {string}
     * @memberof ProducersResponseProducersInner
     */
    addressLine1: string;
    /**
     * ビル名・号室など
     * @type {string}
     * @memberof ProducersResponseProducersInner
     */
    addressLine2: string;
    /**
     * 登録日時 (unixtime)
     * @type {number}
     * @memberof ProducersResponseProducersInner
     */
    createdAt: number;
    /**
     * 更新日時 (unixtime)
     * @type {number}
     * @memberof ProducersResponseProducersInner
     */
    updatedAt: number;
}

/**
 * Check if a given object implements the ProducersResponseProducersInner interface.
 */
export function instanceOfProducersResponseProducersInner(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "coordinatorId" in value;
    isInstance = isInstance && "lastname" in value;
    isInstance = isInstance && "firstname" in value;
    isInstance = isInstance && "lastnameKana" in value;
    isInstance = isInstance && "firstnameKana" in value;
    isInstance = isInstance && "storeName" in value;
    isInstance = isInstance && "thumbnailUrl" in value;
    isInstance = isInstance && "thumbnails" in value;
    isInstance = isInstance && "headerUrl" in value;
    isInstance = isInstance && "headers" in value;
    isInstance = isInstance && "email" in value;
    isInstance = isInstance && "phoneNumber" in value;
    isInstance = isInstance && "postalCode" in value;
    isInstance = isInstance && "prefecture" in value;
    isInstance = isInstance && "city" in value;
    isInstance = isInstance && "addressLine1" in value;
    isInstance = isInstance && "addressLine2" in value;
    isInstance = isInstance && "createdAt" in value;
    isInstance = isInstance && "updatedAt" in value;

    return isInstance;
}

export function ProducersResponseProducersInnerFromJSON(json: any): ProducersResponseProducersInner {
    return ProducersResponseProducersInnerFromJSONTyped(json, false);
}

export function ProducersResponseProducersInnerFromJSONTyped(json: any, ignoreDiscriminator: boolean): ProducersResponseProducersInner {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'coordinatorId': json['coordinatorId'],
        'lastname': json['lastname'],
        'firstname': json['firstname'],
        'lastnameKana': json['lastnameKana'],
        'firstnameKana': json['firstnameKana'],
        'storeName': json['storeName'],
        'thumbnailUrl': json['thumbnailUrl'],
        'thumbnails': ((json['thumbnails'] as Array<any>).map(CoordinatorsResponseCoordinatorsInnerThumbnailsInnerFromJSON)),
        'headerUrl': json['headerUrl'],
        'headers': ((json['headers'] as Array<any>).map(CoordinatorsResponseCoordinatorsInnerHeadersInnerFromJSON)),
        'email': json['email'],
        'phoneNumber': json['phoneNumber'],
        'postalCode': json['postalCode'],
        'prefecture': json['prefecture'],
        'city': json['city'],
        'addressLine1': json['addressLine1'],
        'addressLine2': json['addressLine2'],
        'createdAt': json['createdAt'],
        'updatedAt': json['updatedAt'],
    };
}

export function ProducersResponseProducersInnerToJSON(value?: ProducersResponseProducersInner | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'coordinatorId': value.coordinatorId,
        'lastname': value.lastname,
        'firstname': value.firstname,
        'lastnameKana': value.lastnameKana,
        'firstnameKana': value.firstnameKana,
        'storeName': value.storeName,
        'thumbnailUrl': value.thumbnailUrl,
        'thumbnails': ((value.thumbnails as Array<any>).map(CoordinatorsResponseCoordinatorsInnerThumbnailsInnerToJSON)),
        'headerUrl': value.headerUrl,
        'headers': ((value.headers as Array<any>).map(CoordinatorsResponseCoordinatorsInnerHeadersInnerToJSON)),
        'email': value.email,
        'phoneNumber': value.phoneNumber,
        'postalCode': value.postalCode,
        'prefecture': value.prefecture,
        'city': value.city,
        'addressLine1': value.addressLine1,
        'addressLine2': value.addressLine2,
        'createdAt': value.createdAt,
        'updatedAt': value.updatedAt,
    };
}

