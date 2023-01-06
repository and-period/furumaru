<template>
  <div>
    <v-card-title>コーディネータ管理</v-card-title>

    <div class="d-flex">
      <v-spacer />
      <v-btn outlined color="primary" @click="handleClickAddButton">
        <v-icon left>mdi-plus</v-icon>
        コーディネータ登録
      </v-btn>
    </div>

    <v-alert v-model="isShow" :type="alertType" class="my-2" dismissible>
      {{ alertText }}
    </v-alert>

    <v-dialog v-model="deleteDialog" width="500">
      <v-card>
        <v-card-title
          >{{ selectedItemName }}を本当に削除しますか？</v-card-title
        >
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
          <v-autocomplete
            v-model="search"
            item-text="firstname"
            :items="coordinators"
            label="絞り込み"
          />
          <v-btn type="submit" class="ml-4" small outlined color="primary">
            <v-icon>mdi-search</v-icon>
            検索
          </v-btn>
          <v-spacer />
        </form>
        <v-data-table
          show-select
          :headers="headers"
          :items="coordinators"
          :search="query"
          :no-results-text="noResultsText"
          :server-items-length="totalItems"
          :footer-props="options"
          no-data-text="登録されているコーディネータがいません。"
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
import { useCoordinatorStore } from '~/store/coordinator'
import { CoordinatorsResponseCoordinatorsInner } from '~/types/api'
import { ApiBaseError } from '~/types/exception'
import { Coordinator } from '~/types/props/coordinator'

export default defineComponent({
  setup() {
    const tab = ref<string>('coordinators')
    const tabItems: Coordinator[] = [
      { name: '基本情報', value: 'coordinators' },
      { name: '生産者管理', value: 'customers' },
    ]

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
    const selectedName = ref<string>('')

    const selectedItemName = computed(() => {
      const selectedItem = coordinators.value.find(
        (item) => item.id === selectedId.value
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
      offset,
    } = usePagination()

    watch(itemsPerPage, () => {
      coordinatorStore.fetchCoordinators(itemsPerPage.value, 0)
    })

    const handleUpdatePage = async (page: number) => {
      updateCurrentPage(page)
      await coordinatorStore.fetchCoordinators(itemsPerPage.value, offset.value)
    }

    const openDeleteDialog = (
      item: CoordinatorsResponseCoordinatorsInner
    ): void => {
      selectedId.value = item.id
      selectedName.value = item.firstname
      deleteDialog.value = true
    }

    const { fetchState, fetch } = useFetch(async () => {
      try {
        await coordinatorStore.fetchCoordinators(
          itemsPerPage.value,
          offset.value
        )
      } catch (err) {
        console.log(err)
      }
    })

    const headers: DataTableHeader[] = [
      {
        text: 'サムネイル',
        value: 'thumbnail',
      },
      {
        text: '店舗名',
        value: 'storeName',
      },
      {
        text: 'コーディネータ名',
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
    ]

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
          message: 'コーディネータを削除しました。',
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
      coordinators,
      totalItems,
      search,
      query,
      selectedItemName,
      // 関数
      handleSearch,
      handleUpdatePage,
      handleUpdateItemsPerPage,
      handleClickAddButton,
      handleEdit,
      openDeleteDialog,
      tab,
      tabItems,
      handleClickCancelButton,
      handleClickDeleteButton,
      handleDeleteFormSubmit,
    }
  },
})
</script>
