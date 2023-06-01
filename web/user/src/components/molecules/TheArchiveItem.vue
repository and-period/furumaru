<script setup lang='ts'>
interface Props {
  id: string
  title: string
  imgSrc: string
}

defineProps<Props>()
</script>

<template>
  <div class="text-main min-w-[368px] ">
    <div class="w-full">
      <img
        class="w-full h-[208px] object-cover"
        :src="imgSrc"
        :alt="`archive-${title}-thumbnail`"
      >
    </div>
    <div class="mt-2 flex flex-col gap-2 w-full">
      <div class="text-sm">
        <span class="font-bold px-2 border-2 border-main text-main rounded">アーカイブ配信</span>
      </div>
      <p class="line-clamp-3 break-words">
        {{ title }}
      </p>
    </div>
  </div>
</template>
