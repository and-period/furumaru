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
/**
 * 
 * @export
 * @interface V1CreateProductReviewRequest
 */
export interface V1CreateProductReviewRequest {
    /**
     * 評価（1〜5）
     * @type {number}
     * @memberof V1CreateProductReviewRequest
     */
    rate: number;
    /**
     * タイトル（64文字以内）
     * @type {string}
     * @memberof V1CreateProductReviewRequest
     */
    title: string;
    /**
     * コメント（2000文字以内）
     * @type {string}
     * @memberof V1CreateProductReviewRequest
     */
    comment: string;
}

/**
 * Check if a given object implements the V1CreateProductReviewRequest interface.
 */
export function instanceOfV1CreateProductReviewRequest(value: object): value is V1CreateProductReviewRequest {
    if (!('rate' in value) || value['rate'] === undefined) return false;
    if (!('title' in value) || value['title'] === undefined) return false;
    if (!('comment' in value) || value['comment'] === undefined) return false;
    return true;
}

export function V1CreateProductReviewRequestFromJSON(json: any): V1CreateProductReviewRequest {
    return V1CreateProductReviewRequestFromJSONTyped(json, false);
}

export function V1CreateProductReviewRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): V1CreateProductReviewRequest {
    if (json == null) {
        return json;
    }
    return {
        
        'rate': json['rate'],
        'title': json['title'],
        'comment': json['comment'],
    };
}

export function V1CreateProductReviewRequestToJSON(value?: V1CreateProductReviewRequest | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'rate': value['rate'],
        'title': value['title'],
        'comment': value['comment'],
    };
}

