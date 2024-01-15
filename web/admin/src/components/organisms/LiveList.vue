<script lang="ts" setup>
import { mdiPencil, mdiPlus } from '@mdi/js'
import useVuelidate from '@vuelidate/core'
import dayjs, { unix } from 'dayjs'
import { getResizedImages } from '~/lib/helpers'
import { getErrorMessage } from '~/lib/validations'
import {
  type CreateLiveRequest,
  type Live,
  type Producer,
  type Product,
  type ProductMediaInner,
  type Schedule,
  ScheduleStatus,
  type UpdateLiveRequest
} from '~/types/api'
import type { LiveTime } from '~/types/props'
import {
  CreateLiveValidationRules,
  TimeDataValidationRules,
  UpdateLiveValidationRules
} from '~/types/validations'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  createDialog: {
    type: Boolean,
    default: false
  },
  updateDialog: {
    type: Boolean,
    default: false
  },
  createFormData: {
    type: Object as PropType<CreateLiveRequest>,
    default: (): CreateLiveRequest => ({
      producerId: '',
      productIds: [],
      comment: '',
      startAt: dayjs().unix(),
      endAt: dayjs().unix()
    })
  },
  updateFormData: {
    type: Object as PropType<UpdateLiveRequest>,
    default: (): UpdateLiveRequest => ({
      productIds: [],
      comment: '',
      startAt: dayjs().unix(),
      endAt: dayjs().unix()
    })
  },
  schedule: {
    type: Object as PropType<Schedule>,
    default: (): Schedule => ({
      id: '',
      coordinatorId: '',
      title: '',
      description: '',
      status: ScheduleStatus.UNKNOWN,
      thumbnailUrl: '',
      thumbnails: [],
      imageUrl: '',
      openingVideoUrl: '',
      public: false,
      approved: false,
      startAt: dayjs().unix(),
      endAt: dayjs().unix(),
      createdAt: 0,
      updatedAt: 0
    })
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
      updatedAt: 0
    })
  },
  lives: {
    type: Array<Live>,
    default: () => []
  },
  producers: {
    type: Array<Producer>,
    default: () => []
  },
  products: {
    type: Array<Product>,
    default: () => []
  }
})

const emits = defineEmits<{
  (e: 'click:new'): void;
  (e: 'click:edit', liveId: string): void;
  (e: 'update:live', live: Live): void;
  (e: 'update:create-dialog', val: boolean): void;
  (e: 'update:update-dialog', val: boolean): void;
  (e: 'update:create-form-data', formData: CreateLiveRequest): void;
  (e: 'update:update-form-data', formData: UpdateLiveRequest): void;
  (e: 'search:producer', name: string): void;
  (e: 'search:product', producerId: string, name: string): void;
  (e: 'submit:create'): void;
  (e: 'submit:update'): void;
  (e: 'submit:delete'): void;
}>()

const liveValue = computed({
  get: (): Live => props.live,
  set: (live: Live): void => emits('update:live', live)
})
const createDialogValue = computed({
  get: (): boolean => props.createDialog,
  set: (val: boolean): void => emits('update:create-dialog', val)
})
const updateDialogValue = computed({
  get: (): boolean => props.updateDialog,
  set: (val: boolean): void => emits('update:update-dialog', val)
})
const createFormDataValue = computed({
  get: (): CreateLiveRequest => props.createFormData,
  set: (formData: CreateLiveRequest): void =>
    emits('update:create-form-data', formData)
})
const updateFormDataValue = computed({
  get: (): UpdateLiveRequest => props.updateFormData,
  set: (formData: UpdateLiveRequest): void =>
    emits('update:update-form-data', formData)
})
const createTimeDataValue = computed({
  get: (): LiveTime => ({
    startDate: unix(props.createFormData?.startAt).format('YYYY-MM-DD'),
    startTime: unix(props.createFormData?.startAt).format('HH:mm'),
    endDate: unix(props.createFormData.endAt).format('YYYY-MM-DD'),
    endTime: unix(props.createFormData.endAt).format('HH:mm')
  }),
  set: (timeData: LiveTime): void => {
    const startAt = dayjs(`${timeData.startDate} ${timeData.startTime}`)
    const endAt = dayjs(`${timeData.endDate} ${timeData.endTime}`)
    createFormDataValue.value.startAt = startAt.unix()
    createFormDataValue.value.endAt = endAt.unix()
  }
})
const updateTimeDataValue = computed({
  get: (): LiveTime => ({
    startDate: unix(props.updateFormData?.startAt).format('YYYY-MM-DD'),
    startTime: unix(props.updateFormData?.startAt).format('HH:mm'),
    endDate: unix(props.updateFormData.endAt).format('YYYY-MM-DD'),
    endTime: unix(props.updateFormData.endAt).format('HH:mm')
  }),
  set: (timeData: LiveTime): void => {
    const startAt = dayjs(`${timeData.startDate} ${timeData.startTime}`)
    const endAt = dayjs(`${timeData.endDate} ${timeData.endTime}`)
    updateFormDataValue.value.startAt = startAt.unix()
    updateFormDataValue.value.endAt = endAt.unix()
  }
})

const createFormDataValidate = useVuelidate(
  CreateLiveValidationRules,
  createFormDataValue
)
const updateFormDataValidate = useVuelidate(
  UpdateLiveValidationRules,
  updateFormDataValue
)
const createTimeDataValidate = useVuelidate(
  TimeDataValidationRules,
  createTimeDataValue
)
const updateTimeDataValidate = useVuelidate(
  TimeDataValidationRules,
  updateTimeDataValue
)

const onChangeCreateStartAt = (): void => {
  const startAt = dayjs(
    `${createTimeDataValue.value.startDate} ${createTimeDataValue.value.startTime}`
  )
  createFormDataValue.value.startAt = startAt.unix()
}

const onChangeCreateEndAt = (): void => {
  const endAt = dayjs(
    `${createTimeDataValue.value.endDate} ${createTimeDataValue.value.endTime}`
  )
  createFormDataValue.value.endAt = endAt.unix()
}

const onChangeCreateProducerId = (): void => {
  onSearchProductFromCreate('')
  createFormDataValue.value.productIds = []
}

const onChangeUpdateStartAt = (): void => {
  const startAt = dayjs(
    `${updateTimeDataValue.value.startDate} ${updateTimeDataValue.value.startTime}`
  )
  updateFormDataValue.value.startAt = startAt.unix()
}

const onChangeUpdateEndAt = (): void => {
  const endAt = dayjs(
    `${updateTimeDataValue.value.endDate} ${updateTimeDataValue.value.endTime}`
  )
  updateFormDataValue.value.endAt = endAt.unix()
}

const getDay = (unixTime: number): string => {
  return unix(unixTime).format('YYYY/MM/DD HH:mm')
}

const getScheduleTerm = (schedule: Schedule): string => {
  return `${getDay(schedule.startAt)} ~ ${getDay(schedule.endAt)}`
}

const getLiveTerm = (live: Live): string => {
  return `${getDay(live.startAt)} ~ ${getDay(live.endAt)}`
}

const getProducer = (live: Live): Producer | undefined => {
  return props.producers.find((producer: Producer): boolean => {
    return producer.id === live?.producerId
  })
}

const getProductsByLive = (live: Live): Product[] => {
  const products: Product[] = []
  props.products.forEach((product: Product): void => {
    if (!live.productIds.includes(product.id)) {
      return
    }
    products.push(product)
  })
  return products
}

const getProductsByProducerId = (producerId: string): Product[] => {
  const products: Product[] = []
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
  if (!producer?.thumbnails) {
    return ''
  }
  return getResizedImages(producer.thumbnails)
}

const getProductInventoryColor = (product: Product): string => {
  return product.inventory > 0 ? '' : 'text-error'
}

const getProductThumbnailUrl = (product: Product): string => {
  const thumbnail = product.media?.find((media: ProductMediaInner) => {
    return media.isThumbnail
  })
  return thumbnail ? thumbnail.url : ''
}

const getProductThumbnails = (product: Product): string => {
  const thumbnail = product.media?.find((media: ProductMediaInner) => {
    return media.isThumbnail
  })
  return thumbnail ? getResizedImages(thumbnail.images) : ''
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
  const timeDataValid = await createTimeDataValidate.value.$validate()
  if (!formDataValid || !timeDataValid) {
    return
  }

  emits('submit:create')
}

const onClickEdit = (liveId: string): void => {
  emits('click:edit', liveId)
}

const onClickCloseUpdateDialog = (): void => {
  updateDialogValue.value = false
}

const onSubmitUpdate = async (): Promise<void> => {
  const formDataValid = await updateFormDataValidate.value.$validate()
  const timeDataValid = await updateTimeDataValidate.value.$validate()
  if (!formDataValid || !timeDataValid) {
    return
  }

  emits('submit:update')
}

const onSubmitDelete = (): void => {
  emits('submit:delete')
}
</script>

<template>
  <v-dialog v-model="createDialogValue" width="500">
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
            v-model="createTimeDataValidate.startDate.$model"
            :error-messages="
              getErrorMessage(createTimeDataValidate.startDate.$errors)
            "
            type="date"
            variant="outlined"
            density="compact"
            class="mr-md-2"
            @update:model-value="onChangeCreateStartAt"
          />
          <v-text-field
            v-model="createTimeDataValidate.startTime.$model"
            :error-messages="
              getErrorMessage(createTimeDataValidate.startTime.$errors)
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
            v-model="createTimeDataValidate.endDate.$model"
            :error-messages="
              getErrorMessage(createTimeDataValidate.endDate.$errors)
            "
            type="date"
            variant="outlined"
            density="compact"
            class="mr-md-2"
            @update:model-value="onChangeCreateEndAt"
          />
          <v-text-field
            v-model="createTimeDataValidate.endTime.$model"
            :error-messages="
              getErrorMessage(createTimeDataValidate.endTime.$errors)
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
          :error-message="
            getErrorMessage(createFormDataValidate.comment.$errors)
          "
          label="概要"
          maxlength="2000"
        />
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn color="" variant="text" @click="onClickCloseCreateDialog">
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

  <v-dialog v-model="updateDialogValue" width="500">
    <v-card>
      <v-card-title class="text-h6 primaryLight">
        スケジュール更新
      </v-card-title>
      <v-card-text>
        <p class="text-subtitle-2 text-grey pb-2">
          ライブ配開始日時
        </p>
        <div class="d-flex flex-column flex-md-row justify-center">
          <v-text-field
            v-model="updateTimeDataValidate.startDate.$model"
            :error-messages="
              getErrorMessage(updateTimeDataValidate.startDate.$errors)
            "
            type="date"
            variant="outlined"
            density="compact"
            class="mr-md-2"
            @update:model-value="onChangeUpdateStartAt"
          />
          <v-text-field
            v-model="updateTimeDataValidate.startTime.$model"
            :error-messages="
              getErrorMessage(updateTimeDataValidate.startTime.$errors)
            "
            type="time"
            variant="outlined"
            density="compact"
            @update:model-value="onChangeUpdateStartAt"
          />
        </div>
        <p class="text-subtitle-2 text-grey pb-2">
          ライブ配終了日時
        </p>
        <div class="d-flex flex-column flex-md-row justify-center">
          <v-text-field
            v-model="updateTimeDataValidate.endDate.$model"
            :error-messages="
              getErrorMessage(updateTimeDataValidate.endDate.$errors)
            "
            type="date"
            variant="outlined"
            density="compact"
            class="mr-md-2"
            @update:model-value="onChangeUpdateEndAt"
          />
          <v-text-field
            v-model="updateTimeDataValidate.endTime.$model"
            :error-messages="
              getErrorMessage(updateTimeDataValidate.endTime.$errors)
            "
            type="time"
            variant="outlined"
            density="compact"
            @update:model-value="onChangeUpdateEndAt"
          />
        </div>
        <v-autocomplete
          v-model="liveValue.producerId"
          label="生産者"
          :items="producers"
          item-title="username"
          item-value="id"
          readonly
        />
        <v-autocomplete
          v-model="updateFormDataValidate.productIds.$model"
          :error-messages="
            getErrorMessage(updateFormDataValidate.productIds.$errors)
          "
          label="関連する商品"
          :items="getProductsByProducerId(liveValue.producerId)"
          item-title="name"
          item-value="id"
          chips
          closable-chips
          clearable
          multiple
          density="comfortable"
          @update:search="onSearchProductFromUpdate"
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
          v-model="updateFormDataValidate.comment.$model"
          :error-message="
            getErrorMessage(updateFormDataValidate.comment.$errors)
          "
          label="概要"
          maxlength="2000"
        />
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn color="" variant="text" @click="onClickCloseUpdateDialog">
          キャンセル
        </v-btn>
        <v-btn
          :loading="loading"
          color="error"
          variant="text"
          @click="onSubmitDelete"
        >
          削除
        </v-btn>
        <v-btn
          :loading="loading"
          color="primary"
          variant="text"
          @click="onSubmitUpdate"
        >
          更新
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
      <v-card v-for="(item, i) in props.lives" :key="`live-${i}`" class="mb-4">
        <v-card-title>
          <v-list-item>
            <template #prepend>
              <v-avatar>
                <v-img
                  cover
                  :src="getProducerThumbnailUrl(item)"
                  :srcset="getProducerThumbnails(item)"
                />
              </v-avatar>
            </template>
            <v-list-item-title>{{ getProducerName(item) }}</v-list-item-title>
            <v-list-item-subtitle>{{ getLiveTerm(item) }}</v-list-item-subtitle>
            <template #append>
              <v-btn
                variant="outlined"
                color="primary"
                size="small"
                @click.stop="onClickEdit(item.id)"
              >
                <v-icon size="small" :icon="mdiPencil" />
              </v-btn>
            </template>
          </v-list-item>
        </v-card-title>

        <v-card-text>
          <v-row>
            <v-col sm="12">
              <p class="text-subtitle-2 text-grey pb-2">
                概要
              </p>
              <p class="text-subtitle-2" v-html="item.comment" />
            </v-col>
            <v-col sm="12">
              <p class="text-subtitle-2 text-grey pb-2">
                関連商品
              </p>
              <v-table>
                <thead>
                  <tr>
                    <th />
                    <th>商品名</th>
                    <th>価格</th>
                    <th>在庫</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="(product, j) in getProductsByLive(item)"
                    :key="`product-${j}`"
                  >
                    <td>
                      <v-img
                        aspect-ratio="1/1"
                        :max-height="56"
                        :max-width="80"
                        :src="getProductThumbnailUrl(product)"
                        :srcset="getProductThumbnails(product)"
                      />
                    </td>
                    <td>{{ product.name }}</td>
                    <td>{{ product.price }}</td>
                    <td :class="getProductInventoryColor(product)">
                      {{ product.inventory }}
                    </td>
                  </tr>
                </tbody>
              </v-table>
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
    </v-col>
    <v-col sm="12">
      <v-btn block variant="outlined" color="primary" @click="onClickNew">
        <v-icon :icon="mdiPlus" />
        生産者と商品を追加
      </v-btn>
    </v-col>
  </v-row>
</template>
