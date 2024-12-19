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
import type { StorageMethodType } from './StorageMethodType';
import {
    StorageMethodTypeFromJSON,
    StorageMethodTypeFromJSONTyped,
    StorageMethodTypeToJSON,
} from './StorageMethodType';
import type { ExperienceMediaInner } from './ExperienceMediaInner';
import {
    ExperienceMediaInnerFromJSON,
    ExperienceMediaInnerFromJSONTyped,
    ExperienceMediaInnerToJSON,
} from './ExperienceMediaInner';
import type { ProductStatus } from './ProductStatus';
import {
    ProductStatusFromJSON,
    ProductStatusFromJSONTyped,
    ProductStatusToJSON,
} from './ProductStatus';
import type { DeliveryType } from './DeliveryType';
import {
    DeliveryTypeFromJSON,
    DeliveryTypeFromJSONTyped,
    DeliveryTypeToJSON,
} from './DeliveryType';
import type { ProductRate } from './ProductRate';
import {
    ProductRateFromJSON,
    ProductRateFromJSONTyped,
    ProductRateToJSON,
} from './ProductRate';

/**
 * 商品情報
 * @export
 * @interface Product
 */
export interface Product {
    /**
     * 商品ID
     * @type {string}
     * @memberof Product
     */
    id: string;
    /**
     * 商品名
     * @type {string}
     * @memberof Product
     */
    name: string;
    /**
     * 商品説明
     * @type {string}
     * @memberof Product
     */
    description: string;
    /**
     * 
     * @type {ProductStatus}
     * @memberof Product
     */
    status: ProductStatus;
    /**
     * コーディネータID
     * @type {string}
     * @memberof Product
     */
    coordinatorId: string;
    /**
     * 生産者ID
     * @type {string}
     * @memberof Product
     */
    producerId: string;
    /**
     * 商品種別ID
     * @type {string}
     * @memberof Product
     */
    categoryId: string;
    /**
     * 品目ID
     * @type {string}
     * @memberof Product
     */
    productTypeId: string;
    /**
     * 商品タグ一覧
     * @type {Array<string>}
     * @memberof Product
     */
    productTagIds: Array<string>;
    /**
     * サムネイルURL
     * @type {string}
     * @memberof Product
     */
    thumbnailUrl: string;
    /**
     * 
     * @type {Array<ExperienceMediaInner>}
     * @memberof Product
     */
    media: Array<ExperienceMediaInner>;
    /**
     * 販売価格（税込）
     * @type {number}
     * @memberof Product
     */
    price: number;
    /**
     * 在庫数
     * @type {number}
     * @memberof Product
     */
    inventory: number;
    /**
     * 重量(kg:少数第1位まで)
     * @type {number}
     * @memberof Product
     */
    weight: number;
    /**
     * 数量単位
     * @type {string}
     * @memberof Product
     */
    itemUnit: string;
    /**
     * 数量単位説明
     * @type {string}
     * @memberof Product
     */
    itemDescription: string;
    /**
     * 
     * @type {DeliveryType}
     * @memberof Product
     */
    deliveryType: DeliveryType;
    /**
     * おすすめポイント1(128文字まで)
     * @type {string}
     * @memberof Product
     */
    recommendedPoint1: string;
    /**
     * おすすめポイント2(128文字まで)
     * @type {string}
     * @memberof Product
     */
    recommendedPoint2: string;
    /**
     * おすすめポイント3(128文字まで)
     * @type {string}
     * @memberof Product
     */
    recommendedPoint3: string;
    /**
     * 賞味期限(単位:日)
     * @type {number}
     * @memberof Product
     */
    expirationDate: number;
    /**
     * 
     * @type {StorageMethodType}
     * @memberof Product
     */
    storageMethodType: StorageMethodType;
    /**
     * 箱の占有率(サイズ:60)
     * @type {number}
     * @memberof Product
     */
    box60Rate: number;
    /**
     * 箱の占有率(サイズ:80)
     * @type {number}
     * @memberof Product
     */
    box80Rate: number;
    /**
     * 箱の占有率(サイズ:100)
     * @type {number}
     * @memberof Product
     */
    box100Rate: number;
    /**
     * 原産地(都道府県)
     * @type {string}
     * @memberof Product
     */
    originPrefecture: string;
    /**
     * 原産地(市区町村)
     * @type {string}
     * @memberof Product
     */
    originCity: string;
    /**
     * 
     * @type {ProductRate}
     * @memberof Product
     */
    rate: ProductRate;
    /**
     * 販売開始日時 (unixtime)
     * @type {number}
     * @memberof Product
     */
    startAt: number;
    /**
     * 販売終了日時 (unixtime)
     * @type {number}
     * @memberof Product
     */
    endAt: number;
}



/**
 * Check if a given object implements the Product interface.
 */
export function instanceOfProduct(value: object): value is Product {
    if (!('id' in value) || value['id'] === undefined) return false;
    if (!('name' in value) || value['name'] === undefined) return false;
    if (!('description' in value) || value['description'] === undefined) return false;
    if (!('status' in value) || value['status'] === undefined) return false;
    if (!('coordinatorId' in value) || value['coordinatorId'] === undefined) return false;
    if (!('producerId' in value) || value['producerId'] === undefined) return false;
    if (!('categoryId' in value) || value['categoryId'] === undefined) return false;
    if (!('productTypeId' in value) || value['productTypeId'] === undefined) return false;
    if (!('productTagIds' in value) || value['productTagIds'] === undefined) return false;
    if (!('thumbnailUrl' in value) || value['thumbnailUrl'] === undefined) return false;
    if (!('media' in value) || value['media'] === undefined) return false;
    if (!('price' in value) || value['price'] === undefined) return false;
    if (!('inventory' in value) || value['inventory'] === undefined) return false;
    if (!('weight' in value) || value['weight'] === undefined) return false;
    if (!('itemUnit' in value) || value['itemUnit'] === undefined) return false;
    if (!('itemDescription' in value) || value['itemDescription'] === undefined) return false;
    if (!('deliveryType' in value) || value['deliveryType'] === undefined) return false;
    if (!('recommendedPoint1' in value) || value['recommendedPoint1'] === undefined) return false;
    if (!('recommendedPoint2' in value) || value['recommendedPoint2'] === undefined) return false;
    if (!('recommendedPoint3' in value) || value['recommendedPoint3'] === undefined) return false;
    if (!('expirationDate' in value) || value['expirationDate'] === undefined) return false;
    if (!('storageMethodType' in value) || value['storageMethodType'] === undefined) return false;
    if (!('box60Rate' in value) || value['box60Rate'] === undefined) return false;
    if (!('box80Rate' in value) || value['box80Rate'] === undefined) return false;
    if (!('box100Rate' in value) || value['box100Rate'] === undefined) return false;
    if (!('originPrefecture' in value) || value['originPrefecture'] === undefined) return false;
    if (!('originCity' in value) || value['originCity'] === undefined) return false;
    if (!('rate' in value) || value['rate'] === undefined) return false;
    if (!('startAt' in value) || value['startAt'] === undefined) return false;
    if (!('endAt' in value) || value['endAt'] === undefined) return false;
    return true;
}

export function ProductFromJSON(json: any): Product {
    return ProductFromJSONTyped(json, false);
}

export function ProductFromJSONTyped(json: any, ignoreDiscriminator: boolean): Product {
    if (json == null) {
        return json;
    }
    return {
        
        'id': json['id'],
        'name': json['name'],
        'description': json['description'],
        'status': ProductStatusFromJSON(json['status']),
        'coordinatorId': json['coordinatorId'],
        'producerId': json['producerId'],
        'categoryId': json['categoryId'],
        'productTypeId': json['productTypeId'],
        'productTagIds': json['productTagIds'],
        'thumbnailUrl': json['thumbnailUrl'],
        'media': ((json['media'] as Array<any>).map(ExperienceMediaInnerFromJSON)),
        'price': json['price'],
        'inventory': json['inventory'],
        'weight': json['weight'],
        'itemUnit': json['itemUnit'],
        'itemDescription': json['itemDescription'],
        'deliveryType': DeliveryTypeFromJSON(json['deliveryType']),
        'recommendedPoint1': json['recommendedPoint1'],
        'recommendedPoint2': json['recommendedPoint2'],
        'recommendedPoint3': json['recommendedPoint3'],
        'expirationDate': json['expirationDate'],
        'storageMethodType': StorageMethodTypeFromJSON(json['storageMethodType']),
        'box60Rate': json['box60Rate'],
        'box80Rate': json['box80Rate'],
        'box100Rate': json['box100Rate'],
        'originPrefecture': json['originPrefecture'],
        'originCity': json['originCity'],
        'rate': ProductRateFromJSON(json['rate']),
        'startAt': json['startAt'],
        'endAt': json['endAt'],
    };
}

export function ProductToJSON(value?: Product | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'id': value['id'],
        'name': value['name'],
        'description': value['description'],
        'status': ProductStatusToJSON(value['status']),
        'coordinatorId': value['coordinatorId'],
        'producerId': value['producerId'],
        'categoryId': value['categoryId'],
        'productTypeId': value['productTypeId'],
        'productTagIds': value['productTagIds'],
        'thumbnailUrl': value['thumbnailUrl'],
        'media': ((value['media'] as Array<any>).map(ExperienceMediaInnerToJSON)),
        'price': value['price'],
        'inventory': value['inventory'],
        'weight': value['weight'],
        'itemUnit': value['itemUnit'],
        'itemDescription': value['itemDescription'],
        'deliveryType': DeliveryTypeToJSON(value['deliveryType']),
        'recommendedPoint1': value['recommendedPoint1'],
        'recommendedPoint2': value['recommendedPoint2'],
        'recommendedPoint3': value['recommendedPoint3'],
        'expirationDate': value['expirationDate'],
        'storageMethodType': StorageMethodTypeToJSON(value['storageMethodType']),
        'box60Rate': value['box60Rate'],
        'box80Rate': value['box80Rate'],
        'box100Rate': value['box100Rate'],
        'originPrefecture': value['originPrefecture'],
        'originCity': value['originCity'],
        'rate': ProductRateToJSON(value['rate']),
        'startAt': value['startAt'],
        'endAt': value['endAt'],
    };
}

