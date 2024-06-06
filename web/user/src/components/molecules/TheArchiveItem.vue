<script setup lang="ts">
interface Props {
  id: string | undefined
  title: string | undefined
  imgSrc: string | undefined
  width: number | undefined
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
  <div class="w-full text-main" @click="handleClick">
    <div class="w-full">
      <nuxt-img
        provider="cloudFront"
        class="h-[208px] w-full object-cover"
        :src="imgSrc"
        :alt="`archive-${title}-thumbnail`"
        fit="cover"
        height="208px"
      />
    </div>
    <div class="mt-2 flex w-full flex-col gap-2">
      <div class="flex items-center justify-between text-sm">
        <span
          class="rounded border border-main px-1 text-[10px] font-bold tracking-[10%] text-main"
        >
          {{ archivedStreamText }}
        </span>
      </div>
      <p class="line-clamp-3 break-words text-[14px] tracking-[10%]">
        {{ title }}
      </p>
    </div>
  </div>
</template>
