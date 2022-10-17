<template>
  <the-promotion-edit-form-page
    :form-data="formData"
    :form-data-loading="fetchState.pending"
    @submit="handleSubmit"
  />
</template>

<script lang="ts">
import { useFetch, useRoute, useRouter } from '@nuxtjs/composition-api'
import { defineComponent, reactive } from '@vue/composition-api'
import dayjs from 'dayjs'

import { usePromotionStore } from '~/store/promotion'
import { PromotionResponse } from '~/types/api'

export default defineComponent({
  setup() {
    const router = useRouter()
    const route = useRoute()
    const id = route.value.params.id
    const promotionStore = usePromotionStore()

    const { getPromotion } = usePromotionStore()

    const formData = reactive<PromotionResponse>({
          id,
          title: '',
          description: '',
          public: false,
          publishedAt: dayjs().unix(),
          discountType: 1,
          discountRate: 0,
          code: '',
          startAt: dayjs().unix(),
          endAt: dayjs().unix(),
          createdAt: -1,
          updatedAt: -1,
    })

    const { fetchState } = useFetch(async () => {
      const promotion = await getPromotion(id)
      formData.title = promotion.title
      formData.description = promotion.description
      formData.public = promotion.public
      formData.publishedAt = promotion.publishedAt
      formData.discountType = promotion.discountType
      formData.discountRate = promotion.discountRate
      formData.code = promotion.code
      formData.startAt = promotion.startAt
      formData.endAt = promotion.endAt
    })

    const handleSubmit = async () => {
      try {
        await promotionStore.editPromotion(id, formData)
        router.push('/promotions')
      } catch (error) {
        console.log(error)
      }
    }

    return {
      fetchState,
      formData,
      handleSubmit
    }
  },
})
</script>
