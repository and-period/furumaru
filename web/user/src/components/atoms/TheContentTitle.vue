<script lang="ts" setup>
interface Props {
  value: string
}

defineProps<Props>()
</script>

<template>
  <div
    class="w-72 rounded-3xl border-2 border-main bg-white py-2 text-center text-xl font-bold uppercase text-main"
  >
    {{ value }}
  </div>
</template>
