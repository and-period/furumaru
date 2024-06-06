<script lang="ts" setup>
interface Props {
  withBadge?: boolean
}

interface Emits {
  (e: 'click'): void
}

withDefaults(defineProps<Props>(), {
  withBadge: false,
})

const emits = defineEmits<Emits>()

const handleClick = () => {
  emits('click')
}
</script>

<template>
  <button
    class="rounded-full p-2"
    @click="handleClick"
  >
    <div class="relative">
      <the-badge
        v-show="withBadge"
        class="absolute right-0 top-0"
      />
      <slot />
    </div>
  </button>
</template>
