<script setup lang="ts">
import type { LiveTimeLineItem } from '~/types/props/schedule'

interface Props {
  items: LiveTimeLineItem[]
}

defineProps<Props>()
</script>

<template>
  <div class="bg-white py-7 pl-8 pr-4">
    <ol class="relative ml-[48px] border-l-[2px] border-orange">
      <the-live-timeline-item
        v-for="(liveTimeline, i) in items"
        :key="i"
        :start-at="liveTimeline.startAt"
        :thumbnail-url="liveTimeline.producer?.thumbnailUrl"
        :username="liveTimeline.producer?.username"
        :comment="liveTimeline.comment"
        :items="liveTimeline.products"
      />
    </ol>
  </div>
</template>
