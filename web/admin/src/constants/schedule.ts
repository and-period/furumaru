import { ScheduleStatus } from '~/types/api/v1'

export const scheduleStatuses = [
  { title: '非公開', value: ScheduleStatus.ScheduleStatusPrivate },
  { title: '申請中', value: ScheduleStatus.ScheduleStatusInProgress },
  { title: '開催前', value: ScheduleStatus.ScheduleStatusWaiting },
  { title: '開催中', value: ScheduleStatus.ScheduleStatusLive },
  { title: '終了(アーカイブ)', value: ScheduleStatus.ScheduleStatusClosed },
]
