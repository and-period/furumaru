<script lang="ts" setup>
import Hls from 'hls.js'
import dayjs from 'dayjs'

interface Props {
  videoSrc: string
  title: string
  startAt: number
  endAt: number
  isLiveStreaming: boolean
  isArchive: boolean
  marcheName: string
  description: string
  address: string
  cnName: string
  cnImgSrc: string
  cordinatorId: string
}

const props = defineProps<Props>()

interface Emits {
  (e: 'click:cordinator', id: string): void
}

const emits = defineEmits<Emits>()

const videoRef = ref<HTMLVideoElement | null>(null)
const showDetail = ref<boolean>(false)
const hls = ref<Hls | null>(null)

const formattedStartAt = computed(() => {
  return dayjs.unix(props.startAt).format('YYYY/MM/DD HH:mm')
})

const handleClickShowDetailButton = () => {
  showDetail.value = !showDetail.value
}

const handleCLickCorodinator = () => {
  emits('click:cordinator', props.cordinatorId)
}

onMounted(() => {
  if (videoRef.value) {
    const video = videoRef.value
    const src = props.videoSrc

    if (props.isArchive) {
      return
    }

    if (Hls.isSupported()) {
      hls.value = new Hls({ enableWorker: false })
      hls.value.loadSource(src)
      hls.value.attachMedia(video)
      videoRef.value.play()
    } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
      video.src = src
      video.play()
    }
  }
})

onUnmounted(() => {
  if (hls.value) {
    hls.value.destroy()
  }
  if (videoRef.value) {
    videoRef.value.pause()
    videoRef.value.src = ''
  }
})
</script>

<template>
  <div>
    <video
      ref="videoRef"
      :src="videoSrc"
      class="aspect-video w-full"
      controls
      playsinline
    />
    <div class="mt-2 px-4">
      <div class="flex items-center gap-2">
        <template v-if="isArchive">
          <div
            class="flex max-w-fit items-center justify-center rounded border-2 border-main px-2 font-bold text-main"
          >
            アーカイブ
          </div>
        </template>
        <template v-else-if="isLiveStreaming">
          <div
            class="flex max-w-fit items-center justify-center rounded border-2 border-orange bg-orange px-2 font-bold text-white"
          >
            <div class="mr-2 pt-[2px]">
              <the-live-icon />
            </div>
            <div class="align-middle">LIVE</div>
          </div>
        </template>

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
          class="h-10 w-10 rounded-full hover:cursor-pointer"
          :alt="`${cnName}のプロフィール画像`"
          @click="handleCLickCorodinator"
        />
        <div class="text-[12px] tracking-[1.2px]">
          <p class="mb-1">{{ marcheName }}/{{ address }}</p>
          <p>
            コーディネーター：
            <span
              class="cursor-pointer hover:underline"
              @click="handleCLickCorodinator"
              >{{ cnName }}</span
            >
          </p>
        </div>
      </div>

      <div>
        <p
          v-if="showDetail"
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
