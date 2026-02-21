<script lang="ts" setup>
import type { Promotion } from '~/types/api/v1'
import { useNotificationDisplay } from '~/composables/useNotificationDisplay'

defineProps({
  promotion: {
    type: Object as PropType<Promotion>,
    default: undefined,
  },
  showTitle: {
    type: Boolean,
    default: false,
  },
})

const { getPromotionTerm, getPromotionDiscount } = useNotificationDisplay()
</script>

<template>
  <v-card
    v-if="promotion"
    variant="outlined"
    class="pa-4"
  >
    <v-card-title
      v-if="showTitle"
      class="text-h6 pa-0 mb-3"
    >
      {{ promotion.title }}
    </v-card-title>

    <div class="d-flex flex-column ga-2">
      <div class="d-flex justify-space-between align-center">
        <span class="text-body-2 text-grey-darken-1">割引コード</span>
        <v-chip
          size="small"
          :color="promotion.code ? 'primary' : 'grey'"
          variant="outlined"
        >
          {{ promotion.code || '未設定' }}
        </v-chip>
      </div>

      <div class="d-flex justify-space-between align-center">
        <span class="text-body-2 text-grey-darken-1">割引額</span>
        <v-chip
          size="small"
          color="success"
          variant="tonal"
        >
          {{ getPromotionDiscount(promotion) }}
        </v-chip>
      </div>

      <div class="d-flex justify-space-between align-center">
        <span class="text-body-2 text-grey-darken-1">使用期間</span>
        <span class="text-body-2">{{ getPromotionTerm(promotion) }}</span>
      </div>
    </div>
  </v-card>

  <v-card
    v-else
    variant="outlined"
    class="pa-4 text-center"
  >
    <v-icon
      icon="mdi-information-outline"
      class="mb-2 text-grey"
    />
    <div class="text-body-2 text-grey-darken-1">
      プロモーション情報を選択してください
    </div>
  </v-card>
</template>
