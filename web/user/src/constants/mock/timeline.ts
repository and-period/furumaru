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
    time: '19:00',
    cnImgSrc: '/img/recommend/cn.png',
    marchName: '大崎上島農園',
    description:
      '広島県の離島で柑橘を育てています。防腐剤やワックスを使わず、安心の美味しさをお届けできるよう心がけています。',
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
      '広島県の離島で柑橘を育てています。防腐剤やワックスを使わず、安心の美味しさをお届けできるよう心がけています。',
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
