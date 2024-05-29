import { type PrefecturesListItem, prefecturesList } from '~/constants'
import type { Prefecture } from '~/types/api'

export function findPrefecture(code?: Prefecture): PrefecturesListItem | undefined {
  if (!code) {
    return undefined
  }
  return prefecturesList.find((item: PrefecturesListItem): boolean => {
    return item.value === code
  })
}
