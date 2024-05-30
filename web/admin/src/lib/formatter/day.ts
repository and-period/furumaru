import { unix } from 'dayjs'

/**
 * unixtime表記の数値をYYYY/MM/DD HH:mm表記文字列に変換する関数
 * @param unixtime
 * @returns
 */
export function dateTimeFormatter(unixtime: number): string {
  return unix(unixtime).format('YYYY/MM/DD HH:mm')
}
