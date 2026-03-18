<script setup lang="ts">
import { type Component, computed } from 'vue'

interface Props {
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost'
  size?: 'sm' | 'md' | 'lg'
  disabled?: boolean
  loading?: boolean
  type?: 'button' | 'submit' | 'reset'
  as?: string | Component
}

interface Emits {
  (e: 'click', event: MouseEvent): void
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
  disabled: false,
  loading: false,
  type: 'button',
  as: 'button',
})

const emits = defineEmits<Emits>()

const isDisabled = computed(() => props.disabled || props.loading)
const isNativeButton = computed(() => props.as === 'button')

const handleClick = (e: MouseEvent) => {
  if (isDisabled.value) {
    e.preventDefault()
    e.stopPropagation()
    return
  }
  emits('click', e)
}

const variantClasses: Record<string, string> = {
  primary: 'bg-orange text-white hover:bg-orange/90 disabled:bg-orange/40',
  secondary: 'bg-main text-white hover:bg-main/90 disabled:bg-main/40',
  danger: 'bg-error text-white hover:bg-error/90 disabled:bg-error/40',
  ghost: 'bg-transparent text-main border border-main hover:bg-main/5 disabled:opacity-40',
}

const sizeClasses: Record<string, string> = {
  sm: 'px-3 py-1.5 text-sm',
  md: 'px-4 py-2.5 text-base',
  lg: 'px-6 py-3 text-lg',
}
</script>

<template>
  <component
    :is="as"
    :type="isNativeButton ? type : undefined"
    :disabled="isNativeButton ? isDisabled : undefined"
    :aria-disabled="!isNativeButton && isDisabled ? true : undefined"
    :tabindex="!isNativeButton && isDisabled ? -1 : undefined"
    :class="[
      'inline-flex items-center justify-center font-medium transition-all duration-200 ease-in-out',
      isDisabled ? 'cursor-not-allowed opacity-40' : '',
      variantClasses[variant],
      sizeClasses[size],
    ]"
    @click="handleClick"
  >
    <svg
      v-if="loading"
      class="mr-2 h-4 w-4 animate-spin"
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
      aria-hidden="true"
    >
      <circle
        class="opacity-25"
        cx="12"
        cy="12"
        r="10"
        stroke="currentColor"
        stroke-width="4"
      />
      <path
        class="opacity-75"
        fill="currentColor"
        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
      />
    </svg>
    <slot />
  </component>
</template>
