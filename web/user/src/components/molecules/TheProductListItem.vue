<script lang="ts" setup>
import { Coordinator, ProductMediaInner } from '~/types/api'

interface Props {
  id: string
  name: string
  inventory: number
  price: number
  hasStock: boolean
  originCity: string
  coordinator: Coordinator | undefined
  thumbnail: ProductMediaInner | undefined
}

interface Emits {
  (e: 'click:item', id: string): void
  (e: 'click:addCart', id: string, quantity: number): void
}

const props = defineProps<Props>()

const emits = defineEmits<Emits>()

const quantity = ref<number>(1)

const priceString = computed<string>(() => {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
  }).format(props.price)
})

const handleClickItem = () => {
  emits('click:item', props.id)
}

const handleClickAddCartButton = () => {
  emits('click:addCart', props.id, quantity.value)
}
</script>

<template>
  <div class="text-main">
    <div class="relative">
      <div
        v-if="!hasStock"
        class="absolute inset-0 flex items-center justify-center bg-black/50"
      >
        <p class="text-lg font-semibold text-white">在庫なし</p>
      </div>
      <picture
        v-if="thumbnail"
        class="w-full cursor-pointer"
        @click="handleClickItem"
      >
        <img
          :src="thumbnail.url"
          :alt="`${name}のサムネイル画像`"
          class="aspect-square w-full"
        />
      </picture>
    </div>

    <p class="mt-2">
      {{ name }}
    </p>

    <p
      class="my-4 text-xl after:ml-2 after:text-[16px] after:content-['(税込)']"
    >
      {{ priceString }}
    </p>

    <div class="flex h-8 items-center gap-2 text-sm">
      <div class="inline-flex items-center">
        <label class="mr-2 block text-[10px] xl:text-[14px]">数量</label>
        <select
          v-model="quantity"
          class="h-full border-[1px] border-main px-2"
          :disabled="!hasStock"
        >
          <option
            v-for="(_, i) in Array.from({
              length: inventory < 10 ? inventory : 10,
            })"
            :key="i + 1"
            :value="i + 1"
          >
            {{ i + 1 }}
          </option>
        </select>
      </div>
      <button
        :disabled="!hasStock"
        class="flex h-full grow items-center justify-center bg-main p-1 text-[10px] text-white lg:px-4 lg:text-[14px]"
        @click="handleClickAddCartButton"
      >
        <the-cart-icon id="add-cart-icon" class="mr-1 h-3 w-3 lg:h-4 lg:w-4" />
        カゴに入れる
      </button>
    </div>
    <div v-if="coordinator" class="mt-4 flex items-center gap-x-4 text-xs">
      <div class="grow">
        <span
          class="inline-block whitespace-pre-wrap"
          v-text="`${coordinator.marcheName}\n${originCity}`"
        >
        </span>
        <hr class="my-2 border-dashed border-main" />
        <span class="before:content-['CN：']">
          {{ coordinator.username }}
        </span>
      </div>
      <div class="h-14 w-14">
        <img
          :src="coordinator.thumbnailUrl"
          :alt="`${coordinator.username}のサムネイル画像`"
          class="aspect-square rounded-full"
        />
      </div>
    </div>
  </div>
</template>
