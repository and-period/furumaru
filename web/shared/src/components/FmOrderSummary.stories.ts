import type { Meta, StoryObj } from '@storybook/vue3';
import FmOrderSummary from './FmOrderSummary.vue';

const meta = {
  title: 'Components/FmOrderSummary',
  component: FmOrderSummary,
  parameters: {
    layout: 'padded',
  },
  tags: ['autodocs'],
  argTypes: {
    isLoading: {
      control: 'boolean',
    },
    enableCoupon: {
      control: 'boolean',
    },
    appliedCouponCode: {
      control: 'text',
    },
    invalidCoupon: {
      control: 'boolean',
    },
    subtotal: {
      control: 'number',
    },
    discount: {
      control: 'number',
    },
    shippingFee: {
      control: 'number',
    },
    total: {
      control: 'number',
    },
  },
} satisfies Meta<typeof FmOrderSummary>;

export default meta;
type Story = StoryObj<typeof meta>;

// Mock data matching the design in the provided image
const mockItems = [
  {
    productId: '1',
    quantity: 1,
    product: {
      id: '1',
      name: 'たまねぎ 500g',
      price: 3000,
      thumbnail: {
        url: 'https://via.placeholder.com/56x56/90EE90/000000?text=玉ねぎ',
        isThumbnail: true,
      },
    },
  },
  {
    productId: '2',
    quantity: 1,
    product: {
      id: '2',
      name: 'レモン 500g',
      price: 3000,
      thumbnail: {
        url: 'https://via.placeholder.com/56x56/FFD700/000000?text=レモン',
        isThumbnail: true,
      },
    },
  },
  {
    productId: '3',
    quantity: 1,
    product: {
      id: '3',
      name: '卵 500g',
      price: 3000,
      thumbnail: {
        url: 'https://via.placeholder.com/56x56/F5DEB3/000000?text=卵',
        isThumbnail: true,
      },
    },
  },
];

const mockCoordinator = {
  id: '1',
  marcheName: '大崎上島マルシェ',
  username: '藤中 拓弥',
  prefecture: '広島県',
  city: '豊田郡大崎上島町',
};

const mockCarts = [
  { id: '1' },
  { id: '2' },
];

// Default story
export const Default: Story = {
  args: {
    items: mockItems,
    coordinator: mockCoordinator,
    carts: mockCarts,
    subtotal: 9000,
    discount: 0,
    total: 9000,
    enableCoupon: true,
    isLoading: false,
  },
};

// With Applied Coupon
export const WithAppliedCoupon: Story = {
  args: {
    ...Default.args,
    appliedCouponCode: 'SAVE10',
    discount: 900,
    total: 8100,
  },
};

// With Invalid Coupon
export const WithInvalidCoupon: Story = {
  args: {
    ...Default.args,
    invalidCoupon: true,
  },
};

// With Shipping Fee
export const WithShippingFee: Story = {
  args: {
    ...Default.args,
    shippingFee: 500,
    total: 9500,
  },
};

// Loading State
export const Loading: Story = {
  args: {
    ...Default.args,
    isLoading: true,
  },
};

// Without Coupon Functionality
export const WithoutCoupon: Story = {
  args: {
    ...Default.args,
    enableCoupon: false,
  },
};

// Single Item
export const SingleItem: Story = {
  args: {
    items: [mockItems[0]],
    coordinator: mockCoordinator,
    carts: [{ id: '1' }],
    subtotal: 3000,
    discount: 0,
    total: 3000,
    enableCoupon: true,
  },
};