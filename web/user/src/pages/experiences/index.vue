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

// GoogleMapコンポーネントへの参照
const googleMapRef = ref<any>(null)

// 中心座標の種類
const centerPositionType = ref<'init' | 'geo' | 'search' | 'manual'>('init')

// 現在地の取得
const { coords, error: geoLocationError } = useGeolocation()

// 検索用キーワード
const searchText = ref<string>('')
// 検索結果
const searchResults = ref<GoogleMapSearchResult[]>([])
// 検索結果の座標保存用
const searchResultPosition = ref<{ lat: number, lng: number }>({ lat: 0, lng: 0 })

// マップの現在中心座標を保持する
const mapCenter = ref<{ lat: number, lng: number }>({ lat: 35.681167, lng: 139.7673068 })

// 中心座標の算出プロパティ
const center = computed(() => {
  // マップ操作後は現在のマップ中心を維持
  if (centerPositionType.value === 'manual') {
    return mapCenter.value
  }

  // 検索結果
  if (centerPositionType.value === 'search') {
    return searchResultPosition.value
  }

  // 現在地（かつマップ操作をしていない場合）
  if (centerPositionType.value === 'geo' && coords.value.latitude !== Infinity) {
    return {
      lat: coords.value.latitude,
      lng: coords.value.longitude,
    }
  }

  // 初期値
  return { lat: 35.681167, lng: 139.7673068 }
})

// 前回のマニュアル操作フラグ
const hasManuallyMoved = ref<boolean>(false)

// 現在地情報が更新されたときに中心を自動更新するかどうかのフラグ
const autoUpdateCenterOnGeoChange = ref<boolean>(true)

// マップの中心座標を更新する関数
const updateMapCenter = () => {
  if (googleMapRef.value && googleMapRef.value.map) {
    const gmap = googleMapRef.value.map
    const center = gmap.getCenter()
    if (center) {
      mapCenter.value = {
        lat: center.lat(),
        lng: center.lng(),
      }
    }
  }
}

const renderer = ref<
  undefined | { render: (obj: { count: number, position: any }) => any }
    >(undefined)

const handleClickSpot = (id: string) => {
  router.push(`/experiences/${id}`)
}

// 同一座標の体験に小さなオフセットを追加して個別に表示する
const experiencesWithOffset = computed(() => {
  const coordGroups = new Map<string, Array<{ experience: any, index: number }>>()
  
  // 同一座標の体験をグループ化
  experiences.value.forEach((experience, index) => {
    const coordKey = `${experience.hostLatitude},${experience.hostLongitude}`
    if (!coordGroups.has(coordKey)) {
      coordGroups.set(coordKey, [])
    }
    coordGroups.get(coordKey)!.push({ experience, index })
  })
  
  // オフセットを適用
  return experiences.value.map((experience) => {
    const coordKey = `${experience.hostLatitude},${experience.hostLongitude}`
    const group = coordGroups.get(coordKey)!
    
    if (group.length === 1) {
      // 単独の場合はオフセットなし
      return experience
    }
    
    // 複数の場合は円形にオフセットを配置
    const itemIndex = group.findIndex(item => item.experience.id === experience.id)
    const angleStep = (2 * Math.PI) / group.length
    const radius = 0.0005 // 約50メートルのオフセット
    const angle = angleStep * itemIndex
    
    const offsetLat = radius * Math.cos(angle)
    const offsetLng = radius * Math.sin(angle)
    
    return {
      ...experience,
      hostLatitude: experience.hostLatitude + offsetLat,
      hostLongitude: experience.hostLongitude + offsetLng,
    }
  })
})

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

// 現在地情報の変更を監視
watch(coords, () => {
  // マップを操作していない場合か、自動更新フラグがtrueの場合のみ現在地に中心を移動
  if ((centerPositionType.value === 'geo' || !hasManuallyMoved.value) && autoUpdateCenterOnGeoChange.value) {
    if (coords.value.latitude !== Infinity && coords.value.longitude !== Infinity) {
      mapCenter.value = {
        lat: coords.value.latitude,
        lng: coords.value.longitude,
      }
    }
  }
})

// ユーザーが位置情報の取得を許可した場合に中心座標を現在地に変更
watch(geoLocationError, () => {
  if (geoLocationError.value === null && !hasManuallyMoved.value) {
    centerPositionType.value = 'geo'
    autoUpdateCenterOnGeoChange.value = true
  }
})

// 体験情報の取得
onMounted(() => {
  refetchExperiences()

  // ユーザーが位置情報の取得を許可した場合に中心座標を現在地に変更
  if (geoLocationError.value === null && !hasManuallyMoved.value) {
    centerPositionType.value = 'geo'
    autoUpdateCenterOnGeoChange.value = true
  }

  // マップコンポーネントが準備できた後に実行される処理
  nextTick(() => {
    // GoogleMapコンポーネントのインスタンスが利用可能になるまで待機
    const checkGoogleMapReady = setInterval(() => {
      if (googleMapRef.value && googleMapRef.value.map) {
        clearInterval(checkGoogleMapReady)
        // 初期状態のマップ中心を保存
        updateMapCenter()
      }
    }, 100)
  })

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

/**
 * マップがドラッグされた時の処理
 */
const onMapDragStart = (e: any) => {
  // マップの現在中心を取得
  updateMapCenter()
  centerPositionType.value = 'manual'
  hasManuallyMoved.value = true
  // マップ操作後は現在地への自動更新を無効化
  autoUpdateCenterOnGeoChange.value = false
}

/**
 * マップの操作が完了した時の処理
 */
const onMapIdle = (e: any) => {
  // マップの操作（ズームやパンなど）が完了した時に現在位置を保存
  updateMapCenter()

  // ユーザーがマップを操作した場合は、centerPositionTypeを'manual'に設定
  if (hasManuallyMoved.value) {
    centerPositionType.value = 'manual'
  }
}

/**
 * 現在地に中心を移動する
 */
const centerToCurrentLocation = () => {
  if (coords.value.latitude !== Infinity && coords.value.longitude !== Infinity) {
    centerPositionType.value = 'geo'
    hasManuallyMoved.value = false
    autoUpdateCenterOnGeoChange.value = true
    mapCenter.value = {
      lat: coords.value.latitude,
      lng: coords.value.longitude,
    }
  }
}

/**
 * マップのズームが変更された時の処理
 */
const onMapZoomChanged = () => {
  hasManuallyMoved.value = true
  autoUpdateCenterOnGeoChange.value = false
  centerPositionType.value = 'manual'
}

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
        ref="googleMapRef"
        :api-key="config.public.GOOGLE_MAPS_API_KEY"
        style="width: 100%; height: 700px"
        :center="center"
        :zoom="12"
        :map-type-control-options="mapTypeControlOptions"
        :clickable-icons="false"
        :gesture-handling="'greedy'"
        @dragstart="onMapDragStart"
        @idle="onMapIdle"
        @zoom_changed="onMapZoomChanged"
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

        <CustomControl position="RIGHT_BOTTOM">
          <button
            aria-label="現在地に移動"
            class="bg-white rounded-full p-2 shadow-md m-3"
            @click="centerToCurrentLocation"
          >
            <div class="h-6 w-6">
              <the-my-location-icon />
            </div>
          </button>
        </CustomControl>
        <MarkerCluster
          :options="
            {
              renderer,
            }"
        >
          <template
            v-for="experience in experiencesWithOffset"
            :key="experience.id"
          >
            <the-experience-marker
              :id="experience.id"
              :longitude="experience.hostLongitude"
              :latitude="experience.hostLatitude"
              :title="experience.title"
              :description="experience.description"
              :thumbnail-url="experience.thumbnailUrl"
              :promotion-video-url="experience.promotionVideoUrl"
              @click:name="handleClickSpot"
            />
          </template>
        </MarkerCluster>
      </GoogleMap>
    </ClientOnly>
  </div>
</template>
