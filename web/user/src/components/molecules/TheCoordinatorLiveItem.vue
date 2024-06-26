<script setup lang="ts">
import dayjs from 'dayjs'
import { ScheduleStatus } from '~/types/api'

interface Props {
  id: string
  title: string
  imgSrc: string
  startAt: number
  isLiveStatus: ScheduleStatus
}

const props = defineProps<Props>()

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

interface Emits {
  (e: 'click'): void
}

const emits = defineEmits<Emits>()

const handleClick = () => {
  emits('click')
}
</script>

<template>
  <div @click="handleClick">
    <div class="flex justify-center">
      <nuxt-img
        provider="cloudFront"
        height="208px"
        fit="contain"
        class="aspect-video max-h-[208px] cursor-pointer object-cover"
        :src="imgSrc"
        :alt="`live-${title}-thumbnail`"
      />
    </div>

    <div class="ml-0 mt-4 flex w-full flex-col gap-2 md:ml-4">
      <div class="mx-4 flex items-center text-sm md:mx-0">
        <div class="grow">
          <span
            :class="{
              'rounded px-2 font-bold': true,
              'border-2 border-orange bg-white text-orange':
                isLiveStreaming(isLiveStatus),
              'border-2 border-main text-main': !isLiveStreaming(isLiveStatus),
            }"
          >
            {{ isLiveStreaming(isLiveStatus) ? "配信中" : "配信予定" }}
          </span>
          <span class="ml-2 text-main after:content-['〜']">{{
            formattedStartAt
          }}</span>
        </div>
      </div>

      <p
        class="mx-4 line-clamp-3 break-words text-[14px] text-main md:mx-0 md:text-[16px]"
      >
        {{ title }}
      </p>
    </div>
  </div>
</template>
