<script lang="ts" setup>
interface Props {
  id: string
  title: string
  priceAdult: number
  thumbnailUrl: string
}

const props = defineProps<Props>()

interface Emits {
  (e: 'click:item', id: string): void
}

const emits = defineEmits<Emits>()

const priceString = computed<string>(() => {
  return new Intl.NumberFormat('ja-JP', {
    style: 'currency',
    currency: 'JPY',
  }).format(props.priceAdult)
})

const handleClickItem = () => {
  emits('click:item', props.id)
}
</script>

<template>
  <div class="mx-auto flex flex-col text-main">
    <div class="relative mx-auto max-w-[144px]">
      <!-- 体験バッジ -->
      <div class="absolute top-1 left-1 z-10 bg-blue-500 text-white px-1 py-0.5 rounded text-[10px] font-bold">
        体験
      </div>
      <picture
        class="w-full hover:cursor-pointer"
        @click="handleClickItem"
      >
        <nuxt-img
          provider="cloudFront"
          :src="thumbnailUrl"
          :alt="`${title}のサムネイル画像`"
          class="aspect-square w-full"
        />
      </picture>
    </div>

    <p
      class="mt-2 line-clamp-3 max-w-[144px] grow text-[14px] tracking-[1.4px] hover:cursor-pointer hover:underline md:text-[14px] md:tracking-[1.6px]"
      @click="handleClickItem"
    >
      {{ title }}
    </p>

    <p
      class="my-2 mb-4 text-[16px] font-bold tracking-[1.6px] md:text-[16px] md:tracking-[2.0px]"
    >
      {{ priceString }}〜
    </p>
  </div>
</template>
