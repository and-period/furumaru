<script setup lang="ts">
import { computed, ref } from 'vue';

interface Props {
	name: string;
	price: number
  thumbnailUrl: string
  stock: number
  soldOutText?: string
	addToCartButtonText?: string;
	selectLabelText?: string
}

interface Emits {
  (e: 'click:addCart', quantity: number): void
}

const props = withDefaults(defineProps<Props>(), {
	addToCartButtonText: 'カゴに入れる',
	selectLabelText: '数量',
  soldOutText: '在庫なし',
});

const emits = defineEmits<Emits>()

const quantity = ref<number>(1)

const priceString = computed<string>(() => {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
  }).format(props.price)
})

/**
 * 商品をカートに追加できるかを示す
 */
const canAddCart = computed<boolean>(() => {
  if (props.stock > 0) {
    return true
  }
  return false
})

/**
 * 数量の選択肢を生成する
 * 在庫が0の場合は1を返す
 * 在庫が10以上の場合は10を返す
 * それ以外は在庫数を返す
 */
const stokeValues = computed<number>(() => {
  if (props.stock == 0) {
    return 1
  }
  if (props.stock > 10) {
    return 10
  }
	return props.stock
})

/**
 * サムネイルのURLが動画かどうかを判定する
 * 動画の場合はtrue、画像の場合はfalseを返す
 */
const thumbnailIsVideo = computed<boolean>(() => {
  try {
    const url = new URL(props.thumbnailUrl)

    // クエリパラメータとハッシュを削除
    url.search = ''
    url.hash = ''

    return url.toString().endsWith('.mp4')
  } catch {
    return false
  }
})

/**
 * カートに追加するボタンをクリックしたときのイベントハンドラ
 */
const handleClickAddCartButton = () => {
  emits('click:addCart', quantity.value)
}
</script>

<template>
  <div class="flex flex-col text-main w-full font-semibold gap-2">
    <div 
      class="relative"
    >
      <div 
        v-if="!canAddCart"
        class="absolute inset-0 flex items-center justify-center bg-black/50"
      >
        <p class="text-lg font-semibold text-white">
          {{ soldOutText }}
        </p>
      </div>
      <div class="block w-full">
        <template v-if="thumbnailIsVideo">
          <video
            :src="thumbnailUrl"
            class="aspect-square w-full"
            :alt="`video of ${name}`"
            :title="name"
            preload="metadata"
            autoplay
            muted
            webkit-playsinline
            playsinline
            loop
          />
        </template>
        <template v-else>
          <div class="aspect-square">
            <img
              class="w-full h-full object-contain"
              :src="thumbnailUrl"
              :alt="`thumbnail of ${name}`"
            >
          </div>
        </template> 
      </div>
    </div>
    <p
      class="line-clamp-3 grow text-[14px] tracking-[1.4px] md:text-[16px] md:tracking-[1.6px]"
    >
      {{ name }}
    </p>

    <p
      class="text-[16px] tracking-[1.6px] md:text-[20px] md:tracking-[2.0px]"
    >
      {{ priceString }}
    </p>

    <div class="inline-flex gap-2 text-xs items-center w-full">
      <div class="inline-flex gap-1 items-center">
        <label
          for="select"
          class="text-nowrap"
        >
          {{ selectLabelText }}
        </label>
        <select
          id="select"
          v-model="quantity"
          class="border border-main py-1 pl-1"
        >
          <template
            v-for="n in stokeValues"
            :key="n"
          >
            <option :value="n">
              {{ n }}
            </option>
          </template>
        </select>
      </div>
      <button
        class="bg-orange text-white py-1 px-2 w-full disabled:bg-main/60 disabled:cursor-not-allowed cursor-pointer"
        :disabled="!canAddCart"
        @click="handleClickAddCartButton"
      >
        {{ addToCartButtonText }}
      </button>
    </div>
  </div>
</template>
