<script setup lang="ts">
import type { LiveTimeLineItem } from '~/types/props/schedule'

interface Props {
  items: LiveTimeLineItem[]
}

defineProps<Props>()

interface Emits {
  (e: 'click:addCart', name: string, id: string, quantity: number): void
}

const emits = defineEmits<Emits>()

const handleClickAddCart = (name: string, id: string, quantity: number) => {
  emits('click:addCart', name, id, quantity)
}
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
        @click:add-cart="handleClickAddCart"
      />
    </ol>
  </div>
</template>
