<script setup lang="ts">
const route = useRoute()

const page = computed(() => {
  return route.query.from
})

useSeoMeta({
  title: 'アカウント更新完了',
})
</script>

<template>
  <div class="container mx-auto p-4 md:p-0">
    <TheAccountEditCompleteCard class="mt-6" :page="page" />
  </div>
</template>
