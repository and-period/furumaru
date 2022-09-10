import { Ref } from '@vue/composition-api'
import {
  helpers,
  MessageProps,
  required as _required,
  email as _email,
  minLength as _minLength,
  maxLength as _maxLength,
} from '@vuelidate/validators'

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

/**
 * 電話番号を判定するカスタムバリデーション
 */
const tel = helpers.withMessage(
  '電話番号は0から始まる数字のみの9桁または10桁で入力してください。',
  helpers.regex(/^0[0-9]{9,10}$/)
)

/**
 * 郵便番号を判定するカスタムバリデーション
 */
const postalCode = helpers.withMessage(
  '郵便番号は数字のみの7桁で入力してください。',
  helpers.regex(/[0-9]{7}/)
)

/**
 * 必須バリデーションルールをラップする関数
 */
const required = helpers.withMessage(
  (_: MessageProps) => 'この項目は必須項目です。',
  _required
)

/**
 * メールアドレスバリデーションルールをラップする関数
 */
const email = helpers.withMessage(
  'メールアドレスの形式で入力してください。',
  _email
)

/**
 * 文字列の最小の長さのバリデーションルールをラップする関数
 * @param max
 * @returns
 */
const minLength = (max: number | Ref<number>) =>
  helpers.withMessage(
    ({ $params }: MessageProps) => `${$params.max}文字以上入力してください。`,
    _minLength(max)
  )

/**
 * 文字列の最大の長さのバリデーションルールをラップする関数
 * @param max
 * @returns
 */
const maxLength = (max: number | Ref<number>) =>
  helpers.withMessage(
    ({ $params }: MessageProps) => `${$params.max}文字までです。`,
    _maxLength(max)
  )

export { kana, tel, required, email, postalCode, minLength, maxLength }
