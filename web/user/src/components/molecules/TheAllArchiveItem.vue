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
      <img
        class="aspect-video max-h-[208px] cursor-pointer object-cover"
        :src="imgSrc"
        :alt="`live-${title}-thumbnail`"
      />
    </div>

    <div class="ml-0 mt-4 flex w-full flex-col gap-2">
      <div class="mx-4 flex items-center text-sm md:mx-0">
        <div class="grow">
          <span
            class="rounded border-2 border-main px-2 font-bold text-main"
          >
          アーカイブ配信
          </span>
        </div>
      </div>

      <p
        class="mx-4 line-clamp-3 break-words text-[14px] text-main md:mx-0 md:text-[16px]"
      >
        {{ title }}
      </p>
    </div>
  </div>
</template>
