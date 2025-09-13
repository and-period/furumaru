/**
 * 注文集計期間
 * @export
 * @enum {string}
 */
export const TopOrderPeriodType = {
  /**
   * 1日単位
   */
  DAY: 'day',
  /**
   * 1週間単位
   */
  WEEK: 'week',
  /**
   * 1ヶ月単位
   */
  MONTH: 'month',
} as const

export type TopOrderPeriodType = (typeof TopOrderPeriodType)[keyof typeof TopOrderPeriodType]

/**
 * 文字コード種別
 * @export
 * @enum {string}
 */
export const CharacterEncodingType = {
  /**
   * UTF-8
   */
  UTF8: 0,
  /**
   * Shift-JIS
   */
  ShiftJIS: 1,
} as const

export type CharacterEncodingType = (typeof CharacterEncodingType)[keyof typeof CharacterEncodingType]

/**
 * 都道府県コード
 * @export
 * @enum {string}
 */
export const Prefecture = {
  /**
   * 不明
   */
  UNKNOWN: 0,
  /**
   * 北海道
   */
  HOKKAIDO: 1,
  /**
   * 青森県
   */
  AOMORI: 2,
  /**
   * 岩手県
   */
  IWATE: 3,
  /**
   * 宮城県
   */
  MIYAGI: 4,
  /**
   * 秋田県
   */
  AKITA: 5,
  /**
   * 山形県
   */
  YAMAGATA: 6,
  /**
   * 福島県
   */
  FUKUSHIMA: 7,
  /**
   * 茨城県
   */
  IBARAKI: 8,
  /**
   * 栃木県
   */
  TOCHIGI: 9,
  /**
   * 群馬県
   */
  GUNMA: 10,
  /**
   * 埼玉県
   */
  SAITAMA: 11,
  /**
   * 千葉県
   */
  CHIBA: 12,
  /**
   * 東京都
   */
  TOKYO: 13,
  /**
   * 神奈川県
   */
  KANAGAWA: 14,
  /**
   * 新潟県
   */
  NIIGATA: 15,
  /**
   * 富山県
   */
  TOYAMA: 16,
  /**
   * 石川県
   */
  ISHIKAWA: 17,
  /**
   * 福井県
   */
  FUKUI: 18,
  /**
   * 山梨県
   */
  YAMANASHI: 19,
  /**
   * 長野県
   */
  NAGANO: 20,
  /**
   * 岐阜県
   */
  GIFU: 21,
  /**
   * 静岡県
   */
  SHIZUOKA: 22,
  /**
   * 愛知県
   */
  AICHI: 23,
  /**
   * 三重県
   */
  MIE: 24,
  /**
   * 滋賀県
   */
  SHIGA: 25,
  /**
   * 京都府
   */
  KYOTO: 26,
  /**
   * 大坂府
   */
  OSAKA: 27,
  /**
   * 兵庫県
   */
  HYOGO: 28,
  /**
   * 奈良県
   */
  NARA: 29,
  /**
   * 和歌山県
   */
  WAKAYAMA: 30,
  /**
   * 鳥取県
   */
  TOTTORI: 31,
  /**
   * 島根県
   */
  SHIMANE: 32,
  /**
   * 岡山県
   */
  OKAYAMA: 33,
  /**
   * 広島県
   */
  HIROSHIMA: 34,
  /**
   * 山口県
   */
  YAMAGUCHI: 35,
  /**
   * 徳島県
   */
  TOKUSHIMA: 36,
  /**
   * 香川県
   */
  KAGAWA: 37,
  /**
   * 愛媛県
   */
  EHIME: 38,
  /**
   * 高知県
   */
  KOCHI: 39,
  /**
   * 福岡県
   */
  FUKUOKA: 40,
  /**
   * 佐賀県
   */
  SAGA: 41,
  /**
   * 長崎県
   */
  NAGASAKI: 42,
  /**
   * 熊本県
   */
  KUMAMOTO: 43,
  /**
   * 大分県
   */
  OITA: 44,
  /**
   * 宮崎県
   */
  MIYAZAKI: 45,
  /**
   * 鹿児島県
   */
  KAGOSHIMA: 46,
  /**
   * 沖縄県
   */
  OKINAWA: 47,
} as const

export type Prefecture = (typeof Prefecture)[keyof typeof Prefecture]
