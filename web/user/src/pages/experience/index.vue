<script setup lang="ts">
import type { Snackbar } from '~/types/props'
import { ExperienceStatus } from '~/types/api'

const i18n = useI18n()

const dt = (str: keyof I18n['items']['experiences']) => {
  return i18n.t(`items.experiences.${str}`)
}

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

const priceString = (price: number) => {
  if (price) {
    return new Intl.NumberFormat('ja-JP', {
      style: 'currency',
      currency: 'JPY',
    }).format(price)
  }
  else {
    return ''
  }
}

const convertToTimeString = (time: string): string => {
  if (time.length === 4) {
    const hour = time.slice(0, 2)
    const minute = time.slice(2, 4)
    return `${hour}:${minute}`
  }
  throw new Error('Invalid input format. Expected a 4-digit string.')
}

const items = {
  experience: {
    id: 'kSByoE6FetnPs5Byk3a9Zx',
    title: '農業体験',
    description: '農業体験の説明 \n explanations of agriculture of experience',
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
                :key="i"
              >
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
                {{ dt("producerLabel") }}:
                <a
                  href="#"
                  class="font-bold underline"
                >
                  {{ items.producer.username }}
                </a>
              </div>
              <div class="text-[12px] tracking-[1.4px] md:text-[14px]">
                {{ items.producer.prefecture }} {{ items.producer.city }}
              </div>
            </div>

            <div
              v-if="items.experience && items.experience.recommendedPoint1"
              class="w-full rounded-2xl bg-base px-[20px] py-[28px] text-main md:mt-8"
            >
              <p
                class="mb-[12px] text-[12px] font-medium tracking-[1.4px] md:text-[14px]"
              >
                {{ dt("highlightsLabel") }}:
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

            <div class="items-center grid grid-cols-12 mt-4">
              <p class="text-[16px] font-medium col-span-5 md:col-span-6">
                {{ dt("adult") }}
              </p>
              <div
                class="col-span-4"
              >
                <div class="flex">
                  <p class="text-[16px] font-medium">
                    {{ priceString(items.experience.priceAdult) }}
                  </p>
                  <p class="pl-2 text-[12px] md:text-[14px] mt-auto">
                    {{ dt("itemPriceTaxIncludedText") }}
                  </p>
                </div>
              </div>
              <div class="col-span-3 md:col-span-2">
                <div
                  v-if="items.experience"
                  class="flex justify-end items-center"
                >
                  <label class="mr-2 hidden md:block text-[14px]">
                    {{ dt("quantityLabel") }}
                  </label>
                  <select
                    class="h-full border-[1px] border-main px-2"
                  >
                    <option
                      v-for="(_, i) in Array.from({
                        length: 10,
                      })"
                      :key="i + 1"
                      :value="i + 1"
                    >
                      {{ i + 1 }}
                    </option>
                  </select>
                </div>
              </div>
            </div>
            <div
              class="flex flex-col divide-y divide-dashed divide-main border-y border-dashed border-main"
            />
            <div class="items-center grid grid-cols-12">
              <p class="text-[16px] font-medium col-span-5 md:col-span-6">
                {{ dt("juniorHighSchoolStudents") }}
              </p>
              <div
                class="col-span-4"
              >
                <div class="flex">
                  <p class="text-[16px] font-medium">
                    {{ priceString(items.experience.priceJuniorHighSchool) }}
                  </p>
                  <p class="pl-2 text-[12px] md:text-[14px] mt-auto">
                    {{ dt("itemPriceTaxIncludedText") }}
                  </p>
                </div>
              </div>
              <div class="col-span-3 md:col-span-2">
                <div
                  v-if="items.experience"
                  class="flex justify-end items-center"
                >
                  <label class="mr-2 hidden md:block text-[14px]">
                    {{ dt("quantityLabel") }}
                  </label>
                  <select
                    class="h-full border-[1px] border-main px-2"
                  >
                    <option
                      v-for="(_, i) in Array.from({
                        length: 10,
                      })"
                      :key="i + 1"
                      :value="i + 1"
                    >
                      {{ i + 1 }}
                    </option>
                  </select>
                </div>
              </div>
            </div>
            <div
              class="flex flex-col divide-y divide-dashed divide-main border-y border-dashed border-main"
            />
            <div class="items-center grid grid-cols-12">
              <p class="text-[16px] font-medium col-span-5 md:col-span-6">
                {{ dt("elementarySchoolStudents") }}
              </p>
              <div
                class="col-span-4"
              >
                <div class="flex">
                  <p class="text-[16px] font-medium">
                    {{ priceString(items.experience.priceElementarySchool) }}
                  </p>
                  <p class="pl-2 text-[12px] md:text-[14px] mt-auto">
                    {{ dt("itemPriceTaxIncludedText") }}
                  </p>
                </div>
              </div>
              <div class="col-span-3 md:col-span-2">
                <div
                  v-if="items.experience"
                  class="flex justify-end items-center"
                >
                  <label class="mr-2 hidden md:block text-[14px]">
                    {{ dt("quantityLabel") }}
                  </label>
                  <select
                    class="h-full border-[1px] border-main px-2"
                  >
                    <option
                      v-for="(_, i) in Array.from({
                        length: 10,
                      })"
                      :key="i + 1"
                      :value="i + 1"
                    >
                      {{ i + 1 }}
                    </option>
                  </select>
                </div>
              </div>
            </div>
            <div
              class="flex flex-col divide-y divide-dashed divide-main border-y border-dashed border-main"
            />
            <div class="items-center grid grid-cols-12">
              <p class="text-[16px] font-medium col-span-5 md:col-span-6">
                {{ dt("preschoolers") }}
              </p>
              <div
                class="col-span-4"
              >
                <div class="flex">
                  <p class="text-[16px] font-medium">
                    {{ priceString(items.experience.pricePreschool) }}
                  </p>
                  <p class="pl-2 text-[12px] md:text-[14px] mt-auto">
                    {{ dt("itemPriceTaxIncludedText") }}
                  </p>
                </div>
              </div>
              <div class="col-span-3 md:col-span-2">
                <div
                  v-if="items.experience"
                  class="flex justify-end items-center"
                >
                  <label class="mr-2 hidden md:block text-[14px]">
                    {{ dt("quantityLabel") }}
                  </label>
                  <select
                    class="h-full border-[1px] border-main px-2"
                  >
                    <option
                      v-for="(_, i) in Array.from({
                        length: 10,
                      })"
                      :key="i + 1"
                      :value="i + 1"
                    >
                      {{ i + 1 }}
                    </option>
                  </select>
                </div>
              </div>
            </div>
            <div
              class="flex flex-col divide-y divide-dashed divide-main border-y border-dashed border-main"
            />
            <div class="items-center grid grid-cols-12">
              <p class="text-[16px] font-medium col-span-5 md:col-span-6">
                {{ dt("senior") }}
              </p>
              <div
                class="col-span-4"
              >
                <div class="flex">
                  <p class="text-[16px] font-medium">
                    {{ priceString(items.experience.priceSenior) }}
                  </p>
                  <p class="pl-2 text-[12px] md:text-[14px] mt-auto">
                    {{ dt("itemPriceTaxIncludedText") }}
                  </p>
                </div>
              </div>
              <div class="col-span-3 md:col-span-2">
                <div
                  v-if="items.experience"
                  class="flex justify-end items-center"
                >
                  <label class="mr-2 hidden md:block text-[14px]">
                    {{ dt("quantityLabel") }}
                  </label>
                  <select
                    class="h-full border-[1px] border-main px-2"
                  >
                    <option
                      v-for="(_, i) in Array.from({
                        length: 10,
                      })"
                      :key="i + 1"
                      :value="i + 1"
                    >
                      {{ i + 1 }}
                    </option>
                  </select>
                </div>
              </div>
            </div>

            <button
              class="mt-2 w-full bg-main py-4 text-center text-white disabled:cursor-not-allowed disabled:bg-main/60 md:mt-8"
              :disabled="!canAddCart"
            >
              {{ dt("addToCartText") }}
            </button>

            <div class="mt-4 inline-flex gap-4">
              <span
                class="rounded-2xl border border-main px-4 py-1 text-[14px] md:text-[16px]"
              >
                {{ items.experienceType?.name }}
              </span>
            </div>
          </div>
          <div class="col-span-2 mt-[40px] pb-10 md:mt-[80px] md:pb-16">
            <article
              class="text-[14px] leading-[32px] tracking-[1.4px] md:text-[16px] md:tracking-[1.6px] whitespace-pre-wrap"
              v-text="items.experience.description"
            />
          </div>
          <div
            class="col-span-2 flex flex-col divide-y divide-dashed divide-main border-y border-dashed border-main text-[14px] md:text-[16px]"
          >
            <div class="grid grid-cols-5 py-4">
              <p class="col-span-2 md:col-span-1">
                {{ dt("estimatedTime") }}
              </p>
              <p class="col-span-3 md:col-span-4">
                {{ items.experience.duration }}分
              </p>
            </div>
            <div class="grid grid-cols-5 py-4">
              <p class="col-span-2 md:col-span-1">
                {{ dt("businessHours") }}
              </p>
              <p class="col-span-3 md:col-span-4">
                {{ convertToTimeString(items.experience.businessOpenTime) }}~{{ convertToTimeString(items.experience.businessCloseTime) }}
              </p>
            </div>
            <div class="grid grid-cols-5 py-4">
              <p class="col-span-2 md:col-span-1">
                {{ dt("locationPostalcode") }}
              </p>
              <p class="col-span-3 md:col-span-4">
                {{ items.experience.hostPostalCode }}
              </p>
            </div>
            <div class="grid grid-cols-5 py-4">
              <p class="col-span-2 md:col-span-1">
                {{ dt("locationAddress") }}
              </p>
              <p class="col-span-3 md:col-span-4">
                {{ items.experience.hostPrefecture }}{{ items.experience.hostCity }}{{ items.experience.hostAddressLine1 }}{{ items.experience.hostAddressLine2 }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- Producer Section -->
    <template v-if="items.producer">
      <div class="w-full">
        <div class="mx-auto mt-[40px] w-full px-4 xl:px-28 max-w-[1440px]">
          <div
            class="flex w-full flex-col rounded-3xl bg-white px-8 py-10 text-main xl:px-16"
          >
            <p
              class="mx-auto w-full rounded-full bg-base py-2 text-center text-[14px] font-bold text-main md:text-[16px]"
            >
              {{ dt("producerInformationTitle") }}
            </p>

            <div
              class="mt-[64px] flex w-full flex-col gap-4 md:flex-row lg:gap-10"
            >
              <div
                class="flex min-w-max flex-col items-center justify-center gap-4 md:flex-row"
              >
                <template v-if="items.producer.thumbnailUrl">
                  <nuxt-img
                    provider="cloudFront"
                    sizes="96px md:120px"
                    fit="cover"
                    :src="items.producer.thumbnailUrl"
                    :alt="`${items.producer.username}`"
                    class="mx-auto block aspect-square w-[96px] rounded-full md:w-[120px] object-cover"
                  />
                </template>
                <template v-else>
                  <img
                    class="mx-auto block aspect-square w-[96px] rounded-full md:w-[120px] object-cover"
                    src="/img/account.png"
                  >
                </template>
                <div
                  class="flex min-w-max grow flex-col items-center gap-2 md:items-start md:gap-2 md:whitespace-nowrap"
                >
                  <p class="text-sm font-[500] tracking-[1.4px]">
                    {{ `${items.producer.prefecture} ${items.producer.city}` }}
                  </p>
                  <div
                    class="flex flex-row items-baseline grow text-[16px] tracking-[1.4px] md:text-[24px]"
                  >
                    <p class="mr-1 text-[14px] font-medium">
                      {{ dt("producerLabel") }}
                    </p>
                    <p>{{ items.producer.username }}</p>
                  </div>
                </div>
              </div>
              <div
                class="pt-2 text-[14px] tracking-[1.4px] md:pt-0 md:text-[16px] md:tracking-[1.6px]"
              >
                {{ items.producer.profile }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>

    <template v-if="items.experience">
      <div class="w-full">
        <div class="mx-auto mt-[40px] w-full px-4 xl:px-28 max-w-[1440px]">
          <div
            class="flex w-full flex-col rounded-3xl bg-white px-8 py-10 text-main xl:px-16"
          >
            <p
              class="mx-auto w-full rounded-full bg-base py-2 text-center text-[14px] font-bold text-main md:text-[16px]"
            >
              {{ dt("accessMethod") }}
            </p>
            <p class="mt-8 text-[14px] md:text-[16px]">
              〒　{{ items.experience.hostPostalCode }}
            </p>
            <div
              class="flex mt-4 text-[14px] md:text-[16px] items-center"
            >
              <img
                src="/img/experience/map.svg"
                class="w-[16px] h-[21px] md:w-[20px] md:h-[42px]"
              >
              <p class="ml-3">
                {{ items.experience.hostPrefecture }}{{ items.experience.hostCity }}{{ items.experience.hostAddressLine1 }}{{ items.experience.hostAddressLine2 }}
              </p>
            </div>
            <div
              class="pt-4 text-[14px] tracking-[1.4px] md:pt-8 md:text-[16px] md:tracking-[1.6px]"
            >
              {{ items.experience.direction }}
            </div>
            <!-- 緯度経度をもとにgoogle mapを表示する -->
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
