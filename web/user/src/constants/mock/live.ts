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
    imgSrc: '/img/live/1.png',
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
  title: '［大崎上島マルシェ］瀬戸内海の島の美味しいものを紹介します！',
  startAt: 1665396001,
  published: true,
  isArchive: false,
  canceled: false,
  description:
    'はじめまして。大崎上島マルシェです。\nレモンをはじめ柑橘類で有名な広島県の離島・大崎上島で採れたおいしいものをご紹介します。\n質問等も大歓迎ですので、コメントはお気軽にどうぞ！\nごゆっくりとお買い物をお楽しみください。',
  marcheName: '大崎上島マルシェ',
  address: '広島県 豊田郡 大崎上島町',
  cnName: '藤中 拓弥',
  cnImgSrc: '/img/recommend/cn.png',
}
