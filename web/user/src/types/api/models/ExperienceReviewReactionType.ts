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
 * 体験レビューのリアクション種別
 * @export
 */
export const ExperienceReviewReactionType = {
    UNKNOWN: 0,
    LIKE: 1,
    DISLIKE: 2
} as const;
export type ExperienceReviewReactionType = typeof ExperienceReviewReactionType[keyof typeof ExperienceReviewReactionType];


export function instanceOfExperienceReviewReactionType(value: any): boolean {
    for (const key in ExperienceReviewReactionType) {
        if (Object.prototype.hasOwnProperty.call(ExperienceReviewReactionType, key)) {
            if (ExperienceReviewReactionType[key as keyof typeof ExperienceReviewReactionType] === value) {
                return true;
            }
        }
    }
    return false;
}

export function ExperienceReviewReactionTypeFromJSON(json: any): ExperienceReviewReactionType {
    return ExperienceReviewReactionTypeFromJSONTyped(json, false);
}

export function ExperienceReviewReactionTypeFromJSONTyped(json: any, ignoreDiscriminator: boolean): ExperienceReviewReactionType {
    return json as ExperienceReviewReactionType;
}

export function ExperienceReviewReactionTypeToJSON(value?: ExperienceReviewReactionType | null): any {
    return value as any;
}

