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
 * マルシェ開催状況
 * @export
 */
export const ScheduleStatus = {
    UNKNOWN: 0,
    WAITING: 1,
    LIVE: 2,
    CLOSED: 3,
    ARCHIVED: 4
} as const;
export type ScheduleStatus = typeof ScheduleStatus[keyof typeof ScheduleStatus];


export function instanceOfScheduleStatus(value: any): boolean {
    for (const key in ScheduleStatus) {
        if (Object.prototype.hasOwnProperty.call(ScheduleStatus, key)) {
            if (ScheduleStatus[key as keyof typeof ScheduleStatus] === value) {
                return true;
            }
        }
    }
    return false;
}

export function ScheduleStatusFromJSON(json: any): ScheduleStatus {
    return ScheduleStatusFromJSONTyped(json, false);
}

export function ScheduleStatusFromJSONTyped(json: any, ignoreDiscriminator: boolean): ScheduleStatus {
    return json as ScheduleStatus;
}

export function ScheduleStatusToJSON(value?: ScheduleStatus | null): any {
    return value as any;
}

