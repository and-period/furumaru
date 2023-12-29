<script lang="ts" setup>
import { mdiDotsVertical } from '@mdi/js'
import { unix } from 'dayjs'
import type { VDataTable } from 'vuetify/lib/components/index.mjs'

import { convertI18nToJapanesePhoneNumber } from '~/lib/formatter'
import type { AlertType } from '~/lib/hooks'
import { type UserOrder, type User, Prefecture, PaymentStatus, UserStatus, AdminRole, type Address, AdminStatus } from '~/types/api'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
  },
  role: {
    type: Number as PropType<AdminRole>,
    default: AdminRole.UNKNOWN
  },
  isAlert: {
    type: Boolean,
    default: false
  },
  alertType: {
    type: String as PropType<AlertType>,
    default: undefined
  },
  alertText: {
    type: String,
    default: ''
  },
  deleteDialog: {
    type: Boolean,
    default: false
  },
  customer: {
    type: Object as PropType<User>,
    default: (): User => ({
      id: '',
      username: '',
      accountId: '',
      lastname: '',
      firstname: '',
      lastnameKana: '',
      firstnameKana: '',
      registered: false,
      status: UserStatus.UNKNOWN,
      email: '',
      phoneNumber: '',
      createdAt: 0,
      updatedAt: 0
    })
  },
  address: {
    type: Object as PropType<Address>,
    default: (): Address => ({
      addressId: '',
      lastname: '',
      firstname: '',
      lastnameKana: '',
      firstnameKana: '',
      postalCode: '',
      prefecture: '',
      prefectureCode: Prefecture.UNKNOWN,
      city: '',
      addressLine1: '',
      addressLine2: '',
      phoneNumber: ''
    })
  },
  orders: {
    type: Array<UserOrder>,
    default: () => []
  },
  orderTotal: {
    type: Number,
    default: 0
  },
  orderAmount: {
    type: Number,
    default: 0
  },
  tableItemsPerPage: {
    type: Number,
    default: 20
  }
})

const emit = defineEmits<{
  (e: 'click:update-page', page: number): void
  (e: 'click:update-items-per-page', page: number): void
  (e: 'click:row', orderId: string): void
  (e: 'update:delete-dialog', v: boolean): void
  (e: 'submit:delete'): void
}>()

const headers: VDataTable['headers'] = [
  {
    title: '注文日時',
    key: 'orderedAt',
    sortable: false
  },
  {
    title: '支払い日時',
    key: 'paidAt',
    sortable: false
  },
  {
    title: '支払い状況',
    key: 'status',
    sortable: false
  },
  {
    title: '支払い合計金額',
    key: 'total',
    sortable: false
  }
]

const activities = [
  {
    eventType: 'notification',
    detail: '注文(#1000)の発想が完了しました。',
    createdAt: '2023/04/12 10:34'
  },
  {
    eventType: 'notification',
    detail: '注文(#1000)の発送済みメールを送りました。',
    createdAt: '2023/04/10 10:34'
  },
  {
    eventType: 'comment',
    username: 'ふるマル管理者',
    detail: '発送準備をコーディネーターに依頼済み',
    createdAt: '2023/04/06 12:00'
  },
  {
    eventType: 'notification',
    detail: '注文(#1000)の支払い完了メールを送りました。',
    createdAt: '2023/04/05 10:34'
  }
]

const deleteDialogValue = computed({
  get: () => props.deleteDialog,
  set: (val: boolean) => emit('update:delete-dialog', val)
})

const isEditable = (): boolean => {
  if (!props.customer || props.customer.status === AdminStatus.DEACTIVATED) {
    return false
  }
  return props.role === AdminRole.ADMINISTRATOR
}

const getName = (): string => {
  if (props.customer?.lastname || props.customer?.firstname) {
    return `${props.customer.lastname} ${props.customer.firstname}`
  }
  return props.customer?.email || ''
}

const getUsername = (): string => {
  return `${props.customer.lastname} ${props.customer.firstname}`
}

const getUsernameKana = (): string => {
  return `${props.customer.lastnameKana} ${props.customer.firstnameKana}`
}

const getPhoneNumber = (): string => {
  return convertI18nToJapanesePhoneNumber(props.customer.phoneNumber)
}

const getAddressArea = (): string => {
  if (!props.address) {
    return ''
  }
  return `${props.address.prefecture} ${props.address.city}`
}

const getCustomerStatus = (): string => {
  switch (props.customer.status) {
    case UserStatus.GUEST:
      return 'ゲスト'
    case UserStatus.PROVISIONAL:
      return '仮登録'
    case UserStatus.VERIFIED:
      return '認証済み'
    case UserStatus.WITH_DRAWAL:
      return '退会済み'
    default:
      return '不明'
  }
}

const getCustomerStatusColor = (): string => {
  switch (props.customer.status) {
    case UserStatus.GUEST:
      return 'secondary'
    case UserStatus.PROVISIONAL:
      return 'warning'
    case UserStatus.VERIFIED:
      return 'primary'
    case UserStatus.WITH_DRAWAL:
      return 'error'
    default:
      return 'unknown'
  }
}

const getPaymentStatus = (status: PaymentStatus): string => {
  switch (status) {
    case PaymentStatus.UNPAID:
      return '未払い'
    case PaymentStatus.AUTHORIZED:
      return 'オーソリ済み'
    case PaymentStatus.PAID:
      return '支払い済み'
    case PaymentStatus.CANCELED:
      return 'キャンセル済み'
    case PaymentStatus.FAILED:
      return '失敗'
    default:
      return '不明'
  }
}

const getPaymentStatusColor = (status:PaymentStatus): string => {
  switch (status) {
    case PaymentStatus.UNPAID:
      return 'secondary'
    case PaymentStatus.AUTHORIZED:
      return 'info'
    case PaymentStatus.PAID:
      return 'primary'
    case PaymentStatus.CANCELED:
      return 'warning'
    case PaymentStatus.FAILED:
      return 'error'
    default:
      return 'unkown'
  }
}

const getOrderedAt = (orderedAt: number): string => {
  if (orderedAt === 0) {
    return '-'
  }
  return unix(orderedAt).format('YYYY/MM/DD HH:mm')
}

const getPaidAt = (paidAt: number): string => {
  if (paidAt === 0) {
    return '-'
  }
  return unix(paidAt).format('YYYY/MM/DD HH:mm')
}

const onClickOpenDeleteDialog = (): void => {
  emit('update:delete-dialog', true)
}

const onClickCloseDeleteDialog = (): void => {
  emit('update:delete-dialog', false)
}

const onClickUpdatePage = (page: number): void => {
  emit('click:update-page', page)
}

const onClickUpdateItemsPerPage = (page: number): void => {
  emit('click:update-items-per-page', page)
}

const onClickRow = (item: UserOrder): void => {
  emit('click:row', item.orderId || '')
}

const onSubmitDelete = (): void => {
  emit('submit:delete')
}
</script>

<template>
  <v-alert v-show="props.isAlert" :type="props.alertType" v-text="props.alertText" />

  <v-dialog v-model="deleteDialogValue" width="500">
    <v-card>
      <v-card-title>
        {{ getName() }}を本当に削除しますか？
      </v-card-title>
      <v-card-actions>
        <v-spacer />
        <v-btn color="error" variant="text" @click="onClickCloseDeleteDialog">
          キャンセル
        </v-btn>
        <v-btn :loading="loading" color="primary" variant="outlined" @click="onSubmitDelete">
          削除
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-row>
    <v-col sm="12" md="12" lg="4" order-lg="2">
      <v-card elevation="0">
        <v-card-title class="d-flex flex-row align-center mx-4 mt-2">
          顧客情報
          <v-spacer />
          <v-menu v-show="isEditable()">
            <template #activator="{ props: item }">
              <v-btn variant="plain" size="small" :icon="mdiDotsVertical" v-bind="item" />
            </template>
            <v-list>
              <v-list-item @click="onClickOpenDeleteDialog">
                削除する
              </v-list-item>
            </v-list>
          </v-menu>
        </v-card-title>
        <v-card-text>
          <v-list>
            <v-list-item class="mb-4">
              <v-list-item-subtitle class="pb-2">
                氏名
              </v-list-item-subtitle>
              <div>{{ getUsername() }}</div>
              <div>{{ getUsernameKana() }}</div>
            </v-list-item>
            <v-list-item class="mb-4">
              <v-list-item-subtitle class="mb-2">
                基本情報
              </v-list-item-subtitle>
              <div>
                ステータス：
                <v-chip size="small" :color="getCustomerStatusColor()">
                  {{ getCustomerStatus() }}
                </v-chip>
              </div>
              <div v-show="customer.username != ''">
                ユーザー名：{{ customer.username }}
              </div>
              <div v-show="customer.accountId != ''">
                アカウントID：{{ customer.accountId }}
              </div>
            </v-list-item>
            <v-list-item class="mb-4">
              <v-list-item-subtitle class="pb-2">
                連絡先情報
              </v-list-item-subtitle>
              <div>メール：{{ props.customer.email }}</div>
              <div>電話番号：{{ getPhoneNumber() }}</div>
            </v-list-item>
            <v-list-item v-show="props.address?.postalCode !== ''">
              <v-list-item-subtitle class="pb-2">
                請求先情報
              </v-list-item-subtitle>
              <div>&#12306; {{ props.address?.postalCode || '' }}</div>
              <div>{{ getAddressArea() }}</div>
              <div>{{ props.address?.addressLine1 || '' }}</div>
              <div>{{ props.address?.addressLine2 || '' }}</div>
            </v-list-item>
          </v-list>
        </v-card-text>
      </v-card>
    </v-col>

    <v-col sm="12" md="12" lg="8" order-lg="1">
      <v-card elevation="0" class="mb-4">
        <v-card-title class="mx-4 mt-2">
          購入情報
        </v-card-title>

        <v-card-text>
          <v-row>
            <v-col>
              <v-card-subtitle class="pb-4">
                支払い金額
              </v-card-subtitle>
              <div class="px-4">
                &yen; {{ props.orderAmount.toLocaleString() }}
              </div>
            </v-col>
            <v-col>
              <v-card-subtitle class="pb-4">
                注文数
              </v-card-subtitle>
              <div class="px-4">
                {{ props.orderTotal.toLocaleString() }} 件
              </div>
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
      <v-card elevation="0" class="mb-4">
        <v-card-text>
          <v-data-table-server
            :headers="headers"
            :loading="loading"
            :items="props.orders"
            :items-per-page="props.tableItemsPerPage"
            :items-length="props.orderTotal"
            no-data-text="注文履歴がありません"
            hover
            @update:page="onClickUpdatePage"
            @update:items-per-page="onClickUpdateItemsPerPage"
            @click:row="(_: any, { item }: any) => onClickRow(item)"
          >
            <template #[`item.status`]="{ item }">
              <v-chip :color="getPaymentStatusColor(item.status)">
                {{ getPaymentStatus(item.status) }}
              </v-chip>
            </template>
            <template #[`item.total`]="{ item }">
              &yen; {{ item.total.toLocaleString() }}
            </template>
            <template #[`item.orderedAt`]="{ item }">
              {{ getOrderedAt(item.orderedAt) }}
            </template>
            <template #[`item.paidAt`]="{ item }">
              {{ getPaidAt(item.paidAt) }}
            </template>
          </v-data-table-server>
        </v-card-text>
      </v-card>

      <div class="pa-4">
        <h4 class="pb-2">
          タイムライン
        </h4>
        <v-divider />

        <v-timeline side="end" density="compact">
          <template v-for="(activity, i) in activities" :key="i">
            <v-timeline-item v-if="activity.eventType === 'notification'" class="mb-4" dot-color="grey" size="small" max-width="75vw">
              <div class="d-flex flex-column flex-lg-row justify-space-between flex-grow-1">
                <div>{{ activity.detail }}</div>
                <div class="flex-shrink-0 text-grey">
                  {{ activity.createdAt }}
                </div>
              </div>
            </v-timeline-item>
            <v-timeline-item v-if="activity.eventType === 'comment'" class="mb-4" dot-color="grey" size="small" max-width="75vw">
              <template #icon>
                <v-avatar image="https://i.pravatar.cc/64" />
              </template>
              <v-card class="elevation-0">
                <v-card-title class="d-lg-flex flex-lg-row align-center">
                  <div class="pr-2">
                    {{ activity.username }}
                  </div>
                  <div class="text-subtitle-2 text-grey">
                    {{ activity.createdAt }}
                  </div>
                </v-card-title>
                <v-card-text>
                  <div>{{ activity.detail }}</div>
                </v-card-text>
              </v-card>
            </v-timeline-item>
          </template>
        </v-timeline>
      </div>
    </v-col>
  </v-row>
</template>
