<script setup lang="ts">
import { useElementHover, useEventListener } from '@vueuse/core'

interface Props {
  src: string
}

defineProps<Props>()

const videoRef = ref<HTMLVideoElement | null>(null)
const videoPaused = ref<boolean>(true)

const elementRef = ref<HTMLElement | null>(null)
const isHovered = useElementHover(elementRef)

const handleClickPlayVideoButton = () => {
  if (videoRef.value) {
    videoRef.value.play()
    videoPaused.value = false
  }
}

const handleClickPauseVideoButton = () => {
  if (videoRef.value) {
    videoRef.value.pause()
    videoPaused.value = true
  }
}

const handleVideoEnded = () => {
  videoPaused.value = true
}

const duration = ref<number>(0)
const currentTime = ref<number>(0)

const progress = computed(() => {
  if (duration.value > 0) {
    return (currentTime.value / duration.value) * 100
  }
  else {
    return 0
  }
})

onMounted(() => {
  if (videoRef.value) {
    useEventListener(videoRef, 'ended', handleVideoEnded)
    useEventListener(videoRef, 'loadedmetadata', () => {
      duration.value = videoRef.value?.duration || 0
    })
    useEventListener(videoRef, 'timeupdate', () => {
      currentTime.value = videoRef.value?.currentTime || 0
    })
  }
})
</script>

<template>
  <div
    ref="elementRef"
    class="relative w-full h-full"
  >
    <video
      ref="videoRef"
      class="w-full h-full border"
      :src="src"
    />
    <div class="absolute bottom-0 h-[3px] w-full bg-main/50 rounded-lg">
      <div
        class="h-full bg-orange/80 rounded-lg"
        :style="{ width: `${progress}%` }"
      />
    </div>

    <div class="absolute top-[50%] w-full text-center z-10">
      <button
        v-if="videoPaused"
        class="bg-black/60 rounded-full p-4"
        @click="handleClickPlayVideoButton"
      >
        <div class="h-8 w-8 text-white">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            stroke="currentColor"
            class="size-6"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M5.25 5.653c0-.856.917-1.398 1.667-.986l11.54 6.347a1.125 1.125 0 0 1 0 1.972l-11.54 6.347a1.125 1.125 0 0 1-1.667-.986V5.653Z"
            />
          </svg>
        </div>
      </button>
      <button
        v-if="!videoPaused && isHovered"
        class="bg-black/60 rounded-full p-4"
        @click="handleClickPauseVideoButton"
      >
        <div class="h-8 w-8 text-white">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            stroke="currentColor"
            class="size-6"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M15.75 5.25v13.5m-7.5-13.5v13.5"
            />
          </svg>
        </div>
      </button>
    </div>
  </div>
</template>
