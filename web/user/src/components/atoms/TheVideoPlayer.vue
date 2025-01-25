<script lang="ts" setup>
interface Props {
  videoSrc: string
}

const props = defineProps<Props>()

const videoRef = ref<HTMLVideoElement | null>(null)

onMounted(() => {
  if (videoRef.value) {
    const video = videoRef.value
    const src = props.videoSrc
    video.src = src
    video.play()
  }
})

onUnmounted(() => {
  if (videoRef.value) {
    videoRef.value.pause()
    videoRef.value.src = ''
  }
})

defineExpose({ videoRef })
</script>

<template>
  <video
    ref="videoRef"
    :src="`${videoSrc}#t=0.1`"
    class="aspect-video w-full bg-black"
    controls
    controlsList="nodownload"
    playsinline
    autoPictureInPicture
  />
</template>
