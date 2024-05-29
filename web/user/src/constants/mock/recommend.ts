export interface RecommendMock {
  id: string
  name: string
  inventory: number
  price: number
  imgSrc: string
  address: string
  cnName: string
  cnImgSrc: string
}

export const MOCK_RECOMMEND_ITEMS: RecommendMock[] = [
  {
    id: '1',
    name: 'たまねぎ 500g',
    inventory: 10,
    price: 3000,
    imgSrc: '/img/recommend/1.png',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '藤中 拓弥',
    cnImgSrc: '/img/recommend/cn.png',
  },
  {
    id: '2',
    name: '【冷凍】黒毛和牛 500g',
    inventory: 10,
    price: 3000,
    imgSrc: '/img/recommend/2.png',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '藤中 拓弥',
    cnImgSrc: '/img/recommend/cn.png',
  },
  {
    id: '3',
    name: 'アスパラ 500g',
    inventory: 10,
    price: 3000,
    imgSrc: '/img/recommend/3.png',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '藤中 拓弥',
    cnImgSrc: '/img/recommend/cn.png',
  },
  {
    id: '4',
    name: '無農薬レモン 500g',
    inventory: 10,
    price: 3000,
    imgSrc: '/img/recommend/4.png',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '藤中 拓弥',
    cnImgSrc: '/img/recommend/cn.png',
  },
  {
    id: '5',
    name: '卵 500g',
    inventory: 10,
    price: 3000,
    imgSrc: '/img/recommend/5.png',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '藤中 拓弥',
    cnImgSrc: '/img/recommend/cn.png',
  },
]
