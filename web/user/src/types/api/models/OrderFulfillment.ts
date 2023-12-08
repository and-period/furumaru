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
import type { FulfillmentStatus } from './FulfillmentStatus';
import {
    FulfillmentStatusFromJSON,
    FulfillmentStatusFromJSONTyped,
    FulfillmentStatusToJSON,
} from './FulfillmentStatus';
import type { Prefecture } from './Prefecture';
import {
    PrefectureFromJSON,
    PrefectureFromJSONTyped,
    PrefectureToJSON,
} from './Prefecture';
import type { ShippingCarrier } from './ShippingCarrier';
import {
    ShippingCarrierFromJSON,
    ShippingCarrierFromJSONTyped,
    ShippingCarrierToJSON,
} from './ShippingCarrier';
import type { ShippingSize } from './ShippingSize';
import {
    ShippingSizeFromJSON,
    ShippingSizeFromJSONTyped,
    ShippingSizeToJSON,
} from './ShippingSize';
import type { ShippingType } from './ShippingType';
import {
    ShippingTypeFromJSON,
    ShippingTypeFromJSONTyped,
    ShippingTypeToJSON,
} from './ShippingType';

/**
 * 注文配送情報
 * @export
 * @interface OrderFulfillment
 */
export interface OrderFulfillment {
    /**
     * 注文配送ID
     * @type {string}
     * @memberof OrderFulfillment
     */
    fulfillmentId: string;
    /**
     * 伝票番号
     * @type {string}
     * @memberof OrderFulfillment
     */
    trackingNumber: string;
    /**
     * 
     * @type {FulfillmentStatus}
     * @memberof OrderFulfillment
     */
    status: FulfillmentStatus;
    /**
     * 
     * @type {ShippingCarrier}
     * @memberof OrderFulfillment
     */
    shippingCarrier: ShippingCarrier;
    /**
     * 
     * @type {ShippingType}
     * @memberof OrderFulfillment
     */
    shippingType: ShippingType;
    /**
     * 箱の通番
     * @type {number}
     * @memberof OrderFulfillment
     */
    boxNumber: number;
    /**
     * 
     * @type {ShippingSize}
     * @memberof OrderFulfillment
     */
    boxSize: ShippingSize;
    /**
     * 箱の占有率
     * @type {number}
     * @memberof OrderFulfillment
     */
    boxRate: number;
    /**
     * 配送日時（unixtime）
     * @type {number}
     * @memberof OrderFulfillment
     */
    shippedAt: number;
    /**
     * 配送先 住所ID
     * @type {string}
     * @memberof OrderFulfillment
     */
    addressId: string;
    /**
     * 配送先 氏名（姓）
     * @type {string}
     * @memberof OrderFulfillment
     */
    lastname: string;
    /**
     * 配送先 氏名（名）
     * @type {string}
     * @memberof OrderFulfillment
     */
    firstname: string;
    /**
     * 配送先 郵便番号
     * @type {string}
     * @memberof OrderFulfillment
     */
    postalCode: string;
    /**
     * 
     * @type {Prefecture}
     * @memberof OrderFulfillment
     */
    prefectureCode: Prefecture;
    /**
     * 配送先 市区町村
     * @type {string}
     * @memberof OrderFulfillment
     */
    city: string;
    /**
     * 配送先 町名・番地
     * @type {string}
     * @memberof OrderFulfillment
     */
    addressLine1: string;
    /**
     * 配送先 ビル名・号室など
     * @type {string}
     * @memberof OrderFulfillment
     */
    addressLine2: string;
    /**
     * 配送先 電話番号
     * @type {string}
     * @memberof OrderFulfillment
     */
    phoneNumber: string;
}

/**
 * Check if a given object implements the OrderFulfillment interface.
 */
export function instanceOfOrderFulfillment(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "fulfillmentId" in value;
    isInstance = isInstance && "trackingNumber" in value;
    isInstance = isInstance && "status" in value;
    isInstance = isInstance && "shippingCarrier" in value;
    isInstance = isInstance && "shippingType" in value;
    isInstance = isInstance && "boxNumber" in value;
    isInstance = isInstance && "boxSize" in value;
    isInstance = isInstance && "boxRate" in value;
    isInstance = isInstance && "shippedAt" in value;
    isInstance = isInstance && "addressId" in value;
    isInstance = isInstance && "lastname" in value;
    isInstance = isInstance && "firstname" in value;
    isInstance = isInstance && "postalCode" in value;
    isInstance = isInstance && "prefectureCode" in value;
    isInstance = isInstance && "city" in value;
    isInstance = isInstance && "addressLine1" in value;
    isInstance = isInstance && "addressLine2" in value;
    isInstance = isInstance && "phoneNumber" in value;

    return isInstance;
}

export function OrderFulfillmentFromJSON(json: any): OrderFulfillment {
    return OrderFulfillmentFromJSONTyped(json, false);
}

export function OrderFulfillmentFromJSONTyped(json: any, ignoreDiscriminator: boolean): OrderFulfillment {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'fulfillmentId': json['fulfillmentId'],
        'trackingNumber': json['trackingNumber'],
        'status': FulfillmentStatusFromJSON(json['status']),
        'shippingCarrier': ShippingCarrierFromJSON(json['shippingCarrier']),
        'shippingType': ShippingTypeFromJSON(json['shippingType']),
        'boxNumber': json['boxNumber'],
        'boxSize': ShippingSizeFromJSON(json['boxSize']),
        'boxRate': json['boxRate'],
        'shippedAt': json['shippedAt'],
        'addressId': json['addressId'],
        'lastname': json['lastname'],
        'firstname': json['firstname'],
        'postalCode': json['postalCode'],
        'prefectureCode': PrefectureFromJSON(json['prefectureCode']),
        'city': json['city'],
        'addressLine1': json['addressLine1'],
        'addressLine2': json['addressLine2'],
        'phoneNumber': json['phoneNumber'],
    };
}

export function OrderFulfillmentToJSON(value?: OrderFulfillment | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'fulfillmentId': value.fulfillmentId,
        'trackingNumber': value.trackingNumber,
        'status': FulfillmentStatusToJSON(value.status),
        'shippingCarrier': ShippingCarrierToJSON(value.shippingCarrier),
        'shippingType': ShippingTypeToJSON(value.shippingType),
        'boxNumber': value.boxNumber,
        'boxSize': ShippingSizeToJSON(value.boxSize),
        'boxRate': value.boxRate,
        'shippedAt': value.shippedAt,
        'addressId': value.addressId,
        'lastname': value.lastname,
        'firstname': value.firstname,
        'postalCode': value.postalCode,
        'prefectureCode': PrefectureToJSON(value.prefectureCode),
        'city': value.city,
        'addressLine1': value.addressLine1,
        'addressLine2': value.addressLine2,
        'phoneNumber': value.phoneNumber,
    };
}

