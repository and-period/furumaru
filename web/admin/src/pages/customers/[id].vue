<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { VDataTable } from 'vuetify/lib/labs/components'
import { useCustomerStore } from '~/store'

const route = useRoute()
const customerStore = useCustomerStore()

const customerId = route.params.id as string

const fetchState = useAsyncData(async () => {
  await customerStore.fetchCustomer(customerId)
})

const { customer } = storeToRefs(customerStore)

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-customer-show :customer="customer" />
</template>
