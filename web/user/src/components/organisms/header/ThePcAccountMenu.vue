<script lang="ts" setup>
import type { LinkItem } from '~/types/props'

interface Props {
  isAuthenticated: boolean
  authenticatedMenuItems: LinkItem[]
  noAuthenticatedMenuItems: LinkItem[]
}

defineProps<Props>()
</script>

<template>
  <the-dropdown-with-icon>
    <template #icon>
      <the-account-icon />
    </template>
    <template #content>
      <div v-if="isAuthenticated">
        <p v-for="(item, i) in authenticatedMenuItems" :key="i">
          {{ item.text }}
        </p>
      </div>
      <div v-else class="flex flex-col">
        <nuxt-link
          v-for="(item, i) in noAuthenticatedMenuItems"
          :key="i"
          :href="item.href"
          class="px-4 py-2 hover:bg-gray-200"
        >
          {{ item.text }}
        </nuxt-link>
      </div>
    </template>
  </the-dropdown-with-icon>
</template>
