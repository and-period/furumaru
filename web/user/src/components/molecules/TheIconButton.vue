<script lang="ts" setup>
interface Props {
  withBadge?: boolean
  ariaLabel?: string
  ariaExpanded?: boolean
  ariaHaspopup?: boolean | string
}

interface Emits {
  (e: 'click'): void
}

withDefaults(defineProps<Props>(), {
  withBadge: false,
  ariaLabel: undefined,
  ariaExpanded: undefined,
  ariaHaspopup: undefined,
})

const emits = defineEmits<Emits>()

const handleClick = () => {
  emits('click')
}
</script>

<template>
  <button
    class="rounded-full p-2"
    :aria-label="ariaLabel || undefined"
    :aria-expanded="ariaExpanded"
    :aria-haspopup="ariaHaspopup || undefined"
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
