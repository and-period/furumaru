<script setup lang="ts">
import { useAuthStore } from '~/store/auth'
import { ApiBaseError } from '~/types/exception'

const router = useRouter()
const { user, updateEmail } = useAuthStore()

const formData = ref<string>('')
const errorMessage = ref<string>('')

const handleSubmit = async () => {
  if (formData.value === '') {
    return
  }
  try {
    await updateEmail(formData.value)
    router.push('/account/edit/complete?from=email')
  } catch (error) {
    if (error instanceof ApiBaseError) {
      errorMessage.value = error.message
    } else {
      errorMessage.value = 'メールアドレスの変更に失敗しました。'
    }
  }
}

useSeoMeta({
  title: 'メールアドレスの変更',
})
</script>

<template>
  <div class="container mx-auto p-4 md:p-0">
    <template v-if="user">
      <the-account-edit-card title="メールアドレスの変更" class="mt-6">
        <!-- エラーメッセージ -->
        <the-alert v-if="errorMessage" class="mt-4 w-full">
          {{ errorMessage }}
        </the-alert>

        <the-alert class="mt-4 w-full">
          この機能は現在準備中で使用できません。
        </the-alert>

        <form class="flex w-full flex-col gap-6" @submit.prevent="handleSubmit">
          <div class="my-10 flex w-full flex-col justify-center gap-14">
            <div>
              <p class="mb-4">現在のメールアドレス</p>
              <p>{{ user.email }}</p>
            </div>

            <the-text-input
              v-model="formData"
              label="新しいメールアドレス"
              class="w-full"
              type="email"
              placeholder="新しいメールアドレス"
              required
            />
          </div>

          <div class="flex w-full flex-col gap-4">
            <button
              class="w-ful bg-main px-4 py-2 text-white disabled:bg-main/70"
              type="submit"
              disabled
            >
              変更する
            </button>
            <nuxt-link
              class="border border-main px-4 py-2 text-center text-main"
              to="/account"
            >
              キャンセル
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
