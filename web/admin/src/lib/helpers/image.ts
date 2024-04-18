interface ImageSize {
  query: string
  size: string
}

const sizes: ImageSize[] = [
  { query: 'width=320&format=webp', size: '1x' },
  { query: 'width=640&format=webp', size: '2x' },
  { query: 'width=960&format=webp', size: '3x' }
]

export function getResizedImages<T extends string> (image: T): string {
  const urls: string[] = sizes.map((size: ImageSize): string => {
    return image ? `${image}?${size.query} ${size.size}` : ''
  })
  return urls.join(', ')
}
