import { VideoStatus } from '~/types/api/v1'

export const videoStatuses = [
  { title: '非公開', value: VideoStatus.VideoStatusPrivate },
  { title: '限定公開', value: VideoStatus.VideoStatusLimited },
  { title: '公開中', value: VideoStatus.VideoStatusPublished },
  { title: '公開予定', value: VideoStatus.VideoStatusWaiting },
]
