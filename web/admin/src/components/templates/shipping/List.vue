<script lang="ts" setup>
import { mdiPlus, mdiDelete } from '@mdi/js'
import type { VDataTable } from 'vuetify/components'

import { dateTimeFormatter } from '~/lib/formatter'
import type { AlertType } from '~/lib/hooks'
import { AdminType } from '~/types/api/v1'
import type { Coordinator, Shipping } from '~/types/api/v1'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  adminType: {
    type: Number as PropType<AdminType>,
    default: AdminType.AdminTypeUnknown,
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
  shippings: {
    type: Array<Shipping>,
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
  (e: 'click:row', shippingId: string): void
  (e: 'click:add'): void
  (e: 'click:delete', shippingId: string): void
  (e: 'update:page', page: number): void
  (e: 'update:items-per-page', page: number): void
  (e: 'submit:delete'): void
}>()

const headers: VDataTable['headers'] = [
  {
    title: '名前',
    key: 'name',
    sortable: false,
  },
  {
    title: 'コーディネーター名',
    key: 'coordinatorId',
    sortable: false,
  },
  {
    title: 'デフォルト設定',
    key: 'isDefault',
    sortable: false,
  },
  {
    title: '更新日時',
    key: 'updatedAt',
    sortable: false,
  },
  {
    title: '',
    key: 'actions',
    sortable: false,
  },
]

const dialogVisible = ref(false)

const isRegisterable = (): boolean => {
  return props.adminType === AdminType.AdminTypeCoordinator
}

const getCoordinatorName = (coordinatorId: string) => {
  const coordinator = props.coordinators.find((coordinator: Coordinator): boolean => {
    return coordinator.id === coordinatorId
  })
  return coordinator ? coordinator.username : ''
}

const getIsDefault = (isDefault: boolean): string => {
  return isDefault ? 'デフォルト' : '-'
}

const getIsDefaultColor = (isDefault: boolean): string => {
  return isDefault ? 'primary' : ''
}

const onClickUpdatePage = (page: number): void => {
  emit('update:page', page)
}

const onClickUpdateItemsPerPage = (page: number): void => {
  emit('update:items-per-page', page)
}

const onClickRow = (shippingId: string): void => {
  emit('click:row', shippingId)
}

const onClickAdd = (): void => {
  emit('click:add')
}

const onClickDelete = (shippingId: string): void => {
  emit('click:delete', shippingId)
}

const onSubmitDelete = (): void => {
  emit('submit:delete')
  dialogVisible.value = false
}
</script>

<template>
  <atoms-app-alert
    :show="props.isAlert"
    :type="props.alertType"
    :text="props.alertText"
  />

  <atoms-app-confirm-dialog
    v-model="dialogVisible"
    title="削除の確認"
    message="本当に削除しますか？"
    :loading="props.loading"
    @confirm="onSubmitDelete"
  />

  <v-card
    class="mt-4"
    flat
  >
    <v-card-title class="d-flex flex-row">
      配送設定一覧
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
        配送情報登録
      </v-btn>
    </v-card-title>

    <v-card-text>
      <v-skeleton-loader
        v-if="loading"
        type="table-heading, table-row-divider@5"
      />
      <v-data-table-server
        v-else
        :headers="headers"
        :items="shippings"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        hover
        no-data-text="登録されている配送設定がありません"
        @update:page="onClickUpdatePage"
        @update:items-per-page="onClickUpdateItemsPerPage"
        @click:row="(_: any, { item }: any) => onClickRow(item.id)"
      >
        <template #[`item.coordinatorId`]="{ item }">
          {{ getCoordinatorName(item.coordinatorId) }}
        </template>
        <template #[`item.isDefault`]="{ item }">
          <v-chip
            size="small"
            :color="getIsDefaultColor(item.isDefault)"
          >
            {{ getIsDefault(item.isDefault) }}
          </v-chip>
        </template>
        <template #[`item.updatedAt`]="{ item }">
          {{ dateTimeFormatter(item.updatedAt) }}
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn
            color="primary"
            size="small"
            variant="outlined"
            @click.stop="onClickDelete(item)"
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
