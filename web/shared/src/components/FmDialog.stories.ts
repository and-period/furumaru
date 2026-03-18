import type { Meta, StoryObj } from '@storybook/vue3';
import { ref } from 'vue';

import FmDialog from './FmDialog.vue';

const meta: Meta = {
  title: 'FmDialog',
  component: FmDialog,
  tags: ['autodocs'],
  argTypes: {
    variant: { control: 'select', options: ['default', 'danger'] },
  },
  args: {},
} satisfies Meta<typeof FmDialog>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    title: '注文を確定しますか？',
    confirmText: '確定する',
    cancelText: 'キャンセル',
    variant: 'default',
  },
  render: (args: Record<string, unknown>) => ({
    components: { FmDialog },
    setup: () => {
      const isOpen = ref(true);
      return { isOpen, args };
    },
    template: `
      <div>
        <button @click="isOpen = true" style="padding: 8px 16px; background: #d97a38; color: white;">ダイアログを開く</button>
        <FmDialog v-bind="args" v-model:open="isOpen">
          <p>この操作を実行してもよろしいですか？</p>
        </FmDialog>
      </div>
    `,
  }),
};

export const Danger: Story = {
  args: {
    title: '商品を削除しますか？',
    confirmText: '削除する',
    cancelText: 'キャンセル',
    variant: 'danger',
  },
  render: (args: Record<string, unknown>) => ({
    components: { FmDialog },
    setup: () => {
      const isOpen = ref(true);
      return { isOpen, args };
    },
    template: `
      <div>
        <button @click="isOpen = true" style="padding: 8px 16px; background: #f44336; color: white;">削除</button>
        <FmDialog v-bind="args" v-model:open="isOpen">
          <p>この商品をカートから削除します。この操作は取り消せません。</p>
        </FmDialog>
      </div>
    `,
  }),
};
