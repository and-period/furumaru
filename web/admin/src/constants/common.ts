import { CharacterEncodingType } from '~/types'
import { TimeWeekday } from '~/types/api/v1'

export const weekdays = [
  { title: '日曜日', value: TimeWeekday.Sunday },
  { title: '月曜日', value: TimeWeekday.Monday },
  { title: '火曜日', value: TimeWeekday.Tuesday },
  { title: '水曜日', value: TimeWeekday.Wednesday },
  { title: '木曜日', value: TimeWeekday.Thursday },
  { title: '金曜日', value: TimeWeekday.Friday },
  { title: '土曜日', value: TimeWeekday.Saturday },
]

export const characterEncodingTypes = [
  { title: 'UTF-8', value: CharacterEncodingType.UTF8 },
  { title: 'Shift-JIS', value: CharacterEncodingType.ShiftJIS },
]
