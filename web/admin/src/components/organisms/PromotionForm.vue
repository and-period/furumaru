<script lang="ts" setup>
import dayjs from 'dayjs'

import { CreatePromotionRequest, DiscountType } from '~/types/api'
import { PromotionTime } from '~/types/props'

const props = defineProps({
  formType: {
    type: String,
    default: 'create',
    validator: (value: string) => {
      return ['create', 'edit'].includes(value)
    }
  },
  formData: {
    type: Object,
    default: (): CreatePromotionRequest => ({
      title: '',
      description: '',
      public: false,
      publishedAt: dayjs().unix(),
      discountType: DiscountType.AMOUNT,
      discountRate: 0,
      code: '',
      startAt: dayjs().unix(),
      endAt: dayjs().unix()
    })
  },
  timeData: {
    type: Object,
    default: (): PromotionTime => ({
      publishedDate: '',
      publishedTime: '',
      startDate: '',
      startTime: '',
      endDate: '',
      endTime: ''
    })
  }
})

const emit = defineEmits<{
  (e: 'update:formData', formData: CreatePromotionRequest): void
  (e: 'update:timeData', timeData: PromotionTime): void
  (e: 'submit'): void
}>()

const formDataValue = computed({
  get: (): CreatePromotionRequest => props.formData as CreatePromotionRequest,
  set: (val: CreatePromotionRequest) => emit('update:formData', val)
})

const timeDataValue = computed({
  get: (): PromotionTime => props.timeData as PromotionTime,
  set: (val: PromotionTime) => emit('update:timeData', val)
})

const btnText = computed(() => {
  return props.formType === 'create' ? '登録' : '更新'
})

const handleSubmit = () => {
  formDataValue.value.publishedAt = dayjs(
    timeDataValue.value.publishedDate + ' ' + timeDataValue.value.publishedTime
  ).unix()
  formDataValue.value.startAt = dayjs(
    timeDataValue.value.startDate + ' ' + timeDataValue.value.startTime
  ).unix()
  formDataValue.value.endAt = dayjs(
    timeDataValue.value.endDate + ' ' + timeDataValue.value.endTime
  ).unix()

  emit('submit')
}

const handleGenerate = () => {
  const code = generateRandomString()
  props.formData.code = code
}

const generateRandomString = (): string => {
  const characters =
    'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  let result = ''
  const charactersLength = characters.length
  for (let i = 0; i < 8; i++) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength))
  }

  return result
}

const getErrorMessage = () => {
  switch (formDataValue.value.discountType) {
    case 1:
      if (formDataValue.value.discountRate >= 0) {
        return ''
      } else {
        return '0以上の値を指定してください'
      }
    case 2:
      if (
        formDataValue.value.discountRate >= 0 &&
        formDataValue.value.discountRate <= 100
      ) {
        return ''
      } else {
        return '0~100の値を指定してください'
      }
    default:
      return ''
  }
}

const statusList = [
  { status: '有効', value: true },
  { status: '無効', value: false }
]

const discountMethodList = [
  { method: '円', value: DiscountType.AMOUNT },
  { method: '%', value: DiscountType.RATE },
  { method: '送料無料', value: DiscountType.FREE_SHIPPING }
]
</script>

<template>
  <form @submit.prevent="handleSubmit">
    <v-card elevation="0">
      <v-card-text>
        <v-text-field
          v-model="formDataValue.title"
          label="タイトル"
          required
          maxlength="200"
        />
        <v-textarea
          v-model="formDataValue.description"
          label="説明"
          maxlength="2000"
        />
        <div class="d-flex align-center">
          <v-text-field
            v-model="formDataValue.code"
            class="mr-4"
            label="割引コード(8文字)"
            required
            maxlength="8"
            :error-messages="
              (props.formData.code.length === 0 || props.formData.code.length === 8) ? '' : '割引コードは8文字です'
            "
          />
          <v-btn variant="outlined" size="small" color="primary" @click="handleGenerate">
            自動生成
          </v-btn>
          <v-spacer />
        </div>
        <v-select
          v-model="formDataValue.public"
          :items="statusList"
          label="ステータス"
          item-title="status"
          item-value="value"
        />
        <div class="d-flex align-center">
          <v-select
            v-model="formDataValue.discountType"
            :items="discountMethodList"
            item-title="method"
            item-value="value"
            label="割引方法"
          />
          <v-text-field
            v-if="props.formData.discountType != 3"
            v-model="formDataValue.discountRate"
            class="ml-4"
            type="number"
            label="割引値"
            :error-messages="getErrorMessage()"
          />
        </div>

        <p class="text-h6">
          投稿開始
        </p>
        <div class="d-flex align-center justify-center">
          <v-text-field
            v-model="timeDataValue.publishedDate"
            type="date"
            variant="outlined"
            required
            class="mr-2"
          />
          <v-text-field
            v-model="timeDataValue.publishedTime"
            type="time"
            required
            variant="outlined"
          />
          <p class="text-h6 mb-6 ml-4">
            〜
          </p>
          <v-spacer />
        </div>

        <p class="text-h6">
          使用期間
        </p>
        <div class="d-flex align-center">
          <v-text-field
            v-model="timeDataValue.startDate"
            type="date"
            variant="outlined"
            required
            class="mr-2"
          />
          <v-text-field
            v-model="timeDataValue.startTime"
            type="time"
            required
            variant="outlined"
          />
          <p class="text-h6 mx-4 mb-6">
            〜
          </p>
          <v-text-field
            v-model="timeDataValue.endDate"
            type="date"
            variant="outlined"
            required
            class="mr-2"
          />
          <v-text-field
            v-model="timeDataValue.endTime"
            type="time"
            required
            variant="outlined"
          />
        </div>
      </v-card-text>
      <v-card-actions>
        <v-btn block variant="outlined" color="primary" type="submit" class="mt-4">
          {{ btnText }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </form>
</template>
