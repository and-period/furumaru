import type { Meta, StoryObj } from '@storybook/vue3';

import FmButton from './FmButton.vue';

const meta: Meta = {
  title: 'FmButton',
  component: FmButton,
  tags: ['autodocs'],
  argTypes: {
    variant: { control: 'select', options: ['primary', 'secondary', 'danger', 'ghost'] },
    size: { control: 'select', options: ['sm', 'md', 'lg'] },
  },
  args: {},
} satisfies Meta<typeof FmButton>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Primary: Story = {
  args: { variant: 'primary' },
  render: args => ({
    components: { FmButton },
    setup: () => ({ args }),
    template: '<FmButton v-bind="args">カゴに入れる</FmButton>',
  }),
};

export const Secondary: Story = {
  args: { variant: 'secondary' },
  render: args => ({
    components: { FmButton },
    setup: () => ({ args }),
    template: '<FmButton v-bind="args">詳しく見る</FmButton>',
  }),
};

export const Danger: Story = {
  args: { variant: 'danger' },
  render: args => ({
    components: { FmButton },
    setup: () => ({ args }),
    template: '<FmButton v-bind="args">削除する</FmButton>',
  }),
};

export const Ghost: Story = {
  args: { variant: 'ghost' },
  render: args => ({
    components: { FmButton },
    setup: () => ({ args }),
    template: '<FmButton v-bind="args">キャンセル</FmButton>',
  }),
};

export const Loading: Story = {
  args: { variant: 'primary', loading: true },
  render: args => ({
    components: { FmButton },
    setup: () => ({ args }),
    template: '<FmButton v-bind="args">送信中...</FmButton>',
  }),
};

export const Disabled: Story = {
  args: { variant: 'primary', disabled: true },
  render: args => ({
    components: { FmButton },
    setup: () => ({ args }),
    template: '<FmButton v-bind="args">送信</FmButton>',
  }),
};

export const Sizes: Story = {
  render: () => ({
    components: { FmButton },
    template: `
      <div style="display: flex; align-items: center; gap: 12px;">
        <FmButton size="sm">小</FmButton>
        <FmButton size="md">中</FmButton>
        <FmButton size="lg">大</FmButton>
      </div>
    `,
  }),
};
