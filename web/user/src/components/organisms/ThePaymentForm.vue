<script setup lang="ts">
import dayjs from 'dayjs'
import type { GuestCheckoutCreditCard } from '~/types/api'
import type { I18n } from '~/types/locales'

const i18n = useI18n()

const ct = (str: keyof I18n['purchase']['confirmation']) => {
  return i18n.t(`purchase.confirmation.${str}`)
}

const checkoutFormData = defineModel<GuestCheckoutCreditCard>({ required: true })

interface Props {
  formId: string
}

defineProps<Props>()

const creditCardMonthValue = computed({
  get: () => {
    if (checkoutFormData.value.month === 0) {
      return '0'
    }
    if (checkoutFormData.value.month < 10) {
      return `0${checkoutFormData.value.month}`
    }
    else {
      return String(checkoutFormData.value.month)
    }
  },
  set: (val: string) => {
    const month = Number(val)
    if (!isNaN(month)) {
      checkoutFormData.value.month = month
    }
  },
})

const handleSubmitCreditCardForm = () => {}
</script>

<template>
  <component
    :is="formId ? 'form' : 'div'"
    :id="formId"
    @submit.prevent="handleSubmitCreditCardForm"
  >
    <div class="mt-4 flex w-full items-center gap-4">
      <the-text-input
        v-model="checkoutFormData.number"
        :placeholder="ct('creditCardNumberPlaceholder')"
        :with-label="false"
        name="cc-number"
        type="text"
        class="w-full"
        pattern="[0-9]*"
        required
      />
      <div
        class="hidden h-[24px] min-w-max items-center gap-2 md:flex"
      >
        <img
          src="~/assets/img/cc/visa.png"
          alt="visa icon"
          class="h-full"
        >
        <img
          src="~/assets/img/cc/jcb.png"
          alt="jcb icon"
          class="h-full"
        >
        <img
          src="~/assets/img/cc/amex.png"
          alt="amex icon"
          class="h-full"
        >
        <img
          src="~/assets/img/cc/master.png"
          alt="master icon"
          class="h-full"
        >
      </div>
    </div>
    <the-text-input
      v-model="checkoutFormData.name"
      :placeholder="ct('cardholderNamePlaceholder')"
      :with-label="false"
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
          {{ ct("expirationMonthPlaceholder") }}
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
        v-model="checkoutFormData.year"
        class="mb-1 block w-full appearance-none rounded-none border-b border-main bg-transparent px-1 py-2 text-inherit focus:outline-none"
      >
        <option
          value="0"
          disabled
        >
          {{ ct("expirationYearPlaceholder") }}
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
    <the-text-input
      v-model="checkoutFormData.verificationValue"
      :placeholder="ct('securityCodePlaceholder')"
      :with-label="false"
      name="cc-csc"
      type="password"
      pattern="[0-9]*"
      class="mt-4 w-1/2"
      required
    />
  </component>
</template>
