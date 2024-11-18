import { VideoStatus } from '~/types/api'

/**
 * 動画のステータス情報を文字列に変換する関数
 * @param status
 * @returns
 */
export function videoStatusToString(status: VideoStatus): string {
  switch (status) {
    case VideoStatus.UNKNOWN:
      return '不明'
    case VideoStatus.PRIVATE:
      return '非公開'
    case VideoStatus.WAITING:
      return '公開前'
    case VideoStatus.LIMITED:
      return '限定公開'
    case VideoStatus.PUBLISHED:
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
    case VideoStatus.UNKNOWN:
      return 'gray'
    case VideoStatus.PRIVATE:
      return 'red'
    case VideoStatus.WAITING:
      return 'info'
    case VideoStatus.LIMITED:
      return 'secondary'
    case VideoStatus.PUBLISHED:
      return 'primary'
    default:
      return 'gray'
  }
}
