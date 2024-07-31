<script setup lang="ts">
const route = useRoute()

const id = computed(() => {
  return route.params.id as string
})

const { data } = await useAsyncData(`volunteer-${id.value}`, () => {
  return $fetch(`/api/cms/volunteer/${id.value}`)
})

const title = computed(() => {
  return data.value.title
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
        class="text-[14px] font-bold tracking-[2px] text-main mb-4 inline-flex gap-1 items-center"
      >
        <TheLeftArrowIcon class="h-3" />
        一覧に戻る
      </nuxt-link>
    </div>

    <div
      class="container mx-auto text-center text-[14px] font-bold tracking-[2px] md:text-[20px]"
    >
      <template v-if="data.title">
        {{ data.title }}
      </template>
    </div>
    <div class="grow container mx-auto">
      <hr class="mt-[40px] mb-[20px]">
      <div
        v-if="data.content"
        class="purose md:prose-lg px-4"
        v-html="data.content"
      />
    </div>
  </div>
</template>
