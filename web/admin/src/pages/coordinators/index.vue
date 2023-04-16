<script lang="ts" setup>
import { mdiPlus, mdiSearchWeb, mdiAccount, mdiPencil, mdiDelete } from '@mdi/js'
import { VDataTable } from 'vuetify/labs/components'

import { useAlert, usePagination } from '~/lib/hooks'
import { useCommonStore, useCoordinatorStore } from '~/store'
import { CoordinatorsResponseCoordinatorsInner, ImageSize } from '~/types/api'
import { ApiBaseError } from '~/types/exception'

const router = useRouter()

const { isShow, alertText, alertType, show } = useAlert('error')

const { addSnackbar } = useCommonStore()

const coordinatorStore = useCoordinatorStore()
const coordinators = computed(() => {
  return coordinatorStore.coordinators
})

const totalItems = computed(() => {
  return coordinatorStore.totalItems
})

const deleteDialog = ref<boolean>(false)
const selectedId = ref<string>('')

const selectedItemName = computed(() => {
  const selectedItem = coordinators.value.find(
    item => item.id === selectedId.value
  )
  return selectedItem
    ? `${selectedItem.lastname} ${selectedItem.firstname}`
    : ''
})

const search = ref<string>('')
const query = ref<string>('')

const noResultsText = computed(() => {
  return `「${query.value}」に一致するデータはありません。`
})

watch(search, () => {
  if (search.value === '') {
    query.value = ''
  }
})

const {
  updateCurrentPage,
  itemsPerPage,
  handleUpdateItemsPerPage,
  options,
  offset
} = usePagination()

watch(itemsPerPage, () => {
  coordinatorStore.fetchCoordinators(itemsPerPage.value, 0)
})

const handleUpdatePage = async (page: number) => {
  updateCurrentPage(page)
  await coordinatorStore.fetchCoordinators(itemsPerPage.value, offset.value)
}

const fetchState = useAsyncData(async () => {
  try {
    await coordinatorStore.fetchCoordinators(itemsPerPage.value, offset.value)
  } catch (err) {
    console.log(err)
  }
})

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

const isLoading = (): boolean => {
  return fetchState?.pending?.value || false
}

const handleClickAddButton = () => {
  router.push('/coordinators/add')
}

const handleSearch = () => {
  query.value = search.value
}

const handleEdit = (item: CoordinatorsResponseCoordinatorsInner) => {
  router.push(`/coordinators/edit/${item.id}`)
}

const handleClickDeleteButton = (
  item: CoordinatorsResponseCoordinatorsInner
): void => {
  selectedId.value = item.id
  deleteDialog.value = true
}

const handleClickCancelButton = () => {
  deleteDialog.value = false
}

const handleDeleteFormSubmit = async () => {
  try {
    await coordinatorStore.deleteCoordinator(selectedId.value)
    addSnackbar({
      color: 'info',
      message: 'コーディネータを削除しました。'
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

const getImages = (coordinator: CoordinatorsResponseCoordinatorsInner): string => {
  if (!coordinator.thumbnails) {
    return ''
  }
  const images: string[] = coordinator.thumbnails.map((thumbnail): string => {
    switch (thumbnail.size) {
      case ImageSize.SMALL:
        return `${thumbnail.url} 1x`
      case ImageSize.MEDIUM:
        return `${thumbnail.url} 2x`
      case ImageSize.LARGE:
        return `${thumbnail.url} 3x`
      default:
        return thumbnail.url
    }
  })
  return images.join(', ')
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
      コーディネータ管理
      <v-spacer />
      <v-btn variant="outlined" color="primary" @click="handleClickAddButton">
        <v-icon start :icon="mdiPlus" />
        コーディネータ登録
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
        <form class="d-flex align-center" @submit.prevent="handleSearch">
          <v-select
            v-model="search"
            item-title="firstname"
            :items="coordinators"
            label="絞り込み"
            variant="underlined"
          />
          <v-btn type="submit" class="ml-4" variant="outlined" color="primary">
            <v-icon :icon="mdiSearchWeb" />
            検索
          </v-btn>
          <v-spacer />
        </form>
        <v-data-table-server
          show-select
          :headers="headers"
          :items="coordinators"
          :search="query"
          :no-results-text="noResultsText"
          :items-length="totalItems"
          :footer-props="options"
          no-data-text="登録されているコーディネータがいません。"
          @update:items-per-page="handleUpdateItemsPerPage"
          @update:page="handleUpdatePage"
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
            <v-btn variant="outlined" color="primary" size="small" @click="handleEdit(item.raw)">
              <v-icon size="small" :icon="mdiPencil" />
              編集
            </v-btn>
            <v-btn
              variant="outlined"
              color="primary"
              size="small"
              @click="handleClickDeleteButton(item.raw)"
            >
              <v-icon size="small" :icon="mdiDelete" />
              削除
            </v-btn>
          </template>
        </v-data-table-server>
      </v-card-text>
    </v-card>
  </div>
</template>
