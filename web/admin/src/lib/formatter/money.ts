/**
 * 数値を価格表記に変換する関数
 * @param n
 * @returns
 */
export function moneyFormat(n: number): string {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
  }).format(n)
}
