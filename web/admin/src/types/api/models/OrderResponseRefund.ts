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
import type { OrderRefundType } from './OrderRefundType';
import {
    OrderRefundTypeFromJSON,
    OrderRefundTypeFromJSONTyped,
    OrderRefundTypeToJSON,
} from './OrderRefundType';

/**
 * 注文キャンセル情報
 * @export
 * @interface OrderResponseRefund
 */
export interface OrderResponseRefund {
    /**
     * 注文キャンセルフラグ
     * @type {boolean}
     * @memberof OrderResponseRefund
     */
    canceled: boolean;
    /**
     * 
     * @type {OrderRefundType}
     * @memberof OrderResponseRefund
     */
    type: OrderRefundType;
    /**
     * 注文キャンセル理由詳細
     * @type {string}
     * @memberof OrderResponseRefund
     */
    reason: string;
    /**
     * 返金金額
     * @type {number}
     * @memberof OrderResponseRefund
     */
    total: number;
}

/**
 * Check if a given object implements the OrderResponseRefund interface.
 */
export function instanceOfOrderResponseRefund(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "canceled" in value;
    isInstance = isInstance && "type" in value;
    isInstance = isInstance && "reason" in value;
    isInstance = isInstance && "total" in value;

    return isInstance;
}

export function OrderResponseRefundFromJSON(json: any): OrderResponseRefund {
    return OrderResponseRefundFromJSONTyped(json, false);
}

export function OrderResponseRefundFromJSONTyped(json: any, ignoreDiscriminator: boolean): OrderResponseRefund {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'canceled': json['canceled'],
        'type': OrderRefundTypeFromJSON(json['type']),
        'reason': json['reason'],
        'total': json['total'],
    };
}

export function OrderResponseRefundToJSON(value?: OrderResponseRefund | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'canceled': value.canceled,
        'type': OrderRefundTypeToJSON(value.type),
        'reason': value.reason,
        'total': value.total,
    };
}

