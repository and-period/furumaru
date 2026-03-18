<script setup lang="ts">
import { type Component } from 'vue'

interface Props {
  title: string
  description?: string
  icon?: string
  actionText?: string
  actionComponent?: string | Component
  actionComponentProps?: Record<string, unknown>
}

interface Emits {
  (e: 'click:action'): void
}

withDefaults(defineProps<Props>(), {
  description: '',
  icon: '📭',
  actionText: '',
  actionComponent: undefined,
  actionComponentProps: undefined,
})

const emits = defineEmits<Emits>()

const handleClickAction = () => {
  emits('click:action')
}
</script>

<template>
  <div class="flex flex-col items-center justify-center px-4 py-12 text-center">
    <div
      class="text-4xl mb-4"
      aria-hidden="true"
    >
      {{ icon }}
    </div>
    <h3 class="text-lg font-semibold text-main mb-2">
      {{ title }}
    </h3>
    <p
      v-if="description"
      class="text-sm text-typography mb-6 max-w-md"
    >
      {{ description }}
    </p>
    <slot name="action">
      <component
        :is="actionComponent || 'button'"
        v-if="actionText"
        :type="!actionComponent ? 'button' : undefined"
        v-bind="actionComponentProps"
        class="bg-orange text-white px-6 py-2.5 font-medium transition-all duration-200 hover:bg-orange/90"
        @click="handleClickAction"
      >
        {{ actionText }}
      </component>
    </slot>
  </div>
</template>
