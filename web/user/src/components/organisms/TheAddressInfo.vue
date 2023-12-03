<script setup lang="ts">
import { convertI18nToJapanesePhoneNumber } from '~/lib/phone-number'
import { Address } from '~/types/api'

interface Props {
  address: Address
}

const props = defineProps<Props>()

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
  <dl class="grid grid-cols-3 gap-2">
    <dt>氏名</dt>
    <dd class="col-span-2">{{ displayName }}</dd>
    <dt>電話番号</dt>
    <dd class="col-span-2">
      {{ convertI18nToJapanesePhoneNumber(address.phoneNumber) }}
    </dd>
    <dt>住所</dt>
    <dd class="col-span-2">{{ displayAddress }}</dd>
  </dl>
</template>
