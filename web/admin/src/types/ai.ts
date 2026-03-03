/** AIチャットパネルに渡すフォームデータの型 */
export interface ProductFormDataForAi {
  name: string
  description: string
  price: number
  cost: number
  inventory: number
  weight: number
  itemUnit: string
  itemDescription: string
  deliveryType: number
  storageMethodType: number
  expirationDate: number
  recommendedPoint1: string
  recommendedPoint2: string
  recommendedPoint3: string
  originPrefectureCode: number
  originCity: string
  scope: number
  box60Rate: number
  box80Rate: number
  box100Rate: number
}

/** Tool Call で返されるフォーム更新データ（全フィールドoptional） */
export type ProductFormUpdate = Partial<ProductFormDataForAi>

/** フォーム更新プレビュー用の変更フィールド */
export interface FormFieldChange {
  field: string
  label: string
  oldValue: unknown
  newValue: unknown
}
