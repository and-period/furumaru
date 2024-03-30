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
 * ファイルアップロードの状態
 * @export
 */
export const UploadStatus = {
    UNKNOWN: 0,
    WAITING: 1,
    SUCEEDED: 2,
    FAILED: 3
} as const;
export type UploadStatus = typeof UploadStatus[keyof typeof UploadStatus];


export function UploadStatusFromJSON(json: any): UploadStatus {
    return UploadStatusFromJSONTyped(json, false);
}

export function UploadStatusFromJSONTyped(json: any, ignoreDiscriminator: boolean): UploadStatus {
    return json as UploadStatus;
}

export function UploadStatusToJSON(value?: UploadStatus | null): any {
    return value as any;
}
