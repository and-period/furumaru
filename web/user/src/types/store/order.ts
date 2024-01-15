import type { Coordinator, Order, OrderItem, Product } from '~/types/api'

export interface OrderHistoryItem extends OrderItem {
  product: Product
}

export interface OrderHistory extends Order {
  coordinator: Coordinator | undefined
}
