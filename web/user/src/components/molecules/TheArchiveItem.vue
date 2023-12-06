<script setup lang="ts">
interface Props {
  id: string | undefined
  title: string | undefined
  imgSrc: string | undefined
}

defineProps<Props>()
</script>

<template>
  <div class="min-w-[368px] text-main">
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
        <button class="h-4 w-4 hover:scale-110">
          <the-ellipsis-vertical-icon class="h-5 w-5" />
        </button>
      </div>
      <p class="line-clamp-3 break-words">
        {{ title }}
      </p>
    </div>
  </div>
</template>
