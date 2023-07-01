<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useCustomerStore } from '~/store'

const route = useRoute()
const customerStore = useCustomerStore()

const customerId = route.params.id as string

const { customer } = storeToRefs(customerStore)

const loading = ref<boolean>(false)

const fetchState = useAsyncData(async () => {
  await customerStore.getCustomer(customerId)
})

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-customer-edit
    :loading="loading"
    :customer="customer"
  />
</template>
