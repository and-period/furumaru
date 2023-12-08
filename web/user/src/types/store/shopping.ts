import type {
  CalcCartResponse,
  Cart,
  CartItem as CartItemInner,
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
  thumbnail: ProductMediaInner | undefined
}

export interface CartItem extends CartItemInner {
  product: ProductItem | undefined
}

export interface ShoppingCart extends Cart {
  boxType: string
  boxSize: number
  useRate: number
  coordinator: Coordinator
  items: CartItem[]
}

export interface CalcCartItem extends CartItem {
  quantity: number
}

export interface CalcCart extends CalcCartResponse {
  items: CalcCartItem[]
}
