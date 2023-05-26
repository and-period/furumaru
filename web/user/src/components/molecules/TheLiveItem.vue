<script setup lang='ts'>
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
      'bg-base': published
    }"
  >
    <div class="w-full relative">
      <div v-if="published" class="h-16 w-16 absolute bg-orange rounded-full flex flex-col justify-center items-center -top-8 -left-8">
        <the-live-icon />
        <div class="text-white uppercase text-xl font-bold">
          live
        </div>
      </div>
      <img :src="imgSrc" :alt="`live-${title}-thumbnail`">
    </div>
    <div class="mt-2 flex flex-col gap-2">
      <div class="inline-flex gap-2 text-sm items-center">
        <span
          :class="{
            'font-bold px-2 rounded': true,
            'bg-white border-2 border-orange text-orange':published,
            'border-2 border-main text-main': !published
          }"
        >
          {{ published? '配信中': '配信予定' }}
        </span>
        <span class="after:content-['〜']">{{ formattedStartAt }}</span>
      </div>
      <p class="line-clamp-3">
        {{ title }}
      </p>
    </div>
  </div>
</template>
