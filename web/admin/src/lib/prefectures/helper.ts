import { prefecturesList, PrefecturesListItem } from '~/constants'
import { CreateShippingRate } from '~/types/api'

interface PrefecturesListSelectItems extends PrefecturesListItem {
  disabled: boolean
}

/**
 * 選択可能な都道府県のリストを返す関数
 * @param items
 * @param index
 * @returns
 */
export function getSelectablePrefecturesList (
  items: CreateShippingRate[],
  index: number
): PrefecturesListSelectItems[] {
  const unselectedPrefecturesList: string[] = [
    ...items
      .filter((_, i) => i !== index)
      .map(item => item.prefectures)
      .flat()
  ]
  return prefecturesList.map((item) => {
    return {
      ...item,
      disabled: unselectedPrefecturesList.includes(item.value)
    }
  })
}
