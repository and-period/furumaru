<script lang="ts" setup>
interface Props {
  title: string
  noItemText: string
  notificationItems: any[]
}

defineProps<Props>()
</script>

<template>
  <the-dropdown-with-icon>
    <template #icon>
      <the-ring-icon />
    </template>
    <template #content>
      <div class="p-4">
        <p class="mb-6 font-semibold">
          {{ title }}
        </p>

        <p v-if="notificationItems.length === 0">
          {{ noItemText }}
        </p>

        <p
          v-for="(item, i) in notificationItems"
          :key="i"
        >
          {{ item }}
        </p>
      </div>
    </template>
  </the-dropdown-with-icon>
</template>
