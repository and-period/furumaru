<script lang="ts" setup>
import {
  ExperienceStatus,
  type ProductMediaInner,
} from '~/types/api'
import type { I18n } from '~/types/locales'
import { experienceStatusToString } from '~/lib/experience'

interface Props {
  id: string
  status: ExperienceStatus
  title: string
  thumbnail: ProductMediaInner | undefined
  thumbnailIsVideo: boolean
}

interface Emits {
  (e: 'click:proceedExperience', id: string): void
}

const props = defineProps<Props>()

const i18n = useI18n()

const lt = (str: keyof I18n['items']['list']) => {
  return i18n.t(`items.list.${str}`)
}

const itemThumbnailAlt = computed<string>(() => {
  return i18n.t('items.list.itemThumbnailAlt', {
    itemName: props.title,
  })
})

const emits = defineEmits<Emits>()

const canAddCart = computed<boolean>(() => {
  if (props.status === ExperienceStatus.ACCEPTING) {
    return true
  }
  return false
})

const handleClickExperienceItem = () => {
  emits('click:proceedExperience', props.id)
}

const handleClickProceedButton = () => {
  emits('click:proceedExperience', props.id)
}
</script>

<template>
  <div class="flex flex-col text-main">
    <div class="relative">
      <div
        v-if="!canAddCart"
        class="absolute inset-0 flex items-center justify-center bg-black/50"
      >
        <p class="text-lg font-semibold text-white">
          {{
            status === ExperienceStatus.SOLD_OUT
              ? lt("soldOutText")
              : experienceStatusToString(status, i18n)
          }}
        </p>
      </div>
      <div
        v-if="thumbnail"
        class="cursor-pointer w-full"
        @click="handleClickExperienceItem"
      >
        <template v-if="thumbnailIsVideo">
          <video
            :src="thumbnail.url"
            class="aspect-square w-full"
            autoplay
            muted
            webkit-playsinline
            playsinline
            loop
          />
        </template>
        <template v-else>
          <div class="aspect-square">
            <picture>
              <nuxt-img
                provider="cloudFront"
                :src="thumbnail.url"
                :alt="itemThumbnailAlt"
                fit="contain"
                sizes="180px md:250px"
                format="png"
                class="w-full h-full object-contain"
              />
            </picture>
          </div>
        </template>
      </div>
    </div>

    <p
      class="mt-2 line-clamp-3 grow text-[14px] tracking-[1.4px] md:text-[16px] md:tracking-[1.6px]"
    >
      {{ title }}
    </p>
    <div class="flex h-6 items-center gap-2 text-[10px]">
      <button
        :disabled="!canAddCart"
        class="flex h-full grow items-center justify-center bg-orange p-1 text-[10px] text-white transition-all duration-200 ease-in-out hover:shadow-lg active:scale-95 disabled:cursor-not-allowed disabled:bg-main/60 lg:px-4 xl:text-[14px]"
        @click="handleClickProceedButton"
      >
        {{ lt("proceedPurchase") }}
      </button>
    </div>
  </div>
</template>
