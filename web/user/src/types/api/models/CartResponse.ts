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
import type { Cart } from './Cart';
import {
    CartFromJSON,
    CartFromJSONTyped,
    CartToJSON,
} from './Cart';
import type { Coordinator } from './Coordinator';
import {
    CoordinatorFromJSON,
    CoordinatorFromJSONTyped,
    CoordinatorToJSON,
} from './Coordinator';
import type { Product } from './Product';
import {
    ProductFromJSON,
    ProductFromJSONTyped,
    ProductToJSON,
} from './Product';

/**
 * 
 * @export
 * @interface CartResponse
 */
export interface CartResponse {
    /**
     * 買い物かご一覧
     * @type {Array<Cart>}
     * @memberof CartResponse
     */
    carts: Array<Cart>;
    /**
     * コーディネータ一覧
     * @type {Array<Coordinator>}
     * @memberof CartResponse
     */
    coordinators: Array<Coordinator>;
    /**
     * 商品一覧
     * @type {Array<Product>}
     * @memberof CartResponse
     */
    products: Array<Product>;
}

/**
 * Check if a given object implements the CartResponse interface.
 */
export function instanceOfCartResponse(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "carts" in value;
    isInstance = isInstance && "coordinators" in value;
    isInstance = isInstance && "products" in value;

    return isInstance;
}

export function CartResponseFromJSON(json: any): CartResponse {
    return CartResponseFromJSONTyped(json, false);
}

export function CartResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): CartResponse {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'carts': ((json['carts'] as Array<any>).map(CartFromJSON)),
        'coordinators': ((json['coordinators'] as Array<any>).map(CoordinatorFromJSON)),
        'products': ((json['products'] as Array<any>).map(ProductFromJSON)),
    };
}

export function CartResponseToJSON(value?: CartResponse | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'carts': ((value.carts as Array<any>).map(CartToJSON)),
        'coordinators': ((value.coordinators as Array<any>).map(CoordinatorToJSON)),
        'products': ((value.products as Array<any>).map(ProductToJSON)),
    };
}

