<script setup lang="ts">
import { datetimeformatterFromUnixtime } from '~/lib/dayjs'
import type { I18n } from '~/types/locales'

interface Props {
  title: string
  description: string
  isArchive: boolean
  isLiveStreaming: boolean
  startAt: number
  marcheName: string
  coordinatorId: string
  coordinatorName: string
  coordinatorImgSrc: string
  coordinatorAddress: string
}

const i18n = useI18n()

const dt = (str: keyof I18n['lives']['details']) => {
  return i18n.t(`lives.details.${str}`)
}

const coordinatorThumbnailAlt = computed<string>(() => {
  return i18n.t('lives.details.coordinatorThumbnailAlt', {
    coordinatorName: props.coordinatorName,
  })
})

const props = defineProps<Props>()

interface Emits {
  (e: 'click:coordinator', id: string): void
}

const emits = defineEmits<Emits>()

const handleCLickCoordinator = () => {
  emits('click:coordinator', props.coordinatorId)
}

const showDetail = ref<boolean>(false)

const handleClickShowDetailButton = () => {
  showDetail.value = !showDetail.value
}
</script>

<template>
  <div class="mt-2 px-4">
    <div class="flex items-center gap-2">
      <template v-if="isArchive">
        <div
          class="flex max-w-fit items-center justify-center rounded border-2 border-main px-2 font-bold text-main"
        >
          {{ dt('archivedStreamText') }}
        </div>
      </template>
      <template v-else-if="isLiveStreaming">
        <div
          class="flex max-w-fit items-center justify-center rounded border-2 border-orange bg-orange px-2 font-bold text-white"
        >
          <div class="mr-2 pt-[2px]">
            <the-live-icon />
          </div>
          <div class="align-middle">
            LIVE
          </div>
        </div>
      </template>

      <div class="text-[14px] tracking-[1.4px] after:content-['〜']">
        {{ datetimeformatterFromUnixtime(startAt) }}
      </div>
    </div>
    <p class="mt-2 line-clamp-1 tracking-[1.6px]">
      {{ title }}
    </p>
    <div class="mt-4 flex items-center gap-2">
      <nuxt-img
        width="40"
        height="40"
        provider="cloudFront"
        :src="coordinatorImgSrc"
        class="h-10 w-10 rounded-full hover:cursor-pointer"
        :alt=coordinatorThumbnailAlt
        @click="handleCLickCoordinator"
      />
      <div class="text-[12px] tracking-[1.2px]">
        <p class="mb-1">
          {{ marcheName }}/{{ coordinatorAddress }}
        </p>
        <p>
          {{ dt('coordinatorLabel') }}:
          <span
            class="cursor-pointer hover:underline"
            @click="handleCLickCoordinator"
          >{{ coordinatorName }}</span>
        </p>
      </div>
    </div>

    <div>
      <p
        v-if="showDetail"
        class="mt-6 whitespace-pre-wrap text-[14px] tracking-[1.4px]"
        v-html="description"
      />
      <button
        class="inline-flex w-full items-center justify-center gap-2 text-[12px] tracking-[1.2px]"
        @click="handleClickShowDetailButton"
      >
        <div>
          {{ showDetail ? dt('hideMarcheDetailsText') : dt('showMarcheDetailsText') }}
        </div>
        <div>
          <the-up-arrow-icon v-if="showDetail" />
          <the-down-arrow-icon v-if="!showDetail" />
        </div>
      </button>
    </div>
  </div>
</template>
