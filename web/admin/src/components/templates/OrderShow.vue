<script lang="ts" setup>
import { mdiDelete, mdiPlus } from '@mdi/js'
import { unix } from 'dayjs'
import type { VDataTable } from 'vuetify/lib/components/index.mjs'

import { convertI18nToJapanesePhoneNumber } from '~/lib/formatter'
import { findPrefecture, getResizedImages } from '~/lib/helpers'
import type { AlertType } from '~/lib/hooks'
import { Prefecture } from '~/types'
import {
  FulfillmentStatus,
  OrderStatus,
  PaymentMethodType,
  PaymentStatus,
  ShippingCarrier,
  ShippingSize,
  ShippingType,
  OrderType,
  RefundType,
} from '~/types/api/v1'
import type { CompleteOrderRequest, Coordinator, Order, OrderItem, OrderFulfillment, Product, ProductMedia, RefundOrderRequest, User } from '~/types/api/v1'
import type { FulfillmentInput } from '~/types/props'

const props = defineProps({
  loading: {
    type: Boolean,
    default: false,
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
  order: {
    type: Object as PropType<Order>,
    default: (): Order => ({
      id: '',
      managementId: 0,
      coordinatorId: '',
      promotionId: '',
      userId: '',
      shippingMessage: '',
      status: OrderStatus.OrderStatusUnknown,
      payment: {
        transactionId: '',
        methodType: 0,
        status: 0,
        subtotal: 0,
        discount: 0,
        shippingFee: 0,
        total: 0,
        addressId: '',
        lastname: '',
        lastnameKana: '',
        firstname: '',
        firstnameKana: '',
        postalCode: '',
        prefecture: '',
        prefectureCode: Prefecture.UNKNOWN,
        city: '',
        addressLine1: '',
        addressLine2: '',
        phoneNumber: '',
        orderedAt: 0,
        paidAt: 0,
      },
      fulfillments: [],
      refund: {
        total: 0,
        type: RefundType.RefundTypeNone,
        reason: '',
        canceled: false,
        canceledAt: 0,
      },
      items: [],
      createdAt: 0,
      updatedAt: 0,
      completedAt: 0,
      type: 0,
      experience: {
        experienceId: '',
        adultCount: 0,
        adultPrice: 0,
        juniorHighSchoolCount: 0,
        juniorHighSchoolPrice: 0,
        elementarySchoolCount: 0,
        elementarySchoolPrice: 0,
        preschoolCount: 0,
        preschoolPrice: 0,
        seniorCount: 0,
        seniorPrice: 0,
        remarks: {
          transportation: '',
          requestedDate: '',
          requestedTime: '',
        },
      },
      metadata: {
        pickupAt: 0,
        pickupLocation: '',
      },
    }),
  },
  coordinator: {
    type: Object as PropType<Coordinator>,
    default: () => ({}),
  },
  customer: {
    type: Object as PropType<User>,
    default: () => ({}),
  },
  products: {
    type: Array<Product>,
    default: () => [],
  },
  completeFormData: {
    type: Object as PropType<CompleteOrderRequest>,
    default: (): CompleteOrderRequest => ({
      shippingMessage: '',
    }),
  },
  refundFormData: {
    type: Object as PropType<RefundOrderRequest>,
    default: (): RefundOrderRequest => ({
      description: '',
    }),
  },
  fulfillmentsFormData: {
    type: Array<FulfillmentInput>,
    default: () => ([]),
  },
  cancelDialog: {
    type: Boolean,
    default: false,
  },
  refundDialog: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits<{
  (e: 'update:complete-form-data', formData: CompleteOrderRequest): void
  (e: 'update:refund-form-data', formData: RefundOrderRequest): void
  (e: 'update:fulfillments-form-data', formData: FulfillmentInput[]): void
  (e: 'update:cancel-dialog', toggle: boolean): void
  (e: 'update:refund-dialog', toggle: boolean): void
  (e: 'submit:capture'): void
  (e: 'submit:draft'): void
  (e: 'submit:complete'): void
  (e: 'submit:update-fulfillment', fulfillmentId: string): void
  (e: 'submit:cancel'): void
  (e: 'submit:refund'): void
}>()

const productHeaders: VDataTable['headers'] = [
  {
    title: '',
    key: 'media',
    width: 80,
    sortable: false,
  },
  {
    title: '商品名',
    key: 'name',
    sortable: false,
  },
  {
    title: '価格',
    key: 'price',
    sortable: false,
  },
  {
    title: '数量',
    key: 'quantity',
    sortable: false,
  },
  {
    title: '小計',
    key: 'total',
    sortable: false,
  },
]

const shippingCarriers = [
  { title: '未選択', value: ShippingCarrier.ShippingCarrierUnknown },
  { title: 'ヤマト運輸', value: ShippingCarrier.ShippingCarrierYamato },
  { title: '佐川急便', value: ShippingCarrier.ShippingCarrierSagawa },
]

const completeFormDataValue = computed({
  get: (): CompleteOrderRequest => props.completeFormData,
  set: (formData: CompleteOrderRequest): void => emit('update:complete-form-data', formData),
})
const refundFormDataValue = computed({
  get: (): RefundOrderRequest => props.refundFormData,
  set: (formData: RefundOrderRequest): void => emit('update:refund-form-data', formData),
})
const fulfillmentsFormDataValue = computed({
  get: (): FulfillmentInput[] => props.fulfillmentsFormData,
  set: (formData: FulfillmentInput[]): void => emit('update:fulfillments-form-data', formData),
})
const cancelDialogValue = computed({
  get: (): boolean => props.cancelDialog,
  set: (val: boolean): void => emit('update:cancel-dialog', val),
})
const refundDialogValue = computed({
  get: (): boolean => props.refundDialog,
  set: (val: boolean): void => emit('update:refund-dialog', val),
})

/**
 * 共通
 */
const getDatetime = (unixtime?: number): string => {
  if (!unixtime || unixtime === 0) {
    return '-'
  }
  return unix(unixtime).format('YYYY/MM/DD HH:mm:ss')
}

const getUserName = (lastname?: string, firstname?: string): string => {
  if (!lastname || lastname === '') {
    return ''
  }
  if (!firstname || firstname === '') {
    return lastname
  }
  return `${lastname} ${firstname}`
}

const isAuthorized = (): boolean => {
  return props.order.status === OrderStatus.OrderStatusWaiting
}

// 発送連絡時のメッセージ下書き保存 - 商品購入時のみ
const isPreservable = (): boolean => {
  if (!props.order || props.order.type !== OrderType.OrderTypeProduct) {
    return false
  }
  const targets: OrderStatus[] = [
    OrderStatus.OrderStatusWaiting,
    OrderStatus.OrderStatusPreparing,
    OrderStatus.OrderStatusShipped,
  ]
  return targets.includes(props.order.status)
}

// 完了通知
const isCompletable = (): boolean => {
  if (!props.order) {
    return false
  }
  const targets: OrderStatus[] = []
  switch (props.order.type) {
    case OrderType.OrderTypeProduct:
      targets.push(OrderStatus.OrderStatusPreparing, OrderStatus.OrderStatusShipped)
      break
    case OrderType.OrderTypeExperience:
      targets.push(OrderStatus.OrderStatusShipped)
      break
  }
  return targets.includes(props.order.status)
}

const isCancelable = (): boolean => {
  const targets: OrderStatus[] = [
    OrderStatus.OrderStatusWaiting,
  ]
  return targets.includes(props.order.status)
}

const isRefundable = (): boolean => {
  const targets: OrderStatus[] = [
    OrderStatus.OrderStatusPreparing,
    OrderStatus.OrderStatusShipped,
    OrderStatus.OrderStatusCompleted,
  ]
  return targets.includes(props.order.status)
}

const isUpdatableFulfillment = (): boolean => {
  const targets: OrderStatus[] = [
    OrderStatus.OrderStatusPreparing,
    OrderStatus.OrderStatusShipped,
    OrderStatus.OrderStatusCompleted,
  ]
  return targets.includes(props.order.status)
}

/**
 * 基本情報
 */
const getCoordinatorName = (): string => {
  return getUserName(props.coordinator?.lastname, props.coordinator?.firstname)
}

const getStatus = (): string => {
  switch (props.order.status) {
    case OrderStatus.OrderStatusUnpaid:
      return '支払い待ち'
    case OrderStatus.OrderStatusWaiting:
      return '受注待ち'
    case OrderStatus.OrderStatusPreparing:
      return props.order.type === OrderType.OrderTypeExperience ? '体験準備中' : '発送準備中'
    case OrderStatus.OrderStatusShipped:
      return '発送完了'
    case OrderStatus.OrderStatusCompleted:
      return '完了'
    case OrderStatus.OrderStatusCanceled:
      return 'キャンセル'
    case OrderStatus.OrderStatusRefunded:
      return '返金'
    case OrderStatus.OrderStatusFailed:
      return '失敗'
    default:
      return '不明'
  }
}

const getStatusColor = (): string => {
  switch (props.order.status) {
    case OrderStatus.OrderStatusUnpaid:
      return 'secondary'
    case OrderStatus.OrderStatusWaiting:
      return 'secondary'
    case OrderStatus.OrderStatusPreparing:
      return 'info'
    case OrderStatus.OrderStatusShipped:
      return 'info'
    case OrderStatus.OrderStatusCompleted:
      return 'primary'
    case OrderStatus.OrderStatusCanceled:
      return 'warning'
    case OrderStatus.OrderStatusRefunded:
      return 'warning'
    case OrderStatus.OrderStatusFailed:
      return 'error'
    default:
      return 'unknown'
  }
}

const getOrderedAt = (): string => {
  return getDatetime(props.order?.payment?.orderedAt)
}

const getCompletedAt = (): string => {
  return getDatetime(props.order?.completedAt)
}

/**
 * 支払い情報
 */
const getAllItems = computed(() => {
  const items: OrderItem[] = []
  if (!props.order) {
    return items
  }
  props.order?.items?.forEach((item: OrderItem): void => {
    const index = items.findIndex((v: OrderItem): boolean => {
      return v.productId === item.productId
    })
    if (index < 0) {
      items.push(item)
    }
    else {
      items[index].quantity += item.quantity
    }
  })
  return items
})

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

const getPaymentMethodType = (): string => {
  switch (props.order?.payment?.methodType) {
    case PaymentMethodType.PaymentMethodTypeCash:
      return '代引支払い'
    case PaymentMethodType.PaymentMethodTypeCreditCard:
      return 'クレジットカード決済'
    case PaymentMethodType.PaymentMethodTypeKonbini:
      return 'コンビニ決済'
    case PaymentMethodType.PaymentMethodTypeBankTransfer:
      return '銀行振込決済'
    case PaymentMethodType.PaymentMethodTypePayPay:
      return 'QR決済（PayPay）'
    case PaymentMethodType.PaymentMethodTypeLinePay:
      return 'QR決済（LINE Pay）'
    case PaymentMethodType.PaymentMethodTypeMerpay:
      return 'QR決済（メルペイ）'
    case PaymentMethodType.PaymentMethodTypeRakutenPay:
      return 'QR決済（楽天ペイ）'
    case PaymentMethodType.PaymentMethodTypeAUPay:
      return 'QR決済（au PAY）'
    case PaymentMethodType.PaymentMethodTypePaidy:
      return 'ペイディ（Paidy）'
    case PaymentMethodType.PaymentMethodTypePayEasy:
      return 'ペイジー（Pay-easy）'
    default:
      return '不明'
  }
}

const getPaidAt = (): string => {
  return getDatetime(props.order?.payment?.paidAt)
}

const getSubtotal = (item: OrderItem): number => {
  return item.price * item.quantity
}

/**
 * 注文情報
 */
const getProduct = (productId: string): Product | undefined => {
  return props.products.find((product: Product): boolean => {
    return product.id === productId
  })
}

const getProductName = (productId: string): string => {
  const product = getProduct(productId)
  return product?.name || ''
}

const getThumbnail = (productId: string): string => {
  const product = getProduct(productId)
  const thumbnail = product?.media.find((media: ProductMedia): boolean => {
    return media.isThumbnail
  })
  return thumbnail?.url || ''
}

const getResizedThumbnails = (productId: string): string => {
  const product = getProduct(productId)
  const thumbnail = product?.media.find((media: ProductMedia) => {
    return media.isThumbnail
  })
  if (!thumbnail) {
    return ''
  }
  return getResizedImages(thumbnail.url)
}

/**
 * 配送情報
 */
const getFulfillmentStatus = (status: FulfillmentStatus): string => {
  switch (status) {
    case FulfillmentStatus.FulfillmentStatusUnfulfilled:
      return '未配送'
    case FulfillmentStatus.FulfillmentStatusFulfilled:
      return '配送済み'
    default:
      return '不明'
  }
}

const getFulfillmentStatusColor = (status: FulfillmentStatus): string => {
  switch (status) {
    case FulfillmentStatus.FulfillmentStatusUnfulfilled:
      return 'error'
    case FulfillmentStatus.FulfillmentStatusFulfilled:
      return 'primary'
    default:
      return 'unkown'
  }
}

const getShippingType = (shippingType: ShippingType): string => {
  switch (shippingType) {
    case ShippingType.ShippingTypeNormal:
      return '常温・冷蔵便'
    case ShippingType.ShippingTypeFrozen:
      return '冷凍便'
    case ShippingType.ShippingTypePickup:
      return '店舗受け取り'
    default:
      return '不明'
  }
}

const getBoxSize = (size: ShippingSize): string => {
  switch (size) {
    case ShippingSize.ShippingSize60:
      return '60'
    case ShippingSize.ShippingSize80:
      return '80'
    case ShippingSize.ShippingSize100:
      return '100'
    default:
      return '不明'
  }
}

/**
 * 顧客情報
 * ※配送先情報は現状すべて同じ情報が返されるため、先頭の値を取得して表示する
 */
const getCustomerName = (): string => {
  return getUserName(props.customer?.lastname, props.customer?.firstname)
}

const getCustomerNameKana = (): string => {
  return getUserName(props.customer?.lastnameKana, props.customer?.firstnameKana)
}

const getShippingAddressName = (): string => {
  return getUserName(props.order?.payment?.lastname, props.order?.payment?.firstname)
}

const getShippingAddressPhoneNumber = (): string => {
  return props.order?.payment?.phoneNumber ? convertI18nToJapanesePhoneNumber(props.order?.payment?.phoneNumber) : ''
}

const getShippingAddressPrefecture = (): string => {
  const prefecture = findPrefecture(props.order?.payment?.prefectureCode)
  return prefecture ? prefecture.text : ''
}

const getFulfillmentAddressName = (): string => {
  if (!props.order || !props.order.fulfillments || props.order.fulfillments.length === 0) {
    return ''
  }
  return getUserName(props.order.fulfillments[0]?.lastname, props.order.fulfillments[0]?.firstname)
}

const getFulfillmentAddressPhoneNumber = (): string => {
  if (!props.order || !props.order.fulfillments || props.order.fulfillments.length === 0) {
    return ''
  }
  return props.order.fulfillments[0]?.phoneNumber ? convertI18nToJapanesePhoneNumber(props.order.fulfillments[0].phoneNumber) : ''
}

const getFulfillmentAddressPrefecture = (): string => {
  if (!props.order || !props.order.fulfillments || props.order.fulfillments.length === 0) {
    return ''
  }
  const prefecture = findPrefecture(props.order?.fulfillments[0]?.prefectureCode)
  return prefecture ? prefecture.text : ''
}

const getRequestDaliveryDay = (fulfillment: OrderFulfillment): string => {
  // TODO: API側の実装ができ次第実装する
  return '未指定'
}

// 店舗受け取りかどうかを判定
const isPickupShipping = (): boolean => {
  if (!props.order || !props.order.fulfillments || props.order.fulfillments.length === 0) {
    return false
  }
  return props.order.fulfillments[0]?.shippingType === ShippingType.ShippingTypePickup
}

// 受け取り日時を取得
const getPickupDate = (): string => {
  const pickupAt = props.order?.metadata?.pickupAt
  if (!pickupAt || pickupAt === 0) {
    return '未指定'
  }
  return unix(pickupAt).format('YYYY年MM月DD日 HH:mm')
}

// 受け取り場所を取得
const getPickupLocation = (): string => {
  return props.order?.metadata?.pickupLocation || '未指定'
}

const getOrderItems = (fulfillmentId: string): OrderItem[] => {
  const items = props.order?.items?.filter((item: OrderItem): boolean => {
    return item.fulfillmentId === fulfillmentId
  })
  return items
}

const showShippingMessage = (): boolean => {
  if (!props.order || props.order.type !== OrderType.OrderTypeProduct) {
    return false
  }
  const targets: OrderStatus[] = [
    OrderStatus.OrderStatusPreparing,
    OrderStatus.OrderStatusShipped,
    OrderStatus.OrderStatusCompleted,
  ]
  return targets.includes(props.order.status)
}

const onClickOpenCancelDialog = (): void => {
  emit('update:cancel-dialog', true)
}

const onClickCloseCancelDialog = (): void => {
  emit('update:cancel-dialog', false)
}

const onClickOpenRefundDialog = (): void => {
  emit('update:refund-dialog', true)
}

const onClickCloseRefundDialog = (): void => {
  emit('update:refund-dialog', false)
}

const onSubmitUpdate = (fulfillmentId: string): void => {
  emit('submit:update-fulfillment', fulfillmentId)
}

const onSubmitSaveDraft = (): void => {
  emit('submit:draft')
}

const onSubmitCapture = (): void => {
  emit('submit:capture')
}

const onSubmitComplete = (): void => {
  emit('submit:complete')
}

const onSubmitCancel = (): void => {
  emit('submit:cancel')
}

const onSubmitRefund = (): void => {
  emit('submit:refund')
}
</script>

<template>
  <v-alert
    v-show="props.isAlert"
    :type="props.alertType"
    v-text="props.alertText"
  />

  <v-dialog
    v-model="cancelDialogValue"
    width="500"
  >
    <v-card>
      <v-card-title>
        本当に注文キャンセルしますか？
      </v-card-title>
      <v-card-actions>
        <v-spacer />
        <v-btn
          color="error"
          variant="text"
          @click="onClickCloseCancelDialog"
        >
          キャンセル
        </v-btn>
        <v-btn
          :loading="loading"
          color="primary"
          variant="outlined"
          @click="onSubmitCancel"
        >
          削除
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog
    v-model="refundDialogValue"
    width="500"
  >
    <v-card>
      <v-card-title>
        返金依頼
      </v-card-title>
      <v-card-text>
        <v-text-field
          v-model="refundFormDataValue.description"
          label="返金理由"
          maxlength="200"
        />
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn
          color="error"
          variant="text"
          @click="onClickCloseRefundDialog"
        >
          キャンセル
        </v-btn>
        <v-btn
          :loading="loading"
          color="primary"
          variant="outlined"
          @click="onSubmitRefund"
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
      lg="8"
    >
      <v-card
        elevation="0"
        class="mb-4"
      >
        <v-card-title class="pb-4">
          基本情報
        </v-card-title>
        <v-card-text>
          <v-row>
            <v-col cols="3">
              注文番号
            </v-col>
            <v-col cols="9">
              {{ order?.id || '' }}
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="3">
              コーディネーター名
            </v-col>
            <v-col cols="9">
              {{ getCoordinatorName() }}
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="3">
              ステータス
            </v-col>
            <v-col cols="9">
              <v-chip
                size="small"
                :color="getStatusColor()"
              >
                {{ getStatus() }}
              </v-chip>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="3">
              購入日時
            </v-col>
            <v-col cols="9">
              {{ getOrderedAt() }}
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="3">
              発送完了日時
            </v-col>
            <v-col cols="9">
              {{ getCompletedAt() }}
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
      <v-card
        elevation="0"
        class="mb-4"
      >
        <v-card-title class="pb-4">
          支払い情報
        </v-card-title>
        <v-card-text>
          <v-row>
            <v-col cols="3">
              ステータス
            </v-col>
            <v-col cols="9">
              <v-chip
                size="small"
                :color="getPaymentStatusColor(order?.payment?.status)"
              >
                {{ getPaymentStatus(order?.payment?.status) }}
              </v-chip>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="3">
              支払い方法
            </v-col>
            <v-col cols="9">
              {{ getPaymentMethodType() }}
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="3">
              支払い日時
            </v-col>
            <v-col cols="9">
              {{ getPaidAt() }}
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12">
              <v-table>
                <tbody class="text-grey">
                  <tr>
                    <td>小計</td>
                    <td>{{ getAllItems.length }}つのアイテム</td>
                    <td>&yen; {{ (order?.payment?.subtotal || 0).toLocaleString() }}</td>
                  </tr>
                  <tr>
                    <td>配送手数料</td>
                    <td>{{ (order?.fulfillments?.length || 0) }}つの箱</td>
                    <td>&yen; {{ (order?.payment?.shippingFee || 0).toLocaleString() }}</td>
                  </tr>
                  <tr>
                    <td>割引金額</td>
                    <td />
                    <td>&yen; {{ (order?.payment?.discount || 0).toLocaleString() }}</td>
                  </tr>
                </tbody>
                <tfoot>
                  <tr>
                    <td>支払い合計（税込み）</td>
                    <td />
                    <td>&yen; {{ (order?.payment?.total || 0).toLocaleString() }}</td>
                  </tr>
                </tfoot>
              </v-table>
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
      <v-card
        v-if="props.order.type === OrderType.OrderTypeProduct"
        elevation="0"
        class="mb-4"
      >
        <v-card-title class="pb-4">
          注文情報
        </v-card-title>
        <v-card-text>
          <v-data-table
            :headers="productHeaders"
            :items="getAllItems"
          >
            <template #[`item.media`]="{ item }">
              <v-img
                aspect-ratio="1/1"
                :max-height="56"
                :max-width="80"
                :src="getThumbnail(item.productId)"
                :srcset="getResizedThumbnails(item.productId)"
              />
            </template>
            <template #[`item.name`]="{ item }">
              {{ getProductName(item.productId) }}
            </template>
            <template #[`item.price`]="{ item }">
              &yen; {{ item.price.toLocaleString() }}
            </template>
            <template #[`item.quantity`]="{ item }">
              {{ item.quantity.toLocaleString() }}
            </template>
            <template #[`item.total`]="{ item }">
              &yen; {{ getSubtotal(item).toLocaleString() }}
            </template>
          </v-data-table>

          <v-list>
            <v-list-item class="mb-4">
              <v-list-item-subtitle class="mb-2">
                顧客情報
              </v-list-item-subtitle>
              <div>{{ getCustomerName() }}</div>
              <div>{{ getCustomerNameKana() }}</div>
              <div class="mt-1">
                &#128231; {{ props.customer.email }}
              </div>
              <div>&phone; {{ props.customer.phoneNumber }}</div>
            </v-list-item>
          </v-list>
        </v-card-text>
      </v-card>
      <v-card
        v-if="props.order.type === OrderType.OrderTypeExperience"
        elevation="0"
        class="mb-4"
      >
        <v-card-title class="pb-4">
          予約情報
        </v-card-title>
        <v-card-text>
          <v-row>
            <v-col cols="3">
              大人:
            </v-col>
            <v-col cols="3">
              {{ props.order?.experience?.adultCount || 0 }}人
            </v-col>
            <v-col cols="6">
              合計: {{ (props.order?.experience?.adultPrice || 0) * (props.order?.experience?.adultCount || 0) }}円
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="3">
              未就学児(3歳〜):
            </v-col>
            <v-col cols="3">
              {{ props.order?.experience?.preschoolCount || 0 }}人
            </v-col>
            <v-col cols="6">
              合計: {{ (props.order?.experience?.preschoolPrice || 0) * (props.order?.experience?.preschoolCount || 0) }}円
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="3">
              小学生:
            </v-col>
            <v-col cols="3">
              {{ props.order?.experience?.elementarySchoolCount || 0 }}人
            </v-col>
            <v-col cols="6">
              合計: {{ (props.order?.experience?.elementarySchoolPrice || 0) * (props.order?.experience?.elementarySchoolCount || 0) }}円
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="3">
              中学生:
            </v-col>
            <v-col cols="3">
              {{ props.order?.experience?.juniorHighSchoolCount || 0 }}人
            </v-col>
            <v-col cols="6">
              合計: {{ (props.order?.experience?.juniorHighSchoolPrice || 0) * (props.order?.experience?.juniorHighSchoolCount || 0) }}円
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="3">
              シニア(65歳〜):
            </v-col>
            <v-col cols="3">
              {{ props.order?.experience?.seniorCount || 0 }}人
            </v-col>
            <v-col cols="6">
              合計: {{ (props.order?.experience?.seniorPrice || 0) * (props.order?.experience?.seniorCount || 0) }}円
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
    </v-col>

    <v-col
      sm="12"
      md="12"
      lg="4"
    >
      <v-card
        elevation="0"
        class="mb-4"
      >
        <v-card-text>
          <v-list>
            <v-list-item class="mb-4">
              <v-list-item-subtitle class="mb-2">
                顧客情報
              </v-list-item-subtitle>
              <div>{{ getCustomerName() }}</div>
              <div>{{ getCustomerNameKana() }}</div>
              <div class="mt-1">
                &#128231; {{ props.customer.email }}
              </div>
              <div>&phone; {{ props.customer.phoneNumber }}</div>
            </v-list-item>
            <v-list-item
              v-if="props.order.type !== OrderType.OrderTypeExperience"
              class="mb-4"
            >
              <v-list-item-subtitle class="pb-2">
                請求先情報
              </v-list-item-subtitle>
              <div>{{ getShippingAddressName() }}</div>
              <div class="mt-1">
                &phone; {{ getShippingAddressPhoneNumber() }}
              </div>
              <div class="mt-1">
                &#12306; {{ props.order?.payment?.postalCode || '' }}
              </div>
              <div>{{ `${getShippingAddressPrefecture()} ${props.order?.payment?.city || ''}` }}</div>
              <div>{{ props.order?.payment?.addressLine1 || '' }}</div>
              <div>{{ props.order?.payment?.addressLine2 || '' }}</div>
            </v-list-item>
            <v-list-item v-if="props.order?.fulfillments?.length > 0">
              <v-list-item-subtitle class="pb-2">
                {{ isPickupShipping() ? '受け取り情報' : '配送先情報' }}
              </v-list-item-subtitle>
              <template v-if="isPickupShipping()">
                <div class="mt-1">
                  📍 受け取り場所: {{ getPickupLocation() }}
                </div>
                <div class="mt-1">
                  📅 受け取り日時: {{ getPickupDate() }}
                </div>
              </template>
              <template v-else>
                <div>{{ getFulfillmentAddressName() }}</div>
                <div class="mt-1">
                  &phone; {{ getFulfillmentAddressPhoneNumber() }}
                </div>
                <div class="mt-1">
                  &#12306; {{ props.order.fulfillments[0]?.postalCode || '' }}
                </div>
                <div>{{ `${getFulfillmentAddressPrefecture()} ${props.order.fulfillments[0]?.city || ''}` }}</div>
                <div>{{ props.order.fulfillments[0]?.addressLine1 || '' }}</div>
                <div>{{ props.order.fulfillments[0]?.addressLine2 || '' }}</div>
              </template>
            </v-list-item>
          </v-list>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>

  <v-row
    v-for="(fulfillment, index) in props.order.fulfillments"
    :key="fulfillment.fulfillmentId"
  >
    <v-col
      sm="12"
      md="12"
      lg="8"
    >
      <v-card
        elevation="0"
        class="mb-4"
      >
        <v-card-title class="pb-4">
          配送詳細 {{ index + 1 }}
        </v-card-title>
        <v-card-text>
          <v-row>
            <v-col cols="3">
              ステータス
            </v-col>
            <v-col cols="9">
              <v-chip
                size="small"
                :color="getFulfillmentStatusColor(fulfillment.status)"
              >
                {{ getFulfillmentStatus(fulfillment.status) }}
              </v-chip>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="3">
              箱のタイプ
            </v-col>
            <v-col cols="9">
              {{ getShippingType(fulfillment.shippingType) }}
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="3">
              箱のサイズ
            </v-col>
            <v-col cols="9">
              {{ getBoxSize(fulfillment.boxSize) }}
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="3">
              箱の占有率
            </v-col>
            <v-col cols="9">
              {{ fulfillment.boxRate }}
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="3">
              配送希望日
            </v-col>
            <v-col cols="9">
              {{ getRequestDaliveryDay(fulfillment) }}
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12">
              <v-data-table
                :headers="productHeaders"
                :items="getOrderItems(fulfillment.fulfillmentId)"
              >
                <template #[`item.media`]="{ item }">
                  <v-img
                    aspect-ratio="1/1"
                    :max-height="56"
                    :max-width="80"
                    :src="getThumbnail(item.productId)"
                    :srcset="getResizedThumbnails(item.productId)"
                  />
                </template>
                <template #[`item.name`]="{ item }">
                  {{ getProductName(item.productId) }}
                </template>
                <template #[`item.price`]="{ item }">
                  &yen; {{ item.price.toLocaleString() }}
                </template>
                <template #[`item.quantity`]="{ item }">
                  {{ item.quantity.toLocaleString() }}
                </template>
                <template #[`item.total`]="{ item }">
                  &yen; {{ getSubtotal(item).toLocaleString() }}
                </template>
              </v-data-table>
            </v-col>
          </v-row>
        </v-card-text>
      </v-card>
    </v-col>
    <v-col
      sm="12"
      md="12"
      lg="4"
    >
      <v-card
        elevation="0"
        class="mb-4"
      >
        <v-card-text>
          <v-select
            v-model="fulfillmentsFormDataValue[index].shippingCarrier"
            label="配送業者"
            :items="shippingCarriers"
            :readonly="!isUpdatableFulfillment()"
          />
          <v-text-field
            v-model="fulfillmentsFormDataValue[index].trackingNumber"
            label="伝票番号"
            :readonly="!isUpdatableFulfillment()"
          />
          <v-btn
            v-show="isUpdatableFulfillment()"
            :loading="loading"
            class="mt-2"
            variant="outlined"
            @click="onSubmitUpdate(fulfillment.fulfillmentId)"
          >
            <v-icon
              start
              :icon="mdiPlus"
            />
            更新
          </v-btn>
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
  <v-row>
    <v-col
      v-show="showShippingMessage()"
      sm="12"
      md="12"
      lg="8"
    >
      <v-card>
        <v-card-title class="pb-4">
          発送連絡
        </v-card-title>
        <v-card-text>
          <v-textarea
            v-model="completeFormDataValue.shippingMessage"
            label="お客様へのメッセージ"
            placeholder="例：ご注文ありがとうございます！商品の発送が完了しました。商品到着まで今しばらくお待ち下さい。"
            :readonly="!isPreservable()"
          />
        </v-card-text>
      </v-card>
    </v-col>
    <v-col
      sm="12"
      md="12"
      lg="8"
    >
      <v-btn
        v-show="isPreservable()"
        :loading="loading"
        variant="outlined"
        color="info"
        class="mr-2"
        @click="onSubmitSaveDraft()"
      >
        <v-icon
          start
          :icon="mdiPlus"
        />
        下書きを保存
      </v-btn>
      <v-btn
        v-show="isAuthorized()"
        :loading="loading"
        variant="outlined"
        color="primary"
        class="mr-2"
        @click="onSubmitCapture()"
      >
        <v-icon
          start
          :icon="mdiPlus"
        />
        注文を確定
      </v-btn>
      <v-btn
        v-show="isCompletable()"
        :loading="loading"
        variant="outlined"
        color="primary"
        class="mr-2"
        @click="onSubmitComplete()"
      >
        <v-icon
          start
          :icon="mdiPlus"
        />
        {{ props.order.type === OrderType.OrderTypeProduct ? '発送完了を通知' : 'レビュー依頼を送信' }}
      </v-btn>
      <v-btn
        v-show="isCancelable()"
        :loading="loading"
        variant="outlined"
        color="error"
        class="mr-2"
        @click="onClickOpenCancelDialog()"
      >
        <v-icon
          start
          :icon="mdiDelete"
        />
        注文をキャンセル
      </v-btn>
      <v-btn
        v-show="isRefundable()"
        :loading="loading"
        variant="outlined"
        color="error"
        @click="onClickOpenRefundDialog"
      >
        <v-icon
          start
          :icon="mdiDelete"
        />
        返金
      </v-btn>
    </v-col>
  </v-row>
</template>
