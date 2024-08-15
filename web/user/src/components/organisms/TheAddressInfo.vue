<script setup lang="ts">
import { convertI18nToJapanesePhoneNumber } from '~/lib/phone-number'
import type { Address } from '~/types/api'
import type { I18n } from '~/types/locales'

interface Props {
  address: Address
}

const i18n = useI18n()
const props = defineProps<Props>()

const at = (str: keyof I18n['purchase']['address']) => {
  return i18n.t(`purchase.address.${str}`)
}

const displayName = computed<string>(() => {
  return `${props.address.lastname} ${props.address.firstname}（${props.address.lastnameKana} ${props.address.firstnameKana}）`
})

const displayAddress = computed<string>(() => {
  return `〒 ${props.address.postalCode} ${props.address.prefecture}${
    props.address.city
  }${props.address.addressLine1}${[props.address.addressLine2]}`
})
</script>

<template>
  <dl class="grid grid-cols-3 gap-2 text-[14px] tracking-[1.4px]">
    <dt>{{ at('nameLabel') }}</dt>
    <dd class="col-span-2">
      {{ displayName }}
    </dd>
    <dt>{{ at('phoneNumberLabel') }}</dt>
    <dd class="col-span-2">
      {{ convertI18nToJapanesePhoneNumber(address.phoneNumber) }}
    </dd>
    <dt>{{ at('addressLabel') }}</dt>
    <dd class="col-span-2">
      {{ displayAddress }}
    </dd>
  </dl>
</template>
