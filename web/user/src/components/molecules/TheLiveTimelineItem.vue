<script setup lang="ts">
import dayjs from 'dayjs'
import type { Product } from '~/types/api'

interface Props {
  startAt: number
  comment: string
  username: string | undefined
  thumbnailUrl: string | undefined
  items: Product[]
}

const props = defineProps<Props>()

const priceFormatter = (price: number) => {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
  }).format(price)
}

const startAtString = computed(() => {
  return dayjs.unix(props.startAt).format('hh:mm')
})
</script>

<template>
  <li class="mb-10 ml-4">
    <div
      class="absolute -left-[9px] mt-2 h-4 w-4 rounded-full border-4 border-orange bg-white"
    />
    <time
      class="absolute -left-14 mt-2 text-[14px] font-medium leading-none tracking-[1.4px]"
    >
      {{ startAtString }}
    </time>

    <div
      class="mt-2 flex flex-col gap-x-12 gap-y-4 pl-6 pt-[24px] md:grid md:grid-cols-4"
    >
      <div class="col-span-1 flex items-center gap-4 md:flex-col">
        <div class="flex w-full flex-col items-center">
          <img
            v-if="thumbnailUrl"
            :src="thumbnailUrl"
            class="mb-2 h-[66px] w-[66px] rounded-full"
          />
          <p
            v-if="username"
            class="text-center text-[14px] font-medium tracking-[1.4px]"
          >
            {{ username }}
          </p>
        </div>
        <div
          class="text-[12px] font-medium tracking-[1.2px]"
          v-html="comment"
        ></div>
      </div>

      <div
        class="flex flex-col gap-4 md:col-span-3 md:grid md:grid-cols-2 md:gap-8"
      >
        <div v-for="item in items" :key="item.id" class="flex gap-[10px]">
          <img :src="item.imgSrc" class="h-20 w-20" />
          <div class="flex flex-col justify-between">
            <div class="text-[12px] tracking-[1.2px]">
              {{ item.name }}
            </div>
            <div>
              <p
                class="mb-2 text-[12px] font-bold after:ml-2 after:content-['(税込)']"
              >
                {{ priceFormatter(item.price) }}
              </p>
              <div class="flex h-6 items-center gap-2 text-[10px]">
                <div class="inline-flex h-full items-center">
                  <select class="h-full border-[1px] border-main px-2">
                    <option value="0">0</option>
                  </select>
                </div>
                <button class="flex h-full bg-main px-4 py-1 text-white">
                  カゴに入れる
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </li>
</template>
