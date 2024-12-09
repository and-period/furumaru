<script setup lang="ts">
import type { Snackbar } from '~/types/props'
import { ExperienceStatus } from '~/types/api'

const i18n = useI18n()

const route = useRoute()

const snackbarItems = ref<Snackbar[]>([])

const selectedMediaIndex = ref<number>(-1)

const selectMediaSrcUrl = computed<string>(() => {
  return selectedMediaIndex.value === -1
    ? items.experience.promotionVideoUrl
    : items.experience.media[selectedMediaIndex.value].url
})

const handleClickMediaItem = (index: number) => {
  selectedMediaIndex.value = index
}

const itemThumbnailAlt = computed<string>(() => {
  return i18n.t('items.list.itemThumbnailAlt', {
    itemName: items.experience.title,
  })
})

const canAddCart = computed<boolean>(() => {
  if (items.experience) {
    return (
      items.experience.status === ExperienceStatus.ACCEPTING
    )
  }
  else {
    return false
  }
})

const items = {
  experience: {
    id: 'kSByoE6FetnPs5Byk3a9Zx',
    title: '農業体験',
    description: '農業体験の説明',
    status: 1,
    coordinatorId: 'kSByoE6FetnPs5Byk3a9Zx',
    producerId: 'kSByoE6FetnPs5Byk3a9Zx',
    experienceTypeId: 'kSByoE6FetnPs5Byk3a9Zx',
    thumbnailUrl: 'https://example.com/image.jpg',
    media: [
      { url: 'https://as2.ftcdn.net/v2/jpg/02/31/56/61/1000_F_231566167_DcxyiS11UCKdIoFpPdkXFpAzeVhh6qFA.jpg', isThumbnail: true },
      { url: 'https://as2.ftcdn.net/v2/jpg/02/31/56/61/1000_F_231566167_DcxyiS11UCKdIoFpPdkXFpAzeVhh6qFA.jpg', isThumbnail: false },
    ],
    priceAdult: 1000,
    priceJuniorHighSchool: 800,
    priceElementarySchool: 600,
    pricePreschool: 400,
    priceSenior: 800,
    recommendedPoint1: 'おすすめポイント1',
    recommendedPoint2: 'おすすめポイント2',
    recommendedPoint3: 'おすすめポイント3',
    promotionVideoUrl: 'https://example.com/promotion.mp4',
    duration: 60,
    direction: '新宿駅から徒歩10分',
    businessOpenTime: '1000',
    businessCloseTime: '1700',
    hostPostalCode: '123-4567',
    hostPrefecture: '東京都',
    hostCity: '千代田区',
    hostAddressLine1: '千代田1-1-1',
    hostAddressLine2: '千代田ビル1F',
    hostLongitude: 139.767052,
    hostLatitude: 35.681167,
    startAt: 1614556800,
    endAt: 1614643199,
  },
  coordinator: {
    id: 'kSByoE6FetnPs5Byk3a9Zx',
    storeName: '&.農園',
    profile: '紹介文です。',
    productTypes: ['kSByoE6FetnPs5Byk3a9Zx'],
    businessDays: [1, 2, 3, 4, 5],
    thumbnailUrl: 'https://and-period.jp/thumbnail.png',
    headerUrl: 'https://and-period.jp/header.png',
    promotionVideoUrl: 'https://and-period.jp/promotion.mp4',
    instagramId: 'instagram-id',
    facebookId: 'facebook-id',
    prefecture: '東京都',
    city: '千代田区',
  },
  producer: {
    id: 'kSByoE6FetnPs5Byk3a9Zx',
    coordinatorId: 'kSByoE6FetnPs5Byk3a9Zx',
    username: '&.農園',
    profile: '紹介文です。',
    thumbnailUrl: 'https://and-period.jp/thumbnail.png',
    headerUrl: 'https://and-period.jp/header.png',
    promotionVideoUrl: 'https://and-period.jp/promotion.mp4',
    instagramId: 'instagram-id',
    facebookId: 'facebook-id',
    prefecture: '東京都',
    city: '千代田区',
  },
  experienceType: {
    id: 'kSByoE6FetnPs5Byk3a9Zx',
    name: '農業体験',
  },
}
</script>

<template>
  <div>
    <!-- Snackbar Items -->
    <template
      v-for="(snackbarItem, i) in snackbarItems"
      :key="i"
    >
      <the-snackbar
        v-model:is-show="snackbarItem.isShow"
        :text="snackbarItem.text"
      />
    </template>

    <!-- Experience Section -->
    <template v-if="items.experience">
      <div class="bg-white w-full">
        <div
          class="gap-10 px-4 pb-6 pt-[40px] text-main md:grid md:grid-cols-2 lg:px-[112px] w-full max-w-[1440px] mx-auto"
        >
          <div class="mx-auto w-full max-w-[100%]">
            <div class="flex aspect-square h-full w-full justify-center">
              <template v-if="items.experience.promotionVideoUrl">
                <the-item-video-player :src="items.experience.promotionVideoUrl" />
              </template>
              <template v-else>
                <nuxt-img
                  provider="cloudFront"
                  fill="contain"
                  class="block h-full w-full object-contain border"
                  :src="selectMediaSrcUrl"
                  :alt="itemThumbnailAlt"
                />
              </template>
            </div>
            <div
              class="hidden-scrollbar mt-2 grid w-full grid-flow-col justify-start gap-2 overflow-x-scroll"
            >
              <template
                v-for="(m, i) in items.experience.media"
                :key
              >
                <template>
                  <nuxt-img
                    width="72px"
                    fill="contain"
                    provider="cloudFront"
                    :src="m.url"
                    :alt="`${itemThumbnailAlt}_${i}`"
                    class="aspect-square w-[72px] h-[72px] cursor-pointer object-contain border block"
                    @click="handleClickMediaItem(i)"
                  />
                </template>
              </template>
            </div>
          </div>

          <div class="mt-4 flex w-full flex-col gap-4">
            <div
              class="break-words text-[16px] tracking-[1.6px] md:text-[24px] md:tracking-[2.4px]"
            >
              {{ items.experience.title }}
            </div>

            <div
              v-if="items.producer"
              class="flex flex-col leading-[32px] md:mt-4"
            >
              <div
                class="text-[14px] tracking-[1.4px] md:text-[16px] md:tracking-[1.6px]"
              >
                生産者
                <a
                  href="#"
                  class="font-bold underline"
                >
                  {{ items.producer.username }}
                </a>
              </div>
              <div class="text-[12px] tracking-[1.4px] md:text-[14px]">
                {{ items.producer.prefecture}} {{ items.producer.city }}
              </div>
            </div>

            <div
              v-if="items.experience && items.experience.recommendedPoint1"
              class="w-full rounded-2xl bg-base px-[20px] py-[28px] text-main md:mt-8"
            >
              <p
                class="mb-[12px] text-[12px] font-medium tracking-[1.4px] md:text-[14px]"
              >
                おすすめポイント
              </p>
              <ol
                class="recommend-list flex flex-col divide-y divide-dashed divide-main px-[4px] pl-[24px]"
              >
                <li
                  v-if="items.experience.recommendedPoint1"
                  class="py-3 text-[14px] font-medium md:text-[16px]"
                >
                  {{ items.experience.recommendedPoint1 }}
                </li>
                <li
                  v-if="items.experience.recommendedPoint2"
                  class="py-3 text-[14px] font-medium md:text-[16px]"
                >
                  {{ items.experience.recommendedPoint2 }}
                </li>
                <li
                  v-if="items.experience.recommendedPoint3"
                  class="py-3 text-[14px] font-medium md:text-[16px]"
                >
                  {{ items.experience.recommendedPoint3 }}
                </li>
              </ol>
            </div>

            <button
              class="mt-2 w-full bg-main py-4 text-center text-white disabled:cursor-not-allowed disabled:bg-main/60 md:mt-8"
              :disabled="!canAddCart"
            >
              カゴに入れる
            </button>

            <div class="mt-4 inline-flex gap-4">
              <span
                class="rounded-2xl border border-main px-4 py-1 text-[14px] md:text-[16px]"
              >
                {{ items.experienceType?.name }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.recommend-list {
  list-style: none;
  counter-reset: li;
  position: relative;
}

.recommend-list li {
  padding-left: 16px;
}

.recommend-list li::before {
  content: counter(li);
  counter-increment: li;
  position: absolute;
  left: 0;
  background-color: #604c3f;
  color: #f9f6ea;
  border-radius: 100%;
  width: 24px;
  height: 24px;
  text-align: center;
}
</style>
