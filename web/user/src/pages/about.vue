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
    <div class="flex flex-wrap items-center justify-center gap-x-7 mt-10">
      <div class="text-main mb-10">
        <p class="text-2xl font-semibold mb-12">
          {{ t('leadSentence') }}
        </p>
        <p class="leading-7" v-html="t('description')" />
      </div>
      <the-concept />
    </div>

    <div class="flex flex-wrap gap-10 my-10 justify-center">
      <the-circle v-for="item, i in circleItems" :key="i" :num="i + 1">
        <p class="font-bold text-2xl mb-6" v-html="item.title" />
        <p class="text-xl leading-9 text-left break-words">
          {{ item.description }}
        </p>
        <div v-if="item.linkItem" class="mt-10 font-semibold text-lg">
          <nuxt-link :to="item.linkItem.href">
            {{ item.linkItem.text }}
          </nuxt-link>
        </div>
      </the-circle>
    </div>
  </div>
</template>
