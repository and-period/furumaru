<script setup lang="ts">
import { useAuthStore } from '~/store/auth'
import type { SignInRequest } from '~/types/api'
import { ApiBaseError } from '~/types/exception'

const authStore = useAuthStore()
const { signIn } = authStore

const router = useRouter()
const route = useRoute()

const loginRequired = computed<boolean>(() => {
  const required = route.query.required
  if (required) {
    return Boolean(required)
  } else {
    return false
  }
})

const fromNewAccounet = computed<boolean>(() => {
  const fromNewAccounetParam = route.query.from_new_accounet
  if (fromNewAccounetParam) {
    return Boolean(fromNewAccounetParam)
  } else {
    return false
  }
})

const coordinatorId = computed<string>(() => {
  const id = route.query.coordinatorId
  if (id) {
    return String(id)
  } else {
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
  } catch (error) {
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
  <div v-if="loginRequired" class="px-4">
    <the-alert
      class="mx-auto my-4 w-full bg-white p-4 lg:w-[768px] xl:w-[1024px]"
    >
      ご購入にはログインが必須です。
    </the-alert>
  </div>

  <div v-if="fromNewAccounet" class="px-4">
    <the-alert
      class="mx-auto my-4 w-full bg-white p-4 lg:w-[768px] xl:w-[1024px]"
    >
      作成したアカウントでログインをしましょう
    </the-alert>
  </div>

  <div
    class="container mx-auto mt-[40px] flex flex-col gap-10 px-4 md:grid md:grid-cols-2 md:gap-0 md:px-0"
  >
    <div
      class="w-full bg-white px-4 py-[40px] tracking-[1.6px] text-main md:mx-auto md:w-[360px] md:px-[40px] lg:w-[480px] xl:w-[560px] xl:px-[80px]"
    >
      <h2 class="mb-[40px] text-center text-[16px] font-bold">
        アカウントをお持ちの方
      </h2>
      <the-alert v-if="authErrorState.hasError" class="mb-4">{{
        authErrorState.errorMessage
      }}</the-alert>
      <the-sign-in-form
        v-model="formData"
        button-text="ログインして購入"
        username-label="メールアドレス"
        password-label="パスワード"
        password-placeholder="パスワード"
        username-placeholder="メールアドレス"
        @submit="handleSubmitSignForm"
      />
      <div v-if="hideV1App" class="mt-[24px] text-center text-[14px] underline">
        パスワードをお忘れの方はこちら
      </div>
    </div>

    <div
      class="w-full bg-white px-4 py-[40px] tracking-[1.6px] text-main md:mx-auto md:w-[360px] md:px-[40px] lg:w-[480px] xl:w-[560px] xl:px-[80px]"
    >
      <h2 class="mb-[40px] text-center text-[16px] font-bold">
        まだ登録されていない方
      </h2>
      <p class="mb-[40px] text-[14px] tracking-[1.4px]">
        アカウントをご登録いただくと毎回のお届け先等の情報入力が不要になり、お買い物がもっと便利になります。
      </p>
      <the-submit-button type="button" @click="handleClickNewAccountButton">
        新規登録する
      </the-submit-button>
      <p class="my-[40px] text-[14px] tracking-[1.4px]">
        アカウントを登録せずにご購入を希望される方はこちらからご利用ください。
      </p>
      <the-submit-without-login-button type="submit" @click="handleSubmitWithoutSignForm" >
        ログインせずに購入
      </the-submit-without-login-button>
    </div>
  </div>
</template>
