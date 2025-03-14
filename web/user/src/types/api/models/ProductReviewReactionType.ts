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


/**
 * 商品レビューのリアクション種別
 * @export
 */
export const ProductReviewReactionType = {
    UNKNOWN: 0,
    LIKE: 1,
    DISLIKE: 2
} as const;
export type ProductReviewReactionType = typeof ProductReviewReactionType[keyof typeof ProductReviewReactionType];


export function instanceOfProductReviewReactionType(value: any): boolean {
    for (const key in ProductReviewReactionType) {
        if (Object.prototype.hasOwnProperty.call(ProductReviewReactionType, key)) {
            if (ProductReviewReactionType[key as keyof typeof ProductReviewReactionType] === value) {
                return true;
            }
        }
    }
    return false;
}

export function ProductReviewReactionTypeFromJSON(json: any): ProductReviewReactionType {
    return ProductReviewReactionTypeFromJSONTyped(json, false);
}

export function ProductReviewReactionTypeFromJSONTyped(json: any, ignoreDiscriminator: boolean): ProductReviewReactionType {
    return json as ProductReviewReactionType;
}

export function ProductReviewReactionTypeToJSON(value?: ProductReviewReactionType | null): any {
    return value as any;
}

