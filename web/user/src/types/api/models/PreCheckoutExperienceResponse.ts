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
import type { Experience } from './Experience';
import {
    ExperienceFromJSON,
    ExperienceFromJSONTyped,
    ExperienceToJSON,
} from './Experience';
import type { Promotion } from './Promotion';
import {
    PromotionFromJSON,
    PromotionFromJSONTyped,
    PromotionToJSON,
} from './Promotion';

/**
 * 
 * @export
 * @interface PreCheckoutExperienceResponse
 */
export interface PreCheckoutExperienceResponse {
    /**
     * 
     * @type {Experience}
     * @memberof PreCheckoutExperienceResponse
     */
    experience: Experience;
    /**
     * 
     * @type {Promotion}
     * @memberof PreCheckoutExperienceResponse
     */
    promotion: Promotion;
    /**
     * 小計(税込:商品合計金額)
     * @type {number}
     * @memberof PreCheckoutExperienceResponse
     */
    subtotal: number;
    /**
     * 割引金額(税込)
     * @type {number}
     * @memberof PreCheckoutExperienceResponse
     */
    discount: number;
    /**
     * 合計金額(税込)
     * @type {number}
     * @memberof PreCheckoutExperienceResponse
     */
    total: number;
    /**
     * 支払い時にAPIへ送信するキー(重複判定用)
     * @type {string}
     * @memberof PreCheckoutExperienceResponse
     */
    requestId?: string;
}

/**
 * Check if a given object implements the PreCheckoutExperienceResponse interface.
 */
export function instanceOfPreCheckoutExperienceResponse(value: object): value is PreCheckoutExperienceResponse {
    if (!('experience' in value) || value['experience'] === undefined) return false;
    if (!('promotion' in value) || value['promotion'] === undefined) return false;
    if (!('subtotal' in value) || value['subtotal'] === undefined) return false;
    if (!('discount' in value) || value['discount'] === undefined) return false;
    if (!('total' in value) || value['total'] === undefined) return false;
    return true;
}

export function PreCheckoutExperienceResponseFromJSON(json: any): PreCheckoutExperienceResponse {
    return PreCheckoutExperienceResponseFromJSONTyped(json, false);
}

export function PreCheckoutExperienceResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): PreCheckoutExperienceResponse {
    if (json == null) {
        return json;
    }
    return {
        
        'experience': ExperienceFromJSON(json['experience']),
        'promotion': PromotionFromJSON(json['promotion']),
        'subtotal': json['subtotal'],
        'discount': json['discount'],
        'total': json['total'],
        'requestId': json['requestId'] == null ? undefined : json['requestId'],
    };
}

export function PreCheckoutExperienceResponseToJSON(value?: PreCheckoutExperienceResponse | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'experience': ExperienceToJSON(value['experience']),
        'promotion': PromotionToJSON(value['promotion']),
        'subtotal': value['subtotal'],
        'discount': value['discount'],
        'total': value['total'],
        'requestId': value['requestId'],
    };
}

