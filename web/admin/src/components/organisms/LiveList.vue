<script lang="ts" setup>
import { mdiPlus } from '@mdi/js'
import useVuelidate from '@vuelidate/core'
import dayjs, { unix } from 'dayjs'
import { getResizedImages } from '~/lib/helpers'
import { getErrorMessage } from '~/lib/validations'
import {

  ScheduleStatus,

} from '~/types/api/v1'
import type { CreateLiveRequest, Live, Producer, Product, ProductMedia, Schedule, UpdateLiveRequest } from '~/types/api/v1'
import type { DateTimeInput } from '~/types/props'
import {
  CreateLiveValidationRules,
  TimeDataValidationRules,
} from '~/types/validations'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  createDialog: {
    type: Boolean,
    default: false,
  },
  createFormData: {
    type: Object as PropType<CreateLiveRequest>,
    default: (): CreateLiveRequest => ({
      producerId: '',
      productIds: [],
      comment: '',
      startAt: dayjs().unix(),
      endAt: dayjs().unix(),
    }),
  },
  schedule: {
    type: Object as PropType<Schedule>,
    default: (): Schedule => ({
      id: '',
      shopId: '',
      coordinatorId: '',
      title: '',
      description: '',
      status: ScheduleStatus.ScheduleStatusUnknown,
      thumbnailUrl: '',
      imageUrl: '',
      openingVideoUrl: '',
      _public: false,
      approved: false,
      startAt: dayjs().unix(),
      endAt: dayjs().unix(),
      createdAt: 0,
      updatedAt: 0,
    }),
  },
  live: {
    type: Object as PropType<Live>,
    default: (): Live => ({
      id: '',
      scheduleId: '',
      producerId: '',
      productIds: [],
      comment: '',
      startAt: dayjs().unix(),
      endAt: dayjs().unix(),
      createdAt: 0,
      updatedAt: 0,
    }),
  },
  lives: {
    type: Array<Live>,
    default: () => [],
  },
  producers: {
    type: Array<Producer>,
    default: () => [],
  },
  products: {
    type: Array<Product>,
    default: () => [],
  },
})

const emits = defineEmits<{
  (e: 'click:new'): void
  (e: 'update:live', live: Live): void
  (e: 'update:create-dialog', val: boolean): void
  (e: 'update:update-dialog', val: boolean): void
  (e: 'update:create-form-data', formData: CreateLiveRequest): void
  (e: 'update:update-form-data', formData: UpdateLiveRequest): void
  (e: 'search:producer', name: string): void
  (e: 'search:product', producerId: string, name: string): void
  (e: 'submit:create'): void
  (e: 'submit:update', liveId: string, formData: UpdateLiveRequest): void
  (e: 'submit:delete', liveId: string): void
}>()

const liveValue = computed({
  get: (): Live => props.live,
  set: (live: Live): void => emits('update:live', live),
})
const createDialogValue = computed({
  get: (): boolean => props.createDialog,
  set: (val: boolean): void => emits('update:create-dialog', val),
})
const createFormDataValue = computed({
  get: (): CreateLiveRequest => props.createFormData,
  set: (formData: CreateLiveRequest): void =>
    emits('update:create-form-data', formData),
})
const createStartTimeDataValue = computed({
  get: (): DateTimeInput => ({
    date: unix(props.createFormData?.startAt).format('YYYY-MM-DD'),
    time: unix(props.createFormData?.startAt).format('HH:mm'),
  }),
  set: (timeData: DateTimeInput): void => {
    const startAt = dayjs(`${timeData.date} ${timeData.time}`)
    createFormDataValue.value.startAt = startAt.unix()
  },
})
const createEndTimeDataValue = computed({
  get: (): DateTimeInput => ({
    date: unix(props.createFormData?.endAt).format('YYYY-MM-DD'),
    time: unix(props.createFormData?.endAt).format('HH:mm'),
  }),
  set: (timeData: DateTimeInput): void => {
    const endAt = dayjs(`${timeData.date} ${timeData.time}`)
    createFormDataValue.value.endAt = endAt.unix()
  },
})

const createFormDataValidate = useVuelidate(
  CreateLiveValidationRules,
  createFormDataValue,
)

const createStartTimeDataValidate = useVuelidate(
  TimeDataValidationRules,
  createStartTimeDataValue,
)
const createEndTimeDataValidate = useVuelidate(
  TimeDataValidationRules,
  createEndTimeDataValue,
)

const onChangeCreateStartAt = (): void => {
  const startAt = dayjs(
    `${createStartTimeDataValue.value.date} ${createStartTimeDataValue.value.time}`,
  )
  createFormDataValue.value.startAt = startAt.unix()
}

const onChangeCreateEndAt = (): void => {
  const endAt = dayjs(
    `${createEndTimeDataValue.value.date} ${createEndTimeDataValue.value.time}`,
  )
  createFormDataValue.value.endAt = endAt.unix()
}

const onChangeCreateProducerId = (): void => {
  onSearchProductFromCreate('')
  createFormDataValue.value.productIds = []
}

const getDay = (unixTime: number): string => {
  return unix(unixTime).format('YYYY/MM/DD HH:mm')
}

const getScheduleTerm = (schedule: Schedule): string => {
  return `${getDay(schedule.startAt)} ~ ${getDay(schedule.endAt)}`
}

const getProducer = (live: Live): Producer | undefined => {
  return props.producers.find((producer: Producer): boolean => {
    return producer.id === live?.producerId
  })
}

const getProductsByLive = (live: Live): Product[] => {
  const products: Product[] = []
  live.productIds.forEach((productId: string): void => {
    const product = props.products.find((product: Product): boolean => {
      return product.id === productId
    })
    if (product) {
      products.push(product)
    }
  })
  return products
}

const getProductsByProducerId = (producerId: string): Product[] => {
  return props.products.filter((product: Product): boolean => {
    return product.producerId === producerId
  })
}

const getProducerName = (live: Live): string => {
  const producer = getProducer(live)
  return producer ? producer.username : ''
}

const getProducerThumbnailUrl = (live: Live): string => {
  const producer = getProducer(live)
  return producer ? producer.thumbnailUrl : ''
}

const getProducerThumbnails = (live: Live): string => {
  const producer = getProducer(live)
  if (!producer?.thumbnailUrl) {
    return ''
  }
  return getResizedImages(producer.thumbnailUrl)
}

const getProductThumbnailUrl = (product: Product): string => {
  const thumbnail = product.media?.find((media: ProductMedia) => {
    return media.isThumbnail
  })
  return thumbnail ? thumbnail.url : ''
}

const onSearchProducer = (name: string): void => {
  emits('search:producer', name)
}

const onSearchProductFromCreate = (name: string): void => {
  emits('search:product', props.createFormData.producerId, name)
}

const onSearchProductFromUpdate = (name: string): void => {
  emits('search:product', props.live.producerId, name)
}

const onClickNew = (): void => {
  emits('click:new')
}

const onClickCloseCreateDialog = (): void => {
  createDialogValue.value = false
}

const onSubmitCreate = async (): Promise<void> => {
  const formDataValid = await createFormDataValidate.value.$validate()
  const startTimeDataValid
    = await createStartTimeDataValidate.value.$validate()
  const endTimeDataValid = await createEndTimeDataValidate.value.$validate()
  if (!formDataValid || !startTimeDataValid || !endTimeDataValid) {
    return
  }

  emits('submit:create')
}

const onSubmitUpdate = async (
  liveId: string,
  formData: UpdateLiveRequest,
): Promise<void> => {
  emits('submit:update', liveId, formData)
}

const onSubmitDelete = (liveId: string): void => {
  emits('submit:delete', liveId)
}
</script>

<template>
  <v-dialog
    v-model="createDialogValue"
    width="500"
  >
    <v-card>
      <v-card-title class="text-h6 primaryLight">
        スケジュール登録
      </v-card-title>
      <v-card-text>
        <p class="text-subtitle-2 text-grey pb-2">
          ライブ配開始日時
        </p>
        <div class="d-flex flex-column flex-md-row justify-center">
          <v-text-field
            v-model="createStartTimeDataValidate.date.$model"
            :error-messages="
              getErrorMessage(createStartTimeDataValidate.date.$errors)
            "
            type="date"
            variant="outlined"
            density="compact"
            class="mr-md-2"
            @update:model-value="onChangeCreateStartAt"
          />
          <v-text-field
            v-model="createStartTimeDataValidate.time.$model"
            :error-messages="
              getErrorMessage(createStartTimeDataValidate.time.$errors)
            "
            type="time"
            variant="outlined"
            density="compact"
            @update:model-value="onChangeCreateStartAt"
          />
        </div>
        <p class="text-subtitle-2 text-grey pb-2">
          ライブ配終了日時
        </p>
        <div class="d-flex flex-column flex-md-row justify-center">
          <v-text-field
            v-model="createEndTimeDataValidate.date.$model"
            :error-messages="
              getErrorMessage(createEndTimeDataValidate.date.$errors)
            "
            type="date"
            variant="outlined"
            density="compact"
            class="mr-md-2"
            @update:model-value="onChangeCreateEndAt"
          />
          <v-text-field
            v-model="createEndTimeDataValidate.time.$model"
            :error-messages="
              getErrorMessage(createEndTimeDataValidate.time.$errors)
            "
            type="time"
            variant="outlined"
            density="compact"
            @update:model-value="onChangeCreateEndAt"
          />
        </div>
        <v-autocomplete
          v-model="createFormDataValidate.producerId.$model"
          :error-messages="
            getErrorMessage(createFormDataValidate.producerId.$errors)
          "
          label="生産者"
          :items="producers"
          item-title="username"
          item-value="id"
          clearable
          @update:search="onSearchProducer"
          @update:model-value="onChangeCreateProducerId"
        />

        <v-autocomplete
          v-model="createFormDataValidate.productIds.$model"
          :error-messages="
            getErrorMessage(createFormDataValidate.productIds.$errors)
          "
          label="関連する商品"
          :items="getProductsByProducerId(createFormDataValue.producerId)"
          item-title="name"
          item-value="id"
          chips
          closable-chips
          clearable
          multiple
          density="comfortable"
          @update:search="onSearchProductFromCreate"
        >
          <template #chip="{ props: val, item }">
            <v-chip
              v-bind="val"
              :prepend-avatar="getProductThumbnailUrl(item.raw)"
              :text="item.raw.name"
              rounded
              class="px-4"
              variant="outlined"
            />
          </template>
          <template #item="{ props: val, item }">
            <v-list-item
              v-bind="val"
              :prepend-avatar="getProductThumbnailUrl(item.raw)"
              :title="item.raw.name"
            />
          </template>
        </v-autocomplete>

        <v-textarea
          v-model="createFormDataValidate.comment.$model"
          :error-messages="
            getErrorMessage(createFormDataValidate.comment.$errors)
          "
          label="概要"
          maxlength="2000"
        />
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn
          color=""
          variant="text"
          @click="onClickCloseCreateDialog"
        >
          キャンセル
        </v-btn>
        <v-btn
          :loading="loading"
          color="primary"
          variant="text"
          @click="onSubmitCreate"
        >
          登録
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-row>
    <v-col sm="12">
      <p class="text-subtitle-2 pb-2">
        マルシェ開催期間
      </p>
      <p class="text-subtitle-2">
        {{ getScheduleTerm(schedule) }}
      </p>
    </v-col>
    <v-col sm="12">
      <!-- 新コンポーネント -->
      <div class="d-flex flex-column ga-2">
        <organisms-live-list-item
          v-for="(item, i) in props.lives"
          :key="`live-${i}`"
          :item="item"
          :producer-thumbnail-url="getProducerThumbnailUrl(item)"
          :producer-thumbnails-srcset="getProducerThumbnails(item)"
          :producer-name="getProducerName(item)"
          :products="getProductsByProducerId(item.producerId)"
          :live-products="getProductsByLive(item)"
          :producers="producers"
          :loading="loading"
          @submit:delete="onSubmitDelete"
          @submit:update="onSubmitUpdate"
        />
      </div>
    </v-col>
    <v-col sm="12">
      <v-btn
        block
        variant="outlined"
        color="primary"
        @click="onClickNew"
      >
        <v-icon :icon="mdiPlus" />
        生産者と商品を追加
      </v-btn>
    </v-col>
  </v-row>
</template>
