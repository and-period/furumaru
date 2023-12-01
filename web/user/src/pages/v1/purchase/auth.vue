<script setup lang="ts">
import { SignInRequest } from '~/types/api'

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

const formData = ref<SignInRequest>({
  username: '',
  password: '',
})

const handleClickNewAccountButton = () => {
  router.push('/signup')
}

useSeoMeta({
  title: 'ログイン',
})
</script>

<template>
  <div v-if="loginRequired" class="px-4">
    <the-alert
      class="mx-auto my-4 w-full bg-white p-4 lg:w-[768px] xl:w-[1024px]"
    >
      ご購入にはログインが必須です。
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
      <the-sign-in-form
        v-model="formData"
        button-text="ログインして購入"
        username-label="メールアドレス"
        password-label="パスワード"
        password-placeholder="パスワード"
        username-placeholder="メールアドレス"
      />
      <div class="mt-[24px] text-center text-[14px] underline">
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
        <br />
        アカウント登録について詳しくはこちら
      </p>
      <the-submit-button type="button" @click="handleClickNewAccountButton">
        新規登録する
      </the-submit-button>
    </div>
  </div>
</template>