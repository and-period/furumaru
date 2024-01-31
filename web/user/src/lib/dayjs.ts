import dayjs from 'dayjs'

/**
 * unixtimeからUIで表示する時刻の形式に変換する関数
 * @param unixtime
 * @returns
 */
export function datetimeformatterFromUnixtime(unixtime: number): string {
  return dayjs.unix(unixtime).format('YYYY/MM/DD HH:mm')
}
