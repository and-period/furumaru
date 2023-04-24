import { ImageSize } from '~/types/api'

class Image {
  size: ImageSize
  url: string

  constructor (size: ImageSize, url: string) {
    this.size = size
    this.url = url
  }
}

export function getResizedImages<T extends Image> (images: T[]): string {
  const strs: string[] = images.map((img: T): string => {
    switch (img.size) {
      case ImageSize.SMALL:
        return `${img.url} 1x`
      case ImageSize.MEDIUM:
        return `${img.url} 2x`
      case ImageSize.LARGE:
        return `${img.url} 3x`
      default:
        return img.url
    }
  })
  return strs.join(', ')
}
