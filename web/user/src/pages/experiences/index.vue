<script setup lang="ts">
import { GoogleMap, MarkerCluster, CustomControl } from 'vue3-google-map'
import { useGeolocation } from '@vueuse/core'
import { useSpotStore } from '~/store/spot'
import type { GoogleMapSearchResult } from '~/types/store'
import { useExperienceStore } from '~/store/experience'

const config = useRuntimeConfig()

const spotStore = useSpotStore()
const { search } = spotStore

const experienceStore = useExperienceStore()
const { fetchExperiences } = experienceStore
const { experiences, experiencesFetchState } = storeToRefs(experienceStore)
const errorMessage = ref<string>('')

const router = useRouter()

// 中心座標の種類
const centerPositionType = ref<'init' | 'geo' | 'search'>('init')

// 現在地の取得
const { coords, error: geoLocationError } = useGeolocation()

// 検索用キーワード
const searchText = ref<string>('')
// 検索結果
const searchResults = ref<GoogleMapSearchResult[]>([])
// 検索結果の座標保存用
const searchResultPosition = ref<{ lat: number, lng: number }>({ lat: 0, lng: 0 })

// 中心座標の算出プロパティ
const center = computed(() => {
  // 検索結果
  if (centerPositionType.value === 'search') {
    return searchResultPosition.value
  }
  // 現在地
  if (centerPositionType.value === 'geo' && coords.value.latitude !== Infinity) {
    return { lat: coords.value.latitude, lng: coords.value.longitude }
  }
  // 初期値
  return { lat: 35.681167, lng: 139.7673068 }
})

const renderer = ref<
  undefined | { render: (obj: { count: number, position: any }) => any }
    >(undefined)

const handleClickSpot = (id: string) => {
  router.push(`/experiences/${id}`)
}

// 検索のエラーメッセージ管理用
const searchResultError = ref<string>('')

const handleSubmitSearchForm = async () => {
  searchResultError.value = ''
  try {
    const results = await search(searchText.value)
    searchResults.value = results
  }
  catch (error) {
    console.error(error)
    if (error instanceof google.maps.MapsRequestError) {
      // 400系のエラー
      searchResultError.value = '検索結果が見つかりません'
    }
    if (error instanceof google.maps.MapsServerError) {
      // 500系のエラー
      errorMessage.value = 'Google Maps API側でエラーが発生しました。'
    }
  }
}

const handleClickSearchResult = (result: GoogleMapSearchResult) => {
  searchResultPosition.value = {
    lat: result.latitude,
    lng: result.longitude,
  }
  centerPositionType.value = 'search'
  searchResults.value = []
}

/**
 * 検索フォームのクリア
 */
const handleClearSearchForm = () => {
  searchResultError.value = ''
  searchResults.value = []
}

const refetchExperiences = async () => {
  try {
    await fetchExperiences(center.value.lng, center.value.lat)
  }
  catch (error) {
    if (error instanceof Error) {
      errorMessage.value = error.message
    }
  }
}

// 中心座標が変更された場合に体験情報を再取得
watch(center, () => {
  refetchExperiences()
})

// ユーザーが位置情報の取得を許可した場合に中心座標を現在地に変更
watch(geoLocationError, () => {
  if (geoLocationError.value === null) {
    centerPositionType.value = 'geo'
  }
})

// 体験情報の取得
onMounted(() => {
  refetchExperiences()
})

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

const mapTypeControlOptions = computed(() => {
  return {
    // 地図の種別の選択位置。10にすると左下に表示される
    position: 10,
  }
})

useSeoMeta({
  title: '体験一覧',
})
</script>

<template>
  <div class="bg-white px-[15px] py-[48px] text-main md:px-[36px]">
    <div class="container mx-auto">
      <div v-if="experiencesFetchState.isLoading">
        <div class="text-center border-t-4 border-main animate-pulse" />
      </div>
    </div>

    <div v-if="errorMessage">
      <div class="border border-orange text-orange bg-white p-4 mb-4">
        <p>{{ errorMessage }}</p>
      </div>
    </div>

    <ClientOnly>
      <GoogleMap
        :api-key="config.public.GOOGLE_MAPS_API_KEY"
        style="width: 100%; height: 700px"
        :center="center"
        :zoom="12"
        :map-type-control-options="mapTypeControlOptions"
        :clickable-icons="false"
      >
        <CustomControl position="TOP_LEFT">
          <div class="relative">
            <the-spot-search-form
              v-model="searchText"
              class="absolute left-2 top-2.5 w-[300px] rounded-full "
              :results="searchResults"
              :error-message="searchResultError"
              @click:result="handleClickSearchResult"
              @clear="handleClearSearchForm"
              @submit="handleSubmitSearchForm"
            />
          </div>
        </CustomControl>
        <MarkerCluster
          :options="
            {
              renderer,
            }"
        >
          <template
            v-for="experience in experiences"
            :key="experience.id"
          >
            <the-experience-marker
              :id="experience.id"
              :longitude="experience.hostLongitude"
              :latitude="experience.hostLatitude"
              :title="experience.title"
              :description="experience.description"
              :thumbnail-url="experience.thumbnailUrl"
              @click:name="handleClickSpot"
            />
          </template>
        </MarkerCluster>
      </GoogleMap>
    </ClientOnly>
  </div>
</template>
