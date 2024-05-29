<script setup lang="ts">
import { useAuthStore } from '~/store/auth'
import { ApiBaseError } from '~/types/exception'

const { user, updateNotificationEnabled } = useAuthStore()

const router = useRouter()
const formData = ref<boolean>(false)

const errorMessage = ref<string>('')

if (user) {
  formData.value = user.notificationEnabled
}

const handleSubmit = async () => {
  try {
    errorMessage.value = ''
    await updateNotificationEnabled(formData.value)
    router.push('/account/edit/complete?from=notification')
  }
  catch (error) {
    if (error instanceof ApiBaseError) {
      errorMessage.value = error.message
    }
    else {
      errorMessage.value = 'メール受信設定の変更に失敗しました。'
    }
  }
}

useSeoMeta({
  title: 'メール受信の変更',
})
</script>

<template>
  <div class="container mx-auto p-4 md:p-0">
    <template v-if="user">
      <the-account-edit-card
        title="メール受信の変更"
        class="mt-6"
      >
        <the-alert
          v-if="errorMessage"
          class="mt-4 w-full"
        >
          {{ errorMessage }}
        </the-alert>

        <form
          class="flex w-full flex-col gap-6"
          @submit.prevent="handleSubmit"
        >
          <div class="my-10 flex w-full flex-col justify-center gap-14">
            <div>
              <p class="mb-4">
                メールアドレス
              </p>
              <p>{{ user.email }}</p>
            </div>

            <div class="flex flex-col gap-4">
              <div class="flex items-center gap-2">
                <input
                  id="true"
                  v-model="formData"
                  type="radio"
                  name="notification"
                  class="h-4 w-4 accent-main"
                  :value="true"
                >
                <label for="true"> お知らせメールを受け取る </label>
              </div>

              <div class="flex items-center gap-2">
                <input
                  id="false"
                  v-model="formData"
                  type="radio"
                  name="notification"
                  class="h-4 w-4 accent-main"
                  :value="false"
                >
                <label for="false"> メール配信を停止する </label>
              </div>
            </div>
          </div>

          <div class="flex w-full flex-col gap-4">
            <button class="w-full bg-main px-4 py-2 text-white">
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
