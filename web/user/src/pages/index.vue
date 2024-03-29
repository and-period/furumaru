<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { MOCK_RECOMMEND_ITEMS } from '~/constants/mock'
import { useTopPageStore } from '~/store/home'
import type { BannerItem } from '~/types/props'

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

useAsyncData('home-content', async () => {
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

const banners: BannerItem[] = [
  {
    imgSrc: '/img/banner2.jpg',
    link: '/about',
    isInternalLink: true,
  },
  { imgSrc: '/img/banner.png', link: '/about', isInternalLink: true },
  {
    imgSrc: '/img/banner3.png',
    link: '/live/p6XURyWhSk2EerwWiYYzDU',
    isInternalLink: true,
  },
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
    <the-carousel :items="banners" />

    <div class="mb-[72px] mt-[76px] flex flex-col gap-y-16">
      <the-content-box title="live" sub-title="配信中・配信予定のマルシェ">
        <template v-if="isInItLoading"> </template>
        <template v-if="lives.length === 0">
          <p class="text-center text-main">
            ただいまライブは開催されていません。
            新しいライブが開催されるまでお待ちください。
          </p>
        </template>
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
                :is-live-status="liveItem.status"
                :marche-name="liveItem.coordinator.marcheName"
                :address="liveItem.coordinator.city"
                :cn-name="liveItem.coordinator.username"
                :cn-img-src="liveItem.coordinator.thumbnailUrl"
                @click="handleClickLiveItem(liveItem.scheduleId)"
              />
            </transition-group>
          </div>
          <div v-if="false" class="mb-4 mt-10 flex w-full justify-center">
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
            class="hidden-scrollbar flex w-full flex-col gap-8 md:flex-row md:overflow-x-scroll"
          >
            <the-archive-item
              v-for="archive in archives"
              :id="archive.scheduleId"
              :key="archive.scheduleId"
              :title="archive.title"
              :img-src="archive.thumbnailUrl"
              :width="368"
              class="cursor-pointer md:min-w-[368px] md:max-w-[368px]"
              @click="handleClickLiveItem(archive.scheduleId)"
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

      <the-content-box
        v-if="false"
        title="recommend"
        sub-title="おすすめの商品"
      >
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
    <div class="flex justify-center">
      <a href="https://lin.ee/49SOeUC"><img src="https://scdn.line-apps.com/n/line_add_friends/btn/ja.png" alt="友だち追加" class="h-[40px] md:h-[65px]"></a>
    </div>
    <p class="mx-2 my-4 flex justify-center text-center text-[16px] md:text-[20px]">LINE友達追加すると、ふるマルの配信・クーポンなどのお得情報をGETできる！</p>

  </div>
</template>
