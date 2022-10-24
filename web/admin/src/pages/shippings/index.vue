<template>
  <div>
    <v-card-title>配送設定一覧</v-card-title>
    <div class="d-flex">
      <v-spacer />
      <v-btn outlined color="primary" @click="handleClickAddButton">
        <v-icon left>mdi-plus</v-icon>
        配送情報登録
      </v-btn>
    </div>
    <v-card class="mt-4" flat :loading="fetchState.pending">
      <v-card-text>
        <v-data-table
          :headers="headers"
          :server-items-length="totalItems"
          :footer-props="options"
          :items="shippings"
          show-expand
          class="elevation-0"
          @update:page="updateCurrentPage"
          @update:items-per-page="handleUpdateItemsPerPage"
        >
          <template #[`item.hasFreeShipping`]="{ item }">
            <v-chip small>
              {{ item.hasFreeShipping ? '有り' : '無し' }}
            </v-chip>
          </template>

          <template #[`item.updatedAt`]="{ item }">
            {{ dateTimeFormatter(item.updatedAt) }}
          </template>

          <template #expanded-item="{ item }">
            <td :colspan="headers.length" class="pa-4">
              <div v-for="n in [60, 80, 100]" :key="n">
                <div class="row my-2">サイズ{{ n }}詳細</div>
                <v-row
                  v-for="(boxRate, i) in item[`box${n}Rates`]"
                  :key="i"
                  class="align-center"
                >
                  <v-col cols="1">
                    {{ boxRate.number }}
                  </v-col>
                  <v-col cols="1">
                    {{ boxRate.name }}
                  </v-col>
                  <v-col cols="1"> {{ moneyFormat(boxRate.price) }} 円 </v-col>
                  <v-col cols="9">
                    <v-select
                      :value="boxRate.prefectures"
                      :items="prefecturesList"
                      :label="`${boxRate.prefectures.length}/${prefecturesList.length}`"
                      multiple
                      hide-details
                    >
                      <template #selection="{ item: selectItem, index }">
                        <v-chip v-if="index < 5" small>
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
        </v-data-table>
      </v-card-text>
    </v-card>
  </div>
</template>

<script lang="ts">
import {
  computed,
  defineComponent,
  useFetch,
  useRouter,
} from '@nuxtjs/composition-api'
import { DataTableHeader } from 'vuetify'

import { prefecturesList } from '~/constants'
import { dateTimeFormatter, moneyFormat } from '~/lib/formatter'
import { usePagination } from '~/lib/hooks'
import { useShippingStore } from '~/store/shippings'
export default defineComponent({
  setup() {
    const shippingStore = useShippingStore()
    const router = useRouter()

    const totalItems = computed(() => {
      return shippingStore.totalItems
    })

    const shippings = computed(() => {
      return shippingStore.shippings
    })

    const headers: DataTableHeader[] = [
      {
        text: '名前',
        value: 'name',
      },
      {
        text: '配送無料オプション',
        value: 'hasFreeShipping',
      },
      {
        text: '更新日',
        value: 'updatedAt',
      },
    ]

    const {
      options,
      offset,
      itemsPerPage,
      updateCurrentPage,
      handleUpdateItemsPerPage,
    } = usePagination()

    const { fetchState } = useFetch(async () => {
      try {
        await shippingStore.fetchShippings(itemsPerPage.value, offset.value)
      } catch (err) {
        console.log(err)
      }
    })

    const handleClickAddButton = () => {
      router.push('/shippings/add')
    }

    return {
      fetchState, // 初回ロード状況
      totalItems, // サーバ上の配送情報の個数
      options, // ページネーションのオプション
      itemsPerPage, // 1ページあたりの表示件数
      headers, // テーブルヘッダー
      shippings, // 配送情報一覧
      prefecturesList, // 都道府県リスト
      // 関数
      dateTimeFormatter,
      moneyFormat,
      updateCurrentPage,
      handleUpdateItemsPerPage,
      handleClickAddButton,
    }
  },
})
</script>
