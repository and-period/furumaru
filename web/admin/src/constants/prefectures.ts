import { Prefecture } from '~/types/api'

export interface PrefecturesListItem {
  id: number
  text: string
  value: Prefecture
}

/**
 * 都道府県一覧
 */
export const prefecturesList: PrefecturesListItem[] = [
  { id: 1, text: '北海道', value: Prefecture.HOKKAIDO },
  { id: 2, text: '青森県', value: Prefecture.AOMORI },
  { id: 3, text: '岩手県', value: Prefecture.IWATE },
  { id: 4, text: '宮城県', value: Prefecture.MIYAGI },
  { id: 5, text: '秋田県', value: Prefecture.AKITA },
  { id: 6, text: '山形県', value: Prefecture.YAMAGATA },
  { id: 7, text: '福島県', value: Prefecture.FUKUSHIMA },
  { id: 8, text: '茨城県', value: Prefecture.IBARAKI },
  { id: 9, text: '栃木県', value: Prefecture.TOCHIGI },
  { id: 10, text: '群馬県', value: Prefecture.GUNMA },
  { id: 11, text: '埼玉県', value: Prefecture.SAITAMA },
  { id: 12, text: '千葉県', value: Prefecture.CHIBA },
  { id: 13, text: '東京都', value: Prefecture.TOKYO },
  { id: 14, text: '神奈川県', value: Prefecture.KANAGAWA },
  { id: 15, text: '新潟県', value: Prefecture.NIIGATA },
  { id: 16, text: '富山県', value: Prefecture.TOYAMA },
  { id: 17, text: '石川県', value: Prefecture.ISHIKAWA },
  { id: 18, text: '福井県', value: Prefecture.FUKUI },
  { id: 19, text: '山梨県', value: Prefecture.YAMANASHI },
  { id: 20, text: '長野県', value: Prefecture.NAGANO },
  { id: 21, text: '岐阜県', value: Prefecture.GIFU },
  { id: 22, text: '静岡県', value: Prefecture.SHIZUOKA },
  { id: 23, text: '愛知県', value: Prefecture.AICHI },
  { id: 24, text: '三重県', value: Prefecture.MIE },
  { id: 25, text: '滋賀県', value: Prefecture.SHIGA },
  { id: 26, text: '京都府', value: Prefecture.KYOTO },
  { id: 27, text: '大阪府', value: Prefecture.OSAKA },
  { id: 28, text: '兵庫県', value: Prefecture.HYOGO },
  { id: 29, text: '奈良県', value: Prefecture.NARA },
  { id: 30, text: '和歌山県', value: Prefecture.WAKAYAMA },
  { id: 31, text: '鳥取県', value: Prefecture.TOTTORI },
  { id: 32, text: '島根県', value: Prefecture.SHIMANE },
  { id: 33, text: '岡山県', value: Prefecture.OKAYAMA },
  { id: 34, text: '広島県', value: Prefecture.HIROSHIMA },
  { id: 35, text: '山口県', value: Prefecture.YAMAGUCHI },
  { id: 36, text: '徳島県', value: Prefecture.TOKUSHIMA },
  { id: 37, text: '香川県', value: Prefecture.KAGAWA },
  { id: 38, text: '愛媛県', value: Prefecture.EHIME },
  { id: 39, text: '高知県', value: Prefecture.KOCHI },
  { id: 40, text: '福岡県', value: Prefecture.FUKUOKA },
  { id: 41, text: '佐賀県', value: Prefecture.SAGA },
  { id: 42, text: '長崎県', value: Prefecture.NAGASAKI },
  { id: 43, text: '熊本県', value: Prefecture.KUMAMOTO },
  { id: 44, text: '大分県', value: Prefecture.OITA },
  { id: 45, text: '宮崎県', value: Prefecture.MIYAZAKI },
  { id: 46, text: '鹿児島県', value: Prefecture.KAGOSHIMA },
  { id: 47, text: '沖縄県', value: Prefecture.OKINAWA }
]
