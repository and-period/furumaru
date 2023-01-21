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
          </v-card-text>
        </v-card>
        <v-tab-item value="tab-orderInformation">
          <v-card elevation="0"> </v-card>
        </v-tab-item>
      </v-tab-item>
    </v-tabs-items>
  </div>
</template>

<script lang="ts">
import { ref, useFetch, useRoute } from '@nuxtjs/composition-api'
import { defineComponent, reactive } from '@vue/composition-api'

import { useOrderStore } from '~/store/orders'
import { OrderResponse, PaymentMethodType, PaymentStatus } from '~/types/api'
import { Order } from '~/types/props/order'

export default defineComponent({
  setup() {
    const route = useRoute()
    const orderStore = useOrderStore()
    const id = route.value.params.id

    const selector = ref<string>('shippingInformation')

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

    const { fetchState } = useFetch(async () => {
      const res = await orderStore.getOrder(id)
      formData.userName = res.userName
      formData.payment = res.payment
    })

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

    const convertPhone = (phoneNumber: string): string => {
      return phoneNumber.replace('+81', '0')
    }

    return {
      items,
      formData,
      selector,
      fetchState,
      getMethodType,
      getPaymentStatus,
      getPaymentStatusColor,
      convertPhone,
    }
  },
})
</script>
