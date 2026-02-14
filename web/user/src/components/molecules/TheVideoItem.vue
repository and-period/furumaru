<script setup lang="ts">
interface Props {
  id: string | undefined
  title: string | undefined
  imgSrc: string | undefined
  width: number | undefined
  archivedStreamText: string
  endAt?: number
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
  <div
    class="w-full text-main rounded-xl overflow-hidden shadow-md hover:shadow-lg transition-shadow duration-300 bg-white"
    @click="handleClick"
  >
    <div class="w-full">
      <nuxt-img
        provider="cloudFront"
        class="h-[208px] w-full object-cover"
        :src="imgSrc"
        :alt="`archive-${title}-thumbnail`"
        fit="contain"
        height="208px"
      />
    </div>
    <div class="mt-2 flex w-full flex-col gap-2 p-3">
      <p class="line-clamp-3 break-words text-[14px] md:text-[16px]">
        {{ title }}
      </p>
    </div>
  </div>
</template>
