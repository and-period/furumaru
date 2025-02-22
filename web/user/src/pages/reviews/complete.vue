<script setup lang="ts">
import type { I18n } from '~/types/locales'

const i18n = useI18n()

const lt = (str: keyof I18n['reviews']) => {
  return i18n.t(`reviews.${str}`)
}

useSeoMeta({
  title: lt('reviewSubmitSuccessPageTitle'),
})
</script>

<template>
  <div class="container mx-auto p-4 xl:p-0 text-main">
    <div class="flex flex-col gap-8 mt-4">
      <div class="flex items-center gap-8 mx-auto justify-center md:flex-row flex-col">
        <div class=" font-semibold text-lg">
          {{ lt('reviewSubmitSuccessMessage') }}
        </div>
        <img src="/img/purchase/complete.svg">
      </div>

      <div class=" text-center">
        <nuxt-link
          to="/"
          class="bg-main py-2 md:w-[400px] inline-block w-full text-white text-center"
        >
          {{ lt('topPageLinkText') }}
        </nuxt-link>
      </div>
    </div>
  </div>
</template>
