import type { Meta, StoryObj } from '@storybook/vue3';

import FmLoadingSpinner from './FmLoadingSpinner.vue';

const meta: Meta = {
  title: 'FmLoadingSpinner',
  component: FmLoadingSpinner,
  tags: ['autodocs'],
  argTypes: {
    size: { control: 'select', options: ['sm', 'md', 'lg'] },
  },
  args: {},
} satisfies Meta<typeof FmLoadingSpinner>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: { size: 'md', label: '読み込み中...' },
};

export const Small: Story = {
  args: { size: 'sm', label: '' },
};

export const Large: Story = {
  args: { size: 'lg', label: '商品情報を読み込み中...' },
};

export const Overlay: Story = {
  args: { size: 'lg', overlay: true, label: '処理中...' },
  decorators: [() => ({ template: '<div style="height: 300px; position: relative;"><story /></div>' })],
};
