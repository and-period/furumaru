<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useCoordinatorStore } from '~/store/coordinator'
import type { I18n } from '~/types/locales'

const i18n = useI18n()

const tt = (str: keyof I18n['base']['top']) => {
  return i18n.t(`base.top.${str}`)
}

const route = useRoute()
const router = useRouter()

const coordinatorStore = useCoordinatorStore()

const { fetchCoordinator } = coordinatorStore
const { coordinatorFetchState } = storeToRefs(coordinatorStore)

const { coordinatorInfo, archives, lives, producers, videos, experiences }
  = storeToRefs(coordinatorStore)

// 商品がある生産者を先に、ない生産者を後に並び替え
const sortedProducers = computed(() => {
  return [...producers.value].sort((a, b) => {
    const aHasProducts = a.products && a.products.length > 0
    const bHasProducts = b.products && b.products.length > 0
    if (aHasProducts && !bHasProducts) return -1
    if (!aHasProducts && bHasProducts) return 1
    return 0
  })
})

const id = computed<string>(() => {
  const ids = route.params.id
  if (Array.isArray(ids)) {
    return ids[0]
  }
  else {
    return ids
  }
})

const handleClickLiveItem = (id: string) => {
  router.push(`/live/${id}`)
}

const handleClickProductItem = (id: string) => {
  router.push(`/items/${id}`)
}

const handleClickVideoItem = (id: string) => {
  router.push(`/video/${id}`)
}

const handleClickExperienceItem = (id: string) => {
  router.push(`/experiences/${id}`)
}

useAsyncData(`coordinator-${id.value}`, async () => {
  await fetchCoordinator(id.value)
  return true
})

useSeoMeta({
  title: coordinatorInfo.value
    ? `${coordinatorInfo.value.marcheName}`
    : 'コーディネーターのページ',
})
</script>

<template>
  <template v-if="coordinatorFetchState.isLoading">
    <div
      class="animate-pulse bg-white px-[112px] pb-6 pt-[40px] text-main md:grid md:grid-cols-2"
    >
      <div class="w-full">
        <div class="mx-auto aspect-square h-[500px] w-[500px] bg-slate-100" />
      </div>
      <div class="flex w-full flex-col gap-4">
        <div class="h-[24px] w-[80%] rounded-md bg-slate-100" />
        <div class="h-[24px] w-[60%] rounded-md bg-slate-100" />
      </div>
    </div>
  </template>
  <template v-else>
    <div class="min-h-screen bg-base">
      <div class="mx-auto w-full max-w-7xl px-4 text-main md:px-8">
        <!-- ヘッダー部分 -->
        <div class="relative h-[160px] w-full overflow-hidden rounded-lg shadow-lg md:h-[320px]">
          <template v-if="coordinatorInfo.headerUrl">
            <img
              class="h-full w-full object-cover transition-transform duration-300 hover:scale-105"
              :src="coordinatorInfo.headerUrl"
              :alt="coordinatorInfo.marcheName + 'のヘッダー画像'"
            >
          </template>
          <template v-else>
            <div class="flex h-full w-full items-center justify-center bg-gradient-to-br from-gray-100 to-gray-300">
              <div class="text-center text-gray-400">
                <p class="text-sm">
                  ヘッダー画像
                </p>
              </div>
            </div>
          </template>
        </div>
        <div
          class="relative bottom-[50px] md:bottom-20 md:grid md:grid-cols-7 md:gap-12 bg-base"
        >
          <div class="col-span-2">
            <div class="flex justify-center relative bottom-[30px] md:bottom-[40px]">
              <img
                :src="coordinatorInfo.thumbnailUrl"
                class="block aspect-square w-[120px] rounded-full border-4 border-white md:w-[168px] shadow-lg"
              >
            </div>
            <p
              class="relative bottom-[20px] md:bottom-[25px] text-center text-[16px] font-bold tracking-[2.0px] md:text-[20px]"
            >
              {{ coordinatorInfo.marcheName }}
            </p>
            <div
              class="relative bottom-[16px] md:bottom-[20px] flex justify-center text-[12px] tracking-[1.4px] md:text-[14px]"
            >
              <p>{{ coordinatorInfo.prefecture }}</p>
              <p class="pl-2">
                {{ coordinatorInfo.city }}
              </p>
            </div>
            <div class="relative bottom-[12px] md:bottom-[16px] flex justify-center tracking-[2.4px]">
              <p class="ml-2 text-[16px] font-bold md:text-[24px]">
                {{ coordinatorInfo.username }}
              </p>
            </div>
            <p
              class="mx-4 text-[14px] tracking-[1.4px] md:mx-0 md:text-[16px] md:tracking-[1.6px]"
            >
              {{ coordinatorInfo.profile }}
            </p>
            <hr class="m-4 border-dashed border-main md:mx-0">
            <div class="mx-4 grid grid-cols-3 md:mx-0">
              <div class="col-span-2 text-[14px] md:text-[16px]">
                SNSでフォローする
              </div>
              <div class="flex justify-end">
                <a
                  :href="
                    'https://www.instagram.com/' + coordinatorInfo.instagramId
                  "
                  target="_blank"
                >
                  <svg
                    class="mr-[16px] h-[24p] w-[24px] md:h-[32px] md:w-[32px]"
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 32 32"
                    fill="none"
                  >
                    <g clip-path="url(#clip0_2557_4584)">
                      <path
                        d="M16 2.88345C20.2725 2.88345 20.778 2.89997 22.4655 2.97682C24.0254 3.04795 24.8727 3.30835 25.4366 3.52747C26.1835 3.81772 26.7164 4.1645 27.2766 4.72404C27.8368 5.28422 28.1835 5.81708 28.4732 6.56399C28.6923 7.12798 28.9527 7.97523 29.0238 9.53509C29.1007 11.222 29.1172 11.7282 29.1172 16.0006C29.1172 20.2731 29.1007 20.7787 29.0238 22.4662C28.9527 24.026 28.6923 24.8733 28.4732 25.4373C28.1829 26.1842 27.8361 26.7171 27.2766 27.2772C26.7164 27.8374 26.1835 28.1842 25.4366 28.4738C24.8727 28.6929 24.0254 28.9533 22.4655 29.0245C20.7787 29.1013 20.2725 29.1178 16 29.1178C11.7275 29.1178 11.2213 29.1013 9.53445 29.0245C7.97459 28.9533 7.12734 28.6929 6.56335 28.4738C5.81645 28.1835 5.28358 27.8368 4.7234 27.2772C4.16323 26.7171 3.81645 26.1842 3.52683 25.4373C3.30772 24.8733 3.04732 24.026 2.97618 22.4662C2.89933 20.7793 2.88282 20.2731 2.88282 16.0006C2.88282 11.7282 2.89933 11.2226 2.97618 9.53509C3.04732 7.97523 3.30772 7.12798 3.52683 6.56399C3.81708 5.81708 4.16386 5.28422 4.7234 4.72404C5.28295 4.16386 5.81645 3.81708 6.56335 3.52747C7.12734 3.30835 7.97459 3.04795 9.53445 2.97682C11.2213 2.89997 11.7275 2.88345 16 2.88345ZM16 0.000635122C11.6545 0.000635122 11.1096 0.0190537 9.40298 0.0971737C7.70022 0.174659 6.53668 0.445221 5.51921 0.840902C4.46681 1.24992 3.57447 1.79676 2.6853 2.68657C1.79549 3.57637 1.24865 4.46809 0.839632 5.52048C0.444586 6.53795 0.174023 7.70086 0.0965386 9.40362C0.0184185 11.1102 0 11.6551 0 16.0006C0 20.3461 0.0184185 20.8911 0.0965386 22.5976C0.174023 24.3004 0.444586 25.464 0.840267 26.4814C1.24929 27.5338 1.79613 28.4262 2.68593 29.3153C3.57574 30.2051 4.46745 30.752 5.51985 31.161C6.53731 31.5567 7.70086 31.8272 9.40362 31.9047C11.1102 31.9828 11.6551 32.0013 16.0006 32.0013C20.3461 32.0013 20.8911 31.9828 22.5976 31.9047C24.3004 31.8272 25.464 31.5567 26.4814 31.161C27.5338 30.752 28.4262 30.2051 29.3153 29.3153C30.2051 28.4255 30.752 27.5338 31.161 26.4814C31.5567 25.464 31.8272 24.3004 31.9047 22.5976C31.9828 20.8911 32.0013 20.3461 32.0013 16.0006C32.0013 11.6551 31.9828 11.1102 31.9047 9.40362C31.8272 7.70086 31.5567 6.53731 31.161 5.51985C30.752 4.46745 30.2051 3.5751 29.3153 2.68593C28.4255 1.79613 27.5338 1.24929 26.4814 0.840267C25.464 0.444586 24.3004 0.174023 22.5976 0.0965386C20.8911 0.0184185 20.3461 0 16.0006 0L16 0.000635122Z"
                        fill="#604C3F"
                      />
                      <path
                        d="M16.0001 7.78516C11.4622 7.78516 7.78418 11.4638 7.78418 16.0011C7.78418 20.5384 11.4628 24.217 16.0001 24.217C20.5374 24.217 24.2161 20.5384 24.2161 16.0011C24.2161 11.4638 20.5374 7.78516 16.0001 7.78516ZM16.0001 21.3349C13.0544 21.3349 10.667 18.9468 10.667 16.0017C10.667 13.0567 13.0551 10.6686 16.0001 10.6686C18.9452 10.6686 21.3332 13.0567 21.3332 16.0017C21.3332 18.9468 18.9452 21.3349 16.0001 21.3349Z"
                        fill="#604C3F"
                      />
                      <path
                        d="M24.5411 9.37901C25.6014 9.37901 26.461 8.51941 26.461 7.45904C26.461 6.39866 25.6014 5.53906 24.5411 5.53906C23.4807 5.53906 22.6211 6.39866 22.6211 7.45904C22.6211 8.51941 23.4807 9.37901 24.5411 9.37901Z"
                        fill="#604C3F"
                      />
                    </g>
                    <defs>
                      <clipPath id="clip0_2557_4584">
                        <rect
                          width="32"
                          height="32"
                          fill="white"
                        />
                      </clipPath>
                    </defs>
                  </svg>
                </a>
                <a
                  :href="'https://www.facebook.com/' + coordinatorInfo.facebookId"
                  target="_blank"
                >
                  <svg
                    class="h-[24px] w-[24px] md:h-[32px] md:w-[32px]"
                    xmlns="http://www.w3.org/2000/svg"
                    width="32"
                    height="32"
                    viewBox="0 0 32 32"
                    fill="none"
                  >
                    <g clip-path="url(#clip0_2557_4588)">
                      <path
                        d="M16 0C7.16344 0 0 7.20722 0 16.0978C0 24.1325 5.85094 30.7924 13.5 32V20.751H9.4375V16.0978H13.5V12.5512C13.5 8.51673 15.8888 6.2882 19.5434 6.2882C21.2941 6.2882 23.125 6.60261 23.125 6.60261V10.5642H21.1075C19.12 10.5642 18.5 11.8051 18.5 13.0782V16.0978H22.9375L22.2281 20.751H18.5V32C26.1491 30.7924 32 24.1328 32 16.0978C32 7.20722 24.8366 0 16 0Z"
                        fill="#604C3F"
                      />
                    </g>
                    <defs>
                      <clipPath id="clip0_2557_4588">
                        <rect
                          width="32"
                          height="32"
                          fill="white"
                        />
                      </clipPath>
                    </defs>
                  </svg>
                </a>
              </div>
            </div>
            <!-- プロモーション動画 -->
            <div
              v-if="coordinatorInfo.promotionVideoUrl"
              class="mx-4 mt-4 md:mx-0"
            >
              <hr class="mb-4 border-dashed border-main">
              <div class="mb-4 text-[14px] md:text-[16px]">
                紹介動画
              </div>
              <video
                class="w-full rounded-lg"
                controls
                :src="coordinatorInfo.promotionVideoUrl"
              >
                お使いのブラウザは動画再生をサポートしていません。
              </video>
            </div>
          </div>
          <div class="static pt-[16px] text-main md:col-span-5 md:pt-[100px]">
            <div class="flex w-full px-4 md:px-0">
              <img
                src="/img/coordinator/marche.svg"
                class="z-10 w-full"
              >
            </div>
            <div
              class="relative bottom-4 z-0 mx-4 bg-white pb-10 pt-[20px] md:pt-[65px] md:bottom-8 md:mx-0 md:w-full"
            >
              <div
                v-if="lives.length > 0"
                class="px-4"
              >
                <div
                  class="flex justify-center rounded-3xl bg-base py-[3px] text-[16px] md:mx-auto"
                >
                  配信中・配信予定のライブ
                </div>
              </div>
              <div
                v-if="lives.length > 0"
                class="mx-4 grid grid-cols-1 gap-8 bg-white pt-4 md:mx-auto md:grid-cols-2"
              >
                <the-coordinator-live-item
                  v-for="liveItem in lives"
                  :id="liveItem.scheduleId"
                  :key="liveItem.scheduleId"
                  :title="liveItem.title"
                  :img-src="liveItem.thumbnailUrl"
                  :start-at="liveItem.startAt"
                  :is-live-status="liveItem.status"
                  class="pt-10"
                  @click="handleClickLiveItem(liveItem.scheduleId)"
                />
              </div>
              <div class="my-8 px-4">
                <div
                  v-if="experiences.length > 0"
                  class="flex items-center justify-center gap-2 rounded-3xl bg-gradient-to-r from-blue-400 to-blue-500 py-3 px-4 text-[16px] text-white font-bold shadow-lg md:gap-3 md:py-4 md:px-6 md:text-[18px]"
                >
                  <svg
                    class="w-6 h-6 md:w-7 md:h-7"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path d="M13 6a3 3 0 11-6 0 3 3 0 016 0zM18 8a2 2 0 11-4 0 2 2 0 014 0zM14 15a4 4 0 00-8 0v3h8v-3zM6 8a2 2 0 11-4 0 2 2 0 014 0zM16 18v-3a5.972 5.972 0 00-.75-2.906A3.005 3.005 0 0119 15v3h-3zM4.75 12.094A5.973 5.973 0 004 15v3H1v-3a3 3 0 013.75-2.906z" />
                  </svg>
                  <span class="tracking-wide">体験をする</span>
                </div>
              </div>
              <div
                v-if="experiences.length > 0"
                class="mx-auto grid grid-cols-1 gap-8 bg-white p-4 md:grid-cols-2"
              >
                <div
                  v-for="experience in experiences"
                  :key="experience.id"
                  class="cursor-pointer rounded-xl overflow-hidden shadow-md hover:shadow-lg transition-shadow duration-300"
                  @click="handleClickExperienceItem(experience.id)"
                >
                  <div class="aspect-video relative overflow-hidden">
                    <img
                      :src="experience.thumbnailUrl"
                      :alt="experience.title"
                      class="w-full h-full object-cover"
                    >
                  </div>
                  <div class="p-4">
                    <h3 class="text-[16px] font-bold tracking-[1.6px] line-clamp-2">
                      {{ experience.title }}
                    </h3>
                    <p class="text-[14px] text-gray-600 mt-2 line-clamp-2">
                      {{ experience.description }}
                    </p>
                  </div>
                </div>
              </div>
              <div class="my-8 px-4">
                <div
                  v-if="videos.length > 0"
                  class="flex items-center justify-center gap-2 rounded-3xl bg-gradient-to-r from-red-400 to-pink-500 py-3 px-4 text-[16px] text-white font-bold shadow-lg md:gap-3 md:py-4 md:px-6 md:text-[18px]"
                >
                  <svg
                    class="w-6 h-6 md:w-7 md:h-7"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path d="M2 6a2 2 0 012-2h6a2 2 0 012 2v8a2 2 0 01-2 2H4a2 2 0 01-2-2V6zM14.553 7.106A1 1 0 0014 8v4a1 1 0 00.553.894l2 1A1 1 0 0018 13V7a1 1 0 00-1.447-.894l-2 1z" />
                  </svg>
                  <span class="tracking-wide">関連動画</span>
                </div>
              </div>
              <div
                v-if="videos.length > 0"
                class="mx-auto grid grid-cols-1 gap-8 bg-white p-4 md:grid-cols-2"
              >
                <the-video-item
                  v-for="video in videos"
                  :id="video.id"
                  :key="video.id"
                  :title="video.title"
                  :img-src="video.thumbnailUrl"
                  :width="368"
                  :end-at="video.publishedAt"
                  :archived-stream-text="tt('archivedStreamText')"
                  class="cursor-pointer rounded-xl overflow-hidden shadow-md hover:shadow-lg transition-shadow duration-300"
                  @click="handleClickVideoItem(video.id)"
                />
              </div>
              <div class="my-8 px-4">
                <div
                  class="flex items-center justify-center gap-2 rounded-3xl bg-gradient-to-r from-purple-400 to-indigo-500 py-3 px-4 text-[16px] text-white font-bold shadow-lg md:gap-3 md:py-4 md:px-6 md:text-[18px]"
                >
                  <svg
                    class="w-6 h-6 md:w-7 md:h-7"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path
                      fill-rule="evenodd"
                      d="M10 18a8 8 0 100-16 8 8 0 000 16zM9.555 7.168A1 1 0 008 8v4a1 1 0 001.555.832l3-2a1 1 0 000-1.664l-3-2z"
                      clip-rule="evenodd"
                    />
                  </svg>
                  <span class="tracking-wide">過去の配信</span>
                </div>
              </div>
              <div
                class="mx-auto grid grid-cols-1 gap-8 bg-white p-4 md:grid-cols-2"
              >
                <the-coordinator-archive-item
                  v-for="archive in archives"
                  :id="archive.scheduleId"
                  :key="archive.scheduleId"
                  :title="archive.title"
                  :img-src="archive.thumbnailUrl"
                  :width="320"
                  class="cursor-pointer rounded-xl overflow-hidden shadow-md hover:shadow-lg transition-shadow duration-300"
                  @click="handleClickLiveItem(archive.scheduleId)"
                />
              </div>
            </div>
            <div
              class="flex w-full flex-nowrap justify-between text-main md:gap-[70px]"
            >
              <img
                class="w-[120px] md:w-[260px]"
                src="/img/coordinator/left.svg"
              >
              <p
                class="whitespace-nowrap pt-5 text-[14px] font-bold md:text-[20px]"
              >
                生産者一覧
              </p>
              <img
                class="w-[120px] md:w-[260px]"
                src="/img/coordinator/right.svg"
              >
            </div>
            <div
              class="grid grid-cols-1 gap-x-4 gap-y-[80px] pt-[80px] md:grid-cols-2 md:pt-[100px] lg:gap-x-6"
            >
              <the-producer-list
                v-for="producer in sortedProducers"
                :id="producer.id"
                :key="producer.id"
                :name="producer.username"
                :profile="producer.profile"
                :img-src="producer.thumbnailUrl"
                :products="producer.products"
                @click:product-item="handleClickProductItem"
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  </template>
</template>
