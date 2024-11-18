<script setup lang="ts">
import { CustomMarker, InfoWindow } from 'vue3-google-map'

interface Props {
  position: {
    lat: number
    lng: number
  }
  title: string
  description: string
  imgSrc: string
}

defineProps<Props>()

const isShowInfoWindow = ref<boolean>(false)
const infoWindowRef = ref<InstanceType<typeof InfoWindow> | null>(null)

const handleClickMarker = () => {
  isShowInfoWindow.value = !isShowInfoWindow.value
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
        class="text-main grid grid-cols-5 gap-2 max-w-sm"
      >
        <div class=" col-span-1">
          <img
            :src="imgSrc"
            class="w-full object-cover"
          >
        </div>
        <div class=" tracking-wider col-span-4 flex flex-col gap-2">
          <div class="font-semibold">
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
