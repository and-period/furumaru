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


/**
 * 配送会社
 * @export
 */
export const ShippingCarrier = {
    UNKNOWN: 0,
    YAMATO: 1,
    SAGAWA: 2
} as const;
export type ShippingCarrier = typeof ShippingCarrier[keyof typeof ShippingCarrier];


export function ShippingCarrierFromJSON(json: any): ShippingCarrier {
    return ShippingCarrierFromJSONTyped(json, false);
}

export function ShippingCarrierFromJSONTyped(json: any, ignoreDiscriminator: boolean): ShippingCarrier {
    return json as ShippingCarrier;
}

export function ShippingCarrierToJSON(value?: ShippingCarrier | null): any {
    return value as any;
}

