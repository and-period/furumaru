<script lang="ts" setup>
import { mdiAccount, mdiDelete, mdiPlus } from '@mdi/js'
import type { VDataTable } from 'vuetify/lib/components/index.mjs'

import { convertI18nToJapanesePhoneNumber } from '~/lib/formatter'
import { getResizedImages } from '~/lib/helpers'
import type { AlertType } from '~/lib/hooks'
import { AdminType, type Coordinator, type Producer } from '~/types/api'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  adminType: {
    type: Number as PropType<AdminType>,
    default: AdminType.UNKNOWN,
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
  return props.adminType === AdminType.COORDINATOR
}

const getCoordinatorName = (coordinatorId: string) => {
  const coordinator = props.coordinators.find((coordinator: Coordinator): boolean => {
    return coordinator.id === coordinatorId
  })
  return coordinator ? coordinator.username : ''
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
  <v-alert
    v-show="props.isAlert"
    :type="props.alertType"
    v-text="props.alertText"
  />

  <v-dialog
    v-model="deleteDialogValue"
    width="500"
  >
    <v-card>
      <v-card-title>
        {{ producerName(selectedItem) }}を本当に削除しますか？
      </v-card-title>
      <v-card-actions>
        <v-spacer />
        <v-btn
          color="error"
          variant="text"
          @click="onClickCloseDeleteDialog"
        >
          キャンセル
        </v-btn>
        <v-btn
          color="primary"
          variant="outlined"
          :loading="loading"
          @click="onClickDelete"
        >
          削除
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-card
    class="mt-4"
    flat
  >
    <v-card-title class="d-flex flex-row">
      生産者管理
      <v-spacer />
      <v-btn
        v-show="isRegisterable()"
        variant="outlined"
        color="primary"
        @click="onClickAdd"
      >
        <v-icon
          start
          :icon="mdiPlus"
        />
        生産者登録
      </v-btn>
    </v-card-title>

    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :loading="loading"
        :items="producers"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        hover
        no-data-text="登録されている生産者がいません。"
        @update:page="onClickUpdatePage"
        @update:items-per-page="onClickUpdateItemsPerPage"
        @click:row="(_: any, { item }: any) => onClickRow(item.id)"
      >
        <template #[`item.thumbnail`]="{ item }">
          <v-avatar>
            <v-img
              v-if="item.thumbnailUrl !== ''"
              cover
              :src="item.thumbnailUrl"
              :srcset="getImages(item)"
            />
            <v-icon
              v-else
              :icon="mdiAccount"
            />
          </v-avatar>
        </template>
        <template #[`item.coordinatorName`]="{ item }">
          {{ getCoordinatorName(item.coordinatorId) }}
        </template>
        <template #[`item.phoneNumber`]="{ item }">
          {{ convertI18nToJapanesePhoneNumber(item.phoneNumber) }}
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn
            color="primary"
            size="small"
            variant="outlined"
            @click.stop="onClickOpenDeleteDialog(item)"
          >
            <v-icon
              size="small"
              :icon="mdiDelete"
            />削除
          </v-btn>
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
