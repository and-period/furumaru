import { helpers } from '@vuelidate/validators'

/**
 * ひらがなを判定する関数
 * @param value
 * @returns
 */
const kanaValidator = (value: string): boolean => {
  const kanaRegex = /^[\u3040-\u309F]+$/
  return kanaRegex.test(value)
}

/**
 * ひらがなを判定するカスタムバリデーションルール
 */
const kana = helpers.withMessage('ひらがなを入力してください。', kanaValidator)

export { kana }
