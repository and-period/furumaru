<script setup lang="ts">
import type { VolunteerBlogItemResponse } from '~/types/cms/volunteer'

const route = useRoute()

const id = computed(() => {
  return route.params.id as string
})

const { data, error } = await useAsyncData<VolunteerBlogItemResponse>(
  `volunteer-${id.value}`,
  () => {
    return $fetch(`/api/cms/volunteer/${id.value}`)
  },
)

const title = computed(() => {
  return data.value?.title || ''
})

if (error.value) {
  throw createError(error.value)
}

useSeoMeta({
  title,
})
</script>

<template>
  <div class="w-full bg-white flex flex-col py-[48px] text-main">
    <div class="container mx-auto px-4">
      <nuxt-link
        to="/volunteer"
        class="md:text-[14px] text-[12px] font-bold tracking-[2px] text-main mb-8 inline-flex gap-1 items-center"
      >
        <the-left-arrow-icon class="h-3" />
        一覧に戻る
      </nuxt-link>
    </div>

    <div class="grow container mx-auto">
      <div
        v-if="data"
        class="prose md:prose-lg px-4 max-w-full"
        v-html="data.content"
      />
    </div>
  </div>
</template>
