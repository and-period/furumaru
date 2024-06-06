<script lang="ts" setup>
import type { BannerItem } from '~/types/props'

interface Props {
  items: BannerItem[]
  lineAddFriendImageUrl: string
  lineAddFriendImageAlt: string
  lineCouponText: string
}

const props = defineProps<Props>()

const router = useRouter()

const activeIdx = ref<number>(0)

const bannerItems = computed(() => {
  return props.items.map((item, i) => {
    return {
      src: item.imgSrc,
      link: item.link,
      isInternalLink: item.isInternalLink,
      // activeはカルーセルの中心に表示される要素の場合にtrueになる
      active: i === activeIdx.value,
      // leftContentはカルーセルの左側に表示される要素の場合にtrueになる
      leftContent:
        activeIdx.value === 0
          ? i === props.items.length - 1
          : i === activeIdx.value - 1,
      // rightContentはカルーセルの右側に表示される要素の場合にtrueになる
      rightContent:
        activeIdx.value === props.items.length - 1
          ? i === 0
          : i === activeIdx.value + 1,
    }
  })
})

const handleClickLeftArrowButton = () => {
  if (activeIdx.value === 0) {
    activeIdx.value = props.items.length - 1
  }
  else {
    activeIdx.value = activeIdx.value - 1
  }
}

const handleClickRightArrowButton = () => {
  if (activeIdx.value === props.items.length - 1) {
    activeIdx.value = 0
  }
  else {
    activeIdx.value = activeIdx.value + 1
  }
}

const handleClickItem = (
  link: string,
  isActive: boolean,
  isInternalLink: boolean,
) => {
  if (isActive) {
    // カルーセルの中心に表示されるアクティブな要素のclickイベントだけをハンドリングする
    if (isInternalLink) {
      // 内部リンクは vue-routerで遷移させる
      router.push(link)
    }
    else {
      // 外部リンクは新規タブで開く
      window.open(link, '_blank', 'noopener,noreferrer')?.focus()
    }
  }
}
</script>

<template>
  <div class="flex items-center justify-center gap-x-4">
    <the-icon-button
      class="z-50 h-10 w-10 bg-white/70 hover:bg-white"
      @click="handleClickLeftArrowButton"
    >
      <the-left-arrow-icon />
    </the-icon-button>

    <div
      class="relative h-[128px] w-[312px] sm:h-[256px] sm:w-[624px] lg:h-[320px] lg:w-[780px]"
    >
      <div
        v-for="(item, i) in bannerItems"
        :key="i"
        :class="{
          'absolute h-full w-full transition-all duration-300': true,
          'z-40 cursor-pointer': item.active,
          'z-0 brightness-75': !item.active,
          'translate-x-[-262px] sm:translate-x-[-624px] lg:translate-x-[-780px]':
            item.leftContent,
          'translate-x-[262px] sm:translate-x-[624px] lg:translate-x-[780px]':
            item.rightContent,
        }"
        :style="`background-position: center; background-size: cover; background-image: url(${item.src});`"
        @click="handleClickItem(item.link, item.active, item.isInternalLink)"
      />
    </div>

    <the-icon-button
      class="z-50 h-10 w-10 bg-white/50 hover:bg-white"
      @click="handleClickRightArrowButton"
    >
      <the-right-arrow-icon />
    </the-icon-button>
  </div>

  <div class="flex justify-center pt-8">
    <a href="https://lin.ee/49SOeUC"><img
      :src="lineAddFriendImageUrl"
      :alt="lineAddFriendImageAlt"
      class="h-[40px] md:h-[65px]"
    >
    </a>
  </div>
  <p class="mt-4 flex justify-center text-center text-[16px] md:text-[20px]">
    {{ lineCouponText }}
  </p>
</template>
