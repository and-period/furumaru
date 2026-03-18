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
  render: () => ({
    components: { FmDialog },
    setup: () => {
      const isOpen = ref(true);
      return { isOpen };
    },
    template: `
      <div>
        <button @click="isOpen = true" style="padding: 8px 16px; background: #d97a38; color: white;">ダイアログを開く</button>
        <FmDialog v-model:open="isOpen" title="注文を確定しますか？" confirm-text="確定する" cancel-text="キャンセル">
          <p>この操作を実行してもよろしいですか？</p>
        </FmDialog>
      </div>
    `,
  }),
};

export const Danger: Story = {
  render: () => ({
    components: { FmDialog },
    setup: () => {
      const isOpen = ref(true);
      return { isOpen };
    },
    template: `
      <div>
        <button @click="isOpen = true" style="padding: 8px 16px; background: #f44336; color: white;">削除</button>
        <FmDialog v-model:open="isOpen" title="商品を削除しますか？" confirm-text="削除する" cancel-text="キャンセル" variant="danger">
          <p>この商品をカートから削除します。この操作は取り消せません。</p>
        </FmDialog>
      </div>
    `,
  }),
};
