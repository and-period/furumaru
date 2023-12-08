<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { MOCK_LIVE_ITEMS, MOCK_RECOMMEND_ITEMS } from '~/constants/mock'
import { useTopPageStore } from '~/store/home'

const router = useRouter()

const topPageStore = useTopPageStore()
const { archives } = storeToRefs(topPageStore)
const { getHomeContent } = topPageStore

const archiveRef = ref<HTMLDivElement | null>(null)
const archiveRefScrollLeft = ref<number>(0)

const updateScrollLeft = () => {
  if (archiveRef.value) {
    archiveRefScrollLeft.value = archiveRef.value.scrollLeft
  }
}

onMounted(() => {
  getHomeContent()
})

onMounted(() => {
  if (archiveRef.value) {
    archiveRef.value.addEventListener('scroll', updateScrollLeft)
  }
})

onUnmounted(() => {
  if (archiveRef.value) {
    archiveRef.value.removeEventListener('scroll', updateScrollLeft)
  }
})

const handleClickArchiveLeftButton = () => {
  if (archiveRef.value) {
    archiveRef.value.scrollTo({
      left: archiveRef.value.scrollLeft - 368,
      behavior: 'smooth',
    })
    updateScrollLeft()
  }
}

const handleClickArchiveRightButton = () => {
  if (archiveRef.value) {
    archiveRef.value.scrollTo({
      left: archiveRef.value.scrollLeft + 368,
      behavior: 'smooth',
    })
    updateScrollLeft()
  }
}

const handleClickLiveItem = (_: string) => {
  router.push('/live')
}

const banners: string[] = [
  '/img/banner.png',
  '/img/banner.png',
  '/img/banner.png',
]

const isOpen = ref<boolean>(false)

const handleClickMoreViewButton = () => {
  isOpen.value = !isOpen.value
}

const liveItems = computed(() => {
  if (isOpen.value) {
    return MOCK_LIVE_ITEMS
  } else {
    return MOCK_LIVE_ITEMS.slice(0, 6)
  }
})

useSeoMeta({
  title: 'トップページ',
})
</script>

<template>
  <div>
    <the-carousel :images="banners" />

    <div class="mb-[72px] mt-[76px] flex flex-col gap-y-16">
      <the-content-box title="live" sub-title="配信中・配信予定のマルシェ">
        <div
          class="mx-auto grid max-w-7xl gap-x-10 gap-y-8 px-2 md:grid-cols-2 lg:grid-cols-3"
        >
          <transition-group
            enter-active-class="duration-300 ease-in-out"
            enter-from-class="opacity-0 h-0"
            enter-to-class="opacity-100 h-full"
            leave-active-class="duration-300 ease-in-out"
            leave-from-class="opacity-100 h-full"
            leave-to-class="opacity-0 h-0"
          >
            <the-live-item
              v-for="liveItem in liveItems"
              :id="liveItem.id"
              :key="liveItem.id"
              :title="liveItem.title"
              :img-src="liveItem.imgSrc"
              :start-at="liveItem.startAt"
              :published="liveItem.published"
              :marche-name="liveItem.marcheName"
              :address="liveItem.address"
              :cn-name="liveItem.cnName"
              :cn-img-src="liveItem.cnImgSrc"
              @click="handleClickLiveItem(id)"
            />
          </transition-group>
        </div>
        <div class="mb-4 mt-10 flex w-full justify-center">
          <button
            class="relative w-60 bg-main py-2 text-white"
            @click="handleClickMoreViewButton"
          >
            もっと見る
            <div class="absolute bottom-3.5 right-4">
              <the-up-arrow-icon v-show="isOpen" fill="white" />
              <the-down-arrow-icon v-show="!isOpen" fill="white" />
            </div>
          </button>
        </div>
      </the-content-box>

      <the-content-box title="archive" sub-title="過去のマルシェ">
        <div class="relative mx-auto flex max-w-[1440px]">
          <div class="absolute left-4 flex h-[208px] items-center">
            <the-icon-button
              class="hidden bg-white/50 hover:bg-white md:block"
              @click="handleClickArchiveLeftButton"
            >
              <the-left-arrow-icon />
            </the-icon-button>
          </div>
          <div
            ref="archiveRef"
            class="hidden-scrollbar flex flex-col gap-8 md:flex-row md:flex-nowrap md:overflow-x-scroll"
          >
            <the-archive-item
              v-for="archive in archives"
              :id="archive.scheduleId"
              :key="archive.scheduleId"
              :title="archive.title"
              :img-src="archive.thumbnailUrl"
            />
          </div>
          <div class="absolute right-4 flex h-[208px] items-center">
            <the-icon-button
              class="hidden bg-white/50 hover:bg-white md:block"
              @click="handleClickArchiveRightButton"
            >
              <the-right-arrow-icon />
            </the-icon-button>
          </div>
        </div>

        <div class="mb-4 mt-10 flex w-full justify-center">
          <button class="w-60 bg-main py-2 text-white">一覧を見る</button>
        </div>
      </the-content-box>

      <the-content-box title="recommend" sub-title="おすすめの商品">
        <div
          class="mx-auto grid max-w-[1440px] grid-cols-2 gap-x-8 gap-y-6 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5"
        >
          <the-product-list-item
            v-for="productItem in MOCK_RECOMMEND_ITEMS"
            :id="productItem.id"
            :key="productItem.id"
            :name="productItem.name"
            :price="productItem.price"
            :img-src="productItem.imgSrc"
            :inventory="productItem.inventory"
            :address="productItem.address"
            :cn-name="productItem.cnName"
            :cn-img-src="productItem.cnImgSrc"
          />
        </div>
      </the-content-box>
    </div>
  </div>
</template>
