<template>
  <the-promotion-create-form-page
    :form-data="formData"
    @submit="handleSubmit"
  />
</template>

<script lang="ts">
import { useRouter } from '@nuxtjs/composition-api'
import { defineComponent, reactive } from '@vue/composition-api'
import dayjs from 'dayjs'

import { usePromotionStore } from '~/store/promotion'
import { CreatePromotionRequest } from '~/types/api'

export default defineComponent({
  setup() {
    const router = useRouter()
    const promotionStore = usePromotionStore()

    const formData = reactive<CreatePromotionRequest>({
      title: '',
      description: '',
      public: false,
      publishedAt: dayjs().unix(),
      discountType: 1,
      discountRate: 0,
      code: '',
      startAt: dayjs().unix(),
      endAt: dayjs().unix(),
    })

    const handleSubmit = async () => {
      try {
        await promotionStore.createPromotion(formData)
        router.push('/promotions')
      } catch (error) {
        console.log(error)
      }
    }

    return {
      formData,
      handleSubmit,
    }
  },
})
</script>
