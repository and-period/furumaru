<script setup lang="ts">
import { ProductStatus, type Product } from '~/types/api'
import { priceFormatter } from '~/lib/price'
import { productStatusToString } from '~/lib/product'
import type { I18n } from '~/types/locales'

interface Props {
  product: Product
}

const props = defineProps<Props>()

interface Emits {
  (e: 'click:item', id: string): void
  (e: 'click:addCart', name: string, id: string, quantity: number): void
}

const emits = defineEmits<Emits>()

const i18n = useI18n()
const dt = (str: keyof I18n['lives']['details']) => {
  return i18n.t(`lives.details.${str}`)
}

const formData = ref<number>(1)

const inventoryCountAboveTen = computed<boolean>(() => {
  return props.product.inventory > 10
})

const thumbnailUrl = computed<string>(() => {
  const thumbnail = props.product.media.find(m => m.isThumbnail)
  if (thumbnail) {
    return thumbnail.url
  }
  else {
    return ''
  }
})

const hasStock = computed<boolean>(() => {
  return props.product.inventory > 0
})

const canAddCart = computed<boolean>(() => {
  if (props.product.status === ProductStatus.FOR_SALE && hasStock.value) {
    return true
  }
  return false
})

const handleClickAddCart = () => {
  emits('click:addCart', props.product.name, props.product.id, formData.value)
}

const handleClickItem = () => {
  if (hasStock.value) {
    emits('click:item', props.product.id)
  }
}
</script>

<template>
  <div class="flex gap-[10px]">
    <template v-if="thumbnailUrl">
      <div class="relative aspect-square h-[80px] w-[80px]">
        <div
          v-if="!canAddCart"
          class="absolute inset-0 flex h-full w-full items-center justify-center bg-black/50"
        >
          <p
            class="font-semibold text-white"
            :class="
              product.status === ProductStatus.PRESALE
                ? 'text-[12px]'
                : 'text-[14px]'
            "
          >
            {{
              product.status === ProductStatus.FOR_SALE
                ? "在庫なし"
                : productStatusToString(product.status, i18n)
            }}
          </p>
        </div>
        <template v-if="thumbnailUrl.endsWith('.mp4')">
          <video
            class="aspect-square h-full w-full"
            :src="thumbnailUrl"
            loop
            autoplay
            muted
            webkit-playsinline
            playsinline
          />
        </template>
        <template v-else>
          <nuxt-img
            :src="thumbnailUrl"
            provider="cloudFront"
            width="80px"
            height="80px"
            class="aspect-square h-full w-full"
            :class="{ 'hover:cursor-pointer hover:underline': canAddCart }"
            @click="handleClickItem"
          />
        </template>
      </div>
    </template>
    <div class="flex flex-col justify-between">
      <div
        class="text-[12px] tracking-[1.2px]"
        :class="{ 'hover:cursor-pointer hover:underline': canAddCart }"
        @click="handleClickItem"
      >
        {{ product.name }}
      </div>
      <div>
        <div class="flex flex-row mb-2 text-[12px] font-bold">
          <p>
            {{ priceFormatter(product.price) }}
          </p>
          <p class="ml-2">
            {{ dt("itemPriceTaxIncludedText") }}
          </p>
        </div>
        <div class="flex h-6 items-center gap-2 text-[10px]">
          <div class="inline-flex h-full items-center">
            <select
              v-model="formData"
              class="h-full border-[1px] border-main px-2"
            >
              <template v-if="inventoryCountAboveTen">
                <option
                  v-for="i in 10"
                  :key="i"
                  :value="i"
                >
                  {{ i }}
                </option>
              </template>
              <template v-else>
                <option
                  v-for="i in product.inventory"
                  :key="i"
                  :value="i"
                >
                  {{ i }}
                </option>
              </template>
            </select>
          </div>
          <button
            class="flex h-full bg-orange px-4 py-1 text-white transition-all duration-200 ease-in-out hover:shadow-lg active:scale-95 disabled:cursor-not-allowed disabled:bg-main/60"
            :disabled="!canAddCart"
            @click.stop="handleClickAddCart"
          >
            {{ dt("addToCartText") }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
