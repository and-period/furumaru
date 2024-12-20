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
 * 割引計算方法
 * @export
 */
export const DiscountType = {
    UNKNOWN: 0,
    AMOUNT: 1,
    RATE: 2,
    FREE_SHIPPING: 3
} as const;
export type DiscountType = typeof DiscountType[keyof typeof DiscountType];


export function instanceOfDiscountType(value: any): boolean {
    for (const key in DiscountType) {
        if (Object.prototype.hasOwnProperty.call(DiscountType, key)) {
            if (DiscountType[key as keyof typeof DiscountType] === value) {
                return true;
            }
        }
    }
    return false;
}

export function DiscountTypeFromJSON(json: any): DiscountType {
    return DiscountTypeFromJSONTyped(json, false);
}

export function DiscountTypeFromJSONTyped(json: any, ignoreDiscriminator: boolean): DiscountType {
    return json as DiscountType;
}

export function DiscountTypeToJSON(value?: DiscountType | null): any {
    return value as any;
}

