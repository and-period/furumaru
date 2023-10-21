<script setup lang="ts">
import { MOCK_LIVE_TIMELINES } from '~/constants/mock'

const liveTimelines = MOCK_LIVE_TIMELINES
</script>

<template>
  <div class="bg-white py-7 pl-8 pr-4">
    <ol class="relative ml-[48px] border-l-[2px] border-orange">
      <the-live-timeline-item
        v-for="liveTimeline in liveTimelines"
        :key="liveTimeline.id"
        :time="liveTimeline.time"
        :cn-img-src="liveTimeline.cnImgSrc"
        :marche-name="liveTimeline.marchName"
        :description="liveTimeline.description"
        :items="liveTimeline.items"
      />
    </ol>
  </div>
</template>
