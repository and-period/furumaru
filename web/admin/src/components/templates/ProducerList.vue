<script lang="ts" setup>
import { mdiAccount, mdiDelete, mdiPlus } from '@mdi/js'
import { VDataTable } from 'vuetify/lib/labs/components'
import { convertI18nToJapanesePhoneNumber } from '~/lib/formatter'
import { getResizedImages } from '~/lib/helpers'
import { AlertType } from '~/lib/hooks'
import { ProducersResponseProducersInner } from '~/types/api'

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
  producers: {
    type: Array<ProducersResponseProducersInner>,
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
  (e: 'click:row', producerId: string): void
  (e: 'click:add'): void
  (e: 'click:add-video', producerId: string): void
  (e: 'click:delete', producerId: string): void
  (e: 'update:delete-dialog', v: boolean): void
}>()

const headers: VDataTable['headers'] = [
  {
    title: 'サムネイル',
    key: 'thumbnail'
  },
  {
    title: '農園名',
    key: 'storeName'
  },
  {
    title: '生産者名',
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

const selectedItem = ref<ProducersResponseProducersInner>()

const deleteDialogValue = computed({
  get: (): boolean => props.deleteDialog,
  set: (val: boolean): void => emit('update:delete-dialog', val)
})

const getName = (producer?: ProducersResponseProducersInner): string => {
  if (!producer) {
    return ''
  }
  return `${producer.lastname} ${producer.firstname}`
}

const getImages = (producer: ProducersResponseProducersInner): string => {
  if (!producer.thumbnails) {
    return ''
  }
  return getResizedImages(producer.thumbnails)
}

const onClickOpenDeleteDialog = (producer: ProducersResponseProducersInner): void => {
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

const onClickAddVideo = (producerId: string): void => {
  emit('click:add-video', producerId)
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
        {{ getName(selectedItem) }}を本当に削除しますか？
      </v-card-title>
      <v-card-actions>
        <v-spacer />
        <v-btn color="error" variant="text" @click="onClickCloseDeleteDialog">
          キャンセル
        </v-btn>
        <v-btn color="primary" variant="outlined" :loading="loading" @click="onClickDelete">
          削除
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-card class="mt-4" flat :loading="props.loading">
    <v-card-title class="d-flex flex-row">
      生産者管理
      <v-spacer />
      <v-btn variant="outlined" color="primary" @click="onClickAdd">
        <v-icon start :icon="mdiPlus" />
        生産者登録
      </v-btn>
    </v-card-title>

    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :items="producers"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        hover
        no-data-text="登録されている生産者がいません。"
        @update:page="onClickUpdatePage"
        @update:items-per-page="onClickUpdateItemsPerPage"
        @click:row="(_: any, { item }:any) => onClickRow(item.raw.id)"
      >
        <template #[`item.thumbnail`]="{ item }">
          <v-avatar>
            <v-img
              v-if="item.raw.thumbnailUrl !== ''"
              aspect-ratio="1/1"
              cover
              :src="item.raw.thumbnailUrl"
              :srcset="getImages(item.raw)"
            />
            <v-icon v-else :icon="mdiAccount" />
          </v-avatar>
        </template>
        <template #[`item.name`]="{ item }">
          {{ getName(item.raw) }}
        </template>
        <template #[`item.phoneNumber`]="{ item }">
          {{ convertI18nToJapanesePhoneNumber(item.raw.phoneNumber) }}
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn
            color="primary"
            size="small"
            variant="outlined"
            class="mr-4"
            @click.stop="onClickOpenDeleteDialog(item.raw)"
          >
            <v-icon size="small" :icon="mdiDelete" />削除
          </v-btn>
          <v-btn variant="outlined" color="primary" size="small" @click.stop="onClickAddVideo(item.raw.id)">
            <v-icon size="small" :icon="mdiPlus" />
            追加
          </v-btn>
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
