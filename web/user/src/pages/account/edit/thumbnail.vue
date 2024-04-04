<script setup lang="ts">
import { useAuthStore } from '~/store/auth'

const { user } = useAuthStore()
</script>

<template>
  <div class="container mx-auto p-4 md:p-0">
    <template v-if="user">
      <the-account-edit-card title="プロフィール写真の変更" class="mt-6">
        <div class="flex w-full flex-col gap-6">
          <div class="mt-10 flex justify-center">
            <img
              v-if="user.thumbnailUrl"
              :src="user.thumbnailUrl"
              alt="サムネイル"
            />
            <the-account-icon
              v-else
              color="white"
              class="h-[120px] w-[120px]"
            />
          </div>

          <div class="text-center underline">
            {{
              user.thumbnailUrl
                ? 'プロフィール写真を変更'
                : 'プロフィール写真を追加'
            }}
          </div>

          <div class="flex w-full flex-col gap-4">
            <button class="w-ful bg-main px-4 py-2 text-white">保存する</button>
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
