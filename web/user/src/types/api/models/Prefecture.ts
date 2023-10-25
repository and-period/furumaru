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
 * 都道府県コード
 * @export
 */
export const Prefecture = {
    UNKNOWN: '',
    HOKKAIDO: 'hokkaido',
    AOMORI: 'aomori',
    IWATE: 'iwate',
    MIYAGI: 'miyagi',
    AKITA: 'akita',
    YAMAGATA: 'yamagata',
    FUKUSHIMA: 'fukushima',
    IBARAKI: 'ibaraki',
    TOCHIGI: 'tochigi',
    GUNMA: 'gunma',
    SAITAMA: 'saitama',
    CHIBA: 'chiba',
    TOKYO: 'tokyo',
    KANAGAWA: 'kanagawa',
    NIIGATA: 'niigata',
    TOYAMA: 'toyama',
    ISHIKAWA: 'ishikawa',
    FUKUI: 'fukui',
    YAMANASHI: 'yamanashi',
    NAGANO: 'nagano',
    GIFU: 'gifu',
    SHIZUOKA: 'shizuoka',
    AICHI: 'aichi',
    MIE: 'mie',
    SHIGA: 'shiga',
    KYOTO: 'kyoto',
    OSAKA: 'osaka',
    HYOGO: 'hyogo',
    NARA: 'nara',
    WAKAYAMA: 'wakayama',
    TOTTORI: 'tottori',
    SHIMANE: 'shimane',
    OKAYAMA: 'okayama',
    HIROSHIMA: 'hiroshima',
    YAMAGUCHI: 'yamaguchi',
    TOKUSHIMA: 'tokushima',
    KAGAWA: 'kagawa',
    EHIME: 'ehime',
    KOCHI: 'kochi',
    FUKUOKA: 'fukuoka',
    SAGA: 'saga',
    NAGASAKI: 'nagasaki',
    KUMAMOTO: 'kumamoto',
    OITA: 'oita',
    MIYAZAKI: 'miyazaki',
    KAGOSHIMA: 'kagoshima',
    OKINAWA: 'okinawa'
} as const;
export type Prefecture = typeof Prefecture[keyof typeof Prefecture];


export function PrefectureFromJSON(json: any): Prefecture {
    return PrefectureFromJSONTyped(json, false);
}

export function PrefectureFromJSONTyped(json: any, ignoreDiscriminator: boolean): Prefecture {
    return json as Prefecture;
}

export function PrefectureToJSON(value?: Prefecture | null): any {
    return value as any;
}

