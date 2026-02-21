<script lang="ts" setup>
import { mdiDelete, mdiPlus, mdiInformation, mdiCalculator, mdiPackageVariant, mdiCalendarStar, mdiCog, mdiAccount, mdiReceipt, mdiStore, mdiTruck, mdiPackageVariantClosed, mdiSend, mdiAlertCircle, mdiCreditCardRefund, mdiCheckCircle, mdiContentSave, mdiCancel, mdiUpdate, mdiCircle, mdiCurrencyJpy, mdiCreditCard, mdiCart, mdiTag, mdiPhone, mdiEmail, mdiMapMarker, mdiCalendar, mdiSchool, mdiAccountSchool, mdiAccountSupervisor, mdiBaby, mdiStoreClock, mdiAlert } from '@mdi/js'
import { unix } from 'dayjs'
import type { VDataTable } from 'vuetify/components'

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
    title: '単価',
    key: 'price',
    sortable: false,
    width: 120,
  },
  {
    title: '数量',
    key: 'quantity',
    sortable: false,
    width: 100,
  },
  {
    title: '小計',
    key: 'total',
    sortable: false,
    width: 120,
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

const getAllItemsCount = computed(() => {
  let count = getAllItems.value.length

  // 体験注文の場合は参加者数も含める
  if (props.order?.type === OrderType.OrderTypeExperience && props.order.experience) {
    count += (props.order.experience.adultCount || 0)
    count += (props.order.experience.preschoolCount || 0)
    count += (props.order.experience.elementarySchoolCount || 0)
    count += (props.order.experience.juniorHighSchoolCount || 0)
    count += (props.order.experience.seniorCount || 0)
  }

  return count
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
  <v-container
    fluid
    class="pa-0"
  >
    <v-alert
      v-show="props.isAlert"
      class="mb-4"
      :type="props.alertType"
      v-text="props.alertText"
    />

    <!-- Order Header -->
    <v-card
      elevation="2"
      class="mb-6"
    >
      <v-card-title class="bg-primary text-white pa-6">
        <v-row align="center">
          <v-col
            cols="12"
            md="8"
          >
            <h2 class="text-h4 font-weight-bold mb-2">
              注文 #{{ order?.managementId || order?.id?.slice(-8) || '' }}
            </h2>
            <p class="text-h6 mb-0">
              {{ getCoordinatorName() || '未設定' }}
            </p>
          </v-col>
          <v-col
            cols="12"
            md="4"
            class="text-md-right"
          >
            <v-chip
              size="large"
              class="bg-white ma-1"
              :color="getStatusColor()"
            >
              <v-icon
                start
                size="small"
                :icon="mdiCircle"
              />
              {{ getStatus() }}
            </v-chip>
          </v-col>
        </v-row>
      </v-card-title>
    </v-card>

    <!-- Summary Cards -->
    <v-row class="mb-6">
      <v-col
        cols="12"
        md="4"
      >
        <v-card elevation="2">
          <v-card-text class="text-center">
            <v-icon
              size="32"
              color="primary"
              class="mb-2"
              :icon="mdiCurrencyJpy"
            />
            <h3 class="text-h4 font-weight-bold mb-1">
              ¥{{ (order?.payment?.total || 0).toLocaleString() }}
            </h3>
            <p class="text-subtitle-1 text-grey">
              合計金額
            </p>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col
        cols="12"
        md="4"
      >
        <v-card elevation="2">
          <v-card-text class="text-center">
            <v-icon
              size="32"
              :color="getPaymentStatusColor(order?.payment.status)"
              :icon="mdiCreditCard"
              class="mb-2"
            />
            <h3 class="text-h5 font-weight-bold mb-1">
              {{ getPaymentStatus(order?.payment.status) }}
            </h3>
            <p class="text-subtitle-1 text-grey">
              支払い状況
            </p>
          </v-card-text>
        </v-card>
      </v-col>
      <v-col
        cols="12"
        md="4"
      >
        <v-card elevation="2">
          <v-card-text class="text-center">
            <v-icon
              size="32"
              color="info"
              class="mb-2"
              :icon="props.order.type === OrderType.OrderTypeExperience ? mdiCalendarStar : mdiPackageVariant"
            />
            <h3 class="text-h5 font-weight-bold mb-1">
              {{ getAllItemsCount }}
            </h3>
            <p class="text-subtitle-1 text-grey">
              {{ props.order.type === OrderType.OrderTypeExperience ? '予約アイテム' : '商品アイテム' }}
            </p>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Main Content -->
    <v-row>
      <v-col
        cols="12"
        lg="8"
      >
        <!-- Basic Information -->
        <v-card
          elevation="2"
          class="mb-4"
        >
          <v-card-title class="bg-grey-lighten-4 py-4">
            <v-icon
              class="mr-2"
              color="primary"
              :icon="mdiInformation"
            />
            基本情報
          </v-card-title>
          <v-card-text class="pa-6">
            <v-row>
              <v-col
                cols="12"
                sm="6"
              >
                <div class="mb-4">
                  <p class="text-subtitle-2 text-grey mb-1">
                    注文番号
                  </p>
                  <p class="text-body-1 font-weight-medium">
                    {{ order?.id || '' }}
                  </p>
                </div>
                <div class="mb-4">
                  <p class="text-subtitle-2 text-grey mb-1">
                    購入日時
                  </p>
                  <p class="text-body-1 font-weight-medium">
                    {{ getOrderedAt() }}
                  </p>
                </div>
              </v-col>
              <v-col
                cols="12"
                sm="6"
              >
                <div class="mb-4">
                  <p class="text-subtitle-2 text-grey mb-1">
                    支払い方法
                  </p>
                  <p class="text-body-1 font-weight-medium">
                    {{ getPaymentMethodType() }}
                  </p>
                </div>
                <div class="mb-4">
                  <p class="text-subtitle-2 text-grey mb-1">
                    {{ props.order.type === OrderType.OrderTypeProduct ? '発送完了日時' : '完了日時' }}
                  </p>
                  <p class="text-body-1 font-weight-medium">
                    {{ getCompletedAt() }}
                  </p>
                </div>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>

        <!-- Payment Breakdown -->
        <v-card
          elevation="2"
          class="mb-4"
        >
          <v-card-title class="bg-grey-lighten-4 py-4">
            <v-icon
              class="mr-2"
              color="primary"
              :icon="mdiCalculator"
            />
            支払い詳細
          </v-card-title>
          <v-card-text class="pa-0">
            <v-table>
              <tbody>
                <tr>
                  <td class="pa-4">
                    <v-icon
                      size="20"
                      class="mr-2"
                      color="grey"
                      :icon="mdiCart"
                    />
                    小計
                  </td>
                  <td class="pa-4 text-grey">
                    {{ getAllItemsCount }}アイテム
                  </td>
                  <td class="pa-4 text-right font-weight-medium">
                    ¥{{ (order?.payment?.subtotal || 0).toLocaleString() }}
                  </td>
                </tr>
                <tr>
                  <td class="pa-4">
                    <v-icon
                      size="20"
                      class="mr-2"
                      color="grey"
                      :icon="mdiTruck"
                    />
                    配送手数料
                  </td>
                  <td class="pa-4 text-grey">
                    {{ (order?.fulfillments?.length || 0) }}箱
                  </td>
                  <td class="pa-4 text-right font-weight-medium">
                    ¥{{ (order?.payment?.shippingFee || 0).toLocaleString() }}
                  </td>
                </tr>
                <tr v-if="order?.payment?.discount > 0">
                  <td class="pa-4">
                    <v-icon
                      size="20"
                      class="mr-2"
                      color="success"
                      :icon="mdiTag"
                    />
                    割引
                  </td>
                  <td class="pa-4" />
                  <td class="pa-4 text-right font-weight-medium text-success">
                    -¥{{ (order?.payment?.discount || 0).toLocaleString() }}
                  </td>
                </tr>
              </tbody>
              <tfoot class="bg-grey-lighten-4">
                <tr>
                  <td class="pa-4 font-weight-bold">
                    合計（税込）
                  </td>
                  <td class="pa-4" />
                  <td class="pa-4 text-right text-h6 font-weight-bold text-primary">
                    ¥{{ (order?.payment?.total || 0).toLocaleString() }}
                  </td>
                </tr>
              </tfoot>
            </v-table>
          </v-card-text>
        </v-card>

        <!-- Product Information -->
        <v-card
          v-if="props.order.type === OrderType.OrderTypeProduct"
          elevation="2"
          class="mb-4"
        >
          <v-card-title class="bg-grey-lighten-4 py-4">
            <v-icon
              class="mr-2"
              color="primary"
              :icon="mdiPackageVariant"
            />
            注文商品
          </v-card-title>
          <v-card-text class="pa-0">
            <v-data-table
              :headers="productHeaders"
              :items="getAllItems"
              hide-default-footer
            >
              <template #[`item.media`]="{ item }">
                <v-avatar
                  size="64"
                  rounded="lg"
                  class="ma-2"
                >
                  <v-img
                    :src="getThumbnail(item.productId)"
                    :srcset="getResizedThumbnails(item.productId)"
                    aspect-ratio="1"
                  />
                </v-avatar>
              </template>
              <template #[`item.name`]="{ item }">
                <div class="font-weight-medium">
                  {{ getProductName(item.productId) }}
                </div>
              </template>
              <template #[`item.price`]="{ item }">
                <span class="font-weight-medium">¥{{ item.price.toLocaleString() }}</span>
              </template>
              <template #[`item.quantity`]="{ item }">
                <v-chip
                  size="small"
                  color="info"
                  variant="tonal"
                >
                  {{ item.quantity.toLocaleString() }}
                </v-chip>
              </template>
              <template #[`item.total`]="{ item }">
                <span class="font-weight-bold text-primary">¥{{ getSubtotal(item).toLocaleString() }}</span>
              </template>
            </v-data-table>
          </v-card-text>
        </v-card>

        <!-- Experience Information -->
        <v-card
          v-if="props.order.type === OrderType.OrderTypeExperience"
          elevation="2"
          class="mb-4"
        >
          <v-card-title class="bg-grey-lighten-4 py-4">
            <v-icon
              class="mr-2"
              color="primary"
              :icon="mdiCalendarStar"
            />
            体験予約詳細
          </v-card-title>
          <v-card-text class="pa-6">
            <organisms-order-experience-details
              :experience="props.order?.experience"
              :loading="loading"
              variant="default"
            />
          </v-card-text>
        </v-card>
      </v-col>

      <!-- Sidebar -->
      <v-col
        cols="12"
        lg="4"
      >
        <!-- Customer Information -->
        <v-card
          elevation="2"
          class="mb-4"
        >
          <v-card-title class="bg-grey-lighten-4 py-4">
            <v-icon
              class="mr-2"
              color="primary"
              :icon="mdiAccount"
            />
            顧客情報
          </v-card-title>
          <v-card-text class="pa-4">
            <div class="mb-4">
              <h4 class="text-h6 font-weight-medium mb-2">
                {{ getCustomerName() }}
              </h4>
              <p class="text-body-2 text-grey mb-2">
                {{ getCustomerNameKana() }}
              </p>
              <div class="d-flex align-center mb-2">
                <v-icon
                  size="16"
                  class="mr-2"
                  color="grey"
                  :icon="mdiEmail"
                />
                <span class="text-body-2">{{ props.customer.email }}</span>
              </div>
              <div class="d-flex align-center">
                <v-icon
                  size="16"
                  class="mr-2"
                  color="grey"
                  :icon="mdiPhone"
                />
                <span class="text-body-2">{{ props.customer.phoneNumber }}</span>
              </div>
            </div>
          </v-card-text>
        </v-card>

        <!-- Billing Address -->
        <v-card
          v-if="props.order.type !== OrderType.OrderTypeExperience && getShippingAddressName()"
          elevation="2"
          class="mb-4"
        >
          <v-card-title class="bg-grey-lighten-4 py-4">
            <v-icon
              class="mr-2"
              color="primary"
              :icon="mdiReceipt"
            />
            請求先情報
          </v-card-title>
          <v-card-text class="pa-4">
            <h4 class="text-subtitle-1 font-weight-medium mb-2">
              {{ getShippingAddressName() }}
            </h4>
            <div class="d-flex align-center mb-1">
              <v-icon
                size="16"
                class="mr-2"
                color="grey"
                :icon="mdiPhone"
              />
              <span class="text-body-2">{{ getShippingAddressPhoneNumber() }}</span>
            </div>
            <div class="d-flex align-start mt-2">
              <v-icon
                size="16"
                class="mr-2 mt-1"
                color="grey"
                :icon="mdiMapMarker"
              />
              <div class="text-body-2">
                <div>〒{{ props.order?.payment?.postalCode || '' }}</div>
                <div>{{ `${getShippingAddressPrefecture()} ${props.order?.payment?.city || ''}` }}</div>
                <div>{{ props.order?.payment?.addressLine1 || '' }}</div>
                <div v-if="props.order?.payment?.addressLine2">
                  {{ props.order?.payment?.addressLine2 }}
                </div>
              </div>
            </div>
          </v-card-text>
        </v-card>

        <!-- Shipping Information -->
        <v-card
          v-if="props.order?.fulfillments?.length > 0"
          elevation="2"
          class="mb-4"
        >
          <v-card-title class="bg-grey-lighten-4 py-4">
            <v-icon
              class="mr-2"
              color="primary"
              :icon="isPickupShipping() ? mdiStore : mdiTruck"
            />
            {{ isPickupShipping() ? '受け取り情報' : '配送先情報' }}
          </v-card-title>
          <v-card-text class="pa-4">
            <template v-if="isPickupShipping()">
              <div class="d-flex align-center mb-3">
                <v-icon
                  size="20"
                  class="mr-3"
                  color="info"
                  :icon="mdiStore"
                />
                <div>
                  <p class="text-subtitle-2 mb-1">
                    受け取り場所
                  </p>
                  <p class="text-body-1 font-weight-medium">
                    {{ getPickupLocation() }}
                  </p>
                </div>
              </div>
              <div class="d-flex align-center">
                <v-icon
                  size="20"
                  class="mr-3"
                  color="info"
                  :icon="mdiCalendar"
                />
                <div>
                  <p class="text-subtitle-2 mb-1">
                    受け取り日時
                  </p>
                  <p class="text-body-1 font-weight-medium">
                    {{ getPickupDate() }}
                  </p>
                </div>
              </div>
            </template>
            <template v-else>
              <h4 class="text-subtitle-1 font-weight-medium mb-2">
                {{ getFulfillmentAddressName() }}
              </h4>
              <div class="d-flex align-center mb-1">
                <v-icon
                  size="16"
                  class="mr-2"
                  color="grey"
                  :icon="mdiPhone"
                />
                <span class="text-body-2">{{ getFulfillmentAddressPhoneNumber() }}</span>
              </div>
              <div class="d-flex align-start mt-2">
                <v-icon
                  size="16"
                  class="mr-2 mt-1"
                  color="grey"
                  :icon="mdiMapMarker"
                />
                <div class="text-body-2">
                  <div>〒{{ props.order.fulfillments[0]?.postalCode || '' }}</div>
                  <div>{{ `${getFulfillmentAddressPrefecture()} ${props.order.fulfillments[0]?.city || ''}` }}</div>
                  <div>{{ props.order.fulfillments[0]?.addressLine1 || '' }}</div>
                  <div v-if="props.order.fulfillments[0]?.addressLine2">
                    {{ props.order.fulfillments[0]?.addressLine2 }}
                  </div>
                </div>
              </div>
            </template>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Fulfillment Details -->
    <v-row
      v-for="(fulfillment, index) in props.order.fulfillments"
      :key="fulfillment.fulfillmentId"
      class="mb-4"
    >
      <v-col cols="12">
        <v-card elevation="2">
          <v-card-title class="bg-grey-lighten-4 py-4">
            <v-icon
              class="mr-2"
              color="primary"
              :icon="mdiPackageVariantClosed"
            />
            配送詳細 {{ index + 1 }}
            <v-spacer />
            <v-chip
              size="small"
              :color="getFulfillmentStatusColor(fulfillment.status)"
            >
              {{ getFulfillmentStatus(fulfillment.status) }}
            </v-chip>
          </v-card-title>
          <v-card-text class="pa-6">
            <!-- Store Pickup Layout -->
            <template v-if="fulfillment.shippingType === ShippingType.ShippingTypePickup">
              <v-row>
                <v-col cols="12">
                  <!-- Pickup Information Highlight -->
                  <v-card
                    variant="tonal"
                    color="success"
                    class="mb-4"
                  >
                    <v-card-title class="pb-2">
                      <v-icon
                        class="mr-2"
                        size="24"
                        :icon="mdiStoreClock"
                      />
                      受け取り予定
                    </v-card-title>
                    <v-card-text>
                      <v-row>
                        <v-col
                          cols="12"
                          md="6"
                        >
                          <div class="mb-4">
                            <p class="text-subtitle-2 text-grey mb-1">
                              受け取り場所
                            </p>
                            <p class="text-h6 font-weight-bold">
                              {{ getPickupLocation() }}
                            </p>
                          </div>
                        </v-col>
                        <v-col
                          cols="12"
                          md="6"
                        >
                          <div class="mb-4">
                            <p class="text-subtitle-2 text-grey mb-1">
                              受け取り日時
                            </p>
                            <p class="text-h6 font-weight-bold">
                              {{ getPickupDate() }}
                            </p>
                          </div>
                        </v-col>
                      </v-row>
                    </v-card-text>
                  </v-card>

                  <!-- Basic Info (Compact) -->
                  <v-row class="mb-4">
                    <v-col
                      cols="12"
                      sm="4"
                    >
                      <div class="text-center">
                        <p class="text-subtitle-2 text-grey mb-1">
                          配送タイプ
                        </p>
                        <v-chip
                          color="success"
                          variant="tonal"
                          size="large"
                        >
                          {{ getShippingType(fulfillment.shippingType) }}
                        </v-chip>
                      </div>
                    </v-col>
                    <v-col
                      cols="12"
                      sm="4"
                    >
                      <div class="text-center">
                        <p class="text-subtitle-2 text-grey mb-1">
                          箱サイズ
                        </p>
                        <v-chip
                          color="info"
                          variant="tonal"
                          size="large"
                        >
                          {{ getBoxSize(fulfillment.boxSize) }}サイズ
                        </v-chip>
                      </div>
                    </v-col>
                    <v-col
                      cols="12"
                      sm="4"
                    >
                      <div class="text-center">
                        <p class="text-subtitle-2 text-grey mb-1">
                          占有率
                        </p>
                        <v-chip
                          color="warning"
                          variant="tonal"
                          size="large"
                        >
                          {{ fulfillment.boxRate }}%
                        </v-chip>
                      </div>
                    </v-col>
                  </v-row>

                  <!-- Products in this fulfillment -->
                  <div v-if="getOrderItems(fulfillment.fulfillmentId).length > 0">
                    <h4 class="text-subtitle-1 font-weight-medium mb-3">
                      受け取り予定商品
                    </h4>
                    <v-data-table
                      :headers="productHeaders"
                      :items="getOrderItems(fulfillment.fulfillmentId)"
                      hide-default-footer
                      class="mb-4"
                    >
                      <template #[`item.media`]="{ item }">
                        <v-avatar
                          size="48"
                          rounded="lg"
                          class="ma-1"
                        >
                          <v-img
                            :src="getThumbnail(item.productId)"
                            :srcset="getResizedThumbnails(item.productId)"
                            aspect-ratio="1"
                          />
                        </v-avatar>
                      </template>
                      <template #[`item.name`]="{ item }">
                        <div class="font-weight-medium">
                          {{ getProductName(item.productId) }}
                        </div>
                      </template>
                      <template #[`item.price`]="{ item }">
                        <span class="font-weight-medium">¥{{ item.price.toLocaleString() }}</span>
                      </template>
                      <template #[`item.quantity`]="{ item }">
                        <v-chip
                          size="small"
                          color="info"
                          variant="tonal"
                        >
                          {{ item.quantity.toLocaleString() }}
                        </v-chip>
                      </template>
                      <template #[`item.total`]="{ item }">
                        <span class="font-weight-bold text-primary">¥{{ getSubtotal(item).toLocaleString() }}</span>
                      </template>
                    </v-data-table>
                  </div>
                </v-col>
              </v-row>
            </template>

            <!-- Regular Shipping Layout -->
            <template v-else>
              <v-row>
                <v-col
                  cols="12"
                  lg="8"
                >
                  <!-- Shipping Info -->
                  <v-row class="mb-4">
                    <v-col
                      cols="12"
                      sm="6"
                      md="3"
                    >
                      <div class="mb-4">
                        <p class="text-subtitle-2 text-grey mb-1">
                          配送タイプ
                        </p>
                        <p class="text-body-1 font-weight-medium">
                          {{ getShippingType(fulfillment.shippingType) }}
                        </p>
                      </div>
                    </v-col>
                    <v-col
                      cols="12"
                      sm="6"
                      md="3"
                    >
                      <div class="mb-4">
                        <p class="text-subtitle-2 text-grey mb-1">
                          箱サイズ
                        </p>
                        <p class="text-body-1 font-weight-medium">
                          {{ getBoxSize(fulfillment.boxSize) }}サイズ
                        </p>
                      </div>
                    </v-col>
                    <v-col
                      cols="12"
                      sm="6"
                      md="3"
                    >
                      <div class="mb-4">
                        <p class="text-subtitle-2 text-grey mb-1">
                          占有率
                        </p>
                        <p class="text-body-1 font-weight-medium">
                          {{ fulfillment.boxRate }}%
                        </p>
                      </div>
                    </v-col>
                    <v-col
                      cols="12"
                      sm="6"
                      md="3"
                    >
                      <div class="mb-4">
                        <p class="text-subtitle-2 text-grey mb-1">
                          配送希望日
                        </p>
                        <p class="text-body-1 font-weight-medium">
                          {{ getRequestDaliveryDay(fulfillment) }}
                        </p>
                      </div>
                    </v-col>
                  </v-row>

                  <!-- Products in this fulfillment -->
                  <div v-if="getOrderItems(fulfillment.fulfillmentId).length > 0">
                    <h4 class="text-subtitle-1 font-weight-medium mb-3">
                      この配送に含まれる商品
                    </h4>
                    <v-data-table
                      :headers="productHeaders"
                      :items="getOrderItems(fulfillment.fulfillmentId)"
                      hide-default-footer
                      class="mb-4"
                    >
                      <template #[`item.media`]="{ item }">
                        <v-avatar
                          size="48"
                          rounded="lg"
                          class="ma-1"
                        >
                          <v-img
                            :src="getThumbnail(item.productId)"
                            :srcset="getResizedThumbnails(item.productId)"
                            aspect-ratio="1"
                          />
                        </v-avatar>
                      </template>
                      <template #[`item.name`]="{ item }">
                        <div class="font-weight-medium">
                          {{ getProductName(item.productId) }}
                        </div>
                      </template>
                      <template #[`item.price`]="{ item }">
                        <span class="font-weight-medium">¥{{ item.price.toLocaleString() }}</span>
                      </template>
                      <template #[`item.quantity`]="{ item }">
                        <v-chip
                          size="small"
                          color="info"
                          variant="tonal"
                        >
                          {{ item.quantity.toLocaleString() }}
                        </v-chip>
                      </template>
                      <template #[`item.total`]="{ item }">
                        <span class="font-weight-bold text-primary">¥{{ getSubtotal(item).toLocaleString() }}</span>
                      </template>
                    </v-data-table>
                  </div>
                </v-col>

                <!-- Tracking Information -->
                <v-col
                  cols="12"
                  lg="4"
                >
                  <v-card
                    variant="tonal"
                    color="info"
                  >
                    <v-card-title class="text-h6 pb-2">
                      配送追跡
                    </v-card-title>
                    <v-card-text>
                      <v-select
                        v-model="fulfillmentsFormDataValue[index].shippingCarrier"
                        label="配送業者"
                        :items="shippingCarriers"
                        variant="outlined"
                        density="comfortable"
                        :readonly="!isUpdatableFulfillment()"
                        class="mb-3"
                      />
                      <v-text-field
                        v-model="fulfillmentsFormDataValue[index].trackingNumber"
                        label="追跡番号"
                        variant="outlined"
                        density="comfortable"
                        :readonly="!isUpdatableFulfillment()"
                        class="mb-3"
                      />
                      <v-btn
                        v-show="isUpdatableFulfillment()"
                        :loading="loading"
                        variant="flat"
                        color="info"
                        size="large"
                        block
                        @click="onSubmitUpdate(fulfillment.fulfillmentId)"
                      >
                        <v-icon
                          start
                          :icon="mdiUpdate"
                        />
                        追跡情報を更新
                      </v-btn>
                    </v-card-text>
                  </v-card>
                </v-col>
              </v-row>
            </template>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Shipping Message -->
    <v-row v-show="showShippingMessage()">
      <v-col cols="12">
        <v-card elevation="2">
          <v-card-title class="bg-grey-lighten-4 py-4">
            <v-icon
              class="mr-2"
              color="primary"
              :icon="mdiSend"
            />
            発送完了メッセージ
          </v-card-title>
          <v-card-text class="pa-6">
            <v-textarea
              v-model="completeFormDataValue.shippingMessage"
              label="お客様へのメッセージ"
              placeholder="例：ご注文ありがとうございます！商品の発送が完了しました。商品到着まで今しばらくお待ち下さい。"
              variant="outlined"
              rows="4"
              :readonly="!isPreservable()"
            />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Dialogs -->
    <v-dialog
      v-model="cancelDialogValue"
      max-width="500"
    >
      <v-card>
        <v-card-title class="bg-error text-white py-4">
          <v-icon
            class="mr-2"
            color="white"
            :icon="mdiAlertCircle"
          />
          注文キャンセル確認
        </v-card-title>
        <v-card-text class="py-6">
          <p class="text-body-1 mb-0">
            本当にこの注文をキャンセルしますか？
          </p>
          <p class="text-body-2 text-grey mt-2">
            この操作は取り消すことができません。
          </p>
        </v-card-text>
        <v-card-actions class="pa-4">
          <v-spacer />
          <v-btn
            color="grey"
            variant="text"
            @click="onClickCloseCancelDialog"
          >
            キャンセル
          </v-btn>
          <v-btn
            :loading="loading"
            color="error"
            variant="flat"
            @click="onSubmitCancel"
          >
            注文をキャンセル
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog
      v-model="refundDialogValue"
      max-width="600"
    >
      <v-card>
        <v-card-title class="bg-warning text-white py-4">
          <v-icon
            class="mr-2"
            color="white"
            :icon="mdiCreditCardRefund"
          />
          返金処理
        </v-card-title>
        <v-card-text class="py-6">
          <p class="text-body-1 mb-4">
            返金理由を入力してください。
          </p>
          <v-textarea
            v-model="refundFormDataValue.description"
            label="返金理由"
            placeholder="返金が必要な理由を詳しく入力してください"
            variant="outlined"
            rows="4"
            maxlength="200"
            counter
          />
        </v-card-text>
        <v-card-actions class="pa-4">
          <v-spacer />
          <v-btn
            color="grey"
            variant="text"
            @click="onClickCloseRefundDialog"
          >
            キャンセル
          </v-btn>
          <v-btn
            :loading="loading"
            color="warning"
            variant="flat"
            :disabled="!refundFormDataValue.description"
            @click="onSubmitRefund"
          >
            返金処理を実行
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Fixed Footer Actions -->
    <v-footer
      app
      color="white"
      elevation="8"
      class="px-6 py-4 fixed-footer-actions"
    >
      <v-container
        fluid
        class="pa-0"
      >
        <div class="d-flex align-center justify-center flex-wrap ga-3">
          <!-- Primary Actions -->
          <v-btn
            v-show="isAuthorized()"
            :loading="loading"
            variant="flat"
            color="success"
            size="large"
            @click="onSubmitCapture()"
          >
            <v-icon
              start
              :icon="mdiCheckCircle"
            />
            注文を確定
          </v-btn>
          <v-btn
            v-show="isCompletable()"
            :loading="loading"
            variant="flat"
            color="primary"
            size="large"
            @click="onSubmitComplete()"
          >
            <v-icon
              start
              :icon="mdiSend"
            />
            {{ props.order.type === OrderType.OrderTypeProduct ? '発送完了を通知' : 'レビュー依頼を送信' }}
          </v-btn>
          <v-btn
            v-show="isPreservable()"
            :loading="loading"
            variant="outlined"
            color="info"
            size="large"
            @click="onSubmitSaveDraft()"
          >
            <v-icon
              start
              :icon="mdiContentSave"
            />
            下書きを保存
          </v-btn>
          <!-- Dangerous Actions -->
          <template v-if="isCancelable() || isRefundable()">
            <v-divider
              vertical
              class="mx-2"
            />
            <v-btn
              v-show="isCancelable()"
              :loading="loading"
              variant="outlined"
              color="error"
              size="large"
              @click="onClickOpenCancelDialog()"
            >
              <v-icon
                start
                :icon="mdiCancel"
              />
              キャンセル
            </v-btn>
            <v-btn
              v-show="isRefundable()"
              :loading="loading"
              variant="outlined"
              color="warning"
              size="large"
              @click="onClickOpenRefundDialog"
            >
              <v-icon
                start
                :icon="mdiCreditCardRefund"
              />
              返金処理
            </v-btn>
          </template>
        </div>
      </v-container>
    </v-footer>
  </v-container>
</template>

<style scoped>
.fixed-footer-actions {
  position: fixed !important;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  border-top: 1px solid rgb(0 0 0 / 12%);
}

.fixed-footer-actions .v-container {
  max-width: none;
}

@media (width <= 1024px) {
  .fixed-footer-actions .d-flex {
    flex-direction: column;
    gap: 8px;
  }

  .fixed-footer-actions .v-btn {
    width: 100%;
  }

  .fixed-footer-actions .v-divider {
    display: none;
  }
}

@media (width <= 640px) {
  .fixed-footer-actions {
    padding: 12px 16px;
  }

  .fixed-footer-actions .v-btn {
    font-size: 14px;
  }
}
</style>
