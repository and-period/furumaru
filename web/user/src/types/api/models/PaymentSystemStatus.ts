/* tslint:disable */
 
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
 * 決済システム状態
 * @export
 */
export const PaymentSystemStatus = {
    UNKNOWN: 0,
    IN_USE: 1,
    OUTAGE: 2
} as const;
export type PaymentSystemStatus = typeof PaymentSystemStatus[keyof typeof PaymentSystemStatus];


export function PaymentSystemStatusFromJSON(json: any): PaymentSystemStatus {
    return PaymentSystemStatusFromJSONTyped(json, false);
}

export function PaymentSystemStatusFromJSONTyped(json: any, ignoreDiscriminator: boolean): PaymentSystemStatus {
    return json as PaymentSystemStatus;
}

export function PaymentSystemStatusToJSON(value?: PaymentSystemStatus | null): any {
    return value as any;
}

