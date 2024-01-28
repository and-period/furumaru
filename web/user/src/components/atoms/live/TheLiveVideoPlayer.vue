<script lang="ts" setup>
import Hls from 'hls.js'

interface Props {
  videoSrc: string
  isArchive: boolean
}

const props = defineProps<Props>()

const videoRef = ref<HTMLVideoElement | null>(null)
const hls = ref<Hls | null>(null)

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
  </div>
</template>
