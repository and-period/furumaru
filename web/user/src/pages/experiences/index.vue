<script setup lang="ts">
import { GoogleMap, MarkerCluster } from 'vue3-google-map'

const config = useRuntimeConfig()

const center = { lat: 34.266422, lng: 132.917558 }

const items = [
  {
    id: 1,
    position: { lat: 34.2684527, lng: 132.91340017 },
    title: 'ふじやファーム',
    description: 'レモン農家',
    imgSrc:
      'https://image-cdn.tabechoku.com/crop/w/126/h/120/cw/120/ch/120/images/d4c7ef52c346e13fbbf63eb936a9b0bda59ec418cdcc299c7a6a1248c3ef9cbb.jpeg',
  },
  {
    id: 2,
    position: { lat: 34.63416837, lng: 132.65951729 },
    title:
      '酒蔵見学＆どぶろく体験ツアー　Brewery tour & Doburoku tasting experience',
    description: '酒蔵見学＆どぶろく体験ツアー',
    imgSrc:
      'https://assets.furumaru.and-period.co.jp/products/media/image/gsztHcK7CvWhDyhYRvpTT4.jpg',
  },
  {
    id: 3,
    position: { lat: 35.62539905972506, lng: 139.5175404502735 },
    title: 'よみうりランド',
    description: '遊園地',
    imgSrc:
      'https://lh5.googleusercontent.com/p/AF1QipMcu3Mp8owZq5nbe2FEAtE21wg9T7LA21tjbwgg=w426-h240-k-no',
  },
]

const renderer = ref<
  undefined | { render: (obj: { count: number, position: any }) => any }
    >(undefined)

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
        // adjust zIndex to be above other markers
        zIndex: Number(google.maps.Marker.MAX_ZINDEX) + count,
      }),
    // new google.maps.marker.AdvancedMarkerElement({
    //   title: String(count),
    //   position,
    // },
    // ),
  }
})
</script>

<template>
  <div class="bg-white px-[15px] py-[48px] text-main md:px-[36px]">
    <div class="container mx-auto">
      <ClientOnly>
        <GoogleMap
          :api-key="config.public.GOOGLE_MAPS_API_KEY"
          style="width: 100%; height: 700px"
          :center="center"
          :zoom="10"
        >
          <MarkerCluster :options="{ renderer: renderer }">
            <template
              v-for="item in items"
              :key="item.id"
            >
              <the-experience-marker
                :position="item.position"
                :title="item.title"
                :description="item.description"
                :img-src="item.imgSrc"
              />
            </template>
          </MarkerCluster>
        </GoogleMap>
      </ClientOnly>
    </div>
  </div>
</template>
