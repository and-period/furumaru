<script lang="ts" setup>
import dayjs from 'dayjs'
import { DataTableHeader } from 'vuetify'

import { usePagination } from '~/lib/hooks'
import { useOrderStore } from '~/store'
import {
  DeliveryType,
  FulfillmentStatus,
  OrderRefundType,
  OrderResponse,
  PaymentMethodType,
  PaymentStatus,
  ShippingCarrier,
  ShippingSize,
} from '~/types/api'
import { Order, OrderItems } from '~/types/props/order'

const route = useRoute()
const orderStore = useOrderStore()
const id = route.params.id

const selector = ref<string>('shippingInformation')

const {
  updateCurrentPage,
  itemsPerPage,
  handleUpdateItemsPerPage,
  options,
  offset,
} = usePagination()

const items: Order[] = [
  { name: '支払い情報', value: 'shippingInformation' },
  { name: '配送情報', value: 'orderInformation' },
]

const formData = reactive<OrderResponse>({
  id: '',
  scheduleId: '',
  promotionId: '',
  userId: '',
  userName: '',
  payment: {
    transactionId: '',
    methodId: '',
    methodType: 0,
    status: 0,
    subtotal: 0,
    discount: 0,
    shippingFee: 0,
    tax: 0,
    total: 0,
    addressId: '',
    lastname: '',
    firstname: '',
    postalCode: '',
    prefecture: '',
    city: '',
    addressLine1: '',
    addressLine2: '',
    phoneNumber: '',
  },
  fulfillment: {
    trackingNumber: '',
    status: 0,
    shippingCarrier: 0,
    shippingMethod: 0,
    boxSize: 0,
    addressId: '',
    lastname: '',
    firstname: '',
    postalCode: '',
    prefecture: '',
    city: '',
    addressLine1: '',
    addressLine2: '',
    phoneNumber: '',
  },
  refund: {
    canceled: false,
    type: 0,
    reason: '',
    total: 0,
  },
  items: [
    {
      productId: '',
      name: '',
      price: 0,
      quantity: 0,
      weight: 0,
      media: [
        {
          url: '',
          isThumbnail: false,
        },
      ],
    },
  ],
  orderedAt: -1,
  paidAt: -1,
  deliveredAt: -1,
  canceledAt: -1,
  createdAt: -1,
  updatedAt: -1,
})

const fetchState = useAsyncData(async () => {
  const res = await orderStore.getOrder(id)
  formData.id = res.id
  formData.userName = res.userName
  formData.payment = res.payment
  formData.fulfillment = res.fulfillment
  formData.refund = res.refund
  formData.items = res.items
  formData.orderedAt = res.orderedAt
  formData.paidAt = res.paidAt
  formData.deliveredAt = res.deliveredAt
  formData.canceledAt = res.canceledAt
})

const headers: DataTableHeader[] = [
  {
    text: 'サムネイル',
    value: 'media',
  },
  {
    text: '商品名',
    value: 'name',
  },
  {
    text: '購入価格',
    value: 'price',
  },
  {
    text: '購入数量',
    value: 'quantity',
  },
  {
    text: '重量',
    value: 'weight',
  },
]

const getDay = (unixTime: number): string => {
  return dayjs.unix(unixTime).format('YYYY/MM/DD HH:mm')
}

const getMethodType = (status: PaymentMethodType): string => {
  switch (status) {
    case PaymentMethodType.CASH:
      return '代引き支払い'
    case PaymentMethodType.CARD:
      return 'クレジットカード払い'
    default:
      return '不明'
  }
}

const getPaymentStatus = (status: PaymentStatus): string => {
  switch (status) {
    case PaymentStatus.UNKNOWN:
      return '不明'
    case PaymentStatus.UNPAID:
      return '未払い'
    case PaymentStatus.PENDING:
      return '保留中'
    case PaymentStatus.AUTHORIZED:
      return 'オーソリ済み'
    case PaymentStatus.PAID:
      return '支払い済み'
    case PaymentStatus.REFUNDED:
      return '返金済み'
    case PaymentStatus.EXPIRED:
      return '期限切れ'
    default:
      return '不明'
  }
}

const getRefundType = (status: OrderRefundType): string => {
  switch (status) {
    default:
      return '不明'
  }
}

const getPaymentStatusColor = (status: PaymentStatus): string => {
  switch (status) {
    case PaymentStatus.UNPAID:
      return 'secondary'
    case PaymentStatus.PENDING:
      return 'secondary'
    case PaymentStatus.AUTHORIZED:
      return 'info'
    case PaymentStatus.PAID:
      return 'primary'
    case PaymentStatus.REFUNDED:
      return 'primary'
    case PaymentStatus.EXPIRED:
      return 'error'
    default:
      return 'unkown'
  }
}

const getRefundStatus = (status: boolean): string => {
  if (status) {
    return 'キャンセル'
  } else {
    return '注文受付済み'
  }
}

const getRefundStatusColor = (status: boolean): string => {
  if (status) {
    return 'error'
  } else {
    return 'primary'
  }
}

const getFulfillmentStatus = (status: FulfillmentStatus): string => {
  switch (status) {
    case FulfillmentStatus.UNFULFILLED:
      return '未配送'
    case FulfillmentStatus.FULFILLED:
      return '配送済み'
    default:
      return '不明'
  }
}

const getFulfillmentStatusColor = (status: FulfillmentStatus): string => {
  switch (status) {
    case FulfillmentStatus.UNFULFILLED:
      return 'error'
    case FulfillmentStatus.FULFILLED:
      return 'primary'
    default:
      return 'unkown'
  }
}

const convertPhone = (phoneNumber: string): string => {
  return phoneNumber.replace('+81', '0')
}

// isThumnailがtrueのものを引っ掛けて商品でサムネイルに設定されているURLを探す
const getThumnail = (medias: OrderItems[]): string => {
  const orderItem: OrderItems[] = medias.filter((item) => item.isThumbnail)
  return orderItem[0].url
}

const getShippingCarrier = (carrier: ShippingCarrier): string => {
  switch (carrier) {
    case ShippingCarrier.YAMATO:
      return 'ヤマト運輸'
    case ShippingCarrier.SAGAWA:
      return '佐川急便'
    default:
      return '不明'
  }
}

const getShippingMethod = (method: DeliveryType): string => {
  switch (method) {
    case DeliveryType.NORMAL:
      return '通常便'
    case DeliveryType.REFRIGERATED:
      return '冷蔵便'
    case DeliveryType.FROZEN:
      return '冷凍便'
    default:
      return '不明'
  }
}

const getBoxSize = (size: ShippingSize): string => {
  switch (size) {
    case ShippingSize.SIZE60:
      return '60'
    case ShippingSize.SIZE80:
      return '80'
    case ShippingSize.SIZE100:
      return '100'
    default:
      return '不明'
  }
}
</script>

<template>
  <div>
    <v-tabs v-model="selector" grow color="dark">
      <v-tabs-slider color="accent"></v-tabs-slider>
      <v-tab
        v-for="item in items"
        :key="item.value"
        :href="`#tab-${item.value}`"
      >
        {{ item.name }}
      </v-tab>
    </v-tabs>
    <v-tabs-items v-model="selector">
      <v-tab-item value="tab-shippingInformation">
        <v-card elevation="0">
          <v-card-text>
            <v-text-field
              name="userName"
              label="注文者名"
              :value="formData.userName"
              readonly
            ></v-text-field>
            <v-text-field
              name="paymentMethodType"
              label="決済手段"
              :value="getMethodType(formData.payment.methodType)"
              readonly
            ></v-text-field>
            <v-container>
              <p class="text-h6">購入情報</p>
              <v-row class="mt-4">
                <span class="mx-4">支払い状況:</span>
                <v-chip
                  small
                  :color="getPaymentStatusColor(formData.payment.status)"
                >
                  {{ getPaymentStatus(formData.payment.status) }}
                </v-chip>
              </v-row>
            </v-container>
            <v-text-field
              v-if="getPaymentStatus(formData.payment.status) == '支払い済み'"
              class="mt-4"
              name="deliveredAt"
              label="支払日時"
              :value="getDay(formData.paidAt)"
              readonly
            ></v-text-field>
            <v-text-field
              class="mt-4"
              name="total"
              label="支払い合計金額"
              :value="formData.payment.total"
              readonly
            >
              <template #append>円</template>
            </v-text-field>
            <div class="d-flex align-center">
              <v-text-field
                class="mr-4"
                name="subTotal"
                label="購入金額"
                :value="formData.payment.subtotal"
                readonly
              >
                <template #append>円</template>
              </v-text-field>
              <v-text-field
                name="discount"
                label="割引金額"
                :value="formData.payment.discount"
                readonly
              >
                <template #append>円</template>
              </v-text-field>
            </div>
            <div class="d-flex align-center mt-4">
              <v-text-field
                class="mr-4"
                name="shippingFee"
                label="配送料金"
                :value="formData.payment.shippingFee"
                readonly
              >
                <template #append>円</template>
              </v-text-field>
              <v-text-field
                name="tax"
                label="消費税"
                :value="formData.payment.tax"
                readonly
              >
                <template #append>円</template>
              </v-text-field>
            </div>
            <p class="text-h6">請求先情報</p>
            <div class="d-flex align-center">
              <v-text-field
                class="mr-4"
                name="lastname"
                label="姓"
                :value="formData.payment.lastname"
                readonly
              ></v-text-field>
              <v-text-field
                name="firstname"
                label="名"
                :value="formData.payment.firstname"
                readonly
              ></v-text-field>
            </div>
            <v-text-field
              name="phoneNumber"
              label="電話番号"
              :value="convertPhone(formData.payment.phoneNumber)"
              readonly
            ></v-text-field>
            <v-text-field
              name="postalCode"
              label="郵便番号"
              :value="formData.payment.postalCode"
              readonly
            ></v-text-field>
            <div class="d-flex align-center">
              <v-text-field
                class="mr-4"
                name="prefecture"
                label="都道府県"
                :value="formData.payment.prefecture"
                readonly
              ></v-text-field>
              <v-text-field
                name="city"
                label="市区町村"
                :value="formData.payment.city"
                readonly
              ></v-text-field>
            </div>
            <v-text-field
              name="addressLine1"
              label="町名・番地"
              :value="formData.payment.addressLine1"
            ></v-text-field>
            <v-text-field
              name="addressLine2"
              label="ビル名・号室など"
              :value="formData.payment.addressLine2"
            ></v-text-field>
            <p class="text-h6">キャンセル情報</p>
            <v-row class="mt-4">
              <span class="mx-4">注文キャンセル状況:</span>
              <v-chip
                small
                :color="getRefundStatusColor(formData.refund.canceled)"
              >
                {{ getRefundStatus(formData.refund.canceled) }}
              </v-chip>
            </v-row>
            <v-container
              v-if="getRefundStatus(formData.refund.canceled) == 'キャンセル'"
            >
              <v-text-field
                class="mt-8"
                name="canceledAt"
                label="注文キャンセル日時"
                :value="getDay(formData.canceledAt)"
                readonly
              ></v-text-field>
              <v-text-field
                name="type"
                label="注文キャンセル理由"
                :value="getRefundType(formData.refund.type)"
                readonly
              ></v-text-field>
              <v-textarea
                name="reason"
                label="注文キャンセル理由詳細"
                :value="formData.refund.reason"
                readonly
              ></v-textarea>
              <v-text-field
                name="refundTotal"
                label="返済金額"
                :value="formData.refund.total"
                readonly
              ></v-text-field>
            </v-container>
          </v-card-text>
        </v-card>
      </v-tab-item>
      <v-tab-item value="tab-orderInformation">
        <v-card-text>
          <v-card elevation="0">
            <p class="text-h6">注文情報</p>
            <v-text-field
              name="id"
              label="注文ID"
              :value="formData.id"
              readonly
            ></v-text-field>
            <v-text-field
              name="orderedAt"
              label="注文日時"
              :value="getDay(formData.orderedAt)"
              readonly
            ></v-text-field>
            <v-row class="my-4">
              <span class="mx-4">配送状況:</span>
              <v-chip
                small
                :color="getFulfillmentStatusColor(formData.fulfillment.status)"
              >
                {{ getFulfillmentStatus(formData.fulfillment.status) }}
              </v-chip>
            </v-row>
            <div class="d-flex align-center">
              <v-text-field
                class="mr-4"
                name="shippingCarrier"
                label="配送会社"
                :value="
                  getShippingCarrier(formData.fulfillment.shippingCarrier)
                "
                readonly
              ></v-text-field>
              <v-text-field
                class="mr-4"
                name="shippingmethod"
                label="配送方法"
                :value="getShippingMethod(formData.fulfillment.shippingMethod)"
                readonly
              ></v-text-field>
              <v-text-field
                name="boxSize"
                label="配送時の箱の大きさ"
                :value="getBoxSize(formData.fulfillment.boxSize)"
                readonly
              ></v-text-field>
            </div>
            <v-text-field
              v-if="
                getFulfillmentStatus(formData.fulfillment.status) == '配送済み'
              "
              name="deliveredAt"
              label="配送日時"
              :value="getDay(formData.deliveredAt)"
              readonly
            ></v-text-field>
            <v-data-table
              :headers="headers"
              :items="formData.items"
              :footer-props="options"
              no-data-text="表示する注文がありません"
            >
              <template #[`item.media`]="{ item }">
                <v-avatar>
                  <img :src="getThumnail(item.media)" />
                </v-avatar>
              </template>
            </v-data-table>
            <v-text-field
              name="trackingNumber"
              label="伝票番号"
              :value="formData.fulfillment.trackingNumber"
              readonly
            ></v-text-field>
            <div>
              <v-text-field
                class="mr-4"
                name="lastname"
                label="姓"
                :value="formData.fulfillment.lastname"
                readonly
              ></v-text-field>
              <v-text-field
                name="firstname"
                label="名"
                :value="formData.fulfillment.firstname"
                readonly
              ></v-text-field>
            </div>
            <v-text-field
              name="phoneNumber"
              label="電話番号"
              :value="convertPhone(formData.fulfillment.phoneNumber)"
              readonly
            ></v-text-field>
            <v-text-field
              name="postalCode"
              label="郵便番号"
              :value="formData.fulfillment.postalCode"
              readonly
            ></v-text-field>
            <div class="d-flex align-center">
              <v-text-field
                class="mr-4"
                name="prefecture"
                label="都道府県"
                :value="formData.fulfillment.prefecture"
                readonly
              ></v-text-field>
              <v-text-field
                name="city"
                label="市区町村"
                :value="formData.fulfillment.city"
                readonly
              ></v-text-field>
            </div>
            <v-text-field
              name="addressLine1"
              label="町名・番地"
              :value="formData.fulfillment.addressLine1"
            ></v-text-field>
            <v-text-field
              name="addressLine2"
              label="ビル名・号室など"
              :value="formData.fulfillment.addressLine2"
            ></v-text-field>
          </v-card>
        </v-card-text>
      </v-tab-item>
    </v-tabs-items>
  </div>
</template>
