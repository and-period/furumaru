export interface LiveItemMock {
  id: string
  title: string
  imgSrc: string
  startAt: number
  published: boolean
  canceled: boolean
  marcheName: string
  address: string
  cnName: string
  cnImgSrc: string
}

export interface LiveVideoMock {
  id: string
  title: string
  videoSrc: string
  imgSrc: string
  startAt: number
  published: boolean
  isArchive: boolean
  canceled: boolean
  description: string
  marcheName: string
  address: string
  cnName: string
  cnImgSrc: string
}

export const MOCK_LIVE_ITEMS: LiveItemMock[] = [
  {
    id: '1',
    imgSrc: '/img/live/demo.jpg',
    title: '[秋の大収穫マルシェ］有機農家が厳選した野菜セット',
    startAt: 1700884800,
    published: true,
    canceled: false,
    marcheName: 'ふるさと掘り起こし隊',
    address: '千葉県 山武市',
    cnName: 'ふるマル太郎',
    cnImgSrc: '/img/recommend/demo.jpeg',
  },
  {
    id: '2',
    imgSrc: '/img/live/2.png',
    title: '［大崎上島マルシェ］瀬戸内海の島の美味しいものを紹介します！',
    startAt: 1665396001,
    published: true,
    canceled: false,
    marcheName: '大崎上島マルシェ',
    address: '広島県 豊田郡 大崎上島町',
    cnName: '藤中 拓弥',
    cnImgSrc: '/img/recommend/cn.png',
  },
  {
    id: '3',
    imgSrc: '/img/live/3.png',
    title: '［大崎上島マルシェ］瀬戸内海の島の美味しいものを紹介します！',
    startAt: 1665396001,
    published: false,
    canceled: false,
    marcheName: '大崎上島マルシェ',
    address: '広島県 豊田郡 大崎上島町',
    cnName: '藤中 拓弥',
    cnImgSrc: '/img/recommend/cn.png',
  },
  {
    id: '4',
    imgSrc: '/img/live/4.png',
    title: '［大崎上島マルシェ］瀬戸内海の島の美味しいものを紹介します！',
    startAt: 1665396001,
    published: false,
    canceled: false,
    marcheName: '大崎上島マルシェ',
    address: '広島県 豊田郡 大崎上島町',
    cnName: '藤中 拓弥',
    cnImgSrc: '/img/recommend/cn.png',
  },
  {
    id: '5',
    imgSrc: '/img/live/5.png',
    title: '［大崎上島マルシェ］瀬戸内海の島の美味しいものを紹介します！',
    startAt: 1665396001,
    published: false,
    canceled: false,
    marcheName: '大崎上島マルシェ',
    address: '広島県 豊田郡 大崎上島町',
    cnName: '藤中 拓弥',
    cnImgSrc: '/img/recommend/cn.png',
  },
  {
    id: '6',
    imgSrc: '/img/live/6.png',
    title: '［大崎上島マルシェ］瀬戸内海の島の美味しいものを紹介します！',
    startAt: 1665396001,
    published: false,
    canceled: false,
    marcheName: '大崎上島マルシェ',
    address: '広島県 豊田郡 大崎上島町',
    cnName: '藤中 拓弥',
    cnImgSrc: '/img/recommend/cn.png',
  },
  {
    id: '7',
    imgSrc: '/img/live/7.png',
    title: '［大崎上島マルシェ］瀬戸内海の島の美味しいものを紹介します！',
    startAt: 1665396001,
    published: false,
    canceled: false,
    marcheName: '大崎上島マルシェ',
    address: '広島県 豊田郡 大崎上島町',
    cnName: '藤中 拓弥',
    cnImgSrc: '/img/recommend/cn.png',
  },
]

export const MOCK_LIVE_VIDEO: LiveVideoMock = {
  id: '1',
  imgSrc: '/img/live/6.png',
  videoSrc: 'https://devimages.apple.com/iphone/samples/bipbop/bipbopall.m3u8',
  title: '[秋の大収穫マルシェ］有機農家が厳選した野菜セット',
  startAt: 1700884800,
  published: true,
  isArchive: false,
  canceled: false,
  description:
    '千葉県の山武市で有機農業を20年営んでいます。\n育てる野菜は年間70品目以上！',
  marcheName: 'ふるさと掘り起こし隊',
  address: '千葉県 山武市',
  cnName: 'ふるマル太郎',
  cnImgSrc: '/img/recommend/demo.jpeg',
}
