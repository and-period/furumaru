import type {
  Cart,
  CartItemsInner,
  Coordinator,
  Product,
  ProductMediaInner,
} from '../api'

export interface ImageItem {
  url: string
  size: number
}

export interface MediaItem {
  url: string
  isThumbnail: boolean
  images: ImageItem[]
}

export interface ProductItem extends Product {
  thumbnail: ProductMediaInner
}

export interface CartItem extends CartItemsInner {
  product: ProductItem
}

export interface ShoppingCart extends Cart {
  boxType: string
  boxSize: number
  useRate: number
  coordinator: Coordinator
  items: CartItem[]
}
