<script lang="ts" setup>
import {
  mdiPlus,
  mdiAccountGroup,
  mdiClock,
  mdiPackageVariant,
  mdiClose,
  mdiCheck,
} from '@mdi/js'
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
  <!-- ライブスケジュール作成ダイアログ -->
  <v-dialog
    v-model="createDialogValue"
    max-width="600"
    scrollable
  >
    <v-card class="create-dialog-card">
      <v-card-title class="d-flex align-center section-header pa-6">
        <v-icon
          :icon="mdiAccountGroup"
          size="24"
          class="mr-3 text-primary"
        />
        <span class="text-h6 font-weight-medium">ライブスケジュール登録</span>
        <v-spacer />
        <v-btn
          :icon="mdiClose"
          variant="text"
          size="small"
          aria-label="閉じる"
          @click="onClickCloseCreateDialog"
        />
      </v-card-title>

      <v-card-text class="pa-6">
        <!-- 開催期間設定 -->
        <div class="mb-6">
          <div class="d-flex align-center mb-4">
            <v-icon
              :icon="mdiClock"
              size="20"
              class="mr-2 text-primary"
            />
            <span class="text-subtitle-1 font-weight-medium">開催期間設定</span>
          </div>

          <div class="mb-4">
            <p class="text-subtitle-2 mb-3 text-grey-darken-1">
              開始日時 *
            </p>
            <v-row>
              <v-col
                cols="12"
                sm="6"
              >
                <v-text-field
                  v-model="createStartTimeDataValidate.date.$model"
                  :error-messages="
                    getErrorMessage(createStartTimeDataValidate.date.$errors)
                  "
                  label="日付"
                  type="date"
                  variant="outlined"
                  density="comfortable"
                  @update:model-value="onChangeCreateStartAt"
                />
              </v-col>
              <v-col
                cols="12"
                sm="6"
              >
                <v-text-field
                  v-model="createStartTimeDataValidate.time.$model"
                  :error-messages="
                    getErrorMessage(createStartTimeDataValidate.time.$errors)
                  "
                  label="時刻"
                  type="time"
                  variant="outlined"
                  density="comfortable"
                  @update:model-value="onChangeCreateStartAt"
                />
              </v-col>
            </v-row>
          </div>

          <div class="mb-4">
            <p class="text-subtitle-2 mb-3 text-grey-darken-1">
              終了日時 *
            </p>
            <v-row>
              <v-col
                cols="12"
                sm="6"
              >
                <v-text-field
                  v-model="createEndTimeDataValidate.date.$model"
                  :error-messages="
                    getErrorMessage(createEndTimeDataValidate.date.$errors)
                  "
                  label="日付"
                  type="date"
                  variant="outlined"
                  density="comfortable"
                  @update:model-value="onChangeCreateEndAt"
                />
              </v-col>
              <v-col
                cols="12"
                sm="6"
              >
                <v-text-field
                  v-model="createEndTimeDataValidate.time.$model"
                  :error-messages="
                    getErrorMessage(createEndTimeDataValidate.time.$errors)
                  "
                  label="時刻"
                  type="time"
                  variant="outlined"
                  density="comfortable"
                  @update:model-value="onChangeCreateEndAt"
                />
              </v-col>
            </v-row>
          </div>
        </div>

        <!-- 生産者・商品設定 -->
        <div class="mb-6">
          <div class="d-flex align-center mb-4">
            <v-icon
              :icon="mdiPackageVariant"
              size="20"
              class="mr-2 text-primary"
            />
            <span class="text-subtitle-1 font-weight-medium">生産者・商品設定</span>
          </div>

          <v-autocomplete
            v-model="createFormDataValidate.producerId.$model"
            :error-messages="
              getErrorMessage(createFormDataValidate.producerId.$errors)
            "
            label="生産者 *"
            :items="producers"
            item-title="username"
            item-value="id"
            variant="outlined"
            density="comfortable"
            clearable
            class="mb-4"
            @update:search="onSearchProducer"
            @update:model-value="onChangeCreateProducerId"
          />

          <v-autocomplete
            v-model="createFormDataValidate.productIds.$model"
            :error-messages="
              getErrorMessage(createFormDataValidate.productIds.$errors)
            "
            label="関連する商品 *"
            :items="getProductsByProducerId(createFormDataValue.producerId)"
            item-title="name"
            item-value="id"
            variant="outlined"
            density="comfortable"
            chips
            closable-chips
            clearable
            multiple
            class="mb-4"
            @update:search="onSearchProductFromCreate"
          >
            <template #chip="{ props: val, item }">
              <v-chip
                v-bind="val"
                :prepend-avatar="getProductThumbnailUrl(item.raw)"
                :text="item.raw.name"
                rounded
                class="px-3"
                variant="outlined"
                color="primary"
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
            label="概要・コメント"
            maxlength="2000"
            variant="outlined"
            density="comfortable"
            rows="3"
            counter
          />
        </div>
      </v-card-text>

      <v-card-actions class="pa-6 pt-0">
        <v-spacer />
        <v-btn
          variant="text"
          size="large"
          @click="onClickCloseCreateDialog"
        >
          <v-icon
            :icon="mdiClose"
            start
          />
          キャンセル
        </v-btn>
        <v-btn
          :loading="loading"
          color="primary"
          variant="elevated"
          size="large"
          @click="onSubmitCreate"
        >
          <v-icon
            :icon="mdiCheck"
            start
          />
          登録
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <!-- メインコンテンツ -->
  <div class="live-list-container">
    <!-- スケジュール情報セクション -->
    <v-card
      class="form-section-card mb-6"
      elevation="2"
    >
      <v-card-title class="d-flex align-center section-header">
        <v-icon
          :icon="mdiClock"
          size="24"
          class="mr-3 text-primary"
        />
        <span class="text-h6 font-weight-medium">マルシェ開催期間</span>
      </v-card-title>
      <v-card-text class="pa-6">
        <div class="schedule-period">
          <v-chip
            color="primary"
            variant="outlined"
            size="large"
            class="pa-4"
          >
            <v-icon
              :icon="mdiClock"
              start
              size="20"
            />
            {{ getScheduleTerm(schedule) }}
          </v-chip>
        </div>
      </v-card-text>
    </v-card>

    <!-- ライブスケジュール一覧セクション -->
    <v-card
      class="form-section-card mb-6"
      elevation="2"
    >
      <v-card-title class="d-flex align-center section-header">
        <v-icon
          :icon="mdiAccountGroup"
          size="24"
          class="mr-3 text-primary"
        />
        <span class="text-h6 font-weight-medium">ライブスケジュール一覧</span>
        <v-spacer />
        <v-btn
          variant="elevated"
          color="primary"
          size="small"
          @click="onClickNew"
        >
          <v-icon
            :icon="mdiPlus"
            start
          />
          追加
        </v-btn>
      </v-card-title>
      <v-card-text class="pa-6">
        <div
          v-if="props.lives.length > 0"
          class="d-flex flex-column ga-4"
        >
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
        <div
          v-else
          class="text-center py-8"
        >
          <v-icon
            :icon="mdiAccountGroup"
            size="64"
            class="text-grey-lighten-1 mb-4"
          />
          <p class="text-body-1 text-grey-darken-1 mb-4">
            ライブスケジュールが登録されていません
          </p>
          <v-btn
            variant="outlined"
            color="primary"
            size="large"
            @click="onClickNew"
          >
            <v-icon
              :icon="mdiPlus"
              start
            />
            最初のライブスケジュールを追加
          </v-btn>
        </div>
      </v-card-text>
    </v-card>
  </div>
</template>

<style scoped>
.form-section-card {
  border-radius: 12px;
  max-width: none;
}

.section-header {
  background: linear-gradient(90deg, rgb(33 150 243 / 5%) 0%, rgb(33 150 243 / 0%) 100%);
  border-bottom: 1px solid rgb(0 0 0 / 5%);
  padding: 20px 24px;
}

.create-dialog-card {
  border-radius: 12px;
}

.live-list-container {
  min-height: 200px;
}

.schedule-period {
  display: flex;
  align-items: center;
  justify-content: center;
}

@media (width <= 600px) {
  .form-section-card {
    border-radius: 8px;
  }

  .section-header {
    padding: 16px 20px;
  }

  .create-dialog-card {
    border-radius: 8px;
    margin: 16px;
  }
}
</style>
