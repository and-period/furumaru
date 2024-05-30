import type {
  MessageProps } from '@vuelidate/validators'
import {
  helpers,
  required as _required,
  email as _email,
  minLength as _minLength,
  maxLength as _maxLength,
  minValue as _minValue,
  maxValue as _maxValue,
  sameAs as _sameAs,
} from '@vuelidate/validators'

import { isRef } from 'vue'

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
  helpers.regex(/^0[0-9]{9,10}$/),
)

/**
 * 郵便番号を判定するカスタムバリデーション
 */
const postalCode = helpers.withMessage(
  '郵便番号は数字のみの7桁で入力してください。',
  helpers.regex(/[0-9]{7}/),
)

/**
 * 必須バリデーションルールをラップする関数
 */
const required = helpers.withMessage(
  (_: MessageProps) => 'この項目は必須項目です。',
  _required,
)

/**
 * メールアドレスバリデーションルールをラップする関数
 */
const email = helpers.withMessage(
  'メールアドレスの形式で入力してください。',
  _email,
)

/**
 * 文字列の最小の長さのバリデーションルールをラップする関数
 * @param min
 * @returns
 */
const minLength = (min: number | Ref<number>) =>
  helpers.withMessage(
    ({ $params }: MessageProps) => `${$params.min}文字以上入力してください。`,
    _minLength(min),
  )

/**
 * 文字列の最大の長さのバリデーションルールをラップする関数
 * @param max
 * @returns
 */
const maxLength = (max: number | Ref<number>) =>
  helpers.withMessage(
    ({ $params }: MessageProps) => `${$params.max}文字までです。`,
    _maxLength(max),
  )

/**
 * 最小値のバリデーションルールをラップする関数
 * @param min
 * @returns
 */
const minValue = (min: string | number | Ref<number> | Ref<string>) =>
  helpers.withMessage(
    ({ $params }: MessageProps) => `${$params.min}以上で入力してください。`,
    _minValue(min),
  )

/**
 * 最大値のバリデーションルールをラップする関数
 * @param min
 * @returns
 */
const maxValue = (max: string | number | Ref<number> | Ref<string>) =>
  helpers.withMessage(
    ({ $params }: MessageProps) => `${$params.max}以下で入力してください。`,
    _maxValue(max),
  )

/**
 * 入力値が別のプロパティとの一致しているかのバリデーションルールをラップする関数
 * @param equalTo
 * @param otherName
 * @returns
 */
const sameAs = (equalTo: unknown, otherName?: string) =>
  helpers.withMessage(
    (_: MessageProps) => '一致しません。',
    _sameAs(equalTo, otherName),
  )

/**
 * 入力値が別のプロパティとの一致していないかのバリデーションルール
 * @param notEqualTo
 * @param otherName
 * @returns
 */
const notSameAs = (notEqualTo: unknown, otherName?: string) =>
  helpers.withMessage(
    `${otherName}と違う値を入力してください`,
    helpers.withParams(
      { type: 'notSameAs', notEqualTo },
      value => value !== notEqualTo,
    ),
  )

/**
 * 配列の最小の長さのバリデーションルール（カスタム）
 * @param minLength 最小の要素数
 * @returns
 */
const minLengthArray = (minLength: number | Ref<number>) => {
  const minLengthValue = isRef(minLength) ? minLength.value : minLength
  return helpers.withMessage(
    `少なくとも${minLengthValue}つ以上選択してください。`,
    helpers.withParams(
      { type: 'minLengthArray', minLength },
      value => Array.isArray(value) && value.length >= minLengthValue,
    ),
  )
}

/**
 * 配列の最代の長さのバリデーションルール（カスタム）
 * @param maxLength 最小の要素数
 * @returns
 */
const maxLengthArray = (maxLength: number | Ref<number>) => {
  const maxLengthValue = isRef(maxLength) ? maxLength.value : maxLength
  return helpers.withMessage(
    `この項目は最大${maxLengthValue}件まです。${maxLengthValue}件以下に修正してください。`,
    helpers.withParams(
      { type: 'minLengthArray', minLength: maxLength },
      value => Array.isArray(value) && value.length <= maxLengthValue,
    ),
  )
}

export {
  kana,
  tel,
  required,
  email,
  postalCode,
  minLength,
  maxLength,
  minValue,
  maxValue,
  sameAs,
  notSameAs,
  minLengthArray,
  maxLengthArray,
}
