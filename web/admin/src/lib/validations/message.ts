import type { ErrorObject } from '@vuelidate/core'

/**
 * バリデーションオブジェクトから指定したキーのエラーメッセージを取得する関数
 * @param errors ErrorObjectの配列
 * @returns
 */
const getErrorMessage = (errors: ErrorObject[]): string => {
  return errors.length > 0 ? errors[0].$message.toString() : ''
}

export { getErrorMessage }
