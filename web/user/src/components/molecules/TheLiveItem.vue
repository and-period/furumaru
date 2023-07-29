<script setup lang="ts">
import * as dayjs from 'dayjs'

interface Props {
  id: string
  title: string
  imgSrc: string
  startAt: number
  published: boolean
}

const props = defineProps<Props>()

const formattedStartAt = computed(() => {
  return dayjs.unix(props.startAt).format('YYYY/MM/DD HH:mm')
})
</script>

<template>
  <div
    :class="{
      'p-4 text-main': true,
      'bg-base': published,
    }"
  >
    <div class="relative w-full">
      <div
        v-if="published"
        class="absolute -left-8 -top-8 flex h-16 w-16 flex-col items-center justify-center rounded-full bg-orange"
      >
        <the-live-icon />
        <div class="text-xl font-bold uppercase text-white">live</div>
      </div>
      <img :src="imgSrc" :alt="`live-${title}-thumbnail`" />
    </div>
    <div class="mt-2 flex flex-col gap-2">
      <div class="flex items-center text-sm">
        <div class="grow">
          <span
            :class="{
              'rounded px-2 font-bold': true,
              'border-2 border-orange bg-white text-orange': published,
              'border-2 border-main text-main': !published,
            }"
          >
            {{ published ? '配信中' : '配信予定' }}
          </span>
          <span class="ml-2 after:content-['〜']">{{ formattedStartAt }}</span>
        </div>
        <button class="h-4 w-4 hover:scale-110">
          <the-ellipsis-vertical-icon class="h-5 w-5" />
        </button>
      </div>
      <p class="line-clamp-3">
        {{ title }}
      </p>
    </div>
  </div>
</template>
