<script setup lang="ts">
import { useBreadcrumbJsonLd } from '~/hooks/seo'

interface BreadcrumbItem {
  name: string
  path: string
}

const { items } = defineProps<{
  items: BreadcrumbItem[]
}>()

useBreadcrumbJsonLd(items)
</script>

<template>
  <nav
    aria-label="パンくずリスト"
    class="px-4 py-3 text-sm text-typography lg:px-[112px]"
  >
    <ol class="flex flex-wrap items-center gap-1">
      <li
        v-for="(item, i) in items"
        :key="item.path"
        class="flex items-center gap-1"
      >
        <NuxtLink
          v-if="i < items.length - 1"
          :to="item.path"
          class="hover:underline"
        >
          {{ item.name }}
        </NuxtLink>
        <span
          v-else
          aria-current="page"
          class="text-main font-medium"
        >
          {{ item.name }}
        </span>
        <span
          v-if="i < items.length - 1"
          aria-hidden="true"
          class="text-gray-400"
        >/</span>
      </li>
    </ol>
  </nav>
</template>
