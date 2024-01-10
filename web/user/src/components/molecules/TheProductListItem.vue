<script lang="ts" setup>
import type { Coordinator, ProductMediaInner } from '~/types/api'

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
  (e: 'click:addCart', name: string, id: string, quantity: number): void
}

const props = defineProps<Props>()

const router = useRouter()

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

const handleClickCoorinator = () => {
  router.push(`/coordinator/${props.coordinator?.id}`)
}

const handleClickAddCartButton = () => {
  emits('click:addCart', props.name, props.id, quantity.value)
}
</script>

<template>
  <div class="flex flex-col text-main">
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

    <p
      class="mt-2 line-clamp-3 grow text-[14px] tracking-[1.4px] md:text-[16px] md:tracking-[1.6px]"
    >
      {{ name }}
    </p>

    <p
      class="my-4 text-[16px] tracking-[1.6px] after:ml-2 after:text-[16px] after:content-['(税込)'] md:text-[20px] md:tracking-[2.0px]"
    >
      {{ priceString }}
    </p>

    <div class="flex h-6 items-center gap-2 text-[10px]">
      <div class="inline-flex h-full items-center">
        <label
          class="mr-2 block whitespace-nowrap text-center text-[8px] md:text-[14px]"
        >
          数量
        </label>
        <select
          v-model="quantity"
          class="h-full border-[1px] border-main px-1"
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
        class="flex h-full grow items-center justify-center bg-main p-1 text-[10px] text-white disabled:cursor-not-allowed disabled:bg-main/60 lg:px-4 xl:text-[14px]"
        @click="handleClickAddCartButton"
      >
        <the-cart-icon id="add-cart-icon" class="mr-1 h-2 w-2 lg:h-4 lg:w-4" />
        カゴに入れる
      </button>
    </div>
    <div
      v-if="coordinator"
      class="mt-4 flex flex-col gap-4 text-xs md:flex-row md:items-center"
    >
      <div class="md:hidden">
        <button @click="handleClickCoorinator">
          <p
            class="mb-2 w-full whitespace-pre-wrap text-[14px] font-bold underline md:text-[15px]"
          >
            {{ coordinator.marcheName }}
          </p>
        </button>
        <p class="text-[11px]">
          {{ coordinator.prefecture }} {{ coordinator.city }}
        </p>
      </div>

      <div class="flex items-center gap-x-4">
        <img
          :src="coordinator.thumbnailUrl"
          :alt="`${coordinator.username}のサムネイル画像`"
          class="block aspect-square h-14 w-14 rounded-full"
        />
        <div>
          <div class="hidden md:block">
            <button @click="handleClickCoorinator">
              <p
                class="mb-2 inline-block whitespace-pre-wrap text-[14px] font-bold underline md:text-[15px]"
              >
                {{ coordinator.marcheName }}
              </p>
            </button>
            <p class="text-[11px]">
              {{ coordinator.prefecture }} {{ coordinator.city }}
            </p>
          </div>
          <div class="mt-[5px] flex flex-col gap-2 md:flex-row">
            <p class="whitespace-nowrap">取扱元:</p>
            <p class="text-[12px] md:text-[14px]">
              {{ coordinator.username }}
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
