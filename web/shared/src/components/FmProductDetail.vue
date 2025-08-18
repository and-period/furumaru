<script setup lang="ts">
import { computed, ref } from 'vue';

// Media item interface
interface MediaItem {
  url: string;
  isThumbnail?: boolean;
}

// Rating information interface
interface ProductRating {
  average: number;
  count: number;
  detail?: { [key: string]: number };
}

// Component props interface
interface Props {
  mediaFiles?: MediaItem[];
  name?: string;
  description?: string;
  originPrefecture?: string;
  originCity?: string;
  rating?: ProductRating;
  recommendedPoint1?: string;
  recommendedPoint2?: string;
  recommendedPoint3?: string;
  expirationDate?: number;
  weight?: number;
  deliveryType?: number;
  storageMethodType?: number;
}

const props = withDefaults(defineProps<Props>(), {
  mediaFiles: () => [],
  name: '',
  description: '',
  originPrefecture: '',
  originCity: '',
  rating: () => ({ average: 0, count: 0 }),
  recommendedPoint1: '',
  recommendedPoint2: '',
  recommendedPoint3: '',
  expirationDate: 0,
  weight: 0,
  deliveryType: 0,
  storageMethodType: 0,
});

const selectedMediaIndex = ref<number>(-1);

// Get the currently selected media URL (thumbnail if none selected)
const selectMediaSrcUrl = computed<string>(() => {
  if (selectedMediaIndex.value === -1) {
    const thumbnail = props.mediaFiles.find(m => m.isThumbnail);
    return thumbnail?.url || props.mediaFiles[0]?.url || '';
  }
  return props.mediaFiles[selectedMediaIndex.value]?.url || '';
});

// Handle media item selection
const handleClickMediaItem = (index: number) => {
  selectedMediaIndex.value = index;
};

// Check if a URL is a video file
const isVideoUrl = (url: string): boolean => {
  try {
    const urlObj = new URL(url);
    urlObj.search = '';
    urlObj.hash = '';
    return urlObj.toString().endsWith('.mp4');
  } catch {
    return false;
  }
};

// Auto-link function for descriptions
const autoLink = (text: string): string => {
  if (!text) return '';
  return text.replace(
    /(https?:\/\/[^\s]+)/g,
    '<a href="$1" target="_blank" rel="noopener noreferrer" class="text-blue-600 underline break-all">$1</a>',
  );
};

const productDescriptionHtml = computed(() => autoLink(props.description));

// Get delivery type text
const getDeliveryTypeText = (type: number): string => {
  switch (type) {
    case 1:
      return '常温';
    case 2:
      return '冷蔵';
    case 3:
      return '冷凍';
    default:
      return '不明';
  }
};

// Get storage method type text
const getStorageMethodTypeText = (type: number): string => {
  switch (type) {
    case 0:
      return '不明';
    case 1:
      return '常温';
    case 2:
      return '冷暗所';
    case 3:
      return '冷蔵';
    case 4:
      return '冷凍';
    default:
      return '不明';
  }
};

// Format expiration date
const expirationDateText = computed<string>(() => {
  if (!props.expirationDate) return '';
  return `${props.expirationDate}日`;
});

// Check if product has recommended points
const hasRecommendedPoints = computed<boolean>(() => {
  return !!(props.recommendedPoint1 || props.recommendedPoint2 || props.recommendedPoint3);
});

// Origin location text
const originLocationText = computed<string>(() => {
  return `${props.originPrefecture} ${props.originCity}`.trim();
});
</script>

<template>
  <div class="bg-white w-full">
    <div class="gap-10 px-4 pb-6 pt-[40px] text-main md:grid md:grid-cols-2 lg:px-[112px] w-full max-w-[1440px] mx-auto">
      <!-- Media Gallery -->
      <div class="mx-auto w-full max-w-[100%]">
        <div class="flex aspect-square h-full w-full justify-center">
          <template v-if="isVideoUrl(selectMediaSrcUrl)">
            <video
              :src="`${selectMediaSrcUrl}#t=0.1`"
              class="block h-full w-full object-contain border"
              controls
              loop
            />
          </template>
          <template v-else>
            <img
              class="block h-full w-full object-contain border"
              :src="selectMediaSrcUrl"
              :alt="`thumbnail of ${name}`"
            >
          </template>
        </div>

        <!-- Media Thumbnails -->
        <div
          v-if="mediaFiles.length > 1"
          class="hidden-scrollbar mt-2 grid w-full grid-flow-col justify-start gap-2 overflow-x-scroll"
        >
          <template
            v-for="(media, index) in mediaFiles"
            :key="index"
          >
            <template v-if="isVideoUrl(media.url)">
              <div
                class="aspect-square w-[72px] h-[72px] cursor-pointer border relative"
                @click="handleClickMediaItem(index)"
              >
                <video
                  :src="`${media.url}#t=0.1`"
                  class="aspect-square w-full object-contain h-full"
                />
                <div class="absolute h-6 w-6 bottom-0 right-0 p-1 bg-main/80 rounded-full text-white">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke-width="1.5"
                    stroke="currentColor"
                    class="size-6"
                  >
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      d="m15.75 10.5 4.72-4.72a.75.75 0 0 1 1.28.53v11.38a.75.75 0 0 1-1.28.53l-4.72-4.72M4.5 18.75h9a2.25 2.25 0 0 0 2.25-2.25v-9a2.25 2.25 0 0 0-2.25-2.25h-9A2.25 2.25 0 0 0 2.25 7.5v9a2.25 2.25 0 0 0 2.25 2.25Z"
                    />
                  </svg>
                </div>
              </div>
            </template>
            <template v-else>
              <img
                width="72px"
                :src="media.url"
                :alt="`${name} image ${index + 1}`"
                class="aspect-square w-[72px] h-[72px] cursor-pointer object-contain border block"
                @click="handleClickMediaItem(index)"
              >
            </template>
          </template>
        </div>
      </div>

      <!-- Product Information -->
      <div class="mt-4 flex w-full flex-col gap-4">
        <!-- Product Name -->
        <div class="break-words text-[16px] tracking-[1.6px] md:text-[24px] md:tracking-[2.4px]">
          {{ name }}
        </div>

        <!-- Origin Location -->
        <div
          v-if="originLocationText"
          class="flex flex-col leading-[32px] md:mt-4"
        >
          <div class="text-[12px] tracking-[1.4px] md:text-[14px]">
            {{ originLocationText }}
          </div>
        </div>

        <!-- Rating Information -->
        <div v-if="rating && rating.count > 0">
          <div class="flex items-center gap-3">
            <div class="inline-flex items-center">
              <!-- Star rating display (simplified) -->
              <div class="flex items-center">
                <span class="text-yellow-500">★</span>
                <span class="text-yellow-500">★</span>
                <span class="text-yellow-500">★</span>
                <span class="text-yellow-500">★</span>
                <span class="text-gray-300">★</span>
              </div>
              <p class="ms-2 text-sm font-bold text-main">
                {{ rating.average }}
              </p>
            </div>
            <div class="text-sm font-medium text-main">
              {{ rating.count }} 件のレビュー
            </div>
          </div>
        </div>

        <!-- Recommended Points -->
        <div
          v-if="hasRecommendedPoints"
          class="w-full rounded-2xl bg-base px-[20px] py-[28px] text-main md:mt-8"
        >
          <p class="mb-[12px] text-[12px] font-medium tracking-[1.4px] md:text-[14px]">
            おすすめポイント
          </p>
          <ol class="recommend-list flex flex-col divide-y divide-dashed divide-main px-[4px] pl-[24px]">
            <li
              v-if="recommendedPoint1"
              class="py-3 text-[14px] font-medium md:text-[16px]"
            >
              {{ recommendedPoint1 }}
            </li>
            <li
              v-if="recommendedPoint2"
              class="py-3 text-[14px] font-medium md:text-[16px]"
            >
              {{ recommendedPoint2 }}
            </li>
            <li
              v-if="recommendedPoint3"
              class="py-3 text-[14px] font-medium md:text-[16px]"
            >
              {{ recommendedPoint3 }}
            </li>
          </ol>
        </div>
      </div>

      <!-- Product Description -->
      <div class="col-span-2 mt-[40px] pb-10 md:mt-[80px] md:pb-16">
        <!-- eslint-disable-next-line vue/no-v-html -->
        <article
          class="text-[14px] leading-[32px] tracking-[1.4px] md:text-[16px] md:tracking-[1.6px] whitespace-pre-wrap"
          v-html="productDescriptionHtml"
        />
      </div>

      <!-- Product Details Table -->
      <div class="col-span-2 flex flex-col divide-y divide-dashed divide-main border-y border-dashed border-main text-[14px] md:text-[16px]">
        <div
          v-if="expirationDate"
          class="grid grid-cols-5 py-4"
        >
          <p class="col-span-2 md:col-span-1">
            賞味期限
          </p>
          <p class="col-span-3 md:col-span-4">
            {{ expirationDateText }}
          </p>
        </div>
        <div
          v-if="weight"
          class="grid grid-cols-5 py-4"
        >
          <p class="col-span-2 md:col-span-1">
            重量
          </p>
          <p class="col-span-3 md:col-span-4">
            {{ weight }}kg
          </p>
        </div>
        <div
          v-if="deliveryType"
          class="grid grid-cols-5 py-4"
        >
          <p class="col-span-2 md:col-span-1">
            配送タイプ
          </p>
          <p class="col-span-3 md:col-span-4">
            {{ getDeliveryTypeText(deliveryType) }}
          </p>
        </div>
        <div
          v-if="storageMethodType"
          class="grid grid-cols-5 py-4"
        >
          <p class="col-span-2 md:col-span-1">
            保存方法
          </p>
          <p class="col-span-3 md:col-span-4">
            {{ getStorageMethodTypeText(storageMethodType) }}
          </p>
        </div>
      </div>
    </div>
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

.hidden-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

.hidden-scrollbar::-webkit-scrollbar {
  display: none;
}

.bg-base {
  background-color: #f9f6ea;
}

.text-main {
  color: #604c3f;
}
</style>