<script setup lang="ts">
import { GoogleMap, MarkerCluster } from 'vue3-google-map'
import { useSpotStore } from '~/store/spot'

const config = useRuntimeConfig()

const spotStore = useSpotStore()
const { spots } = storeToRefs(spotStore)
const { fetchSpots } = spotStore

const router = useRouter()

// const center = { lat: 34.266422, lng: 132.917558 }
const center = { lat: 35.681167, lng: 139.7673068 }

const renderer = ref<
  undefined | { render: (obj: { count: number, position: any }) => any }
    >(undefined)

const { status, error } = useAsyncData('spots', () => {
  return fetchSpots(center.lng, center.lat)
})

const handleClickSpot = (id: string) => {
  router.push(`/experiences/${id}`)
}

onMounted(() => {
  const svg = window.btoa(`
  <svg fill="#604C3F" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 240 240">
    <circle cx="120" cy="120" opacity=".8" r="70" />
  </svg>`)

  renderer.value = {
    render: ({ count, position }: { count: number, position: any }) =>
      new google.maps.Marker({
        label: {
          text: String(count),
          color: 'white',
        },
        position,
        icon: {
          url: `data:image/svg+xml;base64,${svg}`,
          scaledSize: new google.maps.Size(75, 75),
        },
        zIndex: Number(google.maps.Marker.MAX_ZINDEX) + count,
      }),
  }
})
</script>

<template>
  <div class="bg-white px-[15px] py-[48px] text-main md:px-[36px]">
    <div class="container mx-auto">
      <div v-if="status === 'pending'">
        <div class="text-center border-t-4 border-main animate-pulse" />
      </div>
    </div>

    <div v-if="status === 'error'">
      <div class="border border-orange text-orange bg-white">
        <p>{{ error }}</p>
      </div>
    </div>

    <template v-if=" status === 'success' ">
      <ClientOnly>
        <GoogleMap
          :api-key="config.public.GOOGLE_MAPS_API_KEY"
          style="width: 100%; height: 700px"
          :center="center"
          :zoom="10"
        >
          <MarkerCluster :options="{ renderer: renderer }">
            <template
              v-for="spot in spots"
              :key="spot.id"
            >
              <the-experience-marker
                :id="spot.id"
                :longitude="spot.longitude"
                :latitude="spot.latitude"
                :name="spot.name"
                :description="spot.description"
                :thumbnail-url="spot.thumbnailUrl"
                @click:name="handleClickSpot"
              />
            </template>
          </MarkerCluster>
        </GoogleMap>
      </ClientOnly>
    </template>
  </div>
</template>
