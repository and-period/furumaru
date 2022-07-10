<template>
  <div>
    <v-card-title>生産者管理</v-card-title>
    <div class="d-flex">
      <v-spacer />
      <v-btn outlined color="primary" @click="handleClickAddButton">
        <v-icon left>mdi-plus</v-icon>
        生産者登録
      </v-btn>
    </div>
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
        >
          <template #item.thumbnail="{ item }">
            <v-avatar>
              <img
                v-if="item.thumbnailUrl !== ''"
                :src="item.thumbnailUrl"
                :alt="`${item.storeName}-profile`"
              />
              <v-icon v-else>mdi-account</v-icon>
            </v-avatar>
          </template>
          <template #item.name="{ item }">
            {{ `${item.lastname} ${item.firstname}` }}
          </template>
          <template #item.actions="{ item }">
            <v-btn outlined color="primary" small @click="handleEdit(item)">
              <v-icon>mdi-plus</v-icon>
              編集
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

import { useProducerStore } from '~/store/producer'
import { ProducersResponseProducers } from '~/types/api'

export default defineComponent({
  setup() {
    const router = useRouter()
    const producerStore = useProducerStore()
    const producers = computed(() => {
      return producerStore.producers
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
    ]

    const handleClickAddButton = () => {
      router.push('/producers/add')
    }

    const handleSearch = () => {
      query.value = search.value
    }

    const handleEdit = (item: ProducersResponseProducers) => {
      console.log(item)
    }

    const { fetchState } = useFetch(async () => {
      try {
        await producerStore.fetchProducers()
      } catch (err) {
        console.log(err)
      }
    })

    return {
      handleClickAddButton,
      headers,
      producers,
      search,
      query,
      noResultsText,
      fetchState,
      handleSearch,
      handleEdit,
    }
  },
})
</script>
