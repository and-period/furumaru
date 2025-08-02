import type { Meta, StoryObj } from '@storybook/vue3';

import FurumaruItem from './FurumaruItem.vue';

const meta: Meta = {
	title: 'FurumaruItem',
	component: FurumaruItem,
	tags: ['autodocs'],
	argTypes: {},
	args: {},
} satisfies Meta<typeof FurumaruItem>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	args: {
		name: 'Furumaru',
	},
}
