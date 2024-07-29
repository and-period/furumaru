<script setup lang="ts">
const { data } = await useAsyncData('volunteer-list', () => {
  return $fetch('/api/cms/volunteer')
})

useSeoMeta({
  title: 'ひろしま援農計画',
})
</script>

<template>
  <div class="w-full bg-white flex flex-col py-[48px] text-main">
    <div
      class="container mx-auto text-center text-[14px] font-bold tracking-[2px] md:text-[20px]"
    >
      ひろしま援農計画
    </div>
    <div class="grow container mx-auto">
      <hr class="mt-[40px] mb-[20px]">
      <template v-if="data.contents">
        <div class="grid md:grid-cols-3 w-full px-4">
          <div
            v-for="content in data.contents"
            :key="content.id"
          >
            <img
              :src="content.eyecatch.url"
              class="md:w-[400px] aspect-video object-cover w-full"
            >
            <div class="flex flex-col gap-2 mt-2">
              <h2 class="text-[18px] font-semibold tracking-[2px]">
                {{ content.title }}
              </h2>
              <div
                class="inline-flex gap-4 [&>span]:px-3 [&>span]:py-1 [&>span]:border [&>span]:border-main [&>span]:text-main [&>span]:rounded-full text-[14px]"
              >
                <span v-if="content.category1">
                  {{ content.category1.name }}
                </span>
                <span v-if="content.category2">
                  {{ content.category2.name }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>
