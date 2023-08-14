import { CartItemMock } from './cart'

export interface PurchaseInnerItemMock {
  boxType: string
  boxSize: number
  useRate: number
  items: CartItemMock[]
}

export interface PurchaseMock {
  marche: string
  address: string
  sender: string
  cartItems: PurchaseInnerItemMock[]
}

export const MOCK_PURCHASE_ITEMS: PurchaseMock[] = [
  {
    marche: '大崎上島マルシェ',
    address: '広島県豊田郡大崎上島町',
    sender: '藤中 拓弥',
    cartItems: [
      {
        boxType: '常温・冷蔵',
        boxSize: 100,
        useRate: 95,
        items: [
          {
            id: '1',
            name: 'たまねぎ 500g',
            price: 3000,
            inventory: 10,
            imgSrc: '/img/recommend/1.png',
          },
          {
            id: '4',
            name: 'レモン 500g',
            inventory: 10,
            price: 3000,
            imgSrc: '/img/recommend/6.png',
          },
        ],
      },
      {
        boxType: '常温・冷蔵',
        boxSize: 60,
        useRate: 30,
        items: [
          {
            id: '5',
            name: '卵 500g',
            inventory: 10,
            price: 3000,
            imgSrc: '/img/recommend/5.png',
          },
        ],
      },
    ],
  },
  {
    marche: '東広島マルシェ',
    address: '広島県東広島市',
    sender: '藤中 拓弥',
    cartItems: [
      {
        boxType: '常温・冷蔵',
        boxSize: 80,
        useRate: 70,
        items: [
          {
            id: '2',
            name: '【冷凍】黒毛和牛 500g',
            inventory: 10,
            price: 3000,
            imgSrc: '/img/recommend/2.png',
          },
        ],
      },
    ],
  },
]
