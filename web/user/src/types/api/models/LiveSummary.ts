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

import { exists, mapValues } from '../runtime';
import type { LiveProduct } from './LiveProduct';
import {
    LiveProductFromJSON,
    LiveProductFromJSONTyped,
    LiveProductToJSON,
} from './LiveProduct';
import type { ScheduleStatus } from './ScheduleStatus';
import {
    ScheduleStatusFromJSON,
    ScheduleStatusFromJSONTyped,
    ScheduleStatusToJSON,
} from './ScheduleStatus';
import type { Thumbnail } from './Thumbnail';
import {
    ThumbnailFromJSON,
    ThumbnailFromJSONTyped,
    ThumbnailToJSON,
} from './Thumbnail';

/**
 * 開催中・開催予定のマルシェ情報
 * @export
 * @interface LiveSummary
 */
export interface LiveSummary {
    /**
     * 開催スケジュールID
     * @type {string}
     * @memberof LiveSummary
     */
    scheduleId: string;
    /**
     * コーディネータID
     * @type {string}
     * @memberof LiveSummary
     */
    coordinatorId: string;
    /**
     * 
     * @type {ScheduleStatus}
     * @memberof LiveSummary
     */
    status: ScheduleStatus;
    /**
     * タイトル
     * @type {string}
     * @memberof LiveSummary
     */
    title: string;
    /**
     * サムネイルURL
     * @type {string}
     * @memberof LiveSummary
     */
    thumbnailUrl: string;
    /**
     * リサイズ済みサムネイルURL一覧
     * @type {Array<Thumbnail>}
     * @memberof LiveSummary
     */
    thumbnails: Array<Thumbnail>;
    /**
     * マルシェ開始日時 (unixtime)
     * @type {number}
     * @memberof LiveSummary
     */
    startAt: number;
    /**
     * マルシェ終了日時 (unixtime)
     * @type {number}
     * @memberof LiveSummary
     */
    endAt: number;
    /**
     * 販売商品一覧
     * @type {Array<LiveProduct>}
     * @memberof LiveSummary
     */
    products: Array<LiveProduct>;
}

/**
 * Check if a given object implements the LiveSummary interface.
 */
export function instanceOfLiveSummary(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "scheduleId" in value;
    isInstance = isInstance && "coordinatorId" in value;
    isInstance = isInstance && "status" in value;
    isInstance = isInstance && "title" in value;
    isInstance = isInstance && "thumbnailUrl" in value;
    isInstance = isInstance && "thumbnails" in value;
    isInstance = isInstance && "startAt" in value;
    isInstance = isInstance && "endAt" in value;
    isInstance = isInstance && "products" in value;

    return isInstance;
}

export function LiveSummaryFromJSON(json: any): LiveSummary {
    return LiveSummaryFromJSONTyped(json, false);
}

export function LiveSummaryFromJSONTyped(json: any, ignoreDiscriminator: boolean): LiveSummary {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'scheduleId': json['scheduleId'],
        'coordinatorId': json['coordinatorId'],
        'status': ScheduleStatusFromJSON(json['status']),
        'title': json['title'],
        'thumbnailUrl': json['thumbnailUrl'],
        'thumbnails': ((json['thumbnails'] as Array<any>).map(ThumbnailFromJSON)),
        'startAt': json['startAt'],
        'endAt': json['endAt'],
        'products': ((json['products'] as Array<any>).map(LiveProductFromJSON)),
    };
}

export function LiveSummaryToJSON(value?: LiveSummary | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'scheduleId': value.scheduleId,
        'coordinatorId': value.coordinatorId,
        'status': ScheduleStatusToJSON(value.status),
        'title': value.title,
        'thumbnailUrl': value.thumbnailUrl,
        'thumbnails': ((value.thumbnails as Array<any>).map(ThumbnailToJSON)),
        'startAt': value.startAt,
        'endAt': value.endAt,
        'products': ((value.products as Array<any>).map(LiveProductToJSON)),
    };
}
