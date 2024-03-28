<script setup lang="ts">
interface Props {
  title: string
}

defineProps<Props>()
</script>

<template>
  <div
    class="mx-auto flex max-w-[560px] flex-col items-center justify-center bg-white px-4 py-6 tracking-[10%] text-main md:p-10 xl:px-20"
  >
    <div class="text-center text-[16px] font-semibold">
      {{ title }}
    </div>
    <slot />
  </div>
</template>
