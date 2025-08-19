<script setup lang="ts">
import { computed } from 'vue'
import dayjs from 'dayjs'
import FmTextInput from './FmTextInput.vue'
import type { CreditCardData } from '../types/index'

// Use dynamic imports to handle assets in library build
import visaIcon from '../assets/img/cc/visa.png'
import jcbIcon from '../assets/img/cc/jcb.png'
import amexIcon from '../assets/img/cc/amex.png'
import masterIcon from '../assets/img/cc/master.png'

interface Props {
  formId?: string
  creditCardNumberPlaceholder?: string
  cardholderNamePlaceholder?: string
  expirationMonthPlaceholder?: string
  expirationYearPlaceholder?: string
  securityCodePlaceholder?: string
  showCardIcons?: boolean
}

withDefaults(defineProps<Props>(), {
  formId: '',
  creditCardNumberPlaceholder: 'クレジットカード番号',
  cardholderNamePlaceholder: '名義人名',
  expirationMonthPlaceholder: '月',
  expirationYearPlaceholder: '年',
  securityCodePlaceholder: 'セキュリティコード',
  showCardIcons: true,
})

const creditCardData = defineModel<CreditCardData>({ required: true })

const creditCardMonthValue = computed({
  get: () => {
    if (creditCardData.value.month === 0) {
      return '0'
    }
    if (creditCardData.value.month < 10) {
      return `0${creditCardData.value.month}`
    }
    else {
      return String(creditCardData.value.month)
    }
  },
  set: (val: string) => {
    const month = Number(val)
    if (!isNaN(month)) {
      creditCardData.value.month = month
    }
  },
})

const handleSubmit = () => {
  // Let parent handle the submit event
}
</script>

<template>
  <component
    :is="formId ? 'form' : 'div'"
    :id="formId"
    @submit.prevent="handleSubmit"
  >
    <div class="mt-4 flex w-full items-center gap-4">
      <FmTextInput
        v-model="creditCardData.number"
        :placeholder="creditCardNumberPlaceholder"
        name="cc-number"
        type="text"
        class="w-full"
        pattern="[0-9]*"
        required
      />
      <div
        v-if="showCardIcons"
        class="hidden h-[24px] min-w-max items-center gap-2 md:flex"
      >
        <img
          :src="visaIcon"
          alt="visa icon"
          class="h-full"
        >
        <img
          :src="jcbIcon"
          alt="jcb icon"
          class="h-full"
        >
        <img
          :src="amexIcon"
          alt="amex icon"
          class="h-full"
        >
        <img
          :src="masterIcon"
          alt="master icon"
          class="h-full"
        >
      </div>
    </div>
    <FmTextInput
      v-model="creditCardData.name"
      :placeholder="cardholderNamePlaceholder"
      name="cc-name"
      type="text"
      class="mt-2 w-full"
      required
    />
    <div class="mt-2 flex gap-4">
      <select
        v-model="creditCardMonthValue"
        class="mb-1 block w-full appearance-none rounded-none border-b border-main bg-transparent px-1 py-2 text-inherit focus:outline-none"
      >
        <option
          :value="0"
          disabled
        >
          {{ expirationMonthPlaceholder }}
        </option>
        <option
          v-for="i in 12"
          :key="i"
          :value="i < 10 ? `0${i}` : String(i)"
        >
          {{ i }}
        </option>
      </select>

      <select
        v-model="creditCardData.year"
        class="mb-1 block w-full appearance-none rounded-none border-b border-main bg-transparent px-1 py-2 text-inherit focus:outline-none"
      >
        <option
          value="0"
          disabled
        >
          {{ expirationYearPlaceholder }}
        </option>
        <option
          v-for="i in 11"
          :key="i"
          :value="dayjs().year() + i - 1"
        >
          {{ dayjs().year() + i - 1 }}
        </option>
      </select>
    </div>
    <FmTextInput
      v-model="creditCardData.verificationValue"
      :placeholder="securityCodePlaceholder"
      name="cc-csc"
      type="password"
      pattern="[0-9]*"
      class="mt-4 w-1/2"
      required
    />
  </component>
</template>