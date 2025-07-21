<script lang="ts" setup>
import dayjs from 'dayjs'
import { storeToRefs } from 'pinia'
import { useAlert } from '~/lib/hooks'

import { useAuthStore, useCommonStore, usePromotionStore, useShopStore } from '~/store'
import { AdminType } from '~/types/api'
import type { UpdatePromotionRequest } from '~/types/api'

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
  public: false,
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

const isEditable = (): boolean => {
  switch (adminType.value) {
    case AdminType.ADMINISTRATOR:
      return true
    case AdminType.COORDINATOR:
      return shopIds.value.includes(promotion.value.shopId)
    default:
      return false
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
  <div>
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
    <div
      class="position-fixed bottom-0 left-0 w-100 bg-white pa-4 text-right elevation-3"
    >
      <div class="d-inline-flex ga-4">
        <v-btn
          color="secondary"
          variant="outlined"
          @click="$router.back()"
        >
          戻る
        </v-btn>
        <v-btn
          v-show="isEditable()"
          color="primary"
          variant="outlined"
          type="submit"
          form="update-promotion-form"
        >
          更新
        </v-btn>
      </div>
    </div>
  </div>
</template>
