<script lang="ts" setup>
import dayjs from 'dayjs'
import { storeToRefs } from 'pinia'
import { useAlert } from '~/lib/hooks'

import { useTopStore } from '~/store'
import { TopOrderPeriodType } from '~/types'

const router = useRouter()
const topStore = useTopStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const { orders, pendingOrders, pendingOrdersTotal, hasPublishedProduct, hasPublishedExperience } = storeToRefs(topStore)

const loading = ref<boolean>(false)
const startAt = ref<number>(dayjs().add(-7, 'day').unix())
const endAt = ref<number>(dayjs().unix())
const periodType = ref<TopOrderPeriodType>(TopOrderPeriodType.DAY)

const fetchState = useAsyncData('home', async (): Promise<void> => {
  try {
    await Promise.all([
      fetchOrders(),
      topStore.fetchPendingOrders(),
      topStore.fetchPublicationStatus(),
    ])
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
})

watch([startAt, endAt, periodType], (): void => {
  fetchState.refresh()
})

const fetchOrders = async (): Promise<void> => {
  try {
    await topStore.fetchOrders(startAt.value, endAt.value, periodType.value)
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const isLoading = (): boolean => {
  return fetchState?.pending?.value || loading.value
}

const handleClickOrder = (orderId: string): void => {
  router.push(`/orders/${orderId}`)
}

const handleGoOrders = (): void => {
  router.push('/orders')
}

const handleGoProducts = (): void => {
  router.push('/products/new')
}

const handleGoExperiences = (): void => {
  router.push('/experiences/new')
}

try {
  await fetchState.execute()
}
catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-home-top
    v-model:start-at="startAt"
    v-model:end-at="endAt"
    v-model:period-type="periodType"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :orders="orders"
    :pending-orders="pendingOrders"
    :pending-orders-total="pendingOrdersTotal"
    :has-published-product="hasPublishedProduct"
    :has-published-experience="hasPublishedExperience"
    @click:order="handleClickOrder"
    @click:go-orders="handleGoOrders"
    @click:go-products="handleGoProducts"
    @click:go-experiences="handleGoExperiences"
  />
</template>
