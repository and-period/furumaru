<script setup lang="ts">
interface Props {
  id: string | undefined
  title: string | undefined
  imgSrc: string | undefined
  width: number | undefined
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
      <img
        class="h-[208px] w-full object-cover"
        :src="imgSrc"
        :alt="`archive-${title}-thumbnail`"
      />
    </div>
    <div class="mt-2 flex w-full flex-col gap-2">
      <div class="flex items-center justify-between text-sm">
        <span class="rounded border-2 border-main px-2 font-bold text-main">
          アーカイブ配信
        </span>
      </div>
      <p class="line-clamp-3 break-words">
        {{ title }}
      </p>
    </div>
  </div>
</template>
