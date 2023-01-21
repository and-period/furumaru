<template>
  <div>
    <v-card-title>
      生産者管理
      <v-spacer />
      <v-btn outlined color="primary" @click="handleClickAddButton">
        <v-icon left>mdi-plus</v-icon>
        生産者登録
      </v-btn>
    </v-card-title>

    <v-alert v-model="isShow" :type="alertType" class="my-2" dismissible>
      {{ alertText }}
    </v-alert>

    <v-dialog v-model="deleteDialog" width="500">
      <v-card>
        <v-card-title>
          {{ selectedItemName }}を本当に削除しますか？
        </v-card-title>
        <v-card-actions>
          <v-spacer />
          <v-btn color="error" text @click="handleClickCancelButton">
            キャンセル
          </v-btn>
          <v-btn color="primary" outlined @click="handleDeleteFormSubmit">
            削除
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-card class="mt-4" flat :loading="fetchState.pending">
      <v-card-text>
        <form class="d-flex align-center" @submit.prevent="handleSearch">
          <v-text-field v-model="search" label="絞り込み" />
          <v-btn type="submit" class="ml-4" small outlined color="primary">
            <v-icon>mdi-search</v-icon>
            検索
          </v-btn>
          <v-spacer />
        </form>
        <v-data-table
          show-select
          :headers="headers"
          :items="producers"
          :search="query"
          :no-results-text="noResultsText"
          :server-items-length="totalItems"
          :footer-props="options"
          no-data-text="登録されている生産者がいません。"
          @update:items-per-page="handleUpdateItemsPerPage"
          @update:page="handleUpdatePage"
        >
          <template #[`item.thumbnail`]="{ item }">
            <v-avatar>
              <img
                v-if="item.thumbnailUrl !== ''"
                :src="item.thumbnailUrl"
                :alt="`${item.storeName}-profile`"
              />
              <v-icon v-else>mdi-account</v-icon>
            </v-avatar>
          </template>
          <template #[`item.name`]="{ item }">
            {{ `${item.lastname} ${item.firstname}` }}
          </template>
          <template #[`item.phoneNumber`]="{ item }">
            {{ `${item.phoneNumber}`.replace('+81', '0') }}
          </template>
          <template #[`item.actions`]="{ item }">
            <v-btn outlined color="primary" small @click="handleEdit(item)">
              <v-icon small>mdi-pencil</v-icon>
              編集
            </v-btn>
            <v-btn
              outlined
              color="primary"
              small
              @click="handleClickDeleteButton(item)"
            >
              <v-icon small>mdi-delete</v-icon>
              削除
            </v-btn>
          </template>
          <template #[`item.video`]="{ item }">
            <v-btn outlined color="primary" small @click="handleAddVideo(item)">
              <v-icon small>mdi-plus</v-icon>
              追加
            </v-btn>
          </template>
        </v-data-table>
      </v-card-text>
    </v-card>
  </div>
</template>

<script lang="ts">
import {
  computed,
  defineComponent,
  ref,
  useFetch,
  useRouter,
  watch,
} from '@nuxtjs/composition-api'
import { DataTableHeader } from 'vuetify'

import { useAlert, usePagination } from '~/lib/hooks'
import { useCommonStore } from '~/store/common'
import { useProducerStore } from '~/store/producer'
import { ProducersResponseProducersInner } from '~/types/api'
import { ApiBaseError } from '~/types/exception'

export default defineComponent({
  setup() {
    const router = useRouter()

    const { isShow, alertText, alertType, show } = useAlert('error')
    const producerStore = useProducerStore()
    const { addSnackbar } = useCommonStore()

    const search = ref<string>('')
    const query = ref<string>('')

    const deleteDialog = ref<boolean>(false)
    const selectedId = ref<string>('')

    const producers = computed(() => {
      return producerStore.producers
    })

    const totalItems = computed(() => {
      return producerStore.totalItems
    })

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
      offset,
    } = usePagination()

    watch(itemsPerPage, () => {
      producerStore.fetchProducers(itemsPerPage.value, 0, '')
    })

    const selectedItemName = computed(() => {
      const selectedItem = producers.value.find(
        (item) => item.id === selectedId.value
      )
      return selectedItem
        ? `${selectedItem.lastname} ${selectedItem.firstname}`
        : ''
    })

    const handleUpdatePage = async (page: number) => {
      updateCurrentPage(page)
      await producerStore.fetchProducers(itemsPerPage.value, offset.value, '')
    }

    const { fetchState, fetch } = useFetch(async () => {
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

    const headers: DataTableHeader[] = [
      {
        text: 'サムネイル',
        value: 'thumbnail',
      },
      {
        text: '農園名',
        value: 'storeName',
      },
      {
        text: '生産者名',
        value: 'name',
      },
      {
        text: 'Email',
        value: 'email',
      },
      {
        text: '電話番号',
        value: 'phoneNumber',
      },
      {
        text: 'Actions',
        value: 'actions',
        sortable: false,
      },
      {
        text: '動画',
        value: 'video',
        sortable: false,
      },
    ]

    const handleClickAddButton = () => {
      router.push('/producers/add')
    }

    const handleSearch = () => {
      query.value = search.value
    }

    const handleEdit = (item: ProducersResponseProducersInner) => {
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
          message: '生産者を削除しました。',
        })
        fetch()
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

    return {
      // 定数
      headers,
      options,
      noResultsText,
      alertType,
      // 変数
      isShow,
      alertText,
      deleteDialog,
      fetchState,
      producers,
      totalItems,
      search,
      query,
      selectedItemName,
      // 関数
      handleClickAddButton,
      handleUpdatePage,
      handleUpdateItemsPerPage,
      handleSearch,
      handleEdit,
      handleClickDeleteButton,
      handleClickCancelButton,
      handleDeleteFormSubmit,
      handleAddVideo,
    }
  },
})
</script>
