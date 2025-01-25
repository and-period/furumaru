<script setup lang="ts">
interface Props {
  id: string | undefined
  title: string | undefined
  imgSrc: string | undefined
  width: number | undefined
  endAt: number
  archivedStreamText: string
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
    class="w-full text-main"
    @click="handleClick"
  >
    <div class="w-full">
      <nuxt-img
        provider="cloudFront"
        class="h-[208px] w-full object-contain"
        :src="imgSrc"
        :alt="`archive-${title}-thumbnail`"
        fit="contain"
        height="208px"
      />
    </div>
    <div class="mt-2 flex w-full flex-col gap-2">
      <p class="line-clamp-3 break-words text-[14px] tracking-[10%]">
        {{ title }}
      </p>
    </div>
  </div>
</template>
