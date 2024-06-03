<script setup lang="ts">
interface Props {
  id: string
  title: string
  imgSrc: string
}

defineProps<Props>()

interface Emits {
  (e: 'click'): void
}

const emits = defineEmits<Emits>()

const handleClick = () => {
  emits('click')
}
</script>

<template>
  <div @click="handleClick">
    <div class="flex justify-center">
      <nuxt-img
        provider="cloudFront"
        fit="contain"
        height="208px"
        class="aspect-video h-[208px] cursor-pointer object-cover"
        :src="imgSrc"
        :alt="`live-${title}-thumbnail`"
      />
    </div>

    <div class="ml-0 mt-4 flex w-full flex-col gap-2">
      <p
        class="mx-4 line-clamp-3 break-words text-[14px] text-main md:mx-0 md:text-[16px]"
      >
        {{ title }}
      </p>
    </div>
  </div>
</template>
