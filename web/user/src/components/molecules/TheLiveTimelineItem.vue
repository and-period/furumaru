<script setup lang="ts">
import dayjs from 'dayjs'
import type { Product } from '~/types/api'

const { unix } = dayjs

interface Props {
  startAt: number
  comment: string
  username: string | undefined
  thumbnailUrl: string | undefined
  items: Product[]
}

const props = defineProps<Props>()

interface Emits {
  (e: 'click:item', id: string): void
  (e: 'click:addCart', name: string, id: string, quantity: number): void
}

const emits = defineEmits<Emits>()

const startAtString = computed(() => {
  return unix(props.startAt).format('HH:mm')
})

const handleClickItem = (prodictId: string) => {
  emits('click:item', prodictId)
}

const handleClickAddCart = (name: string, id: string, quantity: number) => {
  emits('click:addCart', name, id, quantity)
}
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
      <div class="col-span-1 flex items-center gap-2 md:flex-col">
        <div class="flex w-[80px] flex-col items-center">
          <img
            v-if="thumbnailUrl"
            :src="thumbnailUrl"
            class="mb-2 h-[48px] w-[48px] rounded-full"
          />
          <p
            v-if="username"
            class="text-center text-[12px] font-medium tracking-[1.4px] md:text-[14px]"
          >
            {{ username }}
          </p>
        </div>
        <div
          class="w-[150px] overflow-auto break-words text-[12px] font-medium tracking-[1.2px]"
        >
          {{ comment }}
        </div>
      </div>

      <div
        class="flex flex-col items-start gap-4 md:col-span-3 md:grid md:grid-cols-2 md:gap-8"
      >
        <the-live-timeline-product
          v-for="item in items"
          :key="item.id"
          :product="item"
          @click:item="handleClickItem"
          @click:add-cart="handleClickAddCart"
        />
      </div>
    </div>
  </li>
</template>
