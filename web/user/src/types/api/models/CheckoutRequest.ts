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
import type { CheckoutCreditCard } from './CheckoutCreditCard';
import {
    CheckoutCreditCardFromJSON,
    CheckoutCreditCardFromJSONTyped,
    CheckoutCreditCardToJSON,
} from './CheckoutCreditCard';
import type { PaymentMethodType } from './PaymentMethodType';
import {
    PaymentMethodTypeFromJSON,
    PaymentMethodTypeFromJSONTyped,
    PaymentMethodTypeToJSON,
} from './PaymentMethodType';

/**
 * 
 * @export
 * @interface CheckoutRequest
 */
export interface CheckoutRequest {
    /**
     * 支払いキー(重複判定用)
     * @type {string}
     * @memberof CheckoutRequest
     */
    requestId: string;
    /**
     * コーディネータID
     * @type {string}
     * @memberof CheckoutRequest
     */
    coordinatorId: string;
    /**
     * 箱の通番（箱単位で購入する場合のみ）
     * @type {number}
     * @memberof CheckoutRequest
     */
    boxNumber: number;
    /**
     * 請求先住所ID
     * @type {string}
     * @memberof CheckoutRequest
     */
    billingAddressId: string;
    /**
     * 配送先住所ID
     * @type {string}
     * @memberof CheckoutRequest
     */
    shippingAddressId: string;
    /**
     * プロモーションコード（割引適用時のみ）
     * @type {string}
     * @memberof CheckoutRequest
     */
    promotionCode: string;
    /**
     * 
     * @type {PaymentMethodType}
     * @memberof CheckoutRequest
     */
    paymentMethod: PaymentMethodType;
    /**
     * 決済ページからの遷移先URL
     * @type {string}
     * @memberof CheckoutRequest
     */
    callbackUrl: string;
    /**
     * 支払い合計金額（誤り検出用）
     * @type {number}
     * @memberof CheckoutRequest
     */
    total: number;
    /**
     * 
     * @type {CheckoutCreditCard}
     * @memberof CheckoutRequest
     */
    creditCard: CheckoutCreditCard;
}

/**
 * Check if a given object implements the CheckoutRequest interface.
 */
export function instanceOfCheckoutRequest(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "requestId" in value;
    isInstance = isInstance && "coordinatorId" in value;
    isInstance = isInstance && "boxNumber" in value;
    isInstance = isInstance && "billingAddressId" in value;
    isInstance = isInstance && "shippingAddressId" in value;
    isInstance = isInstance && "promotionCode" in value;
    isInstance = isInstance && "paymentMethod" in value;
    isInstance = isInstance && "callbackUrl" in value;
    isInstance = isInstance && "total" in value;
    isInstance = isInstance && "creditCard" in value;

    return isInstance;
}

export function CheckoutRequestFromJSON(json: any): CheckoutRequest {
    return CheckoutRequestFromJSONTyped(json, false);
}

export function CheckoutRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): CheckoutRequest {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'requestId': json['requestId'],
        'coordinatorId': json['coordinatorId'],
        'boxNumber': json['boxNumber'],
        'billingAddressId': json['billingAddressId'],
        'shippingAddressId': json['shippingAddressId'],
        'promotionCode': json['promotionCode'],
        'paymentMethod': PaymentMethodTypeFromJSON(json['paymentMethod']),
        'callbackUrl': json['callbackUrl'],
        'total': json['total'],
        'creditCard': CheckoutCreditCardFromJSON(json['creditCard']),
    };
}

export function CheckoutRequestToJSON(value?: CheckoutRequest | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'requestId': value.requestId,
        'coordinatorId': value.coordinatorId,
        'boxNumber': value.boxNumber,
        'billingAddressId': value.billingAddressId,
        'shippingAddressId': value.shippingAddressId,
        'promotionCode': value.promotionCode,
        'paymentMethod': PaymentMethodTypeToJSON(value.paymentMethod),
        'callbackUrl': value.callbackUrl,
        'total': value.total,
        'creditCard': CheckoutCreditCardToJSON(value.creditCard),
    };
}

