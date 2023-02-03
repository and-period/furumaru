<script lang="ts" setup>
interface Props {
  num?: number
}

defineProps<Props>()
</script>

<template>
  <div class="py-12 px-28 bg-white rounded-full w-[620px] h-[620px] text-center text-main">
    <p class="text-orange uppercase text-3xl font-bold mb-3">
      Point
    </p>
    <div class="block relative w-[68px] h-[85px] mx-auto mb-8">
      <p class="font-extrabold text-4xl text-white absolute bottom-2 left-[16px]">
        {{ num }}
      </p>
      <the-point />
    </div>
    <div class="flex flex-col">
      <slot />
    </div>
  </div>
</template>
