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
 * 配送方法
 * @export
 */
export const DeliveryType = {
    UNKNOWN: 0,
    NORMAL: 1,
    REFRIGERATED: 2,
    FROZEN: 3
} as const;
export type DeliveryType = typeof DeliveryType[keyof typeof DeliveryType];


export function instanceOfDeliveryType(value: any): boolean {
    for (const key in DeliveryType) {
        if (Object.prototype.hasOwnProperty.call(DeliveryType, key)) {
            if (DeliveryType[key as keyof typeof DeliveryType] === value) {
                return true;
            }
        }
    }
    return false;
}

export function DeliveryTypeFromJSON(json: any): DeliveryType {
    return DeliveryTypeFromJSONTyped(json, false);
}

export function DeliveryTypeFromJSONTyped(json: any, ignoreDiscriminator: boolean): DeliveryType {
    return json as DeliveryType;
}

export function DeliveryTypeToJSON(value?: DeliveryType | null): any {
    return value as any;
}

