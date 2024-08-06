<script setup lang="ts">
const error = useError()

useSeoMeta({
  title: 'エラーページ',
})
</script>

<template>
  <nuxt-layout>
    <div class="container mx-auto p-4 text-main">
      <h2 class="text-[14px] font-bold tracking-[2px] md:text-[20px]">
        {{
          error.statusCode === 404
            ? "指定したページが見つかりません"
            : "エラーが発生しました"
        }}
      </h2>
      <nuxt-link
        to="/"
        class="underline mt-6 inline-block"
      >
        トップページに戻る
      </nuxt-link>
    </div>
  </nuxt-layout>
</template>
