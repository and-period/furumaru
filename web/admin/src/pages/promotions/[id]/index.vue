<script lang="ts" setup>
import dayjs from 'dayjs'
import { storeToRefs } from 'pinia'
import { useAlert } from '~/lib/hooks'

import { useAuthStore, useCommonStore, usePromotionStore, useShopStore } from '~/store'
import type { UpdatePromotionRequest } from '~/types/api/v1'

const router = useRouter()
const route = useRoute()
const commonStore = useCommonStore()
const authStore = useAuthStore()
const promotionStore = usePromotionStore()
const shopStore = useShopStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const promotionId = route.params.id as string

const { shopIds, adminType } = storeToRefs(authStore)
const { promotion } = storeToRefs(promotionStore)
const { shop } = storeToRefs(shopStore)

const loading = ref<boolean>(false)
const formData = ref<UpdatePromotionRequest>({
  title: '',
  description: '',
  _public: false,
  discountType: 1,
  discountRate: 0,
  code: '',
  startAt: dayjs().unix(),
  endAt: dayjs().unix(),
})

const fetchState = useAsyncData(async (): Promise<void> => {
  try {
    await promotionStore.getPromotion(promotionId)
    formData.value = { ...promotion.value }
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
})

const isLoading = (): boolean => {
  return fetchState?.pending?.value || loading.value
}

const handleSubmit = async (): Promise<void> => {
  try {
    loading.value = true
    const req: UpdatePromotionRequest = {
      ...formData.value,
      discountRate: Number(formData.value.discountRate),
    }
    await promotionStore.updatePromotion(promotionId, req)
    commonStore.addSnackbar({
      message: 'セール情報の編集が完了しました',
      color: 'info',
    })
    router.push('/promotions')
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
  }
}

try {
  await fetchState.execute()
}
catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-promotion-edit
    v-model:form-data="formData"
    :loading="isLoading()"
    :shop-ids="shopIds"
    :admin-type="adminType"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :promotion="promotion"
    :shop="shop"
    @submit="handleSubmit"
  />
</template>
