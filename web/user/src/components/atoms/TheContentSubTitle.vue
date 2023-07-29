<script lang="ts" setup>
interface Props {
  value: string
}

defineProps<Props>()
</script>

<template>
  <div class="flex items-center justify-center gap-x-4 text-main">
    <span class="inline-block w-5 border-t border-main" />
    {{ value }}
    <span class="inline-block w-5 border-t border-main" />
  </div>
</template>
