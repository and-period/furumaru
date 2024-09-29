<script setup lang="ts">
import { useAuthStore } from '~/store/auth'
import type { SignInRequest } from '~/types/api'
import { ApiBaseError } from '~/types/exception'
import type { I18n } from '~/types/locales'

const authStore = useAuthStore()
const { signIn } = authStore

const i18n = useI18n()

const router = useRouter()
const route = useRoute()

const at = (str: keyof I18n['purchase']['auth']) => {
  return i18n.t(`purchase.auth.${str}`)
}

const loginRequired = computed<boolean>(() => {
  const required = route.query.required
  if (required) {
    return Boolean(required)
  }
  else {
    return false
  }
})

const fromNewAccounet = computed<boolean>(() => {
  const fromNewAccounetParam = route.query.from_new_accounet
  if (fromNewAccounetParam) {
    return Boolean(fromNewAccounetParam)
  }
  else {
    return false
  }
})

const coordinatorId = computed<string>(() => {
  const id = route.query.coordinatorId
  if (id) {
    return String(id)
  }
  else {
    return ''
  }
})

const cartNumber = computed<number | undefined>(() => {
  const id = route.query.cartNumber
  const idNumber = Number(id)
  if (idNumber === 0) {
    return undefined
  }
  if (isNaN(idNumber)) {
    return undefined
  }
  return idNumber
})

const authErrorState = ref({
  hasError: false,
  errorMessage: '',
})

const formData = ref<SignInRequest>({
  username: '',
  password: '',
})

const handleClickNewAccountButton = () => {
  router.push({
    path: '/signup',
    query: {
      redirect_to_purchase: true,
      coordinatorId: coordinatorId.value,
      cartNumber: cartNumber.value,
    },
  })
}

const handleSubmitSignForm = async () => {
  try {
    await signIn(formData.value)
    router.push({
      path: '/v1/purchase/address',
      query: {
        coordinatorId: coordinatorId.value,
        cartNumber: cartNumber.value,
      },
    })
  }
  catch (error) {
    authErrorState.value.hasError = true
    if (error instanceof ApiBaseError) {
      authErrorState.value.errorMessage = error.message
    }
  }
}

const handleSubmitWithoutSignForm = () => {
  router.push({
    path: '/v1/purchase/guest/address',
    query: {
      coordinatorId: coordinatorId.value,
      cartNumber: cartNumber.value,
    },
  })
}

useSeoMeta({
  title: 'ログイン',
})

const hideV1App = false
</script>

<template>
  <!-- 購入導線改善の検証のためにコメントアウト
  <div
    v-if="loginRequired"
    class="px-4"
  >
    <the-alert
      class="mx-auto my-4 w-full bg-white p-4 lg:w-[768px] xl:w-[1024px]"
    >
      {{ at('loginRequiredMessage') }}
    </the-alert>
  </div>
  -->

  <div
    v-if="fromNewAccounet"
    class="px-4"
  >
    <the-alert
      class="mx-auto my-4 w-full bg-white p-4 lg:w-[768px] xl:w-[1024px]"
    >
      {{ at('loginNewAccountMessage') }}
    </the-alert>
  </div>

  <div
    class="container mx-auto mt-[40px] flex flex-col gap-10 px-4 md:grid md:grid-cols-2 md:gap-0 md:px-0"
  >
    <div
      class="w-full bg-white px-4 py-[40px] tracking-[1.6px] text-main md:mx-auto md:w-[360px] md:px-[40px] lg:w-[480px] xl:w-[560px] xl:px-[80px]"
    >
      <h2 class="mb-[40px] text-center text-[16px] font-bold">
        {{ at('notSignUpTitle') }}
      </h2>
      <p class="my-[40px] text-[14px] tracking-[1.4px]">
        {{ at('checkoutWithoutAccountDescription') }}
      </p>
      <the-submit-without-login-button
        type="submit"
        @click="handleSubmitWithoutSignForm"
      >
        {{ at('checkoutWithoutAccountButtonText') }}
      </the-submit-without-login-button>
    </div>
    <div
      class="w-full bg-white px-4 py-[40px] tracking-[1.6px] text-main md:mx-auto md:w-[360px] md:px-[40px] lg:w-[480px] xl:w-[560px] xl:px-[80px]"
    >
      <h2 class="mb-[40px] text-center text-[16px] font-bold">
        {{ at('withAccountTitle') }}
      </h2>
      <the-alert
        v-if="authErrorState.hasError"
        class="mb-4"
      >
        {{
          authErrorState.errorMessage
        }}
      </the-alert>
      <the-sign-in-form
        v-model="formData"
        :button-text="at('loginAndCheckoutButtonText')"
        :username-label="at('usernameLabel')"
        :password-label="at('passwordLabel')"
        :password-placeholder="at('passwordPlaceholder')"
        :username-placeholder="at('usernamePlaceholder')"
        @submit="handleSubmitSignForm"
      />
      <button @click="handleClickNewAccountButton">
        <p
          class="mt-4 inline-block whitespace-pre-wrap text-[14px] font-bold underline md:text-[15px]"
        >
          {{ at('noAccountButtonText') }}
        </p>
      </button>
      <div
        v-if="hideV1App"
        class="mt-[24px] text-center text-[14px] underline"
      >
        {{ at('forgetPasswordLink') }}
      </div>
    </div>
  </div>
</template>
