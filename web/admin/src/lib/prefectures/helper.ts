import { prefecturesList, TOTAL_PREFECTURE_COUNT } from '~/constants'
import type { PrefecturesListItem } from '~/constants'
import type { UpdateDefaultShippingRate, UpsertShippingRate } from '~/types/api/v1'

export interface PrefecturesListSelectItems extends PrefecturesListItem {
  disabled: boolean
}

/**
 * 選択可能な都道府県のリストを返す関数
 * @param items
 * @param index
 * @returns
 */
export function getSelectablePrefecturesList(
  items: UpdateDefaultShippingRate[] | UpsertShippingRate[],
  index: number,
): PrefecturesListSelectItems[] {
  const unselectedPrefecturesList: number[] = [
    ...items
      .filter((_, i) => i !== index)
      .map(item => item.prefectureCodes)
      .flat(),
  ]
  return prefecturesList.map((item) => {
    return {
      ...item,
      disabled: unselectedPrefecturesList.includes(item.value),
    }
  })
}

/**
 * すべての都道府県がカバーされているかを検証する関数
 * @param items 配送料設定の配列
 * @returns すべての都道府県（47件）がカバーされている場合はtrue
 */
export function hasAllPrefecturesCovered(
  items: UpdateDefaultShippingRate[] | UpsertShippingRate[],
): boolean {
  const allCodes = new Set(items.flatMap(item => item.prefectureCodes))
  return allCodes.size === TOTAL_PREFECTURE_COUNT
}
