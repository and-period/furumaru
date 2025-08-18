import type { Meta, StoryObj } from '@storybook/vue3';

import FmProductDetail from './FmProductDetail.vue';

const meta: Meta = {
	title: 'FmProductDetail',
	component: FmProductDetail,
	tags: ['autodocs'],
	argTypes: {
		mediaFiles: {
			description: '商品の画像ファイルのパスのリスト',
			control: { type: 'object' }
		},
		name: {
			description: '商品名',
			control: { type: 'text' }
		},
		description: {
			description: '商品説明',
			control: { type: 'text' }
		},
		originPrefecture: {
			description: '商品の産地情報（都道府県）',
			control: { type: 'text' }
		},
		originCity: {
			description: '商品の産地情報（市区町村）',
			control: { type: 'text' }
		},
		rating: {
			description: '商品の評価情報',
			control: { type: 'object' }
		},
		recommendedPoint1: {
			description: '商品のおすすめポイント1',
			control: { type: 'text' }
		},
		recommendedPoint2: {
			description: '商品のおすすめポイント2',
			control: { type: 'text' }
		},
		recommendedPoint3: {
			description: '商品のおすすめポイント3',
			control: { type: 'text' }
		},
		expirationDate: {
			description: '商品の賞味期間情報（日数）',
			control: { type: 'number' }
		},
		weight: {
			description: '商品の重さ（kg）',
			control: { type: 'number' }
		},
		deliveryType: {
			description: '商品の配送タイプの情報',
			control: { type: 'select' },
			options: [0, 1, 2, 3]
		},
		storageMethodType: {
			description: '商品の保存方法の情報',
			control: { type: 'select' },
			options: [0, 1, 2, 3, 4]
		},
	},
	args: {},
} satisfies Meta<typeof FmProductDetail>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	args: {
		mediaFiles: [
			{
				url: 'https://assets.furumaru.and-period.co.jp/products/media/image/sUJfjmMoMvGwjht9QgdWZK.jpg',
				isThumbnail: true
			},
			{
				url: 'https://assets.furumaru.and-period.co.jp/products/media/image/j3pDMfg6vKAVijrjJnXghL.jpg',
				isThumbnail: false
			}
		],
		name: 'ふるマル厳選 北海道産アスパラガス',
		description: '北海道の豊かな大地で育った新鮮なアスパラガスをお届けします。\n甘くて柔らかい食感が特徴の春の味覚をぜひご堪能ください。\n\n生産者が愛情を込めて育てた自慢の一品です。',
		originPrefecture: '北海道',
		originCity: '札幌市',
		rating: {
			average: 4.5,
			count: 12,
			detail: {
				'1': 0,
				'2': 1,
				'3': 2,
				'4': 4,
				'5': 5
			}
		},
		recommendedPoint1: '朝採れの新鮮なアスパラガスをその日のうちに発送',
		recommendedPoint2: '甘くて柔らかい食感で、茹でるだけで美味しくお召し上がりいただけます',
		recommendedPoint3: '栄養価が高く、ビタミンやミネラルが豊富に含まれています',
		expirationDate: 5,
		weight: 1.0,
		deliveryType: 2,
		storageMethodType: 3,
	},
};

export const WithVideoThumbnail: Story = {
	args: {
		...Default.args,
		mediaFiles: [
			{
				url: 'https://assets.furumaru.and-period.co.jp/products/media/video/87guGSpjE2nrAMUR8UHxzN.mp4#t=0.1',
				isThumbnail: true
			},
			{
				url: 'https://assets.furumaru.and-period.co.jp/products/media/image/j3pDMfg6vKAVijrjJnXghL.jpg',
				isThumbnail: false
			}
		],
	},
};

export const MinimalInfo: Story = {
	args: {
		mediaFiles: [
			{
				url: 'https://assets.furumaru.and-period.co.jp/products/media/image/sUJfjmMoMvGwjht9QgdWZK.jpg',
				isThumbnail: true
			}
		],
		name: 'シンプル商品',
		description: '基本的な商品情報のみを持つサンプル商品です。',
		originPrefecture: '東京都',
		originCity: '渋谷区',
	},
};

export const FullInfo: Story = {
	args: {
		...Default.args,
		mediaFiles: [
			{
				url: 'https://assets.furumaru.and-period.co.jp/products/media/image/sUJfjmMoMvGwjht9QgdWZK.jpg',
				isThumbnail: true
			},
			{
				url: 'https://assets.furumaru.and-period.co.jp/products/media/image/j3pDMfg6vKAVijrjJnXghL.jpg',
				isThumbnail: false
			},
			{
				url: 'https://assets.furumaru.and-period.co.jp/products/media/video/87guGSpjE2nrAMUR8UHxzN.mp4#t=0.1',
				isThumbnail: false
			}
		],
		name: '瀬戸内・牡蠣セット（殻付き10～15個＆剥き身500g）',
		description: '瀬戸内海で育った新鮮な牡蠣をセットでお届けします。\n殻付きの牡蠣は焼き牡蠣や蒸し牡蠣でお楽しみいただけます。\n剥き身はフライやパスタなど様々な料理にご活用ください。\n\n詳しい調理方法は https://example.com/recipes をご確認ください。',
		originPrefecture: '広島県',
		originCity: '呉市',
		rating: {
			average: 4.8,
			count: 25,
			detail: {
				'1': 0,
				'2': 0,
				'3': 1,
				'4': 6,
				'5': 18
			}
		},
		recommendedPoint1: '瀬戸内海の恵まれた環境で育った高品質な牡蠣',
		recommendedPoint2: '殻付きと剥き身の両方をセットでお得にお届け',
		recommendedPoint3: '生産者直送で鮮度抜群、届いたその日にお召し上がりいただけます',
		expirationDate: 3,
		weight: 2.5,
		deliveryType: 2,
		storageMethodType: 3,
	},
};