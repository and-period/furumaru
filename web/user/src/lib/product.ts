import { ProductStatus } from '~/types/api'

/**
 * 商品のステータスに応じて文言を出しわける関数
 * @param status
 * @returns
 */
export function productStatusToString(status: ProductStatus): string {
  switch (status) {
    case ProductStatus.FOR_SALE:
      return '販売中'
    case ProductStatus.OUT_OF_SALES:
      return '販売終了'
    case ProductStatus.PRESALE:
      return '近日販売開始'
    default:
      return '無効な商品'
  }
}
