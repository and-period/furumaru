<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { MOCK_RECOMMEND_ITEMS } from '~/constants/mock'
import { useTopPageStore } from '~/store/home'
import type { BannerItem } from '~/types/props'
import type { I18n } from '~/types/locales'

const i18n = useI18n()

const router = useRouter()

const topPageStore = useTopPageStore()
const { archives, lives } = storeToRefs(topPageStore)
const { getHomeContent } = topPageStore

const tt = (str: keyof I18n['base']['top']) => {
  return i18n.t(`base.top.${str}`)
}

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

const handleClickLiveMore = () => {
  router.push(`/marches`)
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

const handleClickAllArchive = () => {
  router.push(`/marches`)
}

const handleClickAllItem = () => {
  router.push(`/items`)
}

useSeoMeta({
  title: 'トップページ',
})
</script>

<template>
  <div>
    <!-- 動画部分 -->
    <div class="relative md:h-[calc(100vh-180px)] h-[calc(100vh-140px)] w-full">
      <div class="absolute bg-black/50 w-full h-full" />
      <div
        class="absolute w-full h-full z-10 flex flex-col md:gap-40 justify-center"
      >
        <div
          class="text-white md:text-[48px] text-[28px] font-bold w-full text-center tracking-wider md:grow-0 grow flex flex-col justify-center"
        >
          <p>&#035;Deep Japan 体験型</p>
          <p>ローカル映像メディア</p>
        </div>

        <div
          class="md:py-0 pb-12 px-4 xl:w-[40%] md:w-[80%] grid grid-cols-2 items-center justify-center md:gap-16 gap-4 mx-auto md:text-[18px] text-[14px] text-white font-bold"
        >
          <nuxt-link
            to="/items"
            class="bg-base/50 rounded-xl px-8 py-2 hover:bg-base/60 tracking-wide text-center"
          >
            <span class="md:text-[24px] text-[18px]">商品</span>と出会う
          </nuxt-link>
          <nuxt-link
            to="/marches"
            class="bg-base/50 rounded-xl px-8 py-2 hover:bg-base/60 tracking-wide text-center"
          >
            <span class="md:text-[24px] text-[18px]">体験</span>と出会う
          </nuxt-link>
        </div>
      </div>
      <video
        webkit-playsinline
        playsinline
        muted
        autoplay
        loop
        class="h-full object-cover w-full"
      >
        <source
          src="/video/furumaru.webm"
          type="video/webm"
        >
      </video>
    </div>

    <the-carousel
      v-if="false"
      :items="banners"
      :line-add-friend-image-url="tt('lineAddFriendImageUrl')"
      :line-add-friend-image-alt="tt('lineAddFriendImageAlt')"
      :line-coupon-text="tt('lineCouponText')"
    />

    <div class="mb-[72px] mt-8 flex flex-col gap-y-16 md:mt-[76px]">
      <the-content-box
        title="live"
        :sub-title="tt('marcheListSubTitle')"
      >
        <template v-if="isInItLoading" />
        <template v-if="lives.length === 0">
          <div class="flex justify-center">
            <img
              src="~/assets/img/furuneko-sleep.png"
              alt="furuneko sleep"
              width="120"
              height="136"
              class="block"
            >
          </div>
          <div class="mt-8 text-center text-[14px] text-main md:text-[16px]">
            <p>{{ tt("noMarcheItemFirstText") }}</p>
            <p class="md:mt-4">
              {{ tt("noMarcheItemSecondText") }}
            </p>
          </div>
          <div
            class="my-4 grid w-full justify-center md:mt-10 md:flex md:gap-x-16"
          >
            <button
              class="w-60 bg-main py-2 text-white"
              @click="handleClickAllArchive"
            >
              {{ tt("pastMarcheLinkText") }}
            </button>
            <button
              class="mt-4 w-60 bg-main py-2 text-white md:mt-0"
              @click="handleClickAllItem"
            >
              {{ tt("productsLinkText") }}
            </button>
          </div>
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
                :live-streaming-text="tt('liveStreamingText')"
                :live-upcoming-text="tt('liveUpcomingText')"
                @click="handleClickLiveItem(liveItem.scheduleId)"
              />
            </transition-group>
          </div>
          <div
            v-if="false"
            class="mb-4 mt-10 flex w-full justify-center"
          >
            <button
              class="relative w-60 bg-main py-2 text-white"
              @click="handleClickMoreViewButton"
            >
              {{ tt("viewMoreText") }}
              <div class="absolute bottom-3.5 right-4">
                <the-up-arrow-icon
                  v-show="isOpen"
                  fill="white"
                />
                <the-down-arrow-icon
                  v-show="!isOpen"
                  fill="white"
                />
              </div>
            </button>
          </div>
        </template>
      </the-content-box>

      <the-content-box
        title="archive"
        :sub-title="tt('archiveListSubTitle')"
      >
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
              :start-at="archive.startAt"
              :end-at="archive.endAt"
              :width="368"
              :archived-stream-text="tt('archivedStreamText')"
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
          <button
            class="w-60 bg-main py-2 text-white"
            @click="handleClickLiveMore"
          >
            {{ tt("archivesLinkText") }}
          </button>
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
  </div>
</template>
