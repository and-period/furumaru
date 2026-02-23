<script lang="ts" setup>
import type { AlertType } from '~/lib/hooks'

interface Props {
  show: boolean
  type: AlertType | undefined
  text: string
}

defineProps<Props>()
</script>

<template>
  <v-alert
    v-show="show"
    :type="type"
    aria-live="polite"
    v-text="text"
  />
</template>
