<script setup lang="ts">
import { CustomMarker, InfoWindow } from 'vue3-google-map'

interface Props {
  longitude: number
  latitude: number
  id: string
  title: string
  description: string
  thumbnailUrl: string
  promotionVideoUrl: string
}

const props = defineProps<Props>()

interface Emits {
  (e: 'click:name', id: string): void
}

const emits = defineEmits<Emits>()

const isShowInfoWindow = ref<boolean>(false)
const infoWindowRef = ref<InstanceType<typeof InfoWindow> | null>(null)

const position = computed(() => {
  return {
    lat: props.latitude,
    lng: props.longitude,
  }
})

const handleClickMarker = () => {
  isShowInfoWindow.value = !isShowInfoWindow.value
}

const handleClickName = () => {
  emits('click:name', props.id)
}

let initChange = false

watch(isShowInfoWindow, (newValue) => {
  if (initChange) {
    return
  }
  else if (newValue) {
    // 初回の変更時だけInfoWindowを強制的に閉じる
    initChange = true
    if (infoWindowRef.value) {
      isShowInfoWindow.value = false
      infoWindowRef.value.close()
    }
  }
})
</script>

<template>
  <CustomMarker
    :options="{ position, anchorPoint: 'TOP_CENTER' }"
    @click.prevent="handleClickMarker"
  >
    <div
      class="bg-base w-10 h-10 rounded-full flex items-center justify-center border-main p-1 border-2"
    >
      <the-furuneko-icon />
    </div>
    <InfoWindow
      ref="infoWindowRef"
      v-model="isShowInfoWindow"
      :options="{ position }"
    >
      <div
        v-show="isShowInfoWindow"
        class="text-main grid grid-cols-1 md:grid-cols-5 gap-2 max-w-sm"
      >
        <div class="col-span-1 md:col-span-1">
          <div v-if="promotionVideoUrl">
            <video
              :src="promotionVideoUrl"
              class="w-full max-h-[150px] mt-2"
              :title="`${title}の動画`"
              muted
              playsinline
              autoplay
              loop
            />
          </div>
          <div v-else>
            <img
              v-if="thumbnailUrl"
              :src="thumbnailUrl"
              class="w-full max-h-[150px] mt-2"
              :title="`${title}の画像`"
            >
          </div>
        </div>

        <div
          class="tracking-wider md:col-span-4 flex flex-col gap-2 ml-0 md:ml-4"
        >
          <div
            class="font-semibold cursor-pointer hover:underline"
            @click="handleClickName"
          >
            {{ title }}
          </div>
          <div>
            {{ description }}
          </div>
        </div>
      </div>
    </InfoWindow>
  </CustomMarker>
</template>
