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
 * 過去のマルシェ情報
 * @export
 * @interface ArchiveSummary
 */
export interface ArchiveSummary {
    /**
     * 開催スケジュールID
     * @type {string}
     * @memberof ArchiveSummary
     */
    scheduleId: string;
    /**
     * コーディネータID
     * @type {string}
     * @memberof ArchiveSummary
     */
    coordinatorId: string;
    /**
     * タイトル
     * @type {string}
     * @memberof ArchiveSummary
     */
    title: string;
    /**
     * マルシェ開始日時 (unixtime)
     * @type {number}
     * @memberof ArchiveSummary
     */
    startAt: number;
    /**
     * マルシェ終了日時 (unixtime)
     * @type {number}
     * @memberof ArchiveSummary
     */
    endAt: number;
    /**
     * サムネイルURL
     * @type {string}
     * @memberof ArchiveSummary
     */
    thumbnailUrl: string;
}

/**
 * Check if a given object implements the ArchiveSummary interface.
 */
export function instanceOfArchiveSummary(value: object): value is ArchiveSummary {
    if (!('scheduleId' in value) || value['scheduleId'] === undefined) return false;
    if (!('coordinatorId' in value) || value['coordinatorId'] === undefined) return false;
    if (!('title' in value) || value['title'] === undefined) return false;
    if (!('startAt' in value) || value['startAt'] === undefined) return false;
    if (!('endAt' in value) || value['endAt'] === undefined) return false;
    if (!('thumbnailUrl' in value) || value['thumbnailUrl'] === undefined) return false;
    return true;
}

export function ArchiveSummaryFromJSON(json: any): ArchiveSummary {
    return ArchiveSummaryFromJSONTyped(json, false);
}

export function ArchiveSummaryFromJSONTyped(json: any, ignoreDiscriminator: boolean): ArchiveSummary {
    if (json == null) {
        return json;
    }
    return {
        
        'scheduleId': json['scheduleId'],
        'coordinatorId': json['coordinatorId'],
        'title': json['title'],
        'startAt': json['startAt'],
        'endAt': json['endAt'],
        'thumbnailUrl': json['thumbnailUrl'],
    };
}

export function ArchiveSummaryToJSON(value?: ArchiveSummary | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'scheduleId': value['scheduleId'],
        'coordinatorId': value['coordinatorId'],
        'title': value['title'],
        'startAt': value['startAt'],
        'endAt': value['endAt'],
        'thumbnailUrl': value['thumbnailUrl'],
    };
}

