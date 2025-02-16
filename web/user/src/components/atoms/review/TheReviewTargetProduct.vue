<script setup lang="ts">
import type { Product } from '~/types/api'
import type { I18n } from '~/types/locales'

const i18n = useI18n()

const lt = (str: keyof I18n['reviews']) => {
  return i18n.t(`reviews.${str}`)
}

interface Props {
  product: Product
}

defineProps<Props>()
</script>

<template>
  <div class="flex flex-col gap-2">
    <div

      class=" font-semibold text-[16px] tracking-[2px]"
    >
      {{ lt('reviewTargetLabel') }}
    </div>
    <div class="flex gap-4">
      <nuxt-img
        provider="cloudFront"
        :src="product.thumbnailUrl"
        class="w-16 h-16 object-contain md:w-32 md:h-32 border "
        fit="contain"
        sizes="64px md:128px"
        format="png"
      />
      <div class="font-semibold tracking-[2px] text-[14px]">
        {{ product.name }}
      </div>
    </div>
  </div>
</template>
