<script setup lang="ts">
const router = useRouter()

const handleClickBackButton = () => {
  router.back()
}
</script>

<template>
  <div
    class="mx-auto flex max-w-[560px] flex-col items-center justify-center bg-white px-4 py-6 tracking-[10%] text-main md:p-10 xl:px-20"
  >
    <p class="mb-12 text-[16px] font-semibold">この操作はログインが必要です</p>
    <nuxt-link
      to="/signin"
      class="mb-6 w-full bg-main px-4 py-2 text-center text-white"
    >
      ログイン画面へ
    </nuxt-link>

    <button
      class="inline-flex items-center gap-1 text-[12px]"
      @click="handleClickBackButton"
    >
      <TheLeftArrowIcon class="h-3 w-3" />
      前の画面に戻る
    </button>
  </div>
</template>
