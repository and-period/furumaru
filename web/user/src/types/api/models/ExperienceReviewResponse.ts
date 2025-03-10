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
import type { ExperienceReview } from './ExperienceReview';
import {
    ExperienceReviewFromJSON,
    ExperienceReviewFromJSONTyped,
    ExperienceReviewToJSON,
} from './ExperienceReview';

/**
 * 
 * @export
 * @interface ExperienceReviewResponse
 */
export interface ExperienceReviewResponse {
    /**
     * 
     * @type {ExperienceReview}
     * @memberof ExperienceReviewResponse
     */
    review: ExperienceReview;
}

/**
 * Check if a given object implements the ExperienceReviewResponse interface.
 */
export function instanceOfExperienceReviewResponse(value: object): value is ExperienceReviewResponse {
    if (!('review' in value) || value['review'] === undefined) return false;
    return true;
}

export function ExperienceReviewResponseFromJSON(json: any): ExperienceReviewResponse {
    return ExperienceReviewResponseFromJSONTyped(json, false);
}

export function ExperienceReviewResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): ExperienceReviewResponse {
    if (json == null) {
        return json;
    }
    return {
        
        'review': ExperienceReviewFromJSON(json['review']),
    };
}

export function ExperienceReviewResponseToJSON(value?: ExperienceReviewResponse | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'review': ExperienceReviewToJSON(value['review']),
    };
}

