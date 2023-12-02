/**
 * 国際規格の電話番号を日本の電話番号に変換する関数
 * @param phoneNumber 国際規格の電話番号
 * @returns 日本の電話番号
 */
export const convertI18nToJapanesePhoneNumber = (phoneNumber: string) => {
  return phoneNumber.replace('+81', '0')
}

/**
 * 日本の電話番号を国際規格の電話番号に変換する関数
 * @param phoneNumber 日本の電話番号
 * @returns 国際規格の電話番号
 */
export const convertJapaneseToI18nPhoneNumber = (phoneNumber: string) => {
  return phoneNumber.replace('0', '+81')
}
