import dayjs, { unix } from 'dayjs'
import type { CreateNotificationRequest, UpdateNotificationRequest } from '~/types/api/v1'
import type { DateTimeInput } from '~/types/props'

export const useNotificationForm = (
  formData: Ref<CreateNotificationRequest | UpdateNotificationRequest>,
) => {
  const timeDataValue = computed({
    get: (): DateTimeInput => ({
      date: unix(formData.value.publishedAt).format('YYYY-MM-DD'),
      time: unix(formData.value.publishedAt).format('HH:mm'),
    }),
    set: (timeData: DateTimeInput): void => {
      const publishedAt = dayjs(`${timeData.date} ${timeData.time}`)
      formData.value.publishedAt = publishedAt.unix()
    },
  })

  const onChangePublishedAt = (): void => {
    const publishedAt = dayjs(`${timeDataValue.value.date} ${timeDataValue.value.time}`)
    formData.value.publishedAt = publishedAt.unix()
  }

  return {
    timeDataValue,
    onChangePublishedAt,
  }
}
