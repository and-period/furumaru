<script lang="ts" setup>
import {
  ProductStatus,
  type Coordinator,
  type ProductMediaInner,
} from '~/types/api'
import type { I18n } from '~/types/locales'
import { productStatusToString } from '~/lib/product'

interface Props {
  id: string
  status: ProductStatus
  name: string
  inventory: number
  price: number
  hasStock: boolean
  originCity: string
  coordinator: Coordinator | undefined
  thumbnail: ProductMediaInner | undefined
  thumbnailIsVideo: boolean
}

interface Emits {
  (e: 'click:item', id: string): void
  (e: 'click:addCart', name: string, id: string, quantity: number): void
}

const props = defineProps<Props>()

const i18n = useI18n()

const router = useRouter()

const lt = (str: keyof I18n['items']['list']) => {
  return i18n.t(`items.list.${str}`)
}

const itemThumbnailAlt = computed<string>(() => {
  return i18n.t('items.list.itemThumbnailAlt', {
    itemName: props.name,
  })
})

const coordinatorThumbnailAlt = computed<string>(() => {
  return i18n.t('items.list.coordinatorThumbnailAlt', {
    coordinatorName: props.coordinator?.username,
  })
})

const emits = defineEmits<Emits>()

const quantity = ref<number>(1)

const priceString = computed<string>(() => {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
  }).format(props.price)
})

const canAddCart = computed<boolean>(() => {
  if (props.status === ProductStatus.FOR_SALE && props.hasStock) {
    return true
  }
  return false
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
        v-if="!canAddCart"
        class="absolute inset-0 flex items-center justify-center bg-black/50"
      >
        <p class="text-lg font-semibold text-white">
          {{
            status === ProductStatus.FOR_SALE
              ? lt("soldOutText")
              : productStatusToString(status, i18n)
          }}
        </p>
      </div>
      <div
        v-if="thumbnail"
        class="cursor-pointer w-full"
        @click="handleClickItem"
      >
        <template v-if="thumbnailIsVideo">
          <video
            :src="thumbnail.url"
            class="aspect-square w-full"
            autoplay
            muted
            webkit-playsinline
            playsinline
            loop
          />
        </template>
        <template v-else>
          <div class="aspect-square">
            <picture>
              <nuxt-img
                provider="cloudFront"
                :src="thumbnail.url"
                :alt="itemThumbnailAlt"
                fit="contain"
                sizes="180px md:250px"
                format="png"
                class="w-full h-full object-contain"
              />
            </picture>
          </div>
        </template>
      </div>
    </div>

    <p
      class="mt-2 line-clamp-3 grow text-[14px] tracking-[1.4px] md:text-[16px] md:tracking-[1.6px]"
    >
      {{ name }}
    </p>

    <p
      class="my-4 text-[16px] tracking-[1.6px] md:text-[20px] md:tracking-[2.0px]"
    >
      {{ priceString }}{{ lt("itemPriceTaxIncludedText") }}
    </p>

    <div class="flex h-6 items-center gap-2 text-[10px]">
      <div class="inline-flex h-full items-center">
        <label
          class="mr-2 block whitespace-nowrap text-center text-[8px] md:text-[14px]"
        >
          {{ lt("quantityLabel") }}
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
        :disabled="!canAddCart"
        class="flex h-full grow items-center justify-center bg-main p-1 text-[10px] text-white disabled:cursor-not-allowed disabled:bg-main/60 lg:px-4 xl:text-[14px]"
        @click="handleClickAddCartButton"
      >
        <the-cart-icon
          id="add-cart-icon"
          class="mr-1 h-2 w-2 lg:h-4 lg:w-4"
        />
        {{ lt("addToCartText") }}
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
          :alt="coordinatorThumbnailAlt"
          class="block aspect-square h-14 w-14 rounded-full object-cover"
        >
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
            <p class="whitespace-nowrap">
              {{ lt("coordinatorLabel") }}:
            </p>
            <p class="text-[12px] md:text-[14px]">
              {{ coordinator.username }}
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
