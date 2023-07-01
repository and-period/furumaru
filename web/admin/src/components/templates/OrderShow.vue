<script lang="ts" setup>
import { unix } from 'dayjs'
import { prefecturesList } from '~/constants'
import { convertI18nToJapanesePhoneNumber } from '~/lib/formatter'
import { AlertType } from '~/lib/hooks'
import { DeliveryType, FulfillmentStatus, OrderRefundType, OrderResponse, PaymentMethodType, PaymentStatus, Prefecture, ShippingCarrier, ShippingSize } from '~/types/api'
import { Order, OrderItems } from '~/types/props'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false
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
  order: {
    type: Object as PropType<OrderResponse>,
    default: () => ({
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
        prefecture: Prefecture.UNKNOWN,
        city: '',
        addressLine1: '',
        addressLine2: '',
        phoneNumber: ''
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
        prefecture: Prefecture.UNKNOWN,
        city: '',
        addressLine1: '',
        addressLine2: '',
        phoneNumber: ''
      },
      refund: {
        canceled: false,
        type: 0,
        reason: '',
        total: 0
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
              isThumbnail: false
            }
          ]
        }
      ],
      orderedAt: -1,
      paidAt: -1,
      deliveredAt: -1,
      canceledAt: -1,
      createdAt: -1,
      updatedAt: -1
    })
  }
})

const items: Order[] = [
  { name: '支払い情報', value: 'shippingInformation' },
  { name: '配送情報', value: 'orderInformation' }
]

const selector = ref<string>('shippingInformation')

const orderValue = computed((): OrderResponse => {
  return props.order
})
const paymentPhoneNumber = computed((): string => {
  return convertI18nToJapanesePhoneNumber(props.order.payment.phoneNumber)
})
const paymentMethodType = computed((): string => {
  return getPaymentMethodType(props.order.payment.methodType)
})
const refundType = computed((): string => {
  return getRefundType(props.order.refund.type)
})
const fulfillmentPhoneNumber = computed((): string => {
  return convertI18nToJapanesePhoneNumber(props.order.fulfillment.phoneNumber)
})
const fulfillmentShippingCarrier = computed((): string => {
  return getShippingCarrier(props.order.fulfillment.shippingCarrier)
})
const fulfillmentShippingMethod = computed((): string => {
  return getShippingMethod(props.order.fulfillment.shippingMethod)
})
const fulfillmentShippingSize = computed((): string => {
  return getBoxSize(props.order.fulfillment.boxSize)
})
const paidAt = computed((): string => {
  return getDay(props.order.paidAt)
})
const canceledAt = computed((): string => {
  return getDay(props.order.canceledAt)
})
const orderedAt = computed((): string => {
  return getDay(props.order.orderedAt)
})
const deliveredAt = computed((): string => {
  return getDay(props.order.deliveredAt)
})

const getDay = (unixTime: number): string => {
  return unix(unixTime).format('YYYY/MM/DD HH:mm')
}

const getPaymentMethodType = (status: PaymentMethodType): string => {
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

// isThumnailがtrueのものを引っ掛けて商品でサムネイルに設定されているURLを探す
const getThumbnail = (medias: OrderItems[]): string => {
  const orderItem: OrderItems[] = medias.filter(item => item.isThumbnail)
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
  <v-alert v-show="props.isAlert" :type="props.alertType" v-text="props.alertText" />

  <v-card>
    <v-tabs v-model="selector" grow color="dark">
      <v-tab v-for="item in items" :key="item.value" :value="item.value">
        {{ item.name }}
      </v-tab>
    </v-tabs>

    <v-window v-model="selector">
      <v-window-item value="shippingInformation">
        <v-card elevation="0">
          <v-card-text>
            <v-text-field
              v-model="orderValue.userName"
              name="userName"
              label="注文者名"
              readonly
            />
            <v-text-field
              v-model="paymentMethodType"
              name="paymentMethodType"
              label="決済手段"
              readonly
            />
            <v-container>
              <p class="text-h6">
                購入情報
              </p>
              <v-row class="mt-4">
                <span class="mx-4">支払い状況:</span>
                <v-chip
                  size="small"
                  :color="getPaymentStatusColor(props.order.payment.status)"
                >
                  {{ getPaymentStatus(props.order.payment.status) }}
                </v-chip>
              </v-row>
            </v-container>
            <v-text-field
              v-if="getPaymentStatus(props.order.payment.status) == '支払い済み'"
              v-model="paidAt"
              class="mt-4"
              name="deliveredAt"
              label="支払日時"
              readonly
            />
            <v-text-field
              v-value="orderValue.payment.total"
              class="mt-4"
              name="total"
              label="支払い合計金額"
              readonly
            >
              <template #append>
                円
              </template>
            </v-text-field>
            <div class="d-flex align-center">
              <v-text-field
                v-model="orderValue.payment.subtotal"
                class="mr-4"
                name="subTotal"
                label="購入金額"
                readonly
              >
                <template #append>
                  円
                </template>
              </v-text-field>
              <v-text-field
                v-model="orderValue.payment.discount"
                name="discount"
                label="割引金額"
                readonly
              >
                <template #append>
                  円
                </template>
              </v-text-field>
            </div>
            <div class="d-flex align-center mt-4">
              <v-text-field
                v-model="orderValue.payment.shippingFee"
                class="mr-4"
                name="shippingFee"
                label="配送料金"
                readonly
              >
                <template #append>
                  円
                </template>
              </v-text-field>
              <v-text-field
                v-model="orderValue.payment.tax"
                name="tax"
                label="消費税"
                readonly
              >
                <template #append>
                  円
                </template>
              </v-text-field>
            </div>
            <p class="text-h6">
              請求先情報
            </p>
            <div class="d-flex align-center">
              <v-text-field
                v-model="orderValue.payment.lastname"
                class="mr-4"
                name="lastname"
                label="姓"
                readonly
              />
              <v-text-field
                v-model="orderValue.payment.firstname"
                name="firstname"
                label="名"
                readonly
              />
            </div>
            <v-text-field
              v-model="paymentPhoneNumber"
              name="phoneNumber"
              label="電話番号"
              readonly
            />
            <v-text-field
              v-model="orderValue.payment.postalCode"
              name="postalCode"
              label="郵便番号"
              readonly
            />
            <div class="d-flex align-center">
              <v-text-field
                v-model="orderValue.payment.prefecture"
                :items="prefecturesList"
                item-title="text"
                item-value="value"
                class="mr-4"
                name="prefecture"
                label="都道府県"
                readonly
              />
              <v-text-field
                v-model="orderValue.payment.city"
                name="city"
                label="市区町村"
                readonly
              />
            </div>
            <v-text-field
              v-model="orderValue.payment.addressLine1"
              name="addressLine1"
              label="町名・番地"
            />
            <v-text-field
              v-model="orderValue.payment.addressLine2"
              name="addressLine2"
              label="ビル名・号室など"
            />
            <p class="text-h6">
              キャンセル情報
            </p>
            <v-row class="mt-4">
              <span class="mx-4">注文キャンセル状況:</span>
              <v-chip
                size="small"
                :color="getRefundStatusColor(props.order.refund.canceled)"
              >
                {{ getRefundStatus(props.order.refund.canceled) }}
              </v-chip>
            </v-row>
            <v-container
              v-if="getRefundStatus(props.order.refund.canceled) == 'キャンセル'"
            >
              <v-text-field
                v-model="canceledAt"
                class="mt-8"
                name="canceledAt"
                label="注文キャンセル日時"
                readonly
              />
              <v-text-field
                v-model="refundType"
                name="type"
                label="注文キャンセル理由"
                readonly
              />
              <v-textarea
                v-model="orderValue.refund.reason"
                name="reason"
                label="注文キャンセル理由詳細"
                readonly
              />
              <v-text-field
                v-model="orderValue.refund.total"
                name="refundTotal"
                label="返済金額"
                readonly
              />
            </v-container>
          </v-card-text>
        </v-card>
      </v-window-item>

      <v-window-item value="orderInformation">
        <v-card elevation="0">
          <v-card-text>
            <v-text-field
              v-model="orderValue.id"
              name="id"
              label="注文ID"
              readonly
            />
            <v-text-field
              v-model="orderedAt"
              name="orderedAt"
              label="注文日時"
              readonly
            />
            <v-row class="my-4">
              <span class="mx-4">配送状況:</span>
              <v-chip
                size="small"
                :color="getFulfillmentStatusColor(props.order.fulfillment.status)"
              >
                {{ getFulfillmentStatus(props.order.fulfillment.status) }}
              </v-chip>
            </v-row>
            <div class="d-flex align-center">
              <v-text-field
                v-model="fulfillmentShippingCarrier"
                class="mr-4"
                name="shippingCarrier"
                label="配送会社"
                readonly
              />
              <v-text-field
                v-model="fulfillmentShippingMethod"
                class="mr-4"
                name="shippingmethod"
                label="配送方法"
                readonly
              />
              <v-text-field
                v-model="fulfillmentShippingSize"
                name="boxSize"
                label="配送時の箱の大きさ"
                readonly
              />
            </div>
            <v-text-field
              v-if="getFulfillmentStatus(props.order.fulfillment.status) == '配送済み'"
              v-model="deliveredAt"
              name="deliveredAt"
              label="配送日時"
              readonly
            />
            <v-data-table-server
              :headers="items"
              :items="props.order.items"
              no-data-text="表示する注文がありません"
            >
              <template #[`item.media`]="{ item }">
                <v-avatar>
                  <img :src="getThumbnail(item.raw.media)">
                </v-avatar>
              </template>
            </v-data-table-server>
            <v-text-field
              v-model="orderValue.fulfillment.trackingNumber"
              name="trackingNumber"
              label="伝票番号"
              readonly
            />
            <div>
              <v-text-field
                v-model="orderValue.fulfillment.lastname"
                class="mr-4"
                name="lastname"
                label="姓"
                readonly
              />
              <v-text-field
                v-model="orderValue.fulfillment.firstname"
                name="firstname"
                label="名"
                readonly
              />
            </div>
            <v-text-field
              v-model="fulfillmentPhoneNumber"
              name="phoneNumber"
              label="電話番号"
              readonly
            />
            <v-text-field
              v-model="orderValue.fulfillment.postalCode"
              name="postalCode"
              label="郵便番号"
              readonly
            />
            <div class="d-flex align-center">
              <v-text-field
                v-model="orderValue.fulfillment.prefecture"
                :items="prefecturesList"
                item-title="text"
                item-value="value"
                class="mr-4"
                name="prefecture"
                label="都道府県"
                readonly
              />
              <v-text-field
                v-model="orderValue.fulfillment.city"
                name="city"
                label="市区町村"
                readonly
              />
            </div>
            <v-text-field
              v-model="orderValue.fulfillment.addressLine1"
              name="addressLine1"
              label="町名・番地"
            />
            <v-text-field
              v-model="orderValue.fulfillment.addressLine2"
              name="addressLine2"
              label="ビル名・号室など"
            />
          </v-card-text>
        </v-card>
      </v-window-item>
    </v-window>
  </v-card>
</template>
