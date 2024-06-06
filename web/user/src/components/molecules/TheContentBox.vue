<script lang="ts" setup>
interface Props {
  title: string
  subTitle: string
}

defineProps<Props>()
</script>

<template>
  <div class="relative mt-6 min-h-[350px] bg-white">
    <div class="absolute -top-6 w-full">
      <the-content-title
        :value="title"
        class="mx-auto"
      />
    </div>
    <the-content-sub-title
      :value="subTitle"
      class="mb-11 mt-12"
    />
    <div class="px-4 py-6">
      <slot />
    </div>
  </div>
</template>
