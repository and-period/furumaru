import type { Meta, StoryObj } from '@storybook/vue3';

import FmTextInput from './FmTextInput.vue';

const meta: Meta = {
  title: 'FmTextInput',
  component: FmTextInput,
  tags: ['autodocs'],
  argTypes: {},
  args: {},
} satisfies Meta<typeof FmTextInput>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    placeholder: 'Enter text',
  },
}

export const Secret: Story = {
  args: {
    type: 'password',
    placeholder: 'Enter secret',
  },
}

export const Labelled: Story = {
  args: {
    label: 'Username',
    placeholder: 'Enter your username',
  },
}

export const WithMessage: Story = {
  args: {
    label: 'Email',
    type: 'email',
    placeholder: 'Enter your email',
    message: 'Please enter a valid email address.',
  },
}


export const HasError: Story = {
  args: {
    label: 'Email',
    type: 'email',
    placeholder: 'Enter your email',
    message: 'Please enter a valid email address.',
    error: true,
    errorMessage: 'Invalid email address.',
  },
}