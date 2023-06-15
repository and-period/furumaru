<script lang="ts" setup>
import { mdiPlus, mdiAccount, mdiDelete } from '@mdi/js'
import { VDataTable } from 'vuetify/labs/components'

import { getResizedImages } from '~/lib/helpers'
import { AlertType } from '~/lib/hooks'
import { CoordinatorsResponseCoordinatorsInner } from '~/types/api'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  deleteDialog: {
    type: Boolean,
    default: false
  },
  isAlert: {
    type: Boolean,
    default: false
  },
  alertType: {
    type: String as PropType<AlertType>,
    default: undefined
  },
  alertText: {
    type: String,
    default: ''
  },
  sortBy: {
    type: Array as PropType<VDataTable['sortBy']>,
    default: () => []
  },
  coordinators: {
    type: Array<CoordinatorsResponseCoordinatorsInner>,
    default: () => []
  },
  tableItemsPerPage: {
    type: Number,
    default: 20
  },
  tableItemsTotal: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits<{
  (e: 'click:update-page', page: number): void
  (e: 'click:update-items-per-page', page: number): void
  (e: 'click:row', coordinatorId: string): void
  (e: 'click:add'): void
  (e: 'click:delete', coordinatorId: string): void
  (e: 'update:delete-dialog', v: boolean): void
  (e: 'update:sort-by', sortBy: VDataTable['sortBy']): void
}>()

const headers: VDataTable['headers'] = [
  {
    title: 'サムネイル',
    key: 'thumbnail'
  },
  {
    title: '店舗名',
    key: 'storeName'
  },
  {
    title: 'コーディネータ名',
    key: 'name'
  },
  {
    title: 'Email',
    key: 'email'
  },
  {
    title: '電話番号',
    key: 'phoneNumber'
  },
  {
    title: 'Actions',
    key: 'actions',
    sortable: false
  }
]

const selectedItem = ref<CoordinatorsResponseCoordinatorsInner>()

const deleteDialogValue = computed({
  get: () => props.deleteDialog,
  set: (val: boolean) => emit('update:delete-dialog', val)
})

const coordinatorName = (coordinator?: CoordinatorsResponseCoordinatorsInner): string => {
  if (!coordinator) {
    return ''
  }
  return `${coordinator.lastname} ${coordinator.firstname}`
}

const getImages = (coordinator: CoordinatorsResponseCoordinatorsInner): string => {
  if (!coordinator.thumbnails) {
    return ''
  }
  return getResizedImages(coordinator.thumbnails)
}

const onClickUpdatePage = (page: number): void => {
  emit('click:update-page', page)
}

const onClickUpdateItemsPerPage = (page: number): void => {
  emit('click:update-items-per-page', page)
}

const onClickUpdateSortBy = (sortBy: VDataTable['sortBy']): void => {
  emit('update:sort-by', sortBy)
}

const onClickRow = (coordinatorId: string): void => {
  emit('click:row', coordinatorId)
}

const onClickOpen = (coordinator: CoordinatorsResponseCoordinatorsInner): void => {
  selectedItem.value = coordinator
  deleteDialogValue.value = true
}

const onClickClose = (): void => {
  deleteDialogValue.value = false
}

const onClickAdd = (): void => {
  emit('click:add')
}

const onClickDelete = (): void => {
  emit('click:delete', selectedItem?.value?.id || '')
}
</script>

<template>
  <v-alert v-show="props.isAlert" :type="props.alertType" v-text="props.alertText" />

  <v-dialog v-model="deleteDialogValue" width="500">
    <v-card>
      <v-card-title>
        {{ coordinatorName(selectedItem) }}を本当に削除しますか？
      </v-card-title>
      <v-card-actions>
        <v-spacer />
        <v-btn color="error" variant="text" @click="onClickClose">
          キャンセル
        </v-btn>
        <v-btn :loading="loading" color="primary" variant="outlined" @click="onClickDelete">
          削除
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-card class="mt-4" flat :loading="loading">
    <v-card-title class="d-flex flex-row">
      コーディネータ管理
      <v-spacer />
      <v-btn variant="outlined" color="primary" @click="onClickAdd">
        <v-icon start :icon="mdiPlus" />
        コーディネータ登録
      </v-btn>
    </v-card-title>

    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :items="coordinators"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        :sort-by="props.sortBy"
        :multi-sort="true"
        hover
        no-data-text="登録されているコーディネータがいません。"
        @update:page="onClickUpdatePage"
        @update:items-per-page="onClickUpdateItemsPerPage"
        @update:sort-by="onClickUpdateSortBy"
        @update:sort-desc="onClickUpdateSortBy"
        @click:row="(_:any, {item}: any) => onClickRow(item.raw.id)"
      >
        <template #[`item.thumbnail`]="{ item }">
          <v-avatar>
            <v-img
              v-if="item.raw.thumbnailUrl !== ''"
              cover
              :src="item.raw.thumbnailUrl"
              :srcset="getImages(item.raw)"
              :alt="`${item.raw.storeName}-profile`"
            />
            <v-icon v-else :icon="mdiAccount" />
          </v-avatar>
        </template>
        <template #[`item.name`]="{ item }">
          {{ `${item.raw.lastname} ${item.raw.firstname}` }}
        </template>
        <template #[`item.phoneNumber`]="{ item }">
          {{ `${item.raw.phoneNumber}`.replace('+81', '0') }}
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn
            variant="outlined"
            color="primary"
            size="small"
            @click.stop="onClickOpen(item.raw)"
          >
            <v-icon size="small" :icon="mdiDelete" />
            削除
          </v-btn>
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
