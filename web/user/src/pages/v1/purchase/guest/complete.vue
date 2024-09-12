<script setup lang="ts">
import { getOperationResultFromOrderStatus } from '~/lib/order'
import { useCheckoutStore } from '~/store/checkout'
import type { GuestCheckoutStateResponse } from '~/types/api'
import type { I18n } from '~/types/locales'
import { ApiBaseError } from '~/types/exception'

const i18n = useI18n()
const config = useRuntimeConfig()
const router = useRouter()
const route = useRoute()

const checkoutStore = useCheckoutStore()
const { guestCheckTransactionStatus } = checkoutStore

const ct = (str: keyof I18n['purchase']['complete']) => {
  return i18n.t(`purchase.complete.${str}`)
}

const orderIDMessage = computed<string>(() => {
  return i18n.t('purchase.complete.orderIDMessage', {
    orderId: checkoutStatus?.value?.orderId,
  })
})

const isLoading = ref<boolean>(true)
const hasError = ref<boolean>(false)
const errorMessage = ref<string>('')

const sessionId = computed<string>(() => {
  const id = route.query.session_id
  console.log('debug', 'sessionId', id)
  if (id) {
    return String(id)
  }
  else {
    return ''
  }
})

const checkoutStatus = ref<GuestCheckoutStateResponse | undefined>(undefined)

const operationResult = computed<string>(() => {
  if (checkoutStatus.value) {
    return getOperationResultFromOrderStatus(checkoutStatus.value.status)
  }
  else {
    return 'unknown'
  }
})

const handleBackTopPageButton = () => {
  router.push('/')
}

const handleBackCartPageButton = () => {
  router.push('/purchase')
}

onMounted(async () => {
  console.log('debug', 'onMounted', sessionId.value)
  if (sessionId.value) {
    try {
      isLoading.value = true
      checkoutStatus.value = await guestCheckTransactionStatus(sessionId.value)
    }
    catch (error) {
      hasError.value = true
      if (error instanceof ApiBaseError) {
        errorMessage.value = error.message
        return
      }
      errorMessage.value = ''
    }
    finally {
      isLoading.value = false
    }
  }
})

useHead({
  script: [
    // 本番環境にだけ Meta Pixel Code を仕込む
    ...(config.public.ENVIRONMENT === 'prd'
      ? [
          {
            key: 'meta-picel-purchase',
            src: '/meta/pixel-code-purchase.js',
            defer: true,
            type: 'text/javascript',
          },
        ]
      : []),
  ],
})

useSeoMeta({
  title: '決済確認',
})
</script>

<template>
  <template v-if="isLoading">
    <div
      class="flex justify-center"
      aria-label="読み込み中"
    >
      <div
        class="h-20 w-20 animate-spin rounded-full border-4 border-main border-t-transparent"
      />
    </div>
  </template>

  <template v-else>
    <template v-if="hasError">
      <div
        class="text-oranges container mx-auto mb-4 border border-orange bg-white p-4 text-orange"
      >
        {{ errorMessage }}
      </div>
    </template>

    <template v-if="operationResult === 'success'">
      <div class="text-main">
        <div class="hidden md:block">
          <div class="flex items-center justify-center gap-x-[40px]">
            <div class="text-xl font-bold tracking-[2px]">
              <p>{{ ct('thanksMessageFirst') }}</p>
              <p>{{ ct('thanksMessageSecond') }}</p>
            </div>
            <img src="/img/purchase/complete.svg">
          </div>
        </div>
        <div class="mt-[40px] block md:hidden">
          <div class="flex justify-center">
            <img src="/img/purchase/complete.svg">
          </div>
          <div class="text-[16px] font-bold tracking-[1.6px]">
            <p class="mt-[40px] flex justify-center">
              {{ ct('thanksMessageFirst') }}
            </p>
            <p class="mt-2 flex justify-center">
              {{ ct('thanksMessageSecond') }}
            </p>
          </div>
        </div>

        <div
          class="mt-[40px] flex flex-col items-center justify-center md:mt-20"
        >
          <p
            class="text-center text-[14px] font-medium tracking-[2px] md:text-xl"
          >
            {{ orderIDMessage }}
          </p>
          <div
            class="flex flex-col items-center justify-center px-[15px] text-[14px] font-medium md:px-0 md:text-[16px]"
          >
            <p class="mt-[40px] md:mt-10">
              {{ ct('completeMessageFirst') }}
            </p>
            <p>
              {{ ct('completeMessageSecond') }}
            </p>
            <p>
              {{ ct('completeMessageThird') }}
            </p>
          </div>
          <div
            class="mb-4 mt-[40px] flex w-full justify-center px-16 md:mt-20 md:px-0"
          >
            <button
              class="w-[400px] bg-main py-2 tracking-[1.6px] text-white"
              @click="handleBackTopPageButton"
            >
              {{ ct('topLinkText') }}
            </button>
          </div>
        </div>
      </div>
    </template>

    <template v-if="operationResult === 'canceled'">
      <div
        class="mt-[40px] flex flex-col items-center justify-center gap-10 text-main"
      >
        <p class="text-[16px] tracking-[1.6px]">
          {{ ct('cancellMessage') }}
        </p>
        <button
          class="w-[400px] bg-main py-2 tracking-[1.6px] text-white"
          @click="handleBackCartPageButton"
        >
          {{ ct('cartLinkText') }}
        </button>
      </div>
    </template>
  </template>
</template>
