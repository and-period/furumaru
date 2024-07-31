<script setup lang="ts">
import type { VolunteerBlogListResponse } from '~/types/cms/volunteer'

const route = useRoute()
const router = useRouter()

const { data } = await useAsyncData<VolunteerBlogListResponse>(
  'volunteer-list',
  () => {
    return $fetch('/api/cms/volunteer')
  },
)

const totalCount = computed(() => {
  return data.value?.totalCount || 0
})

// 現在のページ番号
const currentPage = computed<number>(() => {
  return route.query.page ? Number(route.query.page) : 1
})

// 1ページ当たりに表示する商品数
const pagePerItems = ref<number>(20)

// ページネーション情報
const pagination = computed<{
  limit: number
  offset: number
  pageArray: number[]
}>(() => {
  const totalPage = Math.ceil(totalCount.value / pagePerItems.value)
  const pageArray = Array.from({ length: totalPage }, (_, i) => i + 1)

  return {
    limit: pagePerItems.value,
    offset: pagePerItems.value * (currentPage.value - 1),
    pageArray,
  }
})

const handleClickPage = (page: number) => {
  router.push({
    query: {
      ...route.query,
      page,
    },
  })
}

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
      <template v-if="data">
        <div class="grid md:grid-cols-3 w-full px-4">
          <nuxt-link
            v-for="content in data.contents"
            :key="content.id"
            :to="`/volunteer/${content.id}`"
            class="flex flex-col gap-4"
          >
            <img
              :src="content.eyecatch.url"
              class="w-full aspect-video object-cover w-full"
            >
            <div class="flex flex-col gap-2 mt-2">
              <h2 class="text-[18px] font-semibold tracking-[2px]">
                {{ content.title }}
              </h2>

              <div class="text-gray flex flex-col gap-2 tracking-[10%]">
                <div>{{ content.name }}</div>
                <div class="inline-flex gap-2 items-center">
                  <the-map-pin-icon class="h-4 w-4" />
                  {{ content.location }}
                </div>
              </div>

              <div
                class="inline-flex gap-4 [&>span]:px-3 [&>span]:py-1 [&>span]:border [&>span]:border-main [&>span]:text-main [&>span]:rounded-full text-[14px]"
              >
                <span
                  v-for="category in content.category"
                  :key="category.id"
                >
                  {{ category }}
                </span>
              </div>
            </div>
          </nuxt-link>
        </div>
      </template>
      <the-pagination
        class="mt-8"
        :current-page="currentPage"
        :page-array="pagination.pageArray"
        @change-page="handleClickPage"
      />
    </div>
  </div>
</template>
