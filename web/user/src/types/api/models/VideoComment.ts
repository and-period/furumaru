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
 * @interface VideoComment
 */
export interface VideoComment {
    /**
     * ユーザーID
     * @type {string}
     * @memberof VideoComment
     */
    userId: string;
    /**
     * ユーザー名
     * @type {string}
     * @memberof VideoComment
     */
    username: string;
    /**
     * アカウントID
     * @type {string}
     * @memberof VideoComment
     */
    accountId: string;
    /**
     * サムネイルURL
     * @type {string}
     * @memberof VideoComment
     */
    thumbnailUrl: string;
    /**
     * コメント
     * @type {string}
     * @memberof VideoComment
     */
    comment: string;
    /**
     * 投稿日時
     * @type {number}
     * @memberof VideoComment
     */
    publishedAt: number;
}

/**
 * Check if a given object implements the VideoComment interface.
 */
export function instanceOfVideoComment(value: object): value is VideoComment {
    if (!('userId' in value) || value['userId'] === undefined) return false;
    if (!('username' in value) || value['username'] === undefined) return false;
    if (!('accountId' in value) || value['accountId'] === undefined) return false;
    if (!('thumbnailUrl' in value) || value['thumbnailUrl'] === undefined) return false;
    if (!('comment' in value) || value['comment'] === undefined) return false;
    if (!('publishedAt' in value) || value['publishedAt'] === undefined) return false;
    return true;
}

export function VideoCommentFromJSON(json: any): VideoComment {
    return VideoCommentFromJSONTyped(json, false);
}

export function VideoCommentFromJSONTyped(json: any, ignoreDiscriminator: boolean): VideoComment {
    if (json == null) {
        return json;
    }
    return {
        
        'userId': json['userId'],
        'username': json['username'],
        'accountId': json['accountId'],
        'thumbnailUrl': json['thumbnailUrl'],
        'comment': json['comment'],
        'publishedAt': json['publishedAt'],
    };
}

export function VideoCommentToJSON(value?: VideoComment | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'userId': value['userId'],
        'username': value['username'],
        'accountId': value['accountId'],
        'thumbnailUrl': value['thumbnailUrl'],
        'comment': value['comment'],
        'publishedAt': value['publishedAt'],
    };
}
