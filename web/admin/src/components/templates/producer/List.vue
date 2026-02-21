<script lang="ts" setup>
import { mdiAccount, mdiDelete, mdiPlus, mdiAccountGroup } from '@mdi/js'
import type { VDataTable } from 'vuetify/components'

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

const { dialogVisible, selectedItem, open: openDeleteDialog, close: closeDeleteDialog } = useDeleteDialog<Producer>()

const isRegisterable = (): boolean => {
  return props.adminType === AdminType.AdminTypeCoordinator
}

const isDeletable = (): boolean => {
  const targets: AdminType[] = [AdminType.AdminTypeAdministrator, AdminType.AdminTypeCoordinator]
  return targets.includes(props.adminType)
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
  emit('click:delete', selectedItem.value?.id || '')
  closeDeleteDialog()
}
</script>

<template>
  <v-container class="pa-0">
    <atoms-app-alert
      :show="props.isAlert"
      :type="props.alertType"
      :text="props.alertText"
      class="mb-4"
    />

    <atoms-app-confirm-dialog
      v-model="dialogVisible"
      title="生産者削除の確認"
      :message="`「${producerName(selectedItem)}」を削除しますか？`"
      :loading="loading"
      @confirm="onClickDelete"
    />

    <v-card
      class="mt-6"
      elevation="0"
      rounded="lg"
    >
      <v-card-title class="d-flex flex-column flex-sm-row align-start align-sm-center justify-space-between pa-4 pa-sm-6 pb-4">
        <div class="d-flex align-center mb-3 mb-sm-0">
          <v-icon
            :icon="mdiAccountGroup"
            size="28"
            class="mr-3 text-primary"
          />
          <div>
            <h1 class="text-h6 text-sm-h5 font-weight-bold text-primary">
              生産者管理
            </h1>
            <p class="text-caption text-sm-body-2 text-medium-emphasis ma-0">
              生産者の登録・編集・削除を行います
            </p>
          </div>
        </div>
        <div class="d-flex ga-3 w-100 w-sm-auto">
          <v-btn
            v-show="isRegisterable()"
            variant="elevated"
            color="primary"
            :size="$vuetify.display.smAndDown ? 'default' : 'large'"
            class="w-100 w-sm-auto"
            @click="onClickAdd"
          >
            <v-icon
              start
              :icon="mdiPlus"
            />
            生産者登録
          </v-btn>
        </div>
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
              v-show="isDeletable()"
              variant="outlined"
              color="error"
              size="small"
              :prepend-icon="mdiDelete"
              @click.stop="openDeleteDialog(item)"
            >
              削除
            </v-btn>
          </template>
        </v-data-table-server>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<style scoped>
.producer-avatar {
  border: 2px solid rgb(0 0 0 / 8%);
}
</style>
