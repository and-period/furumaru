<script setup lang="ts">
import { useAuthStore } from '~/store/auth'

const router = useRouter()
const { user, updateAccountId } = useAuthStore()

const formData = ref<string>('')

if (user) {
  formData.value = user.accountId
}

const handleSubmit = async () => {
  if (formData.value === '') {
    return
  }
  await updateAccountId(formData.value)
  router.push('/account/edit/complete?from=accountId')
}

useSeoMeta({
  title: 'ユーザーIDの変更',
})
</script>

<template>
  <div class="container mx-auto p-4 md:p-0">
    <template v-if="user">
      <the-account-edit-card title="ユーザーIDの変更" class="mt-6">
        <form class="flex w-full flex-col gap-6" @submit.prevent="handleSubmit">
          <div class="my-10 flex w-full flex-col justify-center gap-2">
            <the-text-input
              v-model="formData"
              required
              label="ユーザーID"
              class="w-full"
              type="text"
              placeholder="ユーザーID"
              :max-length="32"
            />
            <p>4文字以上32文字以下の必要があります。</p>
            <p>
              使用可能な文字列は半角英数字とハイフン（-）、アンダースコア（_）です。
            </p>
          </div>

          <div class="flex w-full flex-col gap-4">
            <button class="w-ful bg-main px-4 py-2 text-white" type="submit">
              変更する
            </button>
            <nuxt-link
              class="border border-main px-4 py-2 text-center text-main"
              to="/account"
            >
              キャンセル
            </nuxt-link>
          </div>
          <div class="text-center">
            <nuxt-link
              to="/account"
              class="inline-flex items-center gap-1 text-[12px]"
            >
              <TheLeftArrowIcon class="h-3 w-3" />
              マイページに戻る
            </nuxt-link>
          </div>
        </form>
      </the-account-edit-card>
    </template>

    <template v-else>
      <the-auth-error-card />
    </template>
  </div>
</template>
