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
import type { ExperienceType } from './ExperienceType';
import {
    ExperienceTypeFromJSON,
    ExperienceTypeFromJSONTyped,
    ExperienceTypeToJSON,
} from './ExperienceType';
import type { Producer } from './Producer';
import {
    ProducerFromJSON,
    ProducerFromJSONTyped,
    ProducerToJSON,
} from './Producer';
import type { Coordinator } from './Coordinator';
import {
    CoordinatorFromJSON,
    CoordinatorFromJSONTyped,
    CoordinatorToJSON,
} from './Coordinator';

/**
 * 
 * @export
 * @interface ExperienceResponse
 */
export interface ExperienceResponse {
    /**
     * 
     * @type {Experience}
     * @memberof ExperienceResponse
     */
    experience: Experience;
    /**
     * 
     * @type {Coordinator}
     * @memberof ExperienceResponse
     */
    coordinator: Coordinator;
    /**
     * 
     * @type {Producer}
     * @memberof ExperienceResponse
     */
    producer: Producer;
    /**
     * 
     * @type {ExperienceType}
     * @memberof ExperienceResponse
     */
    experienceType: ExperienceType;
}

/**
 * Check if a given object implements the ExperienceResponse interface.
 */
export function instanceOfExperienceResponse(value: object): value is ExperienceResponse {
    if (!('experience' in value) || value['experience'] === undefined) return false;
    if (!('coordinator' in value) || value['coordinator'] === undefined) return false;
    if (!('producer' in value) || value['producer'] === undefined) return false;
    if (!('experienceType' in value) || value['experienceType'] === undefined) return false;
    return true;
}

export function ExperienceResponseFromJSON(json: any): ExperienceResponse {
    return ExperienceResponseFromJSONTyped(json, false);
}

export function ExperienceResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): ExperienceResponse {
    if (json == null) {
        return json;
    }
    return {
        
        'experience': ExperienceFromJSON(json['experience']),
        'coordinator': CoordinatorFromJSON(json['coordinator']),
        'producer': ProducerFromJSON(json['producer']),
        'experienceType': ExperienceTypeFromJSON(json['experienceType']),
    };
}

export function ExperienceResponseToJSON(value?: ExperienceResponse | null): any {
    if (value == null) {
        return value;
    }
    return {
        
        'experience': ExperienceToJSON(value['experience']),
        'coordinator': CoordinatorToJSON(value['coordinator']),
        'producer': ProducerToJSON(value['producer']),
        'experienceType': ExperienceTypeToJSON(value['experienceType']),
    };
}
