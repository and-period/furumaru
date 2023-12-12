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
import type { OrderFulfillment } from './OrderFulfillment';
import {
    OrderFulfillmentFromJSON,
    OrderFulfillmentFromJSONTyped,
    OrderFulfillmentToJSON,
} from './OrderFulfillment';
import type { OrderItem } from './OrderItem';
import {
    OrderItemFromJSON,
    OrderItemFromJSONTyped,
    OrderItemToJSON,
} from './OrderItem';
import type { OrderPayment } from './OrderPayment';
import {
    OrderPaymentFromJSON,
    OrderPaymentFromJSONTyped,
    OrderPaymentToJSON,
} from './OrderPayment';
import type { OrderRefund } from './OrderRefund';
import {
    OrderRefundFromJSON,
    OrderRefundFromJSONTyped,
    OrderRefundToJSON,
} from './OrderRefund';
import type { OrderStatus } from './OrderStatus';
import {
    OrderStatusFromJSON,
    OrderStatusFromJSONTyped,
    OrderStatusToJSON,
} from './OrderStatus';

/**
 * 注文履歴情報
 * @export
 * @interface Order
 */
export interface Order {
    /**
     * 注文履歴ID
     * @type {string}
     * @memberof Order
     */
    id: string;
    /**
     * コーディネータID
     * @type {string}
     * @memberof Order
     */
    coordinatorId: string;
    /**
     * プロモーションID
     * @type {string}
     * @memberof Order
     */
    promotionId: string;
    /**
     * 
     * @type {OrderStatus}
     * @memberof Order
     */
    status: OrderStatus;
    /**
     * 
     * @type {OrderPayment}
     * @memberof Order
     */
    payment: OrderPayment;
    /**
     * 
     * @type {OrderRefund}
     * @memberof Order
     */
    refund: OrderRefund;
    /**
     * 注文配送一覧
     * @type {Array<OrderFulfillment>}
     * @memberof Order
     */
    fulfillments: Array<OrderFulfillment>;
    /**
     * 注文商品一覧
     * @type {Array<OrderItem>}
     * @memberof Order
     */
    items: Array<OrderItem>;
}

/**
 * Check if a given object implements the Order interface.
 */
export function instanceOfOrder(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "coordinatorId" in value;
    isInstance = isInstance && "promotionId" in value;
    isInstance = isInstance && "status" in value;
    isInstance = isInstance && "payment" in value;
    isInstance = isInstance && "refund" in value;
    isInstance = isInstance && "fulfillments" in value;
    isInstance = isInstance && "items" in value;

    return isInstance;
}

export function OrderFromJSON(json: any): Order {
    return OrderFromJSONTyped(json, false);
}

export function OrderFromJSONTyped(json: any, ignoreDiscriminator: boolean): Order {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'coordinatorId': json['coordinatorId'],
        'promotionId': json['promotionId'],
        'status': OrderStatusFromJSON(json['status']),
        'payment': OrderPaymentFromJSON(json['payment']),
        'refund': OrderRefundFromJSON(json['refund']),
        'fulfillments': ((json['fulfillments'] as Array<any>).map(OrderFulfillmentFromJSON)),
        'items': ((json['items'] as Array<any>).map(OrderItemFromJSON)),
    };
}

export function OrderToJSON(value?: Order | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'coordinatorId': value.coordinatorId,
        'promotionId': value.promotionId,
        'status': OrderStatusToJSON(value.status),
        'payment': OrderPaymentToJSON(value.payment),
        'refund': OrderRefundToJSON(value.refund),
        'fulfillments': ((value.fulfillments as Array<any>).map(OrderFulfillmentToJSON)),
        'items': ((value.items as Array<any>).map(OrderItemToJSON)),
    };
}
