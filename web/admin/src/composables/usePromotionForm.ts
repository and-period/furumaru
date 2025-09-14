import dayjs, { unix } from 'dayjs'
import type { CreatePromotionRequest, UpdatePromotionRequest } from '~/types/api/v1'
import type { DateTimeInput } from '~/types/props'

export const usePromotionForm = (
  formData: Ref<CreatePromotionRequest | UpdatePromotionRequest>,
) => {
  const startTimeDataValue = computed({
    get: (): DateTimeInput => ({
      date: unix(formData.value.startAt).format('YYYY-MM-DD'),
      time: unix(formData.value.startAt).format('HH:mm'),
    }),
    set: (timeData: DateTimeInput): void => {
      const startAt = dayjs(`${timeData.date} ${timeData.time}`)
      formData.value.startAt = startAt.unix()
    },
  })

  const endTimeDataValue = computed({
    get: (): DateTimeInput => ({
      date: unix(formData.value.endAt).format('YYYY-MM-DD'),
      time: unix(formData.value.endAt).format('HH:mm'),
    }),
    set: (timeData: DateTimeInput): void => {
      const endAt = dayjs(`${timeData.date} ${timeData.time}`)
      formData.value.endAt = endAt.unix()
    },
  })

  const onChangeStartAt = (): void => {
    const startAt = dayjs(`${startTimeDataValue.value.date} ${startTimeDataValue.value.time}`)
    formData.value.startAt = startAt.unix()
  }

  const onChangeEndAt = (): void => {
    const endAt = dayjs(`${endTimeDataValue.value.date} ${endTimeDataValue.value.time}`)
    formData.value.endAt = endAt.unix()
  }

  const generateRandomCode = (): string => {
    const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
    let result = ''
    const charactersLength = characters.length
    for (let i = 0; i < 8; i++) {
      result += characters.charAt(Math.floor(Math.random() * charactersLength))
    }
    return result
  }

  return {
    startTimeDataValue,
    endTimeDataValue,
    onChangeStartAt,
    onChangeEndAt,
    generateRandomCode,
  }
}
