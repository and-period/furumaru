<script setup lang="ts">
import { useAuthStore } from '~/store/auth'

const router = useRouter()
const { user, updateThumbnail } = useAuthStore()

const formData = ref<File | null>(null)

const previewUrl = computed(() => {
  if (formData.value === null) {
    return ''
  }
  return URL.createObjectURL(formData.value)
})

const handleFileChange = (e: Event) => {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    formData.value = file
  }
}

const handleSubmit = async () => {
  console.log('submit')
  if (formData.value) {
    await updateThumbnail(formData.value)
    router.push('/account/edit/complete?from=thumbnail')
  }
}

useSeoMeta({
  title: 'プロフィール写真の変更',
})
</script>

<template>
  <div class="container mx-auto p-4 md:p-0">
    <template v-if="user">
      <the-account-edit-card title="プロフィール写真の変更" class="mt-6">
        <form class="flex w-full flex-col gap-6" @submit.prevent="handleSubmit">
          <div class="mt-10 flex justify-center">
            <template v-if="previewUrl">
              <img
                :src="previewUrl"
                alt="サムネイル"
                class="h-[120px] w-[120px] rounded-full"
              />
            </template>
            <template v-else>
              <img
                v-if="user.thumbnailUrl"
                :src="user.thumbnailUrl"
                alt="サムネイル"
                class="h-[120px] w-[120px] rounded-full"
              />
              <the-account-icon
                v-else
                color="white"
                class="h-[120px] w-[120px] rounded-full"
              />
            </template>
          </div>

          <div class="text-center">
            <label for="profile" class="cursor-pointer underline">
              {{
                user.thumbnailUrl
                  ? 'プロフィール写真を変更'
                  : 'プロフィール写真を追加'
              }}
            </label>
            <input
              id="profile"
              type="file"
              class="hidden"
              accept="image/*"
              @change="handleFileChange"
            />
          </div>

          <div class="flex w-full flex-col gap-4">
            <button class="w-ful bg-main px-4 py-2 text-white" type="submit">
              保存する
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
