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
import type { OrderRefundType } from './OrderRefundType';
import {
    OrderRefundTypeFromJSON,
    OrderRefundTypeFromJSONTyped,
    OrderRefundTypeToJSON,
} from './OrderRefundType';

/**
 * 注文キャンセル情報
 * @export
 * @interface OrderRefund
 */
export interface OrderRefund {
    /**
     * 返金金額
     * @type {number}
     * @memberof OrderRefund
     */
    total: number;
    /**
     * 
     * @type {OrderRefundType}
     * @memberof OrderRefund
     */
    type: OrderRefundType;
    /**
     * 注文キャンセル理由
     * @type {string}
     * @memberof OrderRefund
     */
    reason: string;
    /**
     * 注文キャンセルフラグ
     * @type {boolean}
     * @memberof OrderRefund
     */
    canceled: boolean;
    /**
     * 注文キャンセル日時（unixtime）
     * @type {number}
     * @memberof OrderRefund
     */
    canceledAt: number;
}

/**
 * Check if a given object implements the OrderRefund interface.
 */
export function instanceOfOrderRefund(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "total" in value;
    isInstance = isInstance && "type" in value;
    isInstance = isInstance && "reason" in value;
    isInstance = isInstance && "canceled" in value;
    isInstance = isInstance && "canceledAt" in value;

    return isInstance;
}

export function OrderRefundFromJSON(json: any): OrderRefund {
    return OrderRefundFromJSONTyped(json, false);
}

export function OrderRefundFromJSONTyped(json: any, ignoreDiscriminator: boolean): OrderRefund {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'total': json['total'],
        'type': OrderRefundTypeFromJSON(json['type']),
        'reason': json['reason'],
        'canceled': json['canceled'],
        'canceledAt': json['canceledAt'],
    };
}

export function OrderRefundToJSON(value?: OrderRefund | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'total': value.total,
        'type': OrderRefundTypeToJSON(value.type),
        'reason': value.reason,
        'canceled': value.canceled,
        'canceledAt': value.canceledAt,
    };
}

