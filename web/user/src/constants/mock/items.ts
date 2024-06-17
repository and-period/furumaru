export interface ProductItemMock {
  id: string
  name: string
  inventory: number
  price: number
  thumbnail: string
  address: string
  cnName: string
  cnthumbnail: string
  status?: number
  coordinatorId?: string
}
export const MOCK_ALL_PRODUCT_ITEMS: ProductItemMock[] = [
  {
    id: '1',
    name: '【7-8月限定 !!】無農薬ブルーベリー　大崎上島産',
    inventory: 10,
    price: 1600,
    thumbnail:
      'web/user/src/pages/items/peach.jpeg',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
    status: 0,
    coordinatorId: '1',
  },
  {
    id: '2',
    name: '【こだわりの減農薬栽培】瀬戸内グリーンレモン（ワックス・防腐剤不使用）',
    inventory: 10,
    price: 3000,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/files/IMG_9161.heic?v=1697199066',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '一房の彩り',
    cnthumbnail: '/img/recommend/cn.png',
    status: 2,
  },
  {
    id: '3',
    name: '【ギフトに最適】大崎上島レモンケーキ「ナチュールシトロン（5個入り）」',
    inventory: 10,
    price: 2000,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/products/DSC00107.jpg?v=1680105559',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
    status: 3,
  },
  {
    id: '4',
    name: '【シャリシャリほろ苦がクセになる！】大崎上島産　紅八朔',
    inventory: 10,
    price: 3000,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/products/DSC03167-2.jpg?v=1680106082',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
    status: 4,
  },
  {
    id: '5',
    name: '【ジューシーな瀬戸内の宝物】大崎上島産レモンピューレ　鮮度抜群　長期保存可能　使いやすいスパウト容器タイプ',
    inventory: 10,
    price: 10800,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/files/image_17886617-5ae5-461a-b426-fc1c340e63d6.jpg?v=1687798346',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '6',
    name: '【予約商品】「大崎上島産グリーンレモン」〜とにかく香りに癒されんさい〜',
    inventory: 10,
    price: 2800,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/files/unnamed_1_a5291f29-038c-4c8f-ab8a-a16a5e8a56c5.jpg?v=1695276263',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '一房の彩り',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '7',
    name: '【予約販売】山梨県産　樹上完熟の極上シャインマスカット「山昭農園」',
    inventory: 10,
    price: 6900,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/files/11_9f84f0d7-399a-4b62-a6a2-3a5756f931e4.jpg?v=1693396099',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '山昭農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '8',
    name: '【予約販売】山梨県産　樹上完熟の極上シャインマスカット「山昭農園」【贈答用】',
    inventory: 10,
    price: 8800,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/files/11_e5519c86-e05a-402c-8383-003d44bd5bce.jpg?v=1693009232',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '一房の彩り',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '9',
    name: '【初夏の爽やかセット】大崎上島産「甘夏」4キロ　大崎上島産レモン1キロ',
    inventory: 10,
    price: 2500,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/products/IMG-1270.jpg?v=1650724864',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '10',
    name: '【大崎上島の隠れた名産品生はちみつ！】無添加/非加熱「谷口さんの純粋蜂蜜」【120g】',
    inventory: 10,
    price: 1500,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/products/E9E7DFE5-CD0B-4B07-8D8E-98BB36B118E7.jpg?v=1671372824',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '11',
    name: '【大崎上島の隠れた名産品！】　無添加/非加熱　百花蜂蜜【250g】',
    inventory: 10,
    price: 1700,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/products/IMG-9289_1_f5768c2e-9c1b-4125-955a-bd2513c5a0a1.jpg?v=1680106556',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '12',
    name: '【大崎上島の隠れた名産品！】　無添加/非加熱　百花蜂蜜【600g】',
    inventory: 10,
    price: 3500,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/products/IMG-7003_88a8ddec-9fff-4012-8e79-4b4ea6e98768.png?v=1680106707',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '13',
    name: '【大崎上島産「甘夏」　5キロ】〜こだわり農家のこれ食べんさい！〜',
    inventory: 10,
    price: 3000,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/products/IMG-7123_1.jpg?v=1680107086',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '14',
    name: '【大崎上島産】　温州みかん　〜食べて応援地方の食材〜',
    inventory: 10,
    price: 1400,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/products/image_84c2b364-3ac3-4a68-bb03-36cc73b4f938.heic?v=1671373072',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '15',
    name: '【季節・数量限定品！！】無添加/非加熱「谷口さんの純粋蜂蜜3本ギフトセット」【120g×3本】',
    inventory: 10,
    price: 3200,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/products/IMG-2570_Square.jpg?v=1680107394',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '16',
    name: '【専用商品】',
    inventory: 10,
    price: 4200,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/products/0BFF103E-1FA1-4B52-8BE9-35C8F0B83B3A_d3065093-b9b6-46d0-b548-96849883f5b8.jpg?v=1676592944',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '17',
    name: '【島の誇り】大崎上島産　極早生みかん　〜シーズン開始！！〜',
    inventory: 10,
    price: 1400,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/products/IMG-9399.jpg?v=1671373426',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '18',
    name: '【希少な島の真っ赤なみかん】田中さんの「プリンス清見」',
    inventory: 10,
    price: 1600,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/products/IMG-0496.jpg?v=1674577423',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '19',
    name: '【数々の賞を受賞】本場ブラジル人オーナーの手がける絶品の「ミナスチーズ」<フレッシュミナスチーズ240g×2個、熟成チーズ220g×2個>',
    inventory: 10,
    price: 6250,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/files/IMG_0895_dab0491e-4097-40d2-9981-8c0a955e470b.jpg?v=1692854009',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: 'ビルミルク',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '20',
    name: '【数々の賞を受賞】本場ブラジル人オーナーの手がける絶品の「ミナスチーズ」<フレッシュミナスチーズ240g、熟成チーズ、クリームチーズ180g>',
    inventory: 10,
    price: 5400,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/files/IMG_0899_f236ab1d-7301-4050-8910-c57006e9e5b4.jpg?v=1692854576',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: 'ビルミルク',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '21',
    name: '【数々の賞を受賞】本場ブラジル人ーオーナーの手がける絶品の「ミナスチーズ」<フレッシュミナスチーズ240g×3個>',
    inventory: 10,
    price: 5200,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/files/IMG_0895.jpg?v=1693014152',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: 'ビルミルク',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '22',
    name: '【数量限定！】新鮮な果実でつくる大崎上島産ブルーベリージャム',
    inventory: 10,
    price: 1750,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/files/20230824_141811.jpg?v=1693405589',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '一房の彩り',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '23',
    name: '【新鮮な果実で作る無添加ジャム】大崎上島産オレンジジャム',
    inventory: 10,
    price: 600,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/products/image_df255074-714b-408f-b060-835e76335608.heic?v=1680107995',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '24',
    name: '【有機農家が厳選】秋野菜セット７種　千葉県山武市産',
    inventory: 10,
    price: 3980,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/files/2023-10-01151719.jpg?v=1696143047',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '三つ豆ファーム',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '25',
    name: '【柑橘の島】大崎上島産 極早生みかん　シーズン到来！！',
    inventory: 10,
    price: 2000,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/files/47fcf562fb8976bee329c23dfc975978.png?v=1697198949',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '26',
    name: '【甘さの爆弾！】大崎上島産　田中さんの「不知火　5キロ」',
    inventory: 10,
    price: 3500,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/products/IMG-0748.jpg?v=1680108296',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '27',
    name: '【繊細かつ華やかな香りの名脇役】神山すだち（ご家庭用B級品）　徳島県神山産',
    inventory: 10,
    price: 3480,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/files/1EDAA0C4-D23D-40DC-AF75-AA26072BF9B7.jpg?v=1692600696',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: 'NPO法人里山みらい',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '28',
    name: '【見た目からは想像できない甘さに驚き！】大崎上島産　はるか',
    inventory: 10,
    price: 3000,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/products/IMG_6933.jpg?v=1680887775',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '29',
    name: '【訳あり！お買い得品】大崎上島産「B級甘夏」　5キロ〜こだわり農家のこれ食べんさい！〜のコピー',
    inventory: 10,
    price: 1600,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/products/image_c03932c7-af02-441d-988c-9fc8c2f889a5.heic?v=1671373603',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '30',
    name: '【防腐剤・ワックス不使用】大崎上島産イエローレモン',
    inventory: 10,
    price: 6000,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/products/14615FED-B6BC-49B7-87D9-69CE2C631FEF.jpg?v=1680108407',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '31',
    name: '【離島の秘宝！】ジャンボ青パパイヤ　1キロ（レモンのおまけつき）',
    inventory: 10,
    price: 3500,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/products/7EADD242-2082-443F-8A3E-DA3865B69D6B.jpg?v=1671373672',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '32',
    name: '【非加熱！無添加！】純大崎上島産ハチミツレモン',
    inventory: 10,
    price: 700,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/products/image_17395ab0-1ebb-4c6f-85e0-37576342731b.heic?v=1680108996',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '33',
    name: '大崎上島産　伊予柑',
    inventory: 10,
    price: 900,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/products/1676119060835_873fa79a-9d14-43c6-8f18-632245731e14.jpg?v=1680109185',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '34',
    name: '大崎上島産　完熟紅八朔　4キロ【贈答用】',
    inventory: 10,
    price: 2800,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/products/image_e6854736-ec7a-4681-a116-4572c34dd976.heic?v=1680109389',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '35',
    name: '大崎上島産「ゴールデン文旦（大橘）」３キロ',
    inventory: 10,
    price: 2100,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/products/IMG_5257_5e796a03-f8fb-4cb0-ad91-e17c34fd4f41.jpg?v=1680109493',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '36',
    name: '大崎上島産はるみ【ちょっぴり贅沢な柑橘の王様】',
    inventory: 10,
    price: 3500,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/products/DSC03102-2.jpg?v=1676687633',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '37',
    name: '小玉はるみ【ちょっぴり贅沢な柑橘の王様】（シーズンラスト）',
    inventory: 10,
    price: 2500,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/products/DSC03076-2.jpg?v=1681536001',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '38',
    name: '旬の朝採り！上山産すだち【ご家庭用】',
    inventory: 10,
    price: 3300,
    thumbnail: NaN,
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: 'NPO法人里山みらい',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '39',
    name: '石地みかん【贈答用】〜毎年人気の大崎上島自慢のみかん〜',
    inventory: 10,
    price: 1800,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/products/IMG-6120.jpg?v=1671373913',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '40',
    name: '訳あり「はれひめ」【ご家庭用】',
    inventory: 10,
    price: 2000,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/products/image_9170f43c-008f-496e-9ab6-74174de25b1a.jpg?v=1672976938',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
  {
    id: '41',
    name: '限定クーポン対象品「谷口さんの純粋蜂蜜」【120g】大崎上島の隠れた名産品生はちみつ！無添加/非加熱',
    inventory: 10,
    price: 1650,
    thumbnail:
      'https://cdn.shopify.com/s/files/1/0593/3151/0447/products/E9E7DFE5-CD0B-4B07-8D8E-98BB36B118E7_5b973b49-af9e-4fe5-b953-e52ce2105dbf.jpg?v=1671373946',
    address: '大崎上島マルシェ\n広島県 豊田郡 大崎上島町',
    cnName: '瀬戸内大崎上島農園',
    cnthumbnail: '/img/recommend/cn.png',
  },
]
