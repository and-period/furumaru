<script lang="ts" setup>
import { mdiAccount, mdiDelete, mdiPlus, mdiAccountGroup } from '@mdi/js'
import type { VDataTable } from 'vuetify/lib/components/index.mjs'

import { convertI18nToJapanesePhoneNumber } from '~/lib/formatter'
import { getResizedImages } from '~/lib/helpers'
import type { AlertType } from '~/lib/hooks'
import { AdminType } from '~/types/api/v1'
import type { Shop, Coordinator, Producer } from '~/types/api/v1'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  adminType: {
    type: Number as PropType<AdminType>,
    default: AdminType.AdminTypeUnknown,
  },
  deleteDialog: {
    type: Boolean,
    default: false,
  },
  isAlert: {
    type: Boolean,
    default: false,
  },
  alertType: {
    type: String as PropType<AlertType>,
    default: undefined,
  },
  alertText: {
    type: String,
    default: '',
  },
  producers: {
    type: Array<Producer>,
    default: () => [],
  },
  shops: {
    type: Array<Shop>,
    default: () => [],
  },
  coordinators: {
    type: Array<Coordinator>,
    default: () => [],
  },
  tableItemsPerPage: {
    type: Number,
    default: 20,
  },
  tableItemsTotal: {
    type: Number,
    default: 0,
  },
})

const emit = defineEmits<{
  (e: 'click:update-page', page: number): void
  (e: 'click:update-items-per-page', page: number): void
  (e: 'click:row', producerId: string): void
  (e: 'click:add'): void
  (e: 'click:delete', producerId: string): void
  (e: 'update:delete-dialog', v: boolean): void
}>()

const headers: VDataTable['headers'] = [
  {
    title: '',
    key: 'thumbnail',
    sortable: false,
  },
  {
    title: '生産者名',
    key: 'username',
    sortable: false,
  },
  {
    title: '担当コーディネーター名',
    key: 'coordinatorName',
    sortable: false,
  },
  {
    title: 'メールアドレス',
    key: 'email',
    sortable: false,
  },
  {
    title: '電話番号',
    key: 'phoneNumber',
    sortable: false,
  },
  {
    title: '',
    key: 'actions',
    sortable: false,
  },
]

const selectedItem = ref<Producer>()

const deleteDialogValue = computed({
  get: (): boolean => props.deleteDialog,
  set: (val: boolean): void => emit('update:delete-dialog', val),
})

const isRegisterable = (): boolean => {
  return props.adminType === AdminType.AdminTypeCoordinator
}

const getShopName = (producer?: Producer) => {
  if (!producer) {
    return ''
  }
  const shops = props.shops.filter((shop: Shop): boolean => {
    return shop.producerIds.includes(producer.id)
  })
  if (shops.length === 0) {
    return ''
  }
  const shopNames = shops.map((shop: Shop): string => {
    return shop.name
  })
  return shopNames.join(', ')
}

const producerName = (producer?: Producer): string => {
  if (!producer) {
    return ''
  }
  return `${producer.lastname} ${producer.firstname}`
}

const getImages = (producer: Producer): string => {
  if (!producer.thumbnailUrl) {
    return ''
  }
  return getResizedImages(producer.thumbnailUrl)
}

const onClickOpenDeleteDialog = (producer: Producer): void => {
  selectedItem.value = producer
  deleteDialogValue.value = true
}

const onClickCloseDeleteDialog = (): void => {
  deleteDialogValue.value = false
}

const onClickUpdatePage = (page: number): void => {
  emit('click:update-page', page)
}

const onClickUpdateItemsPerPage = (page: number): void => {
  emit('click:update-items-per-page', page)
}

const onClickRow = (notificationId: string): void => {
  emit('click:row', notificationId)
}

const onClickAdd = (): void => {
  emit('click:add')
}

const onClickDelete = (): void => {
  emit('click:delete', selectedItem?.value?.id || '')
}
</script>

<template>
  <v-container class="pa-0">
    <v-alert
      v-show="props.isAlert"
      :type="props.alertType"
      class="mb-4"
      v-text="props.alertText"
    />

    <v-dialog
      v-model="deleteDialogValue"
      max-width="500"
    >
      <v-card class="delete-dialog-card">
        <v-card-title class="text-h6 font-weight-medium">
          削除の確認
        </v-card-title>
        <v-card-text class="text-body-1">
          <div class="d-flex align-center mb-3">
            <v-avatar
              v-if="selectedItem?.thumbnailUrl"
              size="48"
              class="mr-3"
            >
              <v-img :src="selectedItem.thumbnailUrl" />
            </v-avatar>
            <v-avatar
              v-else
              color="grey-lighten-3"
              size="48"
              class="mr-3"
            >
              <v-icon
                :icon="mdiAccount"
                color="grey"
              />
            </v-avatar>
            <div>
              <div class="font-weight-medium">
                {{ producerName(selectedItem) }}
              </div>
              <div class="text-caption text-grey">
                {{ selectedItem?.email }}
              </div>
            </div>
          </div>
          <v-alert
            type="warning"
            variant="tonal"
            density="compact"
            class="text-body-2"
          >
            この操作は取り消せません。生産者に関連するすべてのデータが削除されます。
          </v-alert>
        </v-card-text>
        <v-card-actions class="pa-4">
          <v-spacer />
          <v-btn
            variant="text"
            @click="onClickCloseDeleteDialog"
          >
            キャンセル
          </v-btn>
          <v-btn
            color="error"
            variant="elevated"
            :loading="loading"
            @click="onClickDelete"
          >
            削除する
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-card
      class="producer-list-card"
      elevation="2"
    >
      <v-card-title class="d-flex align-center pa-6">
        <div class="d-flex align-center">
          <v-icon
            :icon="mdiAccountGroup"
            size="28"
            class="mr-3 text-primary"
          />
          <div>
            <h2 class="text-h5 font-weight-bold">
              生産者管理
            </h2>
            <p class="text-caption text-grey mb-0">
              {{ props.tableItemsTotal }}件の生産者が登録されています
            </p>
          </div>
        </div>
        <v-spacer />
        <v-btn
          v-show="isRegisterable()"
          color="primary"
          variant="elevated"
          @click="onClickAdd"
        >
          <v-icon
            :icon="mdiPlus"
            start
          />
          生産者登録
        </v-btn>
      </v-card-title>
      <v-divider />

      <v-card-text class="pa-0">
        <v-data-table-server
          :headers="headers"
          :loading="loading"
          :items="producers"
          :items-per-page="props.tableItemsPerPage"
          :items-length="props.tableItemsTotal"
          class="producer-table"
          hover
          no-data-text="登録されている生産者がいません。"
          @update:page="onClickUpdatePage"
          @update:items-per-page="onClickUpdateItemsPerPage"
          @click:row="(_: any, { item }: any) => onClickRow(item.id)"
        >
          <template #[`item.thumbnail`]="{ item }">
            <v-avatar
              size="40"
              class="producer-avatar"
            >
              <v-img
                v-if="item.thumbnailUrl !== ''"
                cover
                :src="item.thumbnailUrl"
                :srcset="getImages(item)"
              />
              <v-icon
                v-else
                :icon="mdiAccount"
                color="grey"
              />
            </v-avatar>
          </template>
          <template #[`item.coordinatorName`]="{ item }">
            <span class="text-body-2">
              {{ getShopName(item) || '未設定' }}
            </span>
          </template>
          <template #[`item.phoneNumber`]="{ item }">
            <span class="text-body-2">
              {{ convertI18nToJapanesePhoneNumber(item.phoneNumber) }}
            </span>
          </template>
          <template #[`item.actions`]="{ item }">
            <v-btn
              color="error"
              size="small"
              variant="text"
              icon
              @click.stop="onClickOpenDeleteDialog(item)"
            >
              <v-icon
                size="small"
                :icon="mdiDelete"
              />
            </v-btn>
          </template>
        </v-data-table-server>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<style scoped>
.producer-list-card {
  border-radius: 12px;
}

.producer-table {
  border-radius: 0 0 12px 12px;
}

/* stylelint-disable-next-line selector-class-pattern */
.producer-table :deep(.v-table__wrapper) {
  border-radius: 0 0 12px 12px;
}

/* stylelint-disable-next-line selector-class-pattern */
.producer-table :deep(tbody tr) {
  cursor: pointer;
  transition: background-color 0.2s ease;
}

/* stylelint-disable-next-line selector-class-pattern */
.producer-table :deep(tbody tr:hover) {
  background-color: rgb(0 0 0 / 2%);
}

.producer-avatar {
  border: 2px solid rgb(0 0 0 / 8%);
}

.delete-dialog-card {
  border-radius: 12px;
}

@media (width <= 600px) {
  .producer-list-card {
    border-radius: 0;
  }

  /* stylelint-disable-next-line selector-class-pattern */
  .producer-table :deep(.v-data-table__td) {
    padding: 8px;
  }
}
</style>
