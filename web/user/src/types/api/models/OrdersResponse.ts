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
import type { Coordinator } from './Coordinator';
import {
    CoordinatorFromJSON,
    CoordinatorFromJSONTyped,
    CoordinatorToJSON,
} from './Coordinator';
import type { Order } from './Order';
import {
    OrderFromJSON,
    OrderFromJSONTyped,
    OrderToJSON,
} from './Order';
import type { Product } from './Product';
import {
    ProductFromJSON,
    ProductFromJSONTyped,
    ProductToJSON,
} from './Product';
import type { Promotion } from './Promotion';
import {
    PromotionFromJSON,
    PromotionFromJSONTyped,
    PromotionToJSON,
} from './Promotion';

/**
 * 
 * @export
 * @interface OrdersResponse
 */
export interface OrdersResponse {
    /**
     * 
     * @type {Array<Order>}
     * @memberof OrdersResponse
     */
    orders: Array<Order>;
    /**
     * 
     * @type {Array<Coordinator>}
     * @memberof OrdersResponse
     */
    coordinators: Array<Coordinator>;
    /**
     * 
     * @type {Array<Promotion>}
     * @memberof OrdersResponse
     */
    promotions: Array<Promotion>;
    /**
     * 
     * @type {Array<Product>}
     * @memberof OrdersResponse
     */
    products: Array<Product>;
    /**
     * 合計数
     * @type {number}
     * @memberof OrdersResponse
     */
    total: number;
}

/**
 * Check if a given object implements the OrdersResponse interface.
 */
export function instanceOfOrdersResponse(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "orders" in value;
    isInstance = isInstance && "coordinators" in value;
    isInstance = isInstance && "promotions" in value;
    isInstance = isInstance && "products" in value;
    isInstance = isInstance && "total" in value;

    return isInstance;
}

export function OrdersResponseFromJSON(json: any): OrdersResponse {
    return OrdersResponseFromJSONTyped(json, false);
}

export function OrdersResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): OrdersResponse {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'orders': ((json['orders'] as Array<any>).map(OrderFromJSON)),
        'coordinators': ((json['coordinators'] as Array<any>).map(CoordinatorFromJSON)),
        'promotions': ((json['promotions'] as Array<any>).map(PromotionFromJSON)),
        'products': ((json['products'] as Array<any>).map(ProductFromJSON)),
        'total': json['total'],
    };
}

export function OrdersResponseToJSON(value?: OrdersResponse | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'orders': ((value.orders as Array<any>).map(OrderToJSON)),
        'coordinators': ((value.coordinators as Array<any>).map(CoordinatorToJSON)),
        'promotions': ((value.promotions as Array<any>).map(PromotionToJSON)),
        'products': ((value.products as Array<any>).map(ProductToJSON)),
        'total': value.total,
    };
}

