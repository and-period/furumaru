<script lang="ts" setup>
interface Props {
  value: string
}

defineProps<Props>()
</script>

<template>
  <div class="w-72 bg-white rounded-3xl border-2 text-main border-main uppercase text-center font-bold text-xl py-2">
    {{ value }}
  </div>
</template>
