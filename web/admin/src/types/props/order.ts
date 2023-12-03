import type { ShippingCarrier } from '../api'

export interface Order {
  name: string
  value: string
}

export interface OrderItems {
  images: []
  isThumbnail: boolean
  url: string
}

export interface FulfillmentInput {
  fulfillmentId: string
  shippingCarrier: ShippingCarrier
  trackingNumber: string
}
