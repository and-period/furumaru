export interface LiveMock {
  id: string
  title: string
  imgSrc: string
  startAt: number
  published: boolean
  canceled: boolean
}

export const MOCK_LIVE_ITEMS: LiveMock[] = [
  {
    id: '1',
    imgSrc: '/img/live/1.png',
    title: '［大崎上島マルシェ］瀬戸内海の島の美味しいものを紹介します！',
    startAt: 1665396001,
    published: true,
    canceled: false
  },
  {
    id: '2',
    imgSrc: '/img/live/2.png',
    title: '［大崎上島マルシェ］瀬戸内海の島の美味しいものを紹介します！',
    startAt: 1665396001,
    published: true,
    canceled: false
  },
  {
    id: '3',
    imgSrc: '/img/live/3.png',
    title: '［大崎上島マルシェ］瀬戸内海の島の美味しいものを紹介します！',
    startAt: 1665396001,
    published: false,
    canceled: false
  },
  {
    id: '4',
    imgSrc: '/img/live/4.png',
    title: '［大崎上島マルシェ］瀬戸内海の島の美味しいものを紹介します！',
    startAt: 1665396001,
    published: false,
    canceled: false
  },
  {
    id: '5',
    imgSrc: '/img/live/5.png',
    title: '［大崎上島マルシェ］瀬戸内海の島の美味しいものを紹介します！',
    startAt: 1665396001,
    published: false,
    canceled: false
  },
  {
    id: '6',
    imgSrc: '/img/live/6.png',
    title: '［大崎上島マルシェ］瀬戸内海の島の美味しいものを紹介します！',
    startAt: 1665396001,
    published: false,
    canceled: false
  }
]
