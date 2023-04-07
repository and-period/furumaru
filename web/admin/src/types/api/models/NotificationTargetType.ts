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
 * 掲載対象
 * @export
 */
export const NotificationTargetType = {
    UNKNOWN: 0,
    USERS: 1,
    PRODUCERS: 2,
    COORDINATORS: 3
} as const;
export type NotificationTargetType = typeof NotificationTargetType[keyof typeof NotificationTargetType];


export function NotificationTargetTypeFromJSON(json: any): NotificationTargetType {
    return NotificationTargetTypeFromJSONTyped(json, false);
}

export function NotificationTargetTypeFromJSONTyped(json: any, ignoreDiscriminator: boolean): NotificationTargetType {
    return json as NotificationTargetType;
}

export function NotificationTargetTypeToJSON(value?: NotificationTargetType | null): any {
    return value as any;
}

