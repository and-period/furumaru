<script setup lang="ts">
import * as dayjs from 'dayjs'
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

const isLiveStreaming = (stasus: ScheduleStatus) => {
  if (stasus === ScheduleStatus.LIVE) {
    return true
  } else {
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
  <div class="min-w-[368px]" @click="handleClick">
      <div class="w-full">
        <img
          class="w-[320px] object-cover"
          :src="imgSrc"
          :alt="`live-${title}-thumbnail`"
        />
      </div>

      <div class="mt-4 flex w-full flex-col gap-2">
        <div class="flex items-center text-sm">
          <div class="grow">
            <span
              :class="{
                'rounded px-2 font-bold': true,
                'border-2 border-orange bg-white text-orange': isLiveStreaming,
                'border-2 border-main text-main': !isLiveStreaming,
              }"
            >
              {{ isLiveStreaming(isLiveStatus) ? '配信中' : '配信予定' }}
            </span>
            <span class="ml-2 after:content-['〜']">{{
              formattedStartAt
            }}</span>
          </div>
        </div>

        <p class="line-clamp-3">
          {{ title }}
        </p>
      </div>
  </div>
</template>
