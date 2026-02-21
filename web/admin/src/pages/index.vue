<script lang="ts" setup>
import dayjs from 'dayjs'
import { storeToRefs } from 'pinia'
import { useAlert } from '~/lib/hooks'

import { useTopStore } from '~/store'
import { TopOrderPeriodType } from '~/types'

const topStore = useTopStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const { orders } = storeToRefs(topStore)

const loading = ref<boolean>(false)
const startAt = ref<number>(dayjs().add(-7, 'day').unix())
const endAt = ref<number>(dayjs().unix())
const periodType = ref<TopOrderPeriodType>(TopOrderPeriodType.DAY)

const fetchState = useAsyncData('home', async (): Promise<void> => {
  try {
    await fetchOrders()
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
  />
</template>
