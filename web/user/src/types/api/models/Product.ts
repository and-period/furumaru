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
import type { DeliveryType } from './DeliveryType';
import {
    DeliveryTypeFromJSON,
    DeliveryTypeFromJSONTyped,
    DeliveryTypeToJSON,
} from './DeliveryType';
import type { ProductMediaInner } from './ProductMediaInner';
import {
    ProductMediaInnerFromJSON,
    ProductMediaInnerFromJSONTyped,
    ProductMediaInnerToJSON,
} from './ProductMediaInner';
import type { ProductStatus } from './ProductStatus';
import {
    ProductStatusFromJSON,
    ProductStatusFromJSONTyped,
    ProductStatusToJSON,
} from './ProductStatus';
import type { StorageMethodType } from './StorageMethodType';
import {
    StorageMethodTypeFromJSON,
    StorageMethodTypeFromJSONTyped,
    StorageMethodTypeToJSON,
} from './StorageMethodType';
import type { Thumbnail } from './Thumbnail';
import {
    ThumbnailFromJSON,
    ThumbnailFromJSONTyped,
    ThumbnailToJSON,
} from './Thumbnail';

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
     * リサイズ済みサムネイルURL一覧
     * @type {Array<Thumbnail>}
     * @memberof Product
     */
    thumbnails: Array<Thumbnail>;
    /**
     * 
     * @type {Array<ProductMediaInner>}
     * @memberof Product
     */
    media: Array<ProductMediaInner>;
    /**
     * 販売価格
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
export function instanceOfProduct(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "name" in value;
    isInstance = isInstance && "description" in value;
    isInstance = isInstance && "status" in value;
    isInstance = isInstance && "coordinatorId" in value;
    isInstance = isInstance && "producerId" in value;
    isInstance = isInstance && "categoryId" in value;
    isInstance = isInstance && "productTypeId" in value;
    isInstance = isInstance && "productTagIds" in value;
    isInstance = isInstance && "thumbnailUrl" in value;
    isInstance = isInstance && "thumbnails" in value;
    isInstance = isInstance && "media" in value;
    isInstance = isInstance && "price" in value;
    isInstance = isInstance && "inventory" in value;
    isInstance = isInstance && "weight" in value;
    isInstance = isInstance && "itemUnit" in value;
    isInstance = isInstance && "itemDescription" in value;
    isInstance = isInstance && "deliveryType" in value;
    isInstance = isInstance && "recommendedPoint1" in value;
    isInstance = isInstance && "recommendedPoint2" in value;
    isInstance = isInstance && "recommendedPoint3" in value;
    isInstance = isInstance && "expirationDate" in value;
    isInstance = isInstance && "storageMethodType" in value;
    isInstance = isInstance && "box60Rate" in value;
    isInstance = isInstance && "box80Rate" in value;
    isInstance = isInstance && "box100Rate" in value;
    isInstance = isInstance && "originPrefecture" in value;
    isInstance = isInstance && "originCity" in value;
    isInstance = isInstance && "startAt" in value;
    isInstance = isInstance && "endAt" in value;

    return isInstance;
}

export function ProductFromJSON(json: any): Product {
    return ProductFromJSONTyped(json, false);
}

export function ProductFromJSONTyped(json: any, ignoreDiscriminator: boolean): Product {
    if ((json === undefined) || (json === null)) {
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
        'thumbnails': ((json['thumbnails'] as Array<any>).map(ThumbnailFromJSON)),
        'media': ((json['media'] as Array<any>).map(ProductMediaInnerFromJSON)),
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
        'startAt': json['startAt'],
        'endAt': json['endAt'],
    };
}

export function ProductToJSON(value?: Product | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'name': value.name,
        'description': value.description,
        'status': ProductStatusToJSON(value.status),
        'coordinatorId': value.coordinatorId,
        'producerId': value.producerId,
        'categoryId': value.categoryId,
        'productTypeId': value.productTypeId,
        'productTagIds': value.productTagIds,
        'thumbnailUrl': value.thumbnailUrl,
        'thumbnails': ((value.thumbnails as Array<any>).map(ThumbnailToJSON)),
        'media': ((value.media as Array<any>).map(ProductMediaInnerToJSON)),
        'price': value.price,
        'inventory': value.inventory,
        'weight': value.weight,
        'itemUnit': value.itemUnit,
        'itemDescription': value.itemDescription,
        'deliveryType': DeliveryTypeToJSON(value.deliveryType),
        'recommendedPoint1': value.recommendedPoint1,
        'recommendedPoint2': value.recommendedPoint2,
        'recommendedPoint3': value.recommendedPoint3,
        'expirationDate': value.expirationDate,
        'storageMethodType': StorageMethodTypeToJSON(value.storageMethodType),
        'box60Rate': value.box60Rate,
        'box80Rate': value.box80Rate,
        'box100Rate': value.box100Rate,
        'originPrefecture': value.originPrefecture,
        'originCity': value.originCity,
        'startAt': value.startAt,
        'endAt': value.endAt,
    };
}

