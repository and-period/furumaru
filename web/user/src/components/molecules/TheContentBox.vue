<script lang="ts" setup>
interface Props{
  title: string
  subTitle: string
}

defineProps<Props>()
</script>

<template>
  <div class="bg-white relative min-h-[350px] mt-6 py-6 px-4">
    <div class="absolute w-full -top-6">
      <the-content-title :value="title" class="mx-auto" />
    </div>
    <the-content-sub-title :value="subTitle" class="mt-12 mb-11" />
    <slot />
  </div>
</template>
