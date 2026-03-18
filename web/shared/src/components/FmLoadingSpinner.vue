<script setup lang="ts">
interface Props {
  size?: 'sm' | 'md' | 'lg'
  overlay?: boolean
  label?: string
}

withDefaults(defineProps<Props>(), {
  size: 'md',
  overlay: false,
  label: '読み込み中...',
})

const sizeClasses: Record<string, string> = {
  sm: 'h-5 w-5 border-2',
  md: 'h-8 w-8 border-4',
  lg: 'h-12 w-12 border-4',
}
</script>

<template>
  <div
    :class="{
      'flex items-center justify-center': true,
      'fixed inset-0 z-50 bg-black/30': overlay,
      'py-4': !overlay,
    }"
    :role="label ? 'status' : 'presentation'"
    :aria-label="label || undefined"
  >
    <div class="flex flex-col items-center gap-3">
      <div
        :class="[
          'animate-spin rounded-full border-main border-t-transparent',
          sizeClasses[size],
        ]"
        aria-hidden="true"
      />
      <p
        v-if="label"
        class="text-sm text-typography"
      >
        {{ label }}
      </p>
    </div>
  </div>
</template>
