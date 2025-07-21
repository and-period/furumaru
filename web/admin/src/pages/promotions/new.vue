<script lang="ts" setup>
import dayjs from 'dayjs'
import { useAlert } from '~/lib/hooks'

import { useCommonStore, usePromotionStore } from '~/store'
import type { CreatePromotionRequest } from '~/types/api'

const router = useRouter()
const commonStore = useCommonStore()
const promotionStore = usePromotionStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const loading = ref<boolean>(false)
const formData = ref<CreatePromotionRequest>({
  title: '',
  description: '',
  public: false,
  discountType: 1,
  discountRate: 0,
  code: '',
  startAt: dayjs().unix(),
  endAt: dayjs().unix(),
})

const handleSubmit = async (): Promise<void> => {
  try {
    loading.value = true
    const req: CreatePromotionRequest = {
      ...formData.value,
      discountRate: Number(formData.value.discountRate),
    }
    await promotionStore.createPromotion(req)
    commonStore.addSnackbar({
      message: `${formData.value.title}を作成しました。`,
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
</script>

<template>
  <div>
    <templates-promotion-new
      v-model:form-data="formData"
      :loading="loading"
      :is-alert="isShow"
      :alert-type="alertType"
      :alert-text="alertText"
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
          color="primary"
          variant="outlined"
          type="submit"
          form="create-promotion-form"
        >
          登録
        </v-btn>
      </div>
    </div>
  </div>
</template>
