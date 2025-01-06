<script setup lang="ts">
import type { GuestCheckoutExperienceRequest } from '~/types/api'

const formData = ref<GuestCheckoutExperienceRequest>({
  requestId: '',
  billingAddressId: '',
  promotionCode: '',
  adultCount: 0,
  juniorHighSchoolCount: 0,
  elementarySchoolCount: 0,
  preschoolCount: 0,
  seniorCount: 0,
  transportation: '',
  requestedDate: '',
  requestedTime: '',
  paymentMethod: 0,
  callbackUrl: '',
  total: 0,
  email: '',
  billingAddress: {
    lastname: '',
    firstname: '',
    lastnameKana: '',
    firstnameKana: '',
    postalCode: '',
    prefectureCode: 0,
    city: '',
    addressLine1: '',
    addressLine2: '',
    phoneNumber: '',
  },
  creditCard: {
    name: '',
    number: '',
    month: 0,
    year: 0,
    verificationValue: '',
  },
})
</script>

<template>
  <div>
    購入画面
  </div>
</template>
