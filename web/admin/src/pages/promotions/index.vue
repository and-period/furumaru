<template>
  <div>
    <v-card-title>セール情報</v-card-title>
    <div class="d-flex">
      <v-spacer />
      <v-btn outlined color="primary" @click="handleClickAddButton">
        <v-icon left>mdi-plus</v-icon>
        セール情報登録
      </v-btn>
    </div>
    <v-card class="mt-4" flat>
      <v-data-table
        :headers="headers"
        :items="promotions"
        no-data-text="登録されているセール情報がありません。"
      >
        <template #[`item.title`]="{ item }">
          {{ item.title }}
        </template>
        <template #[`item.public`]="{ item }">
          <v-chip small :color="getStatusColor(item.public)">
            {{ getStatus(item.public) }}
          </v-chip>
        </template>
        <template #[`item.code`]="{ item }">
          {{ item.code }}
        </template>
        <template #[`item.discount`]="{ item }">
          {{ getDiscount(item.discountType, item.discountRate) }}
        </template>
        <template #[`item.startAt`]="{ item }">
          {{ getDay(item.startAt) }}
        </template>
        <template #[`item.endAt`]="{ item }">
          {{ getDay(item.endAt) }}
        </template>
        <template #[`item.actions`]="{ item }">
          <v-btn outlined color="primary" small @click="handleEdit(item)">
            <v-icon small>mdi-pencil</v-icon>
            編集
          </v-btn>
        </template>
      </v-data-table>
    </v-card>
  </div>
</template>

<script lang="ts">
import { computed, useFetch, useRouter } from '@nuxtjs/composition-api'
import { defineComponent } from '@vue/composition-api'
import dayjs from 'dayjs'
import { DataTableHeader } from 'vuetify'

import { usePromotionStore } from '~/store/promotion'

export default defineComponent({
  setup() {
    const router = useRouter()
    const promotionStore = usePromotionStore()
    const promotions = computed(() => {
      return promotionStore.promotions
    })

    const headers: DataTableHeader[] = [
      {
        text: 'タイトル',
        value: 'title',
      },
      {
        text: 'ステータス',
        value: 'public',
      },
      {
        text: '割引コード',
        value: 'code',
      },
      {
        text: '割引方法',
        value: 'discount',
      },
      {
        text: '使用開始',
        value: 'startAt',
      },
      {
        text: '使用終了',
        value: 'endAt',
      },
      {
        text: 'Actions',
        value: 'actions',
        sortable: false,
      },
    ]

    const getDiscount = (
      discountType: number,
      discountRate: number
    ): string => {
      switch (discountType) {
        case 1:
          return '-' + discountRate + '円'
        case 2:
          return '-' + discountRate + '%'
        case 3:
          return '送料無料'
        default:
          return ''
      }
    }

    const handleClickAddButton = () => {
      router.push('/promotions/add')
    }

    const getStatus = (status: boolean): string => {
      if (status) {
        return '有効'
      } else {
        return '無効'
      }
    }

    const getStatusColor = (status: boolean): string => {
      if (status) {
        return 'primary'
      } else {
        return 'accentDarken'
      }
    }

    const getDay = (unixTime: number): string => {
      return dayjs.unix(unixTime).format('YYYY/MM/DD hh:mm')
    }

    const { fetchState } = useFetch(async () => {
      try {
        await promotionStore.fetchPromotions()
      } catch (err) {
        console.log(err)
      }
    })

    return {
      headers,
      promotions,
      fetchState,
      handleClickAddButton,
      getDiscount,
      getStatus,
      getStatusColor,
      getDay,
    }
  },
})
</script>
