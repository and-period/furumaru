<template>
  <the-producer-edit-form-page
    :form-data="formData"
    :form-data-loading="fetchState.pending"
  />
</template>

<script lang="ts">
import {
  defineComponent,
  reactive,
  useFetch,
  useRoute,
} from '@nuxtjs/composition-api'

import { useProducerStore } from '~/store/producer'
import { ProducerResponse } from '~/types/api'

export default defineComponent({
  setup() {
    const route = useRoute()
    const id = route.value.params.id

    const { getProducer } = useProducerStore()

    const formData = reactive<ProducerResponse>({
      id,
      lastname: '',
      lastnameKana: '',
      firstname: '',
      firstnameKana: '',
      addressLine1: '',
      addressLine2: '',
      city: '',
      prefecture: '',
      phoneNumber: '',
      postalCode: '',
      storeName: '',
      headerUrl: '',
      createdAt: -1,
      updatedAt: -1,
      thumbnailUrl: '',
      email: '',
    })

    const { fetchState } = useFetch(async () => {
      const producer = await getProducer(id)
      formData.lastname = producer.lastname
      formData.lastnameKana = producer.lastnameKana
      formData.firstname = producer.firstname
      formData.firstnameKana = producer.firstnameKana
      formData.addressLine1 = producer.addressLine1
      formData.addressLine2 = producer.addressLine2
      formData.city = producer.city
      formData.prefecture = producer.prefecture
      formData.phoneNumber = producer.phoneNumber
      formData.postalCode = producer.postalCode
      formData.storeName = producer.storeName
      formData.headerUrl = producer.headerUrl
      formData.thumbnailUrl = producer.thumbnailUrl
      formData.email = producer.email
      formData.createdAt = producer.createdAt
      formData.updatedAt = producer.updatedAt
    })

    return {
      id,
      fetchState,
      formData,
    }
  },
})
</script>
