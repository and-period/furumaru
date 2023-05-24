<script lang="ts" setup>
import { LinkItem } from '~/types/props'

import { I18n } from '~/types/locales'

interface CircleItem {
  title: string
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
    description: t('firstPointDescription'),
    linkItem: {
      text: t('firstPointLinkText'),
      href: localePath('/items')
    }
  },
  {
    title: t('secondPointTitle'),
    description: t('secondPointDescription')
  },
  {
    title: t('thirdPointTitle'),
    description: t('thirdPointDescription')
  },
  {
    title: t('forthPointTitle'),
    description: t('forthPointDescription'),
    linkItem: {
      text: t('forthPointLinkText'),
      href: localePath('/')
    }
  }
])
</script>

<template>
  <div>
    <div class="flex flex-wrap items-center justify-center gap-x-[120px] pt-[80px]">
      <div class="text-main mb-10">
        <p class="text-[32px] font-semibold mb-12">
          {{ t('leadSentence') }}
        </p>
        <p class="leading-10 text-xl whitespace-pre-wrap" v-html="t('description')" />
      </div>
      <the-concept />
    </div>

    <div class="flex flex-wrap gap-10 my-20 justify-center">
      <the-circle v-for="item, i in circleItems" :key="i" :num="i + 1">
        <p class="font-bold text-2xl mb-6" v-html="item.title" />
        <p class="text-xl leading-9 text-left break-words whitespace-pre-wrap">
          {{ item.description }}
        </p>
        <div v-if="item.linkItem" class="mt-10 font-semibold text-xl">
          <nuxt-link :to="item.linkItem.href">
            <div class="">
              {{ item.linkItem.text }}
              <svg width="10" height="17" viewBox="0 0 10 17" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path fill-rule="evenodd" clip-rule="evenodd" d="M1.49023 16.9707L0.0356148 15.5161L7.06628 8.48542L0.0356165 1.45476L1.49024 0.000140667L9.97552 8.48542L1.49023 16.9707Z" fill="#604C3F"/>
              </svg>
            </div>
          </nuxt-link>
        </div>
      </the-circle>
    </div>
  </div>
</template>
