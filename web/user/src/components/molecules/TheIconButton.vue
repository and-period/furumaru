<script lang="ts" setup>
interface Props {
  withBadge?: boolean
}

interface Emits {
  (e: 'click'): void;
}

withDefaults(defineProps<Props>(), {
  withBadge: false
})

const emits = defineEmits<Emits>()

const handleClick = () => {
  emits('click')
}
</script>

<template>
  <button class="p-2 rounded-full" @click="handleClick">
    <div class="relative">
      <the-badge v-show="withBadge" class="absolute top-0 right-0" />
      <slot />
    </div>
  </button>
</template>
