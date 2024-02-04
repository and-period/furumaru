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
/**
 * 
 * @export
 * @interface LiveComment
 */
export interface LiveComment {
    /**
     * ユーザーID
     * @type {string}
     * @memberof LiveComment
     */
    userId: string;
    /**
     * ユーザー名
     * @type {string}
     * @memberof LiveComment
     */
    username: string;
    /**
     * アカウントID
     * @type {string}
     * @memberof LiveComment
     */
    accountId: string;
    /**
     * サムネイルURL
     * @type {string}
     * @memberof LiveComment
     */
    thumbnailUrl: string;
    /**
     * コメント
     * @type {string}
     * @memberof LiveComment
     */
    comment: string;
    /**
     * 投稿日時
     * @type {number}
     * @memberof LiveComment
     */
    publishedAt: number;
}

/**
 * Check if a given object implements the LiveComment interface.
 */
export function instanceOfLiveComment(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "userId" in value;
    isInstance = isInstance && "username" in value;
    isInstance = isInstance && "accountId" in value;
    isInstance = isInstance && "thumbnailUrl" in value;
    isInstance = isInstance && "comment" in value;
    isInstance = isInstance && "publishedAt" in value;

    return isInstance;
}

export function LiveCommentFromJSON(json: any): LiveComment {
    return LiveCommentFromJSONTyped(json, false);
}

export function LiveCommentFromJSONTyped(json: any, ignoreDiscriminator: boolean): LiveComment {
    if ((json === undefined) || (json === null)) {
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

export function LiveCommentToJSON(value?: LiveComment | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'userId': value.userId,
        'username': value.username,
        'accountId': value.accountId,
        'thumbnailUrl': value.thumbnailUrl,
        'comment': value.comment,
        'publishedAt': value.publishedAt,
    };
}
