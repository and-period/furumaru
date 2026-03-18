import type { Meta, StoryObj } from '@storybook/vue3';
import { ref } from 'vue';

import FmDialog from './FmDialog.vue';

const meta = {
  title: 'FmDialog',
  component: FmDialog,
  tags: ['autodocs'],
  argTypes: {
    variant: { control: 'select', options: ['default', 'danger'] },
  },
} satisfies Meta<typeof FmDialog>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    open: true,
    title: '注文を確定しますか？',
    confirmText: '確定する',
    cancelText: 'キャンセル',
  },
  render: args => ({
    components: { FmDialog },
    setup: () => {
      const isOpen = ref(args.open);
      return { args, isOpen };
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
    open: true,
    title: '商品を削除しますか？',
    confirmText: '削除する',
    cancelText: 'キャンセル',
    variant: 'danger',
  },
  render: args => ({
    components: { FmDialog },
    setup: () => {
      const isOpen = ref(args.open);
      return { args, isOpen };
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
