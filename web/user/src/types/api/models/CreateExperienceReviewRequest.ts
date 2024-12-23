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
 * 
 * @export
 * @interface CreateExperienceReviewRequest
 */
export interface CreateExperienceReviewRequest {
    /**
     * 評価（1〜5）
     * @type {number}
     * @memberof CreateExperienceReviewRequest
     */
    rate: number;
    /**
     * タイトル（64文字以内）
     * @type {string}
     * @memberof CreateExperienceReviewRequest
     */
    title: string;
    /**
     * コメント（2000文字以内）
     * @type {string}
     * @memberof CreateExperienceReviewRequest
     */
    comment: string;
}

/**
 * Check if a given object implements the CreateExperienceReviewRequest interface.
 */
export function instanceOfCreateExperienceReviewRequest(value: object): value is CreateExperienceReviewRequest {
    if (!('rate' in value) || value['rate'] === undefined) return false;
    if (!('title' in value) || value['title'] === undefined) return false;
    if (!('comment' in value) || value['comment'] === undefined) return false;
    return true;
}

export function CreateExperienceReviewRequestFromJSON(json: any): CreateExperienceReviewRequest {
    return CreateExperienceReviewRequestFromJSONTyped(json, false);
}

export function CreateExperienceReviewRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): CreateExperienceReviewRequest {
    if (json == null) {
        return json;
    }
    return {
        
        'rate': json['rate'],
        'title': json['title'],
        'comment': json['comment'],
    };
}

export function CreateExperienceReviewRequestToJSON(value?: CreateExperienceReviewRequest | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'rate': value['rate'],
        'title': value['title'],
        'comment': value['comment'],
    };
}

