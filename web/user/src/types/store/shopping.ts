import type {
  CalcCartResponse,
  Cart,
  CartItem as CartItemInner,
  Coordinator,
  PaymentSystem,
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

// UI上で表示する商品の型定義
export interface ProductItem extends Product {
  thumbnail: ProductMediaInner | undefined
}

// UI上で表示するカートの中身の型定義
export interface CartItem extends CartItemInner {
  product: ProductItem | undefined
}

// UI上で表示するカートの型定義
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

// UI上で表示する支払いシステムの状態の型定義
export interface PaymentSystemStatus extends PaymentSystem {
  methodName: string
}
