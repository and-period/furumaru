<script lang="ts" setup>
interface Props {
  value: string
}

defineProps<Props>()
</script>

<template>
  <div class="text-main flex justify-center items-center gap-x-4">
    <span class="border-main border-t inline-block w-5 " />
    {{ value }}
    <span class="border-main border-t inline-block w-5 " />
  </div>
</template>
