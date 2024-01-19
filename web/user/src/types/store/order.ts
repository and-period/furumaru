import type { Coordinator, Order, OrderItem, Product } from '~/types/api'

export interface OrderHistoryItem extends OrderItem {
  product: Product | undefined
}

export interface OrderHistory extends Order {
  coordinator: Coordinator | undefined
  items: OrderHistoryItem[]
}
