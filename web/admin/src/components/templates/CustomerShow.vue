<script lang="ts" setup>
import { mdiAccountCircle, mdiDelete, mdiDotsVertical, mdiHistory, mdiShopping, mdiTimeline } from '@mdi/js'
import { unix } from 'dayjs'
import type { VDataTable } from 'vuetify/lib/components/index.mjs'

import { convertI18nToJapanesePhoneNumber } from '~/lib/formatter'
import type { AlertType } from '~/lib/hooks'
import { Prefecture } from '~/types'
import { PaymentStatus, UserStatus, AdminType, AdminStatus } from '~/types/api/v1'
import type { UserOrder, User, Address } from '~/types/api/v1'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
  },
  adminType: {
    type: Number as PropType<AdminType>,
    default: AdminType.AdminTypeUnknown,
  },
  isAlert: {
    type: Boolean,
    default: false,
  },
  alertType: {
    type: String as PropType<AlertType>,
    default: undefined,
  },
  alertText: {
    type: String,
    default: '',
  },
  deleteDialog: {
    type: Boolean,
    default: false,
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
      status: UserStatus.UserStatusUnknown,
      email: '',
      phoneNumber: '',
      thumbnailUrl: '',
      createdAt: 0,
      updatedAt: 0,
    }),
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
      phoneNumber: '',
    }),
  },
  orders: {
    type: Array<UserOrder>,
    default: () => [],
  },
  totalOrderCount: {
    type: Number,
    default: 0,
  },
  totalPaymentCount: {
    type: Number,
    default: 0,
  },
  totalProductAmount: {
    type: Number,
    default: 0,
  },
  totalPaymentAmount: {
    type: Number,
    default: 0,
  },
  tableItemsPerPage: {
    type: Number,
    default: 20,
  },
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
    sortable: false,
  },
  {
    title: '支払い日時',
    key: 'paidAt',
    sortable: false,
  },
  {
    title: '支払い状況',
    key: 'status',
    sortable: false,
  },
  {
    title: '支払い合計金額',
    key: 'total',
    sortable: false,
  },
]

// TODO: Replace with API data
const activities = computed(() => {
  // Placeholder data - should be replaced with real API call
  return [
    {
      eventType: 'notification',
      detail: '注文(#1000)の発想が完了しました。',
      createdAt: '2023/04/12 10:34',
    },
    {
      eventType: 'notification',
      detail: '注文(#1000)の発送済みメールを送りました。',
      createdAt: '2023/04/10 10:34',
    },
    {
      eventType: 'comment',
      username: 'ふるマル管理者',
      detail: '発送準備をコーディネーターに依頼済み',
      createdAt: '2023/04/06 12:00',
    },
    {
      eventType: 'notification',
      detail: '注文(#1000)の支払い完了メールを送りました。',
      createdAt: '2023/04/05 10:34',
    },
  ]
})

const deleteDialogValue = computed({
  get: () => props.deleteDialog,
  set: (val: boolean) => emit('update:delete-dialog', val),
})

const isEditable = (): boolean => {
  if (!props.customer || props.customer.status === AdminStatus.AdminStatusDeactivated) {
    return false
  }
  return props.adminType === AdminType.AdminTypeAdministrator
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
    case UserStatus.UserStatusGuest:
      return 'ゲスト'
    case UserStatus.UserStatusProvisional:
      return '仮登録'
    case UserStatus.UserStatusVerified:
      return '認証済み'
    case UserStatus.UserStatusDeactivated:
      return '退会済み'
    default:
      return '不明'
  }
}

const getCustomerStatusColor = (): string => {
  switch (props.customer.status) {
    case UserStatus.UserStatusGuest:
      return 'secondary'
    case UserStatus.UserStatusProvisional:
      return 'warning'
    case UserStatus.UserStatusVerified:
      return 'primary'
    case UserStatus.UserStatusDeactivated:
      return 'error'
    default:
      return 'unknown'
  }
}

const getPaymentStatus = (status: PaymentStatus): string => {
  switch (status) {
    case PaymentStatus.PaymentStatusUnpaid:
      return '未払い'
    case PaymentStatus.PaymentStatusAuthorized:
      return 'オーソリ済み'
    case PaymentStatus.PaymentStatusPaid:
      return '支払い済み'
    case PaymentStatus.PaymentStatusCanceled:
      return 'キャンセル済み'
    case PaymentStatus.PaymentStatusFailed:
      return '失敗'
    default:
      return '不明'
  }
}

const getPaymentStatusColor = (status: PaymentStatus): string => {
  switch (status) {
    case PaymentStatus.PaymentStatusUnpaid:
      return 'secondary'
    case PaymentStatus.PaymentStatusAuthorized:
      return 'info'
    case PaymentStatus.PaymentStatusPaid:
      return 'primary'
    case PaymentStatus.PaymentStatusCanceled:
      return 'warning'
    case PaymentStatus.PaymentStatusFailed:
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
  <v-alert
    v-show="props.isAlert"
    :type="props.alertType"
    v-text="props.alertText"
  />

  <v-dialog
    v-model="deleteDialogValue"
    width="500"
  >
    <v-card>
      <v-card-title>
        {{ getName() }}を本当に削除しますか？
      </v-card-title>
      <v-card-actions>
        <v-spacer />
        <v-btn
          color="error"
          variant="text"
          @click="onClickCloseDeleteDialog"
        >
          キャンセル
        </v-btn>
        <v-btn
          :loading="loading"
          color="primary"
          variant="outlined"
          @click="onSubmitDelete"
        >
          削除
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-row>
    <v-col
      sm="12"
      md="12"
      lg="4"
      order-lg="2"
    >
      <v-card
        elevation="0"
        class="customer-info-card"
      >
        <v-card-title class="d-flex flex-row align-center pa-6">
          <v-icon
            :icon="mdiAccountCircle"
            size="24"
            class="mr-2 text-primary"
          />
          顧客情報
          <v-spacer />
          <v-menu v-show="isEditable()">
            <template #activator="{ props: item }">
              <v-btn
                variant="text"
                size="small"
                :icon="mdiDotsVertical"
                class="action-menu-btn"
                v-bind="item"
              />
            </template>
            <v-list>
              <v-list-item
                class="text-error"
                @click="onClickOpenDeleteDialog"
              >
                <template #prepend>
                  <v-icon
                    :icon="mdiDelete"
                    size="20"
                  />
                </template>
                削除する
              </v-list-item>
            </v-list>
          </v-menu>
        </v-card-title>
        <v-card-text class="px-6">
          <v-list>
            <v-list-item class="mb-6 px-0">
              <v-list-item-subtitle class="text-subtitle-2 font-weight-medium text-grey-darken-1 pb-3">
                氏名
              </v-list-item-subtitle>
              <div class="text-h6 font-weight-medium mb-1">
                {{ getUsername() }}
              </div>
              <div class="text-body-2 text-grey-darken-2">
                {{ getUsernameKana() }}
              </div>
            </v-list-item>
            <v-list-item class="mb-6 px-0">
              <v-list-item-subtitle class="text-subtitle-2 font-weight-medium text-grey-darken-1 pb-3">
                基本情報
              </v-list-item-subtitle>
              <div class="d-flex align-center mb-2">
                <span class="text-body-2 text-grey-darken-2 mr-2">ステータス：</span>
                <v-chip
                  size="small"
                  :color="getCustomerStatusColor()"
                  variant="flat"
                >
                  {{ getCustomerStatus() }}
                </v-chip>
              </div>
              <div
                v-show="customer.username != ''"
                class="text-body-2 mb-1"
              >
                <span class="text-grey-darken-2">ユーザー名：</span>
                <span>{{ customer.username }}</span>
              </div>
              <div
                v-show="customer.accountId != ''"
                class="text-body-2"
              >
                <span class="text-grey-darken-2">アカウントID：</span>
                <span>{{ customer.accountId }}</span>
              </div>
            </v-list-item>
            <v-list-item class="mb-6 px-0">
              <v-list-item-subtitle class="text-subtitle-2 font-weight-medium text-grey-darken-1 pb-3">
                連絡先情報
              </v-list-item-subtitle>
              <div class="text-body-2 mb-1">
                <span class="text-grey-darken-2">メール：</span>
                <span>{{ props.customer.email }}</span>
              </div>
              <div class="text-body-2">
                <span class="text-grey-darken-2">電話番号：</span>
                <span>{{ getPhoneNumber() }}</span>
              </div>
            </v-list-item>
            <v-list-item
              v-show="props.address?.postalCode !== ''"
              class="px-0"
            >
              <v-list-item-subtitle class="text-subtitle-2 font-weight-medium text-grey-darken-1 pb-3">
                請求先情報
              </v-list-item-subtitle>
              <div class="text-body-2 mb-1">
                <v-icon
                  icon="mdi-map-marker"
                  size="16"
                  class="mr-1"
                />
                〒{{ props.address?.postalCode || '' }}
              </div>
              <div class="text-body-2 mb-1">
                {{ getAddressArea() }}
              </div>
              <div class="text-body-2 mb-1">
                {{ props.address?.addressLine1 || '' }}
              </div>
              <div
                v-show="props.address?.addressLine2"
                class="text-body-2"
              >
                {{ props.address?.addressLine2 || '' }}
              </div>
            </v-list-item>
          </v-list>
        </v-card-text>
      </v-card>
    </v-col>

    <v-col
      sm="12"
      md="12"
      lg="8"
      order-lg="1"
    >
      <v-card
        elevation="0"
        class="mb-4 purchase-stats-card"
      >
        <v-card-title class="pa-6">
          <v-icon
            :icon="mdiShopping"
            size="24"
            class="mr-2 text-primary"
          />
          購入情報
        </v-card-title>

        <v-card-text class="px-6">
          <v-row>
            <v-col
              cols="12"
              sm="6"
            >
              <div class="stats-item">
                <v-card-subtitle class="text-subtitle-2 font-weight-medium text-grey-darken-1 pb-2">
                  購入商品金額（※送料等は除く）
                </v-card-subtitle>
                <div class="text-h5 font-weight-bold text-primary">
                  ¥{{ props.totalProductAmount.toLocaleString() }}
                </div>
              </div>
            </v-col>
            <v-col
              cols="12"
              sm="6"
            >
              <div class="stats-item">
                <v-card-subtitle class="text-subtitle-2 font-weight-medium text-grey-darken-1 pb-2">
                  注文数
                </v-card-subtitle>
                <div class="text-h5 font-weight-bold text-primary">
                  {{ props.totalPaymentCount.toLocaleString() }} 件
                </div>
              </div>
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
      <v-card
        elevation="0"
        class="mb-4 order-history-card"
      >
        <v-card-title class="pa-6">
          <v-icon
            :icon="mdiHistory"
            size="24"
            class="mr-2 text-primary"
          />
          注文履歴
        </v-card-title>
        <v-card-text>
          <v-data-table-server
            :headers="headers"
            :loading="loading"
            :items="props.orders"
            :items-per-page="props.tableItemsPerPage"
            :items-length="props.totalOrderCount"
            no-data-text="注文履歴がありません"
            hover
            class="order-history-table"
            @update:page="onClickUpdatePage"
            @update:items-per-page="onClickUpdateItemsPerPage"
            @click:row="(_: any, { item }: any) => onClickRow(item)"
          >
            <template #[`item.status`]="{ item }">
              <v-chip
                :color="getPaymentStatusColor(item.status)"
                size="small"
                variant="flat"
              >
                {{ getPaymentStatus(item.status) }}
              </v-chip>
            </template>
            <template #[`item.orderedAt`]="{ item }">
              {{ getOrderedAt(item.orderedAt) }}
            </template>
            <template #[`item.paidAt`]="{ item }">
              {{ getPaidAt(item.paidAt) }}
            </template>
            <template #[`item.total`]="{ item }">
              <span class="font-weight-medium">
                ¥{{ item.total.toLocaleString() }}
              </span>
            </template>
          </v-data-table-server>
        </v-card-text>
      </v-card>

      <v-card
        elevation="0"
        class="timeline-card"
      >
        <v-card-title class="pa-6">
          <v-icon
            :icon="mdiTimeline"
            size="24"
            class="mr-2 text-primary"
          />
          タイムライン
        </v-card-title>
        <v-divider />
        <v-card-text class="pa-6">
          <v-timeline
            side="end"
            density="compact"
            class="customer-timeline"
          >
            <template
              v-for="(activity, i) in activities"
              :key="i"
            >
              <v-timeline-item
                v-if="activity.eventType === 'notification'"
                class="mb-4"
                dot-color="info"
                size="small"
                max-width="75vw"
              >
                <template #icon>
                  <v-icon
                    icon="mdi-bell"
                    size="16"
                  />
                </template>
                <div class="timeline-content">
                  <div class="d-flex flex-column flex-lg-row justify-space-between flex-grow-1">
                    <div class="text-body-2">
                      {{ activity.detail }}
                    </div>
                    <div class="flex-shrink-0 text-caption text-grey-darken-1 mt-1 mt-lg-0">
                      {{ activity.createdAt }}
                    </div>
                  </div>
                </div>
              </v-timeline-item>
              <v-timeline-item
                v-if="activity.eventType === 'comment'"
                class="mb-4"
                dot-color="primary"
                size="small"
                max-width="75vw"
              >
                <template #icon>
                  <v-avatar
                    image="https://i.pravatar.cc/64"
                    size="32"
                  />
                </template>
                <v-card
                  class="elevation-0 timeline-comment-card"
                  variant="outlined"
                >
                  <v-card-title class="d-lg-flex flex-lg-row align-center pa-4">
                    <div class="text-subtitle-2 font-weight-medium pr-2">
                      {{ activity.username }}
                    </div>
                    <div class="text-caption text-grey-darken-1">
                      {{ activity.createdAt }}
                    </div>
                  </v-card-title>
                  <v-card-text class="pa-4 pt-0">
                    <div class="text-body-2">
                      {{ activity.detail }}
                    </div>
                  </v-card-text>
                </v-card>
              </v-timeline-item>
            </template>
          </v-timeline>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>

<style scoped>
.customer-info-card,
.purchase-stats-card,
.order-history-card,
.timeline-card {
  border-radius: 12px;
  border: 1px solid rgb(0 0 0 / 5%);
}

.action-menu-btn {
  opacity: 0.7;
  transition: opacity 0.2s;
}

.action-menu-btn:hover {
  opacity: 1;
}

.stats-item {
  padding: 16px;
  background: rgb(33 150 243 / 4%);
  border-radius: 8px;
  border-left: 4px solid rgb(33 150 243);
}

.order-history-table {
  border-radius: 8px;
}

.customer-timeline {
  margin-top: 8px;
}

.timeline-content {
  background: rgb(255 255 255 / 80%);
  border-radius: 6px;
  padding: 8px 12px;
}

.timeline-comment-card {
  background: rgb(243 247 251 / 70%);
  border-radius: 8px;
  border: 1px solid rgb(33 150 243 / 20%);
}

@media (width <= 960px) {
  .stats-item {
    margin-bottom: 16px;
  }

  .timeline-content {
    font-size: 14px;
  }
}
</style>
