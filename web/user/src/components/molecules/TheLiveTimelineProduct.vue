<script setup lang="ts">
import type { Product } from '~/types/api'
import { priceFormatter } from '~/lib/price'

interface Props {
  product: Product
}

const props = defineProps<Props>()

interface Emits {
  (e: 'click:item', id: string): void
  (e: 'click:addCart', name: string, id: string, quantity: number): void
}

const emits = defineEmits<Emits>()

const formData = ref<number>(1)

const inventoryCountAboveTen = computed<boolean>(() => {
  return props.product.inventory > 10
})

const thumbnailUrl = computed<string>(() => {
  const thumbnail = props.product.media.find((m) => m.isThumbnail)
  if (thumbnail) {
    return thumbnail.url
  } else {
    return ''
  }
})

const hasStock = computed<boolean>(() => {
  return props.product.inventory > 0
})

const handleClickAddCart = () => {
  emits('click:addCart', props.product.name, props.product.id, formData.value)
}

const handleClickItemTitle = () => {
  if (hasStock.value) {
    emits('click:item', props.product.id)
  }
}
</script>

<template>
  <div class="flex gap-[10px]">
    <template v-if="thumbnailUrl">
      <div class="relative">
        <div
          v-if="!hasStock"
          class="absolute inset-0 flex items-center justify-center bg-black/50"
        >
          <p class="text-[14px] font-semibold text-white">在庫なし</p>
        </div>
        <img :src="thumbnailUrl" class="h-20 w-20" />
      </div>
    </template>
    <div class="flex flex-col justify-between">
      <div
        class="text-[12px] tracking-[1.2px]"
        :class="{ 'hover:cursor-pointer hover:underline': hasStock }"
        @click="handleClickItemTitle"
      >
        {{ product.name }}
      </div>
      <div>
        <p
          class="mb-2 text-[12px] font-bold after:ml-2 after:content-['(税込)']"
        >
          {{ priceFormatter(product.price) }}
        </p>
        <div class="flex h-6 items-center gap-2 text-[10px]">
          <div class="inline-flex h-full items-center">
            <select
              v-model="formData"
              class="h-full border-[1px] border-main px-2"
            >
              <template v-if="inventoryCountAboveTen">
                <option v-for="i in 10" :key="i" :value="i">
                  {{ i }}
                </option>
              </template>
              <template v-else>
                <option v-for="i in product.inventory" :key="i" :value="i">
                  {{ i }}
                </option>
              </template>
            </select>
          </div>
          <button
            class="flex h-full bg-main px-4 py-1 text-white disabled:cursor-not-allowed disabled:bg-main/60"
            :disabled="!hasStock"
            @click.stop="handleClickAddCart"
          >
            カゴに入れる
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
