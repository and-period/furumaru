<script setup lang="ts">
import dayjs from 'dayjs'
import { ScheduleStatus } from '~/types/api'

interface Props {
  id: string
  title: string
  imgSrc: string
  startAt: number
  isLiveStatus: ScheduleStatus
  marcheName: string
  address: string
  cnName: string
  cnImgSrc: string
  liveStreamingText: string
  liveUpcomingText: string
}

interface Emits {
  (e: 'click'): void
}

const props = defineProps<Props>()

const emits = defineEmits<Emits>()

const formattedStartAt = computed(() => {
  return dayjs.unix(props.startAt).format('YYYY/MM/DD HH:mm')
})

const isLiveStreaming = (status: ScheduleStatus) => {
  if (status === ScheduleStatus.LIVE || status === ScheduleStatus.CLOSED) {
    return true
  }
  else {
    return false
  }
}

const handleClick = () => {
  emits('click')
}
</script>

<template>
  <div
    class="lg:group cursor-pointer bg-base drop-shadow-sm duration-75 ease-in-out lg:hover:z-10 lg:hover:scale-[1.2] lg:hover:bg-white"
    @click="handleClick"
  >
    <div class="relative w-full p-4">
      <div
        v-if="isLiveStreaming(isLiveStatus)"
        class="absolute -left-4 -top-4 z-[1] flex h-16 w-16 flex-col items-center justify-center rounded-full bg-orange xl:-left-8 xl:-top-8"
      >
        <the-live-icon />
        <div class="text-xl font-bold uppercase text-white">
          live
        </div>
      </div>

      <div class="relative">
        <nuxt-img
          provider="cloudFront"
          :src="imgSrc"
          fit="contain"
          :alt="`live-${title}-thumbnail`"
          class="aspect-video w-full object-cover"
          sizes="320px md:368px"
        />
        <div
          v-if="!isLiveStreaming(isLiveStatus)"
          class="absolute bottom-0 flex h-[48px] w-full items-center justify-center bg-black/50 text-[16px] font-bold tracking-[1.6px] text-white"
        >
          {{ formattedStartAt }} 〜 {{ liveUpcomingText }}
        </div>
      </div>

      <div class="mt-4 flex w-full flex-col gap-2">
        <div class="flex items-center text-sm">
          <div class="grow">
            <span
              :class="{
                'rounded px-2 font-bold': true,
                'border-2 border-orange bg-white text-orange':
                  isLiveStreaming(isLiveStatus),
                'border-2 border-main text-main':
                  !isLiveStreaming(isLiveStatus),
              }"
            >
              {{ isLiveStreaming(isLiveStatus) ? liveStreamingText : liveUpcomingText }}
            </span>
            <span class="ml-2 after:content-['〜']">{{
              formattedStartAt
            }}</span>
          </div>
          <button class="h-4 w-4 hover:scale-110">
            <the-ellipsis-vertical-icon class="h-5 w-5" />
          </button>
        </div>

        <p class="line-clamp-3">
          {{ title }}
        </p>

        <div
          class="absolute bottom-[-120px] left-0 hidden h-[120px] w-full bg-white p-4 group-hover:block"
        >
          <hr class="border-dashed">
          <div class="mt-4 flex w-full items-center justify-end gap-4">
            <div class="text-[12px] tracking-[1.2px]">
              <p class="mb-1">
                {{ marcheName }}/{{ address }}
              </p>
              <p class="tracking-[1.3px]">
                コーディネーター：{{ cnName }}
              </p>
            </div>
            <nuxt-img
              provider="cloudFront"
              :src="cnImgSrc"
              :alt="`${cnName}のーアバター画像`"
              width="40"
              height="40"
              class="h-10 w-10 rounded-full"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
