<script lang="ts" setup>
import { mdiDelete, mdiPlus } from '@mdi/js'
import type { VDataTable } from 'vuetify/lib/components/index.mjs'

import { convertI18nToJapanesePhoneNumber } from '~/lib/formatter'
import type { AlertType } from '~/lib/hooks'
import type { Administrator } from '~/types/api'

interface Props {
  loading: boolean
  deleteDialog: boolean
  isAlert: boolean
  alertType: AlertType | undefined
  alertText: string
  administrators: Administrator[]
  tableItemsPerPage: number
  tableItemsTotal: number
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  deleteDialog: false,
  isAlert: false,
  alertText: '',
  administrators: () => [],
  tableItemsPerPage: 20,
  tableItemsTotal: 0,
})

const emit = defineEmits<{
  (e: 'click:update-page', page: number): void
  (e: 'click:update-items-per-page', page: number): void
  (e: 'click:row', notificationId: string): void
  (e: 'click:add'): void
  (e: 'click:delete', notificationId: string): void
  (e: 'update:delete-dialog', v: boolean): void
}>()

const headers: VDataTable['headers'] = [
  {
    title: '管理者名',
    key: 'name',
  },
  {
    title: 'Email',
    key: 'email',
  },
  {
    title: '電話番号',
    key: 'phoneNumber',
  },
  {
    title: 'Actions',
    key: 'actions',
    sortable: false,
  },
]

const selectedItem = ref<Administrator>()

const deleteDialogValue = computed({
  get: (): boolean => props.deleteDialog,
  set: (val: boolean): void => emit('update:delete-dialog', val),
})

const getName = (administrator?: Administrator): string => {
  if (!administrator) {
    return ''
  }
  return `${administrator.lastname} ${administrator.firstname}`
}

const onClickUpdatePage = (page: number): void => {
  emit('click:update-page', page)
}

const onClickUpdateItemsPerPage = (page: number): void => {
  emit('click:update-items-per-page', page)
}

const onClickRow = (promotionId: string): void => {
  emit('click:row', promotionId)
}

const onClickOpenDeleteDialog = (administrator: Administrator): void => {
  selectedItem.value = administrator
  deleteDialogValue.value = true
}

const onClickCloseDeleteDialog = (): void => {
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
      <v-card-title class="text-h7">
        {{ getName(selectedItem) }}を本当に削除しますか？
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
          :loading="loading"
          color="primary"
          variant="outlined"
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
      管理者管理
      <v-spacer />
      <v-btn
        variant="outlined"
        color="primary"
        @click="onClickAdd"
      >
        <v-icon
          start
          :icon="mdiPlus"
        />
        管理者登録
      </v-btn>
    </v-card-title>

    <v-card-text>
      <v-data-table-server
        :headers="headers"
        :loading="loading"
        :items="props.administrators"
        :items-per-page="props.tableItemsPerPage"
        :items-length="props.tableItemsTotal"
        hover
        no-data-text="登録されている管理者がいません。"
        @update:page="onClickUpdatePage"
        @update:items-per-page="onClickUpdateItemsPerPage"
        @click:row="(_: any, { item }: any) => onClickRow(item.id)"
      >
        <template #[`item.name`]="{ item }">
          {{ getName(item) }}
        </template>
        <template #[`item.phoneNumber`]="{ item }">
          {{ convertI18nToJapanesePhoneNumber(item.phoneNumber) }}
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn
            color="primary"
            size="small"
            variant="outlined"
            :prepend-icon="mdiDelete"
            @click.stop="onClickOpenDeleteDialog(item)"
          >
            削除
          </v-btn>
        </template>
      </v-data-table-server>
    </v-card-text>
  </v-card>
</template>
