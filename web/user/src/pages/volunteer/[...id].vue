<script setup lang="ts">
import type { VolunteerBlogItemResponse } from '~/types/cms/volunteer'

const route = useRoute()

const id = computed(() => {
  return route.params.id as string
})

const { data } = await useAsyncData<VolunteerBlogItemResponse>(
  `volunteer-${id.value}`,
  () => {
    return $fetch(`/api/cms/volunteer/${id.value}`)
  },
)

const title = computed(() => {
  return data.value?.title || ''
})

useSeoMeta({
  title,
})
</script>

<template>
  <div class="w-full bg-white flex flex-col py-[48px] text-main">
    <div class="container mx-auto">
      <nuxt-link
        to="/volunteer"
        class="md:text-[14px] text-[12px] font-bold tracking-[2px] text-main md:mb-4 mb-8 inline-flex gap-1 items-center"
      >
        <the-left-arrow-icon class="h-3" />
        一覧に戻る
      </nuxt-link>
    </div>

    <div
      class="container mx-auto text-center text-[14px] font-bold tracking-[2px] md:text-[20px]"
    >
      <template v-if="data">
        {{ data.title }}
      </template>
    </div>
    <div class="grow container mx-auto">
      <hr class="mt-[40px] mb-[20px]">
      <div
        v-if="data"
        class="prose md:prose-lg px-4 max-w-full"
        v-html="data.content"
      />
    </div>
  </div>
</template>
