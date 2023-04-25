import { ImageSize } from '~/types/api'

interface Image {
  size: ImageSize
  url: string
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
