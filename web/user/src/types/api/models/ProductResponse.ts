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
import type { Category } from './Category';
import {
    CategoryFromJSON,
    CategoryFromJSONTyped,
    CategoryToJSON,
} from './Category';
import type { Producer } from './Producer';
import {
    ProducerFromJSON,
    ProducerFromJSONTyped,
    ProducerToJSON,
} from './Producer';
import type { ProductTag } from './ProductTag';
import {
    ProductTagFromJSON,
    ProductTagFromJSONTyped,
    ProductTagToJSON,
} from './ProductTag';
import type { ProductType } from './ProductType';
import {
    ProductTypeFromJSON,
    ProductTypeFromJSONTyped,
    ProductTypeToJSON,
} from './ProductType';
import type { Product } from './Product';
import {
    ProductFromJSON,
    ProductFromJSONTyped,
    ProductToJSON,
} from './Product';
import type { Coordinator } from './Coordinator';
import {
    CoordinatorFromJSON,
    CoordinatorFromJSONTyped,
    CoordinatorToJSON,
} from './Coordinator';

/**
 * 
 * @export
 * @interface ProductResponse
 */
export interface ProductResponse {
    /**
     * 
     * @type {Product}
     * @memberof ProductResponse
     */
    product: Product;
    /**
     * 
     * @type {Coordinator}
     * @memberof ProductResponse
     */
    coordinator: Coordinator;
    /**
     * 
     * @type {Producer}
     * @memberof ProductResponse
     */
    producer: Producer;
    /**
     * 
     * @type {Category}
     * @memberof ProductResponse
     */
    category: Category;
    /**
     * 
     * @type {ProductType}
     * @memberof ProductResponse
     */
    productType: ProductType;
    /**
     * 商品タグ一覧
     * @type {Array<ProductTag>}
     * @memberof ProductResponse
     */
    productTags: Array<ProductTag>;
}

/**
 * Check if a given object implements the ProductResponse interface.
 */
export function instanceOfProductResponse(value: object): value is ProductResponse {
    if (!('product' in value) || value['product'] === undefined) return false;
    if (!('coordinator' in value) || value['coordinator'] === undefined) return false;
    if (!('producer' in value) || value['producer'] === undefined) return false;
    if (!('category' in value) || value['category'] === undefined) return false;
    if (!('productType' in value) || value['productType'] === undefined) return false;
    if (!('productTags' in value) || value['productTags'] === undefined) return false;
    return true;
}

export function ProductResponseFromJSON(json: any): ProductResponse {
    return ProductResponseFromJSONTyped(json, false);
}

export function ProductResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): ProductResponse {
    if (json == null) {
        return json;
    }
    return {
        
        'product': ProductFromJSON(json['product']),
        'coordinator': CoordinatorFromJSON(json['coordinator']),
        'producer': ProducerFromJSON(json['producer']),
        'category': CategoryFromJSON(json['category']),
        'productType': ProductTypeFromJSON(json['productType']),
        'productTags': ((json['productTags'] as Array<any>).map(ProductTagFromJSON)),
    };
}

export function ProductResponseToJSON(value?: ProductResponse | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'product': ProductToJSON(value['product']),
        'coordinator': CoordinatorToJSON(value['coordinator']),
        'producer': ProducerToJSON(value['producer']),
        'category': CategoryToJSON(value['category']),
        'productType': ProductTypeToJSON(value['productType']),
        'productTags': ((value['productTags'] as Array<any>).map(ProductTagToJSON)),
    };
}

