export interface LiveTimelineItemMock {
  id: string
  name: string
  inventory: number
  price: number
  imgSrc: string
}

export interface LiveTimelineMock {
  id: string
  time: string
  cnImgSrc: string
  marchName: string
  description: string
  items: LiveTimelineItemMock[]
}

export const MOCK_LIVE_TIMELINES: LiveTimelineMock[] = [
  {
    id: '1',
    time: '13:00',
    cnImgSrc: '/img/recommend/demo.png',
    marchName: '三つ豆\nファーム ',
    description:
      '千葉県の山武市で有機農業を20年営んでいます。\n育てる野菜は年間70品目以上！',
    items: [
      {
        id: '2',
        name: 'かぶと里芋 500g',
        inventory: 10,
        price: 3000,
        imgSrc: '/img/recommend/8.png',
      },
      {
        id: '2',
        name: 'トマト 500g',
        inventory: 10,
        price: 1500,
        imgSrc: '/img/recommend/7.png',
      },
      {
        id: '3',
        name: 'たまねぎ 500g',
        price: 3000,
        inventory: 10,
        imgSrc: '/img/recommend/1.png',
      },
      {
        id: '4',
        name: 'アスパラ 500g',
        inventory: 10,
        price: 3000,
        imgSrc: '/img/recommend/3.png',
      },
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
    id: '2',
    time: '19:15',
    cnImgSrc: '/img/recommend/cn.png',
    marchName: '大崎上島農園',
    description:
      '千葉県の山武市で有機農業を20年営んでいます。\n育てる野菜は年間70品目以上！',
    items: [
      {
        id: '2',
        name: '無農薬レモン 500g',
        inventory: 10,
        price: 3000,
        imgSrc: '/img/recommend/4.png',
      },
      {
        id: '2',
        name: 'レモン 500g',
        inventory: 10,
        price: 3000,
        imgSrc: '/img/recommend/6.png',
      },
      {
        id: '3',
        name: 'たまねぎ 500g',
        price: 3000,
        inventory: 10,
        imgSrc: '/img/recommend/1.png',
      },
      {
        id: '4',
        name: 'アスパラ 500g',
        inventory: 10,
        price: 3000,
        imgSrc: '/img/recommend/3.png',
      },
    ],
  },
]
