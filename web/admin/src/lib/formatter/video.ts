import { VideoStatus } from '~/types/api/v1'

/**
 * 動画のステータス情報を文字列に変換する関数
 * @param status
 * @returns
 */
export function videoStatusToString(status: VideoStatus): string {
  switch (status) {
    case VideoStatus.VideoStatusUnknown:
      return '不明'
    case VideoStatus.VideoStatusPrivate:
      return '非公開'
    case VideoStatus.VideoStatusWaiting:
      return '公開前'
    case VideoStatus.VideoStatusLimited:
      return '限定公開'
    case VideoStatus.VideoStatusPublished:
      return '公開済み'
    default:
      return '不明'
  }
}

/**
 * 動画のステータス情報を色に変換する関数
 * @param status
 * @returns
 */
export function videoStatusToColor(status: VideoStatus): string {
  switch (status) {
    case VideoStatus.VideoStatusUnknown:
      return 'gray'
    case VideoStatus.VideoStatusPrivate:
      return 'red'
    case VideoStatus.VideoStatusWaiting:
      return 'info'
    case VideoStatus.VideoStatusLimited:
      return 'secondary'
    case VideoStatus.VideoStatusPublished:
      return 'primary'
    default:
      return 'gray'
  }
}
