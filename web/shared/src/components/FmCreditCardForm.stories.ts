import type { Meta, StoryObj } from '@storybook/vue3'
import { ref } from 'vue'

import FmCreditCardForm from './FmCreditCardForm.vue'
import type { CreditCardData } from '../types/index'

const meta: Meta = {
  title: 'FmCreditCardForm',
  component: FmCreditCardForm,
  tags: ['autodocs'],
  argTypes: {
    showCardIcons: {
      control: { type: 'boolean' },
    },
  },
  args: {},
} satisfies Meta<typeof FmCreditCardForm>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  render: (args) => ({
    components: { FmCreditCardForm },
    setup() {
      const creditCardData = ref<CreditCardData>({
        name: '',
        number: '',
        month: 0,
        year: 0,
        verificationValue: '',
      })
      
      return { creditCardData, args }
    },
    template: '<FmCreditCardForm v-model="creditCardData" v-bind="args" />',
  }),
  args: {},
}

export const WithFormId: Story = {
  render: (args) => ({
    components: { FmCreditCardForm },
    setup() {
      const creditCardData = ref<CreditCardData>({
        name: '',
        number: '',
        month: 0,
        year: 0,
        verificationValue: '',
      })
      
      return { creditCardData, args }
    },
    template: '<FmCreditCardForm v-model="creditCardData" v-bind="args" />',
  }),
  args: {
    formId: 'credit-card-form',
  },
}

export const WithoutCardIcons: Story = {
  render: (args) => ({
    components: { FmCreditCardForm },
    setup() {
      const creditCardData = ref<CreditCardData>({
        name: '',
        number: '',
        month: 0,
        year: 0,
        verificationValue: '',
      })
      
      return { creditCardData, args }
    },
    template: '<FmCreditCardForm v-model="creditCardData" v-bind="args" />',
  }),
  args: {
    showCardIcons: false,
  },
}

export const WithCustomPlaceholders: Story = {
  render: (args) => ({
    components: { FmCreditCardForm },
    setup() {
      const creditCardData = ref<CreditCardData>({
        name: '',
        number: '',
        month: 0,
        year: 0,
        verificationValue: '',
      })
      
      return { creditCardData, args }
    },
    template: '<FmCreditCardForm v-model="creditCardData" v-bind="args" />',
  }),
  args: {
    creditCardNumberPlaceholder: 'Enter card number',
    cardholderNamePlaceholder: 'Enter cardholder name',
    expirationMonthPlaceholder: 'Month',
    expirationYearPlaceholder: 'Year',
    securityCodePlaceholder: 'CVV',
  },
}