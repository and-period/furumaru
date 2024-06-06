export interface CartItemMock {
  id: string
  name: string
  inventory: number
  price: number
  imgSrc: string
}

export interface CartMock {
  marche: string
  boxType: string
  boxSize: number
  items: CartItemMock[]
}

export const MOCK_CART_ITEMS: CartMock[] = [
  {
    marche: '大崎上島',
    boxType: '常温・冷蔵',
    boxSize: 100,
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
    marche: '大崎上島',
    boxType: '常温・冷蔵',
    boxSize: 60,
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
  {
    marche: '大崎上島',
    boxType: '常温・冷蔵',
    boxSize: 80,
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
]
