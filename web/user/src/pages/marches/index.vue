<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useLiveStore } from '~/store/live'
import type { I18n } from '~/types/locales'

const i18n = useI18n()

const route = useRoute()
const router = useRouter()

const lt = (str: keyof I18n['lives']['list']) => {
  return i18n.t(`lives.list.${str}`)
}

const liveStore = useLiveStore()

const { fetchArchives } = useLiveStore()
const { archivesFetchState, archiveResponse, totalArchivesCount }
  = storeToRefs(liveStore)

// 1ページ当たりに表示するマルシェ数
const pagePerItems = ref<number>(20)

// 現在のページ番号
const currentPage = computed<number>(() => {
  return route.query.page ? Number(route.query.page) : 1
})

// ページネーション情報
const pagination = computed<{
  limit: number
  offset: number
  pageArray: number[]
}>(() => {
  const totalPage = Math.ceil(totalArchivesCount.value / pagePerItems.value)
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

const handleClickLiveItem = (id: string) => {
  router.push(`/live/${id}`)
}

watch(currentPage, () => {
  fetchArchives(pagePerItems.value, pagination.value.offset)
})

useAsyncData('products', () => {
  return fetchArchives(pagePerItems.value, pagination.value.offset)
})

useSeoMeta({
  title: 'すべてのマルシェ',
})
</script>

<template>
  <div
    class="flex flex-col bg-white px-[15px] py-[48px] text-main md:px-[36px]"
  >
    <div class="w-full">
      <p
        class="text-center text-[14px] font-bold tracking-[2px] md:text-[20px]"
      >
        {{ lt('allMarcheTitle')}}
      </p>
    </div>
    <hr class="mt-[40px]">
    <div
      class="mx-auto mt-[24px] grid max-w-[1440px] gap-x-[19px] gap-y-6 md:grid-cols-3 md:gap-x-8 lg:grid-cols-3 xl:grid-cols-4"
    >
      <template v-if="archivesFetchState.isLoading">
        <div
          v-for="i in [1, 2, 3, 4, 5]"
          :key="i"
          class="w-full animate-pulse"
        >
          <div class="aspect-square w-full bg-slate-200" />
          <div class="mt-2 h-[24px] w-[80%] rounded-lg bg-slate-200" />
          <div class="mt-2 h-[24px] w-[60%] rounded-lg bg-slate-200" />
        </div>
      </template>

      <template v-else>
        <the-all-archive-item
          v-for="archive in archiveResponse.archives"
          :id="archive.scheduleId"
          :key="archive.scheduleId"
          :title="archive.title"
          :img-src="archive.thumbnailUrl"
          class="cursor-pointer"
          @click="handleClickLiveItem(archive.scheduleId)"
        />
      </template>
    </div>
    <the-pagination
      class="mt-8"
      :current-page="currentPage"
      :page-array="pagination.pageArray"
      @change-page="handleClickPage"
    />
  </div>
</template>
