<script lang="ts" setup>
import dayjs from 'dayjs'
import { DataTableHeader } from 'vuetify'

import { usePromotionStore } from '~/store/promotion'
import { PromotionsResponsePromotionsInner } from '~/types/api'

const router = useRouter()
const promotionStore = usePromotionStore()

const deleteDialog = ref<boolean>(false)
const selectedId = ref<string>('')
const selectedName = ref<string>('')

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
    text: '投稿開始',
    value: 'publishedAt',
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

const handleDelete = async (): Promise<void> => {
  try {
    await promotionStore.deletePromotion(selectedId.value)
  } catch (err) {
    console.log(err)
  }
  deleteDialog.value = false
}

const handleEdit = (item: PromotionsResponsePromotionsInner) => {
  router.push(`/promotions/edit/${item.id}`)
}

const openDeleteDialog = (
  item: PromotionsResponsePromotionsInner
): void => {
  selectedId.value = item.id
  selectedName.value = item.title
  deleteDialog.value = true
}

const hideDeleteDialog = () => {
  deleteDialog.value = false
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
    return 'error'
  }
}

const getDay = (unixTime: number): string => {
  return dayjs.unix(unixTime).format('YYYY/MM/DD HH:mm')
}

const fetchState = useAsyncData(async () => {
  try {
    await promotionStore.fetchPromotions()
  } catch (err) {
    console.log(err)
  }
})
</script>

<template>
  <div>
    <v-card-title>
      セール情報
      <v-spacer />
      <v-btn outlined color="primary" @click="handleClickAddButton">
        <v-icon left>mdi-plus</v-icon>
        セール情報登録
      </v-btn>
    </v-card-title>

    <v-dialog v-model="deleteDialog" width="500">
      <v-card>
        <v-card-title class="text-h7">
          {{ selectedName }}を本当に削除しますか？
        </v-card-title>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="error" text @click="hideDeleteDialog">
            キャンセル
          </v-btn>
          <v-btn color="primary" outlined @click="handleDelete"> 削除 </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-card class="mt-4" flat>
      <v-card-text>
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
          <template #[`item.publishedAt`]="{ item }">
            {{ getDay(item.publishedAt) }}
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
            <v-btn
              outlined
              color="primary"
              small
              @click="openDeleteDialog(item)"
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
