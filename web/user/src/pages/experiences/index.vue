<script setup lang="ts">
import { GoogleMap, Marker } from 'vue3-google-map'

const config = useRuntimeConfig()

const center = { lat: 40.689247, lng: -74.044502 }
</script>

<template>
  <div
    class="flex flex-col bg-white px-[15px] py-[48px] text-main md:px-[36px]"
  >
    <GoogleMap
      :api-key="config.public.GOOGLE_MAPS_API_KEY"
      style="width: 100%; height: 500px"
      :center="center"
      :zoom="15"
    >
      <Marker :options="{ position: center }" />
    </GoogleMap>
  </div>
</template>
