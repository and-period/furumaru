<script setup lang="ts">
import { useAuthStore } from '~/store/auth'
import type { UpdateAuthPasswordRequest } from '~/types/api'
import { ApiBaseError } from '~/types/exception'

const router = useRouter()
const { user, updatePassword } = useAuthStore()

const errorMessage = ref<string>('')
const formData = ref<UpdateAuthPasswordRequest>({
  oldPassword: '',
  newPassword: '',
  passwordConfirmation: '',
})

const handleSubmit = async () => {
  try {
    if (formData.value.newPassword !== formData.value.passwordConfirmation) {
      alert('新しいパスワードと新しいパスワード（確認用）が一致しません')
    } else {
      await updatePassword(formData.value)
      router.push('/account/edit/complete?from=password')
    }
  } catch (error) {
    if (error instanceof ApiBaseError) {
      errorMessage.value = error.message
      return
    }
    console.error(error)
  }
}
</script>

<template>
  <div class="container mx-auto p-4 md:p-0">
    <template v-if="user">
      <the-account-edit-card title="パスワードの変更" class="mt-6">
        <the-alert v-if="errorMessage">{{ errorMessage }}</the-alert>

        <form class="flex w-full flex-col gap-6" @submit.prevent="handleSubmit">
          <div class="mt-10 flex w-full flex-col justify-center gap-4">
            <the-text-input
              v-model="formData.oldPassword"
              label="現在のパスワード"
              class="w-full"
              type="password"
              placeholder="現在のパスワード"
              required
            />

            <the-text-input
              v-model="formData.newPassword"
              label="新しいパスワード"
              class="w-full"
              type="password"
              placeholder="新しいパスワード"
              required
            />

            <the-text-input
              v-model="formData.passwordConfirmation"
              label="新しいパスワード（確認用）"
              class="w-full"
              type="password"
              placeholder="新しいパスワード（確認用）"
              required
            />
          </div>

          <div class="my-4 text-center">
            <nuxt-link to="/account/edit/id" class="underline">
              パスワードをお忘れの場合
            </nuxt-link>
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
