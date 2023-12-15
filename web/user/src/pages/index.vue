<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { MOCK_RECOMMEND_ITEMS } from '~/constants/mock'
import { useTopPageStore } from '~/store/home'

const router = useRouter()

const topPageStore = useTopPageStore()
const { archives, lives } = storeToRefs(topPageStore)
const { getHomeContent } = topPageStore

const isInItLoading = ref<boolean>(false)

const archiveRef = ref<HTMLDivElement | null>(null)
const archiveRefScrollLeft = ref<number>(0)

const updateScrollLeft = () => {
  if (archiveRef.value) {
    archiveRefScrollLeft.value = archiveRef.value.scrollLeft
  }
}

onMounted(async () => {
  isInItLoading.value = true
  await getHomeContent()
  isInItLoading.value = false
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

const handleClickLiveItem = (id: string) => {
  router.push(`/live/${id}`)
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

useSeoMeta({
  title: 'トップページ',
})
</script>

<template>
  <div>
    <the-carousel :images="banners" />

    <div class="mb-[72px] mt-[76px] flex flex-col gap-y-16">
      <the-content-box title="live" sub-title="配信中・配信予定のマルシェ">
        <template v-if="isInItLoading"> </template>
        <template v-if="lives.length > 0">
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
                v-for="liveItem in lives"
                :id="liveItem.scheduleId"
                :key="liveItem.scheduleId"
                :title="liveItem.title"
                :img-src="liveItem.thumbnailUrl"
                :start-at="liveItem.startAt"
                :is-live-streaming="liveItem.isLiveStreaming"
                :marche-name="liveItem.coordinator.marcheName"
                :address="liveItem.coordinator.city"
                :cn-name="liveItem.coordinator.username"
                :cn-img-src="liveItem.coordinator.thumbnailUrl"
                @click="handleClickLiveItem(liveItem.scheduleId)"
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
        </template>
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
