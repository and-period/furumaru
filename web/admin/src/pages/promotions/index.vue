<script lang="ts" setup>
import { mdiPlus, mdiDelete } from '@mdi/js'
import dayjs from 'dayjs'
import { VDataTable } from 'vuetify/lib/labs/components'

import { usePromotionStore } from '~/store'
import { PromotionsResponsePromotionsInner } from '~/types/api'

const router = useRouter()
const promotionStore = usePromotionStore()

const deleteDialog = ref<boolean>(false)
const selectedId = ref<string>('')
const selectedName = ref<string>('')

const promotions = computed(() => {
  return promotionStore.promotions
})

const headers: VDataTable['headers'] = [
  {
    title: 'タイトル',
    key: 'title'
  },
  {
    title: 'ステータス',
    key: 'public'
  },
  {
    title: '割引コード',
    key: 'code'
  },
  {
    title: '割引方法',
    key: 'discount'
  },
  {
    title: '投稿開始',
    key: 'publishedAt'
  },
  {
    title: '使用開始',
    key: 'startAt'
  },
  {
    title: '使用終了',
    key: 'endAt'
  },
  {
    title: 'Actions',
    key: 'actions',
    sortable: false
  }
]

const getDiscount = (discountType: number, discountRate: number): string => {
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

const handleClickRow = (item: PromotionsResponsePromotionsInner) => {
  router.push(`/promotions/edit/${item.id}`)
}

const openDeleteDialog = (item: PromotionsResponsePromotionsInner): void => {
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

try {
  await fetchState.execute()
} catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <div>
    <v-card-title class="d-flex flex-row">
      セール情報
      <v-spacer />
      <v-btn variant="outlined" color="primary" @click="handleClickAddButton">
        <v-icon start :icon="mdiPlus" />
        セール情報登録
      </v-btn>
    </v-card-title>

    <v-dialog v-model="deleteDialog" width="500">
      <v-card>
        <v-card-title class="text-h7">
          {{ selectedName }}を本当に削除しますか？
        </v-card-title>
        <v-card-actions>
          <v-spacer />
          <v-btn color="error" variant="text" @click="hideDeleteDialog">
            キャンセル
          </v-btn>
          <v-btn color="primary" variant="outlined" @click="handleDelete">
            削除
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-card class="mt-4" flat>
      <v-card-text>
        <v-data-table
          :headers="headers"
          :items="promotions"
          hover
          no-data-text="登録されているセール情報がありません。"
          @click:row="(_: any, {item}: any) => handleClickRow(item.raw)"
        >
          <template #[`item.title`]="{ item }">
            {{ item.raw.title }}
          </template>
          <template #[`item.public`]="{ item }">
            <v-chip size="small" :color="getStatusColor(item.raw.public)">
              {{ getStatus(item.raw.public) }}
            </v-chip>
          </template>
          <template #[`item.code`]="{ item }">
            {{ item.raw.code }}
          </template>
          <template #[`item.discount`]="{ item }">
            {{ getDiscount(item.raw.discountType, item.raw.discountRate) }}
          </template>
          <template #[`item.publishedAt`]="{ item }">
            {{ getDay(item.raw.publishedAt) }}
          </template>
          <template #[`item.startAt`]="{ item }">
            {{ getDay(item.raw.startAt) }}
          </template>
          <template #[`item.endAt`]="{ item }">
            {{ getDay(item.raw.endAt) }}
          </template>
          <template #[`item.actions`]="{ item }">
            <v-btn
              color="primary"
              size="small"
              variant="outlined"
              @click.stop="openDeleteDialog(item.raw)"
            >
              <v-icon size="small" :icon="mdiDelete" />
              削除
            </v-btn>
          </template>
        </v-data-table>
      </v-card-text>
    </v-card>
  </div>
</template>
