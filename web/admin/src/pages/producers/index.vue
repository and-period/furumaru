<script lang="ts" setup>
import { mdiPlus, mdiAccount, mdiDelete } from '@mdi/js'
import { VDataTable } from 'vuetify/lib/labs/components'

import { getResizedImages } from '~/lib/helpers'
import { useAlert, usePagination } from '~/lib/hooks'
import { useCommonStore, useProducerStore } from '~/store'
import { ProducersResponseProducersInner } from '~/types/api'
import { ApiBaseError } from '~/types/exception'

const router = useRouter()

const { isShow, alertText, alertType, show } = useAlert('error')
const producerStore = useProducerStore()
const { addSnackbar } = useCommonStore()

const deleteDialog = ref<boolean>(false)
const selectedId = ref<string>('')

const producers = computed(() => {
  return producerStore.producers
})

const totalItems = computed(() => {
  return producerStore.totalItems
})

const {
  updateCurrentPage,
  itemsPerPage,
  handleUpdateItemsPerPage,
  options,
  offset
} = usePagination()

watch(itemsPerPage, () => {
  producerStore.fetchProducers(itemsPerPage.value, 0, '')
})

const selectedItemName = computed(() => {
  const selectedItem = producers.value.find(
    item => item.id === selectedId.value
  )
  return selectedItem
    ? `${selectedItem.lastname} ${selectedItem.firstname}`
    : ''
})

const handleUpdatePage = async (page: number) => {
  updateCurrentPage(page)
  await producerStore.fetchProducers(itemsPerPage.value, offset.value, '')
}

const fetchState = useAsyncData(async () => {
  try {
    await producerStore.fetchProducers(itemsPerPage.value, offset.value)
  } catch (error) {
    const errorMessage =
      error instanceof ApiBaseError
        ? error.message
        : '不明なエラーが発生しました。'
    show(errorMessage)
  }
})

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
  },
  {
    title: '動画',
    key: 'video',
    sortable: false
  }
]

const handleClickAddButton = () => {
  router.push('/producers/add')
}

const handleClickRow = (item: ProducersResponseProducersInner) => {
  router.push(`/producers/edit/${item.id}`)
}

const handleClickCancelButton = () => {
  deleteDialog.value = false
}

const handleClickDeleteButton = (item: ProducersResponseProducersInner) => {
  selectedId.value = item.id
  deleteDialog.value = true
}

const handleDeleteFormSubmit = async () => {
  try {
    await producerStore.deleteProducer(selectedId.value)
    addSnackbar({
      color: 'info',
      message: '生産者を削除しました。'
    })
    fetchState.refresh()
  } catch (error) {
    const errorMessage =
      error instanceof ApiBaseError
        ? error.message
        : '不明なエラーが発生しました。'
    show(errorMessage)
  }
  deleteDialog.value = false
}

const handleAddVideo = (item: ProducersResponseProducersInner) => {
  console.log(item)
}

const isLoading = (): boolean => {
  return fetchState?.pending?.value || false
}

const getImages = (producer: ProducersResponseProducersInner): string => {
  if (!producer.thumbnails) {
    return ''
  }
  return getResizedImages(producer.thumbnails)
}

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <div>
    <v-card-title class="d-flex flex-row">
      生産者管理
      <v-spacer />
      <v-btn variant="outlined" color="primary" @click="handleClickAddButton">
        <v-icon start :icon="mdiPlus" />
        生産者登録
      </v-btn>
    </v-card-title>

    <v-alert v-model="isShow" :type="alertType" class="my-2" closable>
      {{ alertText }}
    </v-alert>

    <v-dialog v-model="deleteDialog" width="500">
      <v-card>
        <v-card-title>
          {{ selectedItemName }}を本当に削除しますか？
        </v-card-title>
        <v-card-actions>
          <v-spacer />
          <v-btn color="error" variant="text" @click="handleClickCancelButton">
            キャンセル
          </v-btn>
          <v-btn color="primary" variant="outlined" @click="handleDeleteFormSubmit">
            削除
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-card class="mt-4" flat :loading="isLoading()">
      <v-card-text>
        <v-data-table-server
          :headers="headers"
          :items="producers"
          :items-length="totalItems"
          :footer-props="options"
          hover
          no-data-text="登録されている生産者がいません。"
          @update:items-per-page="handleUpdateItemsPerPage"
          @update:page="handleUpdatePage"
          @click:row="(_:any, {item}:any) => handleClickRow(item.raw)"
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
              color="primary"
              size="small"
              variant="outlined"
              @click.stop="handleClickDeleteButton(item.raw)"
            >
              <v-icon size="small" :icon="mdiDelete" />
              削除
            </v-btn>
          </template>
          <template #[`item.video`]="{ item }">
            <v-btn variant="outlined" color="primary" size="small" @click.stop="handleAddVideo(item.raw)">
              <v-icon size="small" :icon="mdiPlus" />
              追加
            </v-btn>
          </template>
        </v-data-table-server>
      </v-card-text>
    </v-card>
  </div>
</template>
