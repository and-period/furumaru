<script setup lang="ts">
interface Props {
  page?:
    | 'username'
    | 'email'
    | 'password'
    | 'thumbnail'
    | 'notification'
    | 'accountId'
}

const props = defineProps<Props>()

const message = computed(() => {
  switch (props.page) {
    case 'username':
      return 'ユーザー名を変更しました。'
    case 'accountId':
      return 'ユーザーIDを変更しました。'
    case 'email':
      return 'メールアドレスを変更しました。'
    case 'password':
      return 'パスワードを変更しました。'
    case 'thumbnail':
      return 'サムネイルを変更しました。'
    case 'notification':
      return '通知設定を変更しました。'
    default:
      return 'アカウント情報を変更しました。'
  }
})
</script>

<template>
  <div
    class="mx-auto flex max-w-[560px] flex-col items-center justify-center bg-white px-4 py-6 tracking-[10%] text-main md:p-10 xl:px-20"
  >
    <p class="mb-12 text-[16px] font-semibold">
      {{ message }}
    </p>
    <nuxt-link
      to="/"
      class="mb-6 w-full bg-main px-4 py-2 text-center text-white"
    >
      TOPページへ
    </nuxt-link>
  </div>
</template>
