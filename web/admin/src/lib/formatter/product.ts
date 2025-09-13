import type { Experience, Product, ProductMedia } from '~/types/api/v1'

/**
 * 商品のサムネイル情報を取得する関数
 */
export const getProductThumbnailUrl = (product: Product): string => {
  const thumbnail = product.media?.find((media: ProductMedia) => {
    return media.isThumbnail
  })
  return thumbnail ? thumbnail.url : ''
}

/**
 * 体験のサムネイル情報を取得する関数
 * @param experience
 * @returns
 */
export const getExperienceThumbnailUrl = (experience: Experience): string => {
  const thumbnail = experience.media?.find((media: ProductMedia) => {
    return media.isThumbnail
  })
  return thumbnail ? thumbnail.url : ''
}
