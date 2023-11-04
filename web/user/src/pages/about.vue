<script lang="ts" setup>
import { LinkItem } from '~/types/props'

import { I18n } from '~/types/locales'

interface CircleItem {
  title: string
  imgSrc: string
  description: string
  linkItem?: LinkItem
}

const i18n = useI18n()
const localePath = useLocalePath()

const t = (str: keyof I18n['base']['about']) => {
  return i18n.t(`base.about.${str}`)
}

const circleItems = computed<CircleItem[]>(() => [
  {
    title: t('firstPointTitle'),
    imgSrc: '/img/about/1.svg',
    description: t('firstPointDescription'),
    linkItem: {
      text: t('firstPointLinkText'),
      href: localePath('/items'),
    },
  },
  {
    title: t('secondPointTitle'),
    imgSrc: '/img/about/2.svg',
    description: t('secondPointDescription'),
  },
  {
    title: t('thirdPointTitle'),
    imgSrc: '/img/about/3.svg',
    description: t('thirdPointDescription'),
  },
  {
    title: t('forthPointTitle'),
    imgSrc: '/img/about/4.svg',
    description: t('forthPointDescription'),
    linkItem: {
      text: t('forthPointLinkText'),
      href: localePath('/'),
    },
  },
])

useSeoMeta({
  title: 'ふるマルについて',
})
</script>

<template>
  <div>
    <div
      class="flex flex-wrap items-center justify-center gap-x-[120px] pt-[80px] tracking-wider"
    >
      <div class="mb-10 text-main">
        <p class="mb-12 text-[32px] font-semibold">
          {{ t('leadSentence') }}
        </p>
        <p
          class="whitespace-pre-wrap text-xl leading-10"
          v-html="t('description')"
        />
      </div>
      <the-concept />
    </div>

    <div class="my-20 flex flex-wrap justify-center gap-10 tracking-wider">
      <div v-for="(item, i) in circleItems" :key="i" :num="i + 1">
        <div
          class="h-[620px] w-[620px] break-words rounded-full bg-white px-28 py-12 text-center text-main"
        >
          <div class="relative mx-auto mb-16 block h-[73px] w-[98px]">
            <img :src="item.imgSrc" :alt="`about-point-${i + 1}`" />
          </div>
          <p class="mb-6 whitespace-pre text-2xl font-bold">
            {{ item.title }}
          </p>
          <p class="whitespace-pre-wrap text-left text-xl leading-9">
            {{ item.description }}
          </p>
          <div v-if="item.linkItem" class="mt-10 text-xl font-semibold">
            <nuxt-link
              :to="item.linkItem.href"
              class="flex flex-row justify-center"
            >
              {{ item.linkItem.text }}
              <div class="ml-2 mt-2">
                <svg
                  width="10"
                  height="17"
                  viewBox="0 0 10 17"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    fill-rule="evenodd"
                    clip-rule="evenodd"
                    d="M1.49023 16.9707L0.0356148 15.5161L7.06628 8.48542L0.0356165 1.45476L1.49024 0.000140667L9.97552 8.48542L1.49023 16.9707Z"
                    fill="#604C3F"
                  />
                </svg>
              </div>
            </nuxt-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
