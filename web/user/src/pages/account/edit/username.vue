<script setup lang="ts">
import { useAuthStore } from '~/store/auth'

const { user } = useAuthStore()

const formData = ref<string>('')

if (user) {
  formData.value = user.username
}
</script>

<template>
  <div class="container mx-auto p-4 md:p-0">
    <template v-if="user">
      <the-account-edit-card title="ユーザー名の変更" class="mt-6">
        <div class="flex w-full flex-col gap-6">
          <div class="my-10 flex w-full justify-center">
            <the-text-input
              v-model="formData"
              label="ユーザー名"
              class="w-full"
              type="text"
              placeholder="ユーザー名"
            />
          </div>

          <div class="flex w-full flex-col gap-4">
            <button class="w-ful bg-main px-4 py-2 text-white">変更する</button>
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
        </div>
      </the-account-edit-card>
    </template>

    <template v-else>
      <the-auth-error-card />
    </template>
  </div>
</template>
