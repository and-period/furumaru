<script lang="ts" setup>
import Hls from 'hls.js'
import * as dayjs from 'dayjs'

interface Props {
  videoSrc: string
  videoType: string
  title: string
  startAt: number
  endAt: string
  isArchive: boolean
  marcheName: string
  description: string
  address: string
  cnName: string
  cnImgSrc: string
}

const props = defineProps<Props>()

const videoRef = ref<HTMLVideoElement | null>(null)
const showDetail = ref<boolean>(false)

onMounted(() => {
  if (videoRef.value) {
    const video = videoRef.value
    const src = props.videoSrc

    if (Hls.isSupported()) {
      const hls = new Hls({ enableWorker: false })
      hls.loadSource(src)
      hls.attachMedia(video)
      videoRef.value.play()
    } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
      video.src = src
      video.play()
    }
  }
})

const formattedStartAt = computed(() => {
  return dayjs.unix(props.startAt).format('YYYY/MM/DD HH:mm')
})

const handleClickShowDetailButton = () => {
  showDetail.value = !showDetail.value
}
</script>

<template>
  <div>
    <video :src="videoSrc" class="aspect-video w-full" controls playsinline />
    <div class="mt-2 px-4 lg:px-0">
      <div class="flex items-center gap-2">
        <div
          :class="{
            'flex max-w-fit items-center justify-center rounded px-2 font-bold': true,
            'border-2 border-orange bg-orange text-white': isArchive,
            'border-2 border-main text-main': !isArchive,
          }"
        >
          <div class="mr-2 pt-[2px]">
            <the-live-icon />
          </div>
          <div class="align-middle">
            {{ isArchive ? 'LIVE' : '配信予定' }}
          </div>
        </div>
        <div class="text-[14px] tracking-[1.4px] after:content-['〜']">
          {{ formattedStartAt }}
        </div>
      </div>
      <p class="mt-2 line-clamp-1 tracking-[1.6px]">
        {{ title }}
      </p>
      <div class="mt-4 flex items-center gap-2">
        <img
          :src="cnImgSrc"
          class="h-10 w-10 rounded-full"
          :alt="`${cnName}のプロフィール画像`"
        />
        <div class="text-[12px] tracking-[1.2px]">
          <p class="mb-1">{{ marcheName }}/{{ address }}</p>
          <p>コーディネーター：{{ cnName }}</p>
        </div>
      </div>

      <div>
        <p
          v-show="showDetail"
          class="mt-6 whitespace-pre-wrap text-[14px] tracking-[1.4px]"
          v-html="description"
        ></p>
        <button
          class="inline-flex w-full items-center justify-center gap-2 text-[12px] tracking-[1.2px]"
          @click="handleClickShowDetailButton"
        >
          <div>
            {{ showDetail ? 'マルシェの詳細を隠す' : 'マルシェの詳細を見る' }}
          </div>
          <div>
            <the-up-arrow-icon v-if="showDetail" />
            <the-down-arrow-icon v-if="!showDetail" />
          </div>
        </button>
      </div>
    </div>
  </div>
</template>
