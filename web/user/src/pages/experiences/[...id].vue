<script setup lang="ts">
import type { Snackbar } from '~/types/props'
import { ExperienceStatus } from '~/types/api'
import type { I18n } from '~/types/locales'
import { useExperienceStore } from '~/store/experience'
import { useAuthStore } from '~/store/auth'

const i18n = useI18n()

const dt = (str: keyof I18n['items']['experiences']) => {
  return i18n.t(`items.experiences.${str}`)
}

const route = useRoute()

const router = useRouter()

const experienceId = computed<string>(() => {
  const ids = route.params.id
  if (Array.isArray(ids)) {
    return ids[0]
  }
  else {
    return route.params.id as string
  }
})

const authStore = useAuthStore()
const { isAuthenticated } = storeToRefs(authStore)

const experienceStore = useExperienceStore()
const { fetchExperience } = experienceStore

const { data, status } = await useAsyncData('spot', () => {
  return fetchExperience(experienceId.value)
})

const snackbarItems = ref<Snackbar[]>([])

const selectedMediaIndex = ref<number>(0)

const selectMediaSrcUrl = computed<string>(() => {
  if (!data.value || !data.value.experience) {
    return ''
  }

  return selectedMediaIndex.value === -1
    ? data.value.experience.promotionVideoUrl
    : data.value.experience.media[selectedMediaIndex.value].url
})

const handleClickMediaItem = (index: number) => {
  selectedMediaIndex.value = index
}

const itemThumbnailAlt = computed<string>(() => {
  if (!data.value || !data.value.experience) {
    return ''
  }

  return i18n.t('items.list.itemThumbnailAlt', {
    itemName: data.value.experience.title,
  })
})

const formData = ref({
  adultCount: 0,
  juniorHighSchoolCount: 0,
  elementarySchoolCount: 0,
  preschoolCount: 0,
  seniorCount: 0,
})

const canAddCart = computed<boolean>(() => {
  if (!data.value || !data.value.experience) {
    return false
  }

  if (data.value.experience.status !== ExperienceStatus.ACCEPTING) {
    return false
  }

  if (
    formData.value.adultCount === 0
    && formData.value.juniorHighSchoolCount === 0
    && formData.value.elementarySchoolCount === 0
    && formData.value.preschoolCount === 0
    && formData.value.seniorCount === 0
  ) {
    return false
  }

  return true
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

const handleClickApplyButton = () => {
  if (isAuthenticated.value) {
    console.log('authenticated')
    router.push('/experiences/purchase', {
      query: {
        id: experienceId.value,
      },
    })
  }
  else {
    router.push({
      path: '/experiences/purchase/guest',
      query: {
        id: experienceId.value,
        ...formData.value,
      },
    })
  }
}

useSeoMeta({
  title: data.value?.experience?.title || '',
})
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

    <template v-if="status === 'pending'">
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

    <!-- Experience Section -->
    <template v-if="data?.experience">
      <div class="bg-white w-full">
        <div
          class="gap-10 px-4 pb-6 pt-[40px] text-main md:grid md:grid-cols-2 lg:px-[112px] w-full max-w-[1440px] mx-auto"
        >
          <div class="mx-auto w-full max-w-[100%]">
            <div class="flex aspect-square h-full w-full justify-center">
              <template v-if="data.experience.promotionVideoUrl">
                <the-item-video-player :src="data.experience.promotionVideoUrl" />
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
                v-for="(m, i) in data.experience.media"
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
              {{ data.experience.title }}
            </div>

            <div
              v-if="data.producer"
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
                  {{ data.producer.username }}
                </a>
              </div>
              <div class="text-[12px] tracking-[1.4px] md:text-[14px]">
                {{ data.producer.prefecture }} {{ data.producer.city }}
              </div>
            </div>

            <div
              v-if="data.experience && data.experience.recommendedPoint1"
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
                  v-if="data.experience.recommendedPoint1"
                  class="py-3 text-[14px] font-medium md:text-[16px]"
                >
                  {{ data.experience.recommendedPoint1 }}
                </li>
                <li
                  v-if="data.experience.recommendedPoint2"
                  class="py-3 text-[14px] font-medium md:text-[16px]"
                >
                  {{ data.experience.recommendedPoint2 }}
                </li>
                <li
                  v-if="data.experience.recommendedPoint3"
                  class="py-3 text-[14px] font-medium md:text-[16px]"
                >
                  {{ data.experience.recommendedPoint3 }}
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
                    {{ priceString(data.experience.priceAdult) }}
                  </p>
                  <p class="pl-2 text-[12px] md:text-[14px] mt-auto">
                    {{ dt("itemPriceTaxIncludedText") }}
                  </p>
                </div>
              </div>
              <div class="col-span-3 md:col-span-2">
                <div
                  v-if="data.experience"
                  class="flex justify-end items-center"
                >
                  <label class="mr-2 hidden md:block text-[14px]">
                    {{ dt("quantityLabel") }}
                  </label>
                  <select
                    v-model="formData.adultCount"
                    class="h-full border-[1px] border-main px-2"
                  >
                    <option
                      v-for="(_, i) in Array.from({
                        length: 11,
                      })"
                      :key="i"
                      :value="i"
                    >
                      {{ i }}
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
                    {{ priceString(data.experience.priceJuniorHighSchool) }}
                  </p>
                  <p class="pl-2 text-[12px] md:text-[14px] mt-auto">
                    {{ dt("itemPriceTaxIncludedText") }}
                  </p>
                </div>
              </div>
              <div class="col-span-3 md:col-span-2">
                <div
                  v-if="data.experience"
                  class="flex justify-end items-center"
                >
                  <label class="mr-2 hidden md:block text-[14px]">
                    {{ dt("quantityLabel") }}
                  </label>
                  <select
                    v-model="formData.juniorHighSchoolCount"
                    class="h-full border-[1px] border-main px-2"
                  >
                    <option
                      v-for="(_, i) in Array.from({
                        length: 11,
                      })"
                      :key="i"
                      :value="i"
                    >
                      {{ i }}
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
                    {{ priceString(data.experience.priceElementarySchool) }}
                  </p>
                  <p class="pl-2 text-[12px] md:text-[14px] mt-auto">
                    {{ dt("itemPriceTaxIncludedText") }}
                  </p>
                </div>
              </div>
              <div class="col-span-3 md:col-span-2">
                <div
                  v-if="data.experience"
                  class="flex justify-end items-center"
                >
                  <label class="mr-2 hidden md:block text-[14px]">
                    {{ dt("quantityLabel") }}
                  </label>
                  <select
                    v-model="formData.elementarySchoolCount"
                    class="h-full border-[1px] border-main px-2"
                  >
                    <option
                      v-for="(_, i) in Array.from({
                        length: 11,
                      })"
                      :key="i"
                      :value="i"
                    >
                      {{ i }}
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
                    {{ priceString(data.experience.pricePreschool) }}
                  </p>
                  <p class="pl-2 text-[12px] md:text-[14px] mt-auto">
                    {{ dt("itemPriceTaxIncludedText") }}
                  </p>
                </div>
              </div>
              <div class="col-span-3 md:col-span-2">
                <div
                  v-if="data.experience"
                  class="flex justify-end items-center"
                >
                  <label class="mr-2 hidden md:block text-[14px]">
                    {{ dt("quantityLabel") }}
                  </label>
                  <select
                    v-model="formData.preschoolCount"
                    class="h-full border-[1px] border-main px-2"
                  >
                    <option
                      v-for="(_, i) in Array.from({
                        length: 11,
                      })"
                      :key="i"
                      :value="i"
                    >
                      {{ i }}
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
                    {{ priceString(data.experience.priceSenior) }}
                  </p>
                  <p class="pl-2 text-[12px] md:text-[14px] mt-auto">
                    {{ dt("itemPriceTaxIncludedText") }}
                  </p>
                </div>
              </div>
              <div class="col-span-3 md:col-span-2">
                <div
                  v-if="data.experience"
                  class="flex justify-end items-center"
                >
                  <label class="mr-2 hidden md:block text-[14px]">
                    {{ dt("quantityLabel") }}
                  </label>
                  <select
                    v-model="formData.seniorCount"
                    class="h-full border-[1px] border-main px-2"
                  >
                    <option
                      v-for="(_, i) in Array.from({
                        length: 11,
                      })"
                      :key="i"
                      :value="i"
                    >
                      {{ i }}
                    </option>
                  </select>
                </div>
              </div>
            </div>

            <button
              class="mt-2 w-full bg-main py-4 text-center text-white disabled:cursor-not-allowed disabled:bg-main/60 md:mt-8"
              :disabled="!canAddCart"
              @click="handleClickApplyButton"
            >
              {{ dt("applyButtonText") }}
            </button>

            <div class="mt-4 inline-flex gap-4">
              <span
                class="rounded-2xl border border-main px-4 py-1 text-[14px] md:text-[16px]"
              >
                {{ data.experienceType?.name }}
              </span>
            </div>
          </div>
          <div class="col-span-2 mt-[40px] pb-10 md:mt-[80px] md:pb-16">
            <article
              class="text-[14px] leading-[32px] tracking-[1.4px] md:text-[16px] md:tracking-[1.6px] whitespace-pre-wrap"
              v-text="data.experience.description"
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
                {{ data.experience.duration }}分
              </p>
            </div>
            <div class="grid grid-cols-5 py-4">
              <p class="col-span-2 md:col-span-1">
                {{ dt("businessHours") }}
              </p>
              <p class="col-span-3 md:col-span-4">
                {{ convertToTimeString(data.experience.businessOpenTime) }}~{{ convertToTimeString(data.experience.businessCloseTime) }}
              </p>
            </div>
            <div class="grid grid-cols-5 py-4">
              <p class="col-span-2 md:col-span-1">
                {{ dt("locationPostalcode") }}
              </p>
              <p class="col-span-3 md:col-span-4">
                {{ data.experience.hostPostalCode }}
              </p>
            </div>
            <div class="grid grid-cols-5 py-4">
              <p class="col-span-2 md:col-span-1">
                {{ dt("locationAddress") }}
              </p>
              <p class="col-span-3 md:col-span-4">
                {{ data.experience.hostPrefecture }}{{ data.experience.hostCity }}{{ data.experience.hostAddressLine1 }}{{ data.experience.hostAddressLine2 }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- Producer Section -->
    <template v-if="data?.producer">
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
                <template v-if="data.producer.thumbnailUrl">
                  <nuxt-img
                    provider="cloudFront"
                    sizes="96px md:120px"
                    fit="cover"
                    :src="data.producer.thumbnailUrl"
                    :alt="`${data.producer.username}`"
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
                    {{ `${data.producer.prefecture} ${data.producer.city}` }}
                  </p>
                  <div
                    class="flex flex-row items-baseline grow text-[16px] tracking-[1.4px] md:text-[24px]"
                  >
                    <p class="mr-1 text-[14px] font-medium">
                      {{ dt("producerLabel") }}
                    </p>
                    <p>{{ data.producer.username }}</p>
                  </div>
                </div>
              </div>
              <div
                class="pt-2 text-[14px] tracking-[1.4px] md:pt-0 md:text-[16px] md:tracking-[1.6px]"
              >
                {{ data.producer.profile }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>

    <template v-if="data?.experience">
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
            <div class="mt-8 text-[14px] md:text-[16px] flex">
              <p>〒</p>
              <p class="ml-3 md:ml-4">
                {{ data.experience.hostPostalCode }}
              </p>
            </div>
            <div
              class="flex mt-4 text-[14px] md:text-[16px] items-center"
            >
              <img
                src="/img/experience/map.svg"
                class="w-[16px] h-[21px] md:w-[20px] md:h-[42px]"
              >
              <p class="ml-3">
                {{ data.experience.hostPrefecture }}{{ data.experience.hostCity }}{{ data.experience.hostAddressLine1 }}{{ data.experience.hostAddressLine2 }}
              </p>
            </div>
            <div
              class="pt-4 text-[14px] tracking-[1.4px] md:pt-8 md:text-[16px] md:tracking-[1.6px]"
            >
              {{ data.experience.direction }}
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
