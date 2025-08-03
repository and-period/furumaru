import type { Meta, StoryObj } from '@storybook/vue3';

import FmProductItem from './FmProductItem.vue';

const meta: Meta = {
	title: 'FurumaruItem',
	component: FmProductItem,
	tags: ['autodocs'],
	argTypes: {},
	args: {},
} satisfies Meta<typeof FmProductItem>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	args: {
		name: 'Furumaru',
		price: 3000,
		stoke: 15,
		thumbnailUrl: 'https://assets.furumaru.and-period.co.jp/products/media/image/sUJfjmMoMvGwjht9QgdWZK.jpg',
	},
}

export const SoldOut: Story = {
	args: {
		...Default.args,
		stoke: 0
	}
}

export const VideoThumbnail: Story = {
	args: {
		...Default.args,
		thumbnailUrl: 'https://assets.furumaru.and-period.co.jp/products/media/video/87guGSpjE2nrAMUR8UHxzN.mp4#t=0.1',
	}
}
