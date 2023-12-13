import type { Producer, Product } from '../api'

export interface LiveTimeLineItem {
  comment: string
  startAt: number
  endAt: number
  producerId: string
  producer: Producer | undefined
  products: Product[]
  productIds: string[]
}
