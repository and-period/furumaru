<script lang="ts" setup>
import { mdiPlus, mdiPencil } from '@mdi/js'
import { VDataTable } from 'vuetify/lib/labs/components'
import { prefecturesList } from '~/constants'
import { dateTimeFormatter, moneyFormat } from '~/lib/formatter'
import { usePagination } from '~/lib/hooks'
import { useShippingStore } from '~/store'

const shippingStore = useShippingStore()
const router = useRouter()

const totalItems = computed(() => {
  return shippingStore.totalItems
})

const shippings = computed(() => {
  return shippingStore.shippings
})

const headers: VDataTable['headers'] = [
  {
    title: '名前',
    key: 'name'
  },
  {
    title: '配送無料オプション',
    key: 'hasFreeShipping'
  },
  {
    title: '更新日',
    key: 'updatedAt'
  },
  {
    title: '',
    key: 'actions'
  }
]

const {
  options,
  offset,
  itemsPerPage,
  updateCurrentPage,
  handleUpdateItemsPerPage
} = usePagination()

const fetchState = useAsyncData(async () => {
  try {
    await shippingStore.fetchShippings(itemsPerPage.value, offset.value)
  } catch (err) {
    console.log(err)
  }
})

const isLoading = (): boolean => {
  return fetchState?.pending?.value || false
}

const handleClickAddButton = () => {
  router.push('/shippings/add')
}

const handleClickEditButton = (id: string) => {
  router.push(`/shippings/edit/${id}`)
}

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <div>
    <v-card-title>
      配送設定一覧
      <v-spacer />
      <v-btn variant="outlined" color="primary" @click="handleClickAddButton">
        <v-icon start :icon="mdiPlus" />
        配送情報登録
      </v-btn>
    </v-card-title>
    <v-card class="mt-4" flat :loading="isLoading">
      <v-card-text>
        <v-data-table-server
          :headers="headers"
          :items-length="totalItems"
          :footer-props="options"
          :items="shippings"
          show-expand
          class="elevation-0"
          @update:page="updateCurrentPage"
          @update:items-per-page="handleUpdateItemsPerPage"
        >
          <template #[`item.hasFreeShipping`]="{ item }">
            <v-chip size="small">
              {{ item.raw.hasFreeShipping ? '有り' : '無し' }}
            </v-chip>
          </template>

          <template #[`item.updatedAt`]="{ item }">
            {{ dateTimeFormatter(item.raw.updatedAt) }}
          </template>

          <template #[`item.actions`]="{ item }">
            <v-btn
              variant="outlined"
              color="primary"
              size="small"
              @click="handleClickEditButton(item.raw.id)"
            >
              <v-icon :icon="mdiPencil" />
              編集
            </v-btn>
          </template>

          <template #expanded-item="{ item }">
            <td :colspan="headers.length" class="pa-4">
              <div v-for="n in [60, 80, 100]" :key="n">
                <div class="row my-2">
                  サイズ{{ n }}詳細
                </div>
                <v-row
                  v-for="(boxRate, i) in item.raw[`box${n}Rates`]"
                  :key="i"
                  class="align-center"
                >
                  <v-col cols="1">
                    {{ boxRate.number }}
                  </v-col>
                  <v-col cols="1">
                    {{ boxRate.name }}
                  </v-col>
                  <v-col cols="1">
                    {{ moneyFormat(boxRate.price) }} 円
                  </v-col>
                  <v-col cols="9">
                    <v-select
                      v-model="boxRate.prefectures"
                      :items="prefecturesList"
                      :label="`${boxRate.prefectures.length}/${prefecturesList.length}`"
                      hide-details
                    >
                      <template #selection="{ item: selectItem, index }">
                        <v-chip v-if="index < 5" size="small">
                          <span>{{ selectItem.text }}</span>
                        </v-chip>
                        <span
                          v-if="index === 5"
                          class="grey--text text-caption"
                        >
                          (+{{ boxRate.prefectures.length - 5 }} others)
                        </span>
                      </template>
                    </v-select>
                  </v-col>
                </v-row>
              </div>
            </td>
          </template>
        </v-data-table-server>
      </v-card-text>
    </v-card>
  </div>
</template>
