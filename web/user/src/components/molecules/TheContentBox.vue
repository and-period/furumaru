<script lang="ts" setup>
interface Props{
  title: string
  subTitle: string
}

defineProps<Props>()
</script>

<template>
  <div class="bg-white relative min-h-[350px] mt-6">
    <div class="absolute w-full -top-6">
      <the-content-title :value="title" class="mx-auto" />
    </div>
    <the-content-sub-title :value="subTitle" class="mt-12 mb-11" />
    <div class="py-6 px-4">
      <slot />
    </div>
  </div>
</template>
