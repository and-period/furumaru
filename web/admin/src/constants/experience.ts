import { ExperienceStatus } from '~/types/api/v1'

export const experienceStatues = [
  { title: '体験受付開始前', value: ExperienceStatus.ExperienceStatusWaiting },
  { title: '体験受付中', value: ExperienceStatus.ExperienceStatusAccepting },
  { title: '体験受付上限', value: ExperienceStatus.ExperienceStatusSoldOut },
  { title: '非公開', value: ExperienceStatus.ExperienceStatusPrivate },
  { title: '体験受付終了', value: ExperienceStatus.ExperienceStatusFinished },
  { title: '不明', value: ExperienceStatus.ExperienceStatusUnknown },
]

export const experiencePublicationStatuses = [
  { title: '公開', value: true },
  { title: '非公開', value: false },
]

export const experienceSoldStatus = [
  { title: '販売中', value: false },
  { title: '在庫なし', value: true },
]
