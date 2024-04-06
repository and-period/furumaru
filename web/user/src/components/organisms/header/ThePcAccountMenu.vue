<script lang="ts" setup>
import type { AuthUserResponse } from '~/types/api'
import type { LinkItem } from '~/types/props'

interface Props {
  isAuthenticated: boolean
  user: AuthUserResponse | undefined
  authenticatedMenuItems: LinkItem[]
  noAuthenticatedMenuItems: LinkItem[]
}

defineProps<Props>()

interface Emits {
  (e: 'click:myPageButton'): void
  (e: 'click:logoutButton'): void
}

const emits = defineEmits<Emits>()

const dropdownRef = ref<{ close: () => void } | undefined>(undefined)

const handleClickMyPageButton = () => {
  emits('click:myPageButton')
  if (dropdownRef.value) {
    dropdownRef.value.close()
  }
}

const handleClickLogoutButton = () => {
  emits('click:logoutButton')
  if (dropdownRef.value) {
    dropdownRef.value.close()
  }
}
</script>

<template>
  <the-dropdown-with-icon ref="dropdownRef">
    <template #icon>
      <the-account-icon />
    </template>
    <template #content>
      <div
        v-if="isAuthenticated && user"
        class="flex flex-col text-[14px] tracking-[1.4px]"
      >
        <p v-for="(item, i) in authenticatedMenuItems" :key="i">
          {{ item.text }}
        </p>
        <div class="mb-8 px-4">
          <template v-if="user.thumbnailUrl">
            <img
              :src="user.thumbnailUrl"
              class="block h-[40px] w-[40px] rounded-full"
            />
          </template>
          <template v-else>
            <img
              src="~/assets/img/account.png"
              class="block h-[40px] w-[40px] rounded-full"
            />
          </template>
          <div class="mt-4 font-bold">
            <template v-if="user.username"> {{ user.username }} 様 </template>
            <template v-else> ユーザー名が未設定です </template>
          </div>
        </div>

        <button
          class="px-4 py-2 text-left underline hover:bg-gray-200"
          @click="handleClickMyPageButton"
        >
          マイページ
        </button>

        <button
          class="px-4 py-2 text-left underline hover:bg-gray-200"
          @click="handleClickLogoutButton"
        >
          ログアウト
        </button>
      </div>
      <div v-else class="flex flex-col">
        <nuxt-link
          v-for="(item, i) in noAuthenticatedMenuItems"
          :key="i"
          :href="item.href"
          class="px-4 py-2 hover:bg-gray-200"
        >
          {{ item.text }}
        </nuxt-link>
      </div>
    </template>
  </the-dropdown-with-icon>
</template>
