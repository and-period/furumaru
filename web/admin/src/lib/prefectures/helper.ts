import { prefecturesList } from '~/constants'
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
