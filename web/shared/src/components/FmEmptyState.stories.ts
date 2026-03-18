import type { Meta, StoryObj } from '@storybook/vue3';

import FmEmptyState from './FmEmptyState.vue';

const meta = {
  title: 'FmEmptyState',
  component: FmEmptyState,
  tags: ['autodocs'],
} satisfies Meta<typeof FmEmptyState>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    title: '商品が見つかりませんでした',
    description: '検索条件を変更して再度お試しください。',
    icon: '🔍',
    actionText: '商品一覧に戻る',
  },
};

export const EmptyCart: Story = {
  args: {
    title: 'カートは空です',
    description: '気になる商品を探してカートに追加してみましょう。',
    icon: '🛒',
    actionText: '商品を探す',
  },
};

export const NoOrders: Story = {
  args: {
    title: '注文履歴がありません',
    description: '商品を購入すると、こちらに注文履歴が表示されます。',
    icon: '📦',
    actionText: '商品を探す',
  },
};

export const NoReviews: Story = {
  args: {
    title: 'レビューはまだありません',
    description: 'この商品を購入してレビューを投稿してみませんか？',
    icon: '✏️',
  },
};
