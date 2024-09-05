import type { Product, ProductMediaInner } from '~/types/api'

/**
 * 商品のサムネイル情報を取得する関数
 */
export const getProductThumbnailUrl = (product: Product): string => {
  const thumbnail = product.media?.find((media: ProductMediaInner) => {
    return media.isThumbnail
  })
  return thumbnail ? thumbnail.url : ''
}
