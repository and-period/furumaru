<script setup lang="ts">
import { convertI18nToJapanesePhoneNumber } from '~/lib/phone-number'
import type { GuestCheckoutAddress } from '~/types/api'
import { prefecturesList } from '~/constants/prefectures'

interface Props {
  address: GuestCheckoutAddress
}

const props = defineProps<Props>()

const findPrefectureValueById = computed<string>(() => {
  const prefecture = prefecturesList.find(item => item.id === props.address.prefectureCode)
  return prefecture?.text ?? '未設定'
})

const displayName = computed<string>(() => {
  return `${props.address.lastname} ${props.address.firstname}（${props.address.lastnameKana} ${props.address.firstnameKana}）`
})

const displayAddress = computed<string>(() => {
  return `〒 ${props.address.postalCode} ${findPrefectureValueById.value}${props.address.city}
  ${props.address.addressLine1}${props.address.addressLine2 ? ` ${props.address.addressLine2}` : ''}`
})
</script>

<template>
  <dl class="grid grid-cols-3 gap-2 text-[14px] tracking-[1.4px]">
    <dt>氏名</dt>
    <dd class="col-span-2">
      {{ displayName }}
    </dd>
    <dt>電話番号</dt>
    <dd class="col-span-2">
      {{ convertI18nToJapanesePhoneNumber(props.address.phoneNumber) }}
    </dd>
    <dt>住所</dt>
    <dd class="col-span-2">
      {{ displayAddress }}
    </dd>
  </dl>
</template>
