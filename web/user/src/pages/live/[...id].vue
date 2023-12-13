<script setup lang="ts">
import { useScheduleStore } from '~/store/schedule'
import type { ScheduleResponse } from '~/types/api'
import type { LiveTimeLineItem } from '~/types/props/schedule'

const scheduleStore = useScheduleStore()
const { getSchedule } = scheduleStore

const route = useRoute()

const isLoading = ref<boolean>(false)
const schedule = ref<ScheduleResponse | undefined>(undefined)

const scheduleId = computed<string>(() => {
  return route.params.id as string
})

const liveTimeLineItems = computed<LiveTimeLineItem[]>(() => {
  if (schedule.value) {
    return (
      schedule.value.lives.map((live) => {
        const producer = schedule.value?.producers.find(
          (p) => p.id === live.producerId,
        )
        const products =
          schedule.value?.products.filter((p) => {
            return p.id in live.productIds
          }) ?? []
        return {
          ...live,
          producer,
          products,
        }
      }) ?? []
    )
  } else {
    return []
  }
})

onMounted(async () => {
  isLoading.value = true
  const res = await getSchedule(scheduleId.value)
  schedule.value = res
  isLoading.value = false
})

useSeoMeta({
  title: 'ライブ配信',
})
</script>

<template>
  <div
    class="mx-auto grid max-w-[1440px] grid-flow-col auto-rows-max grid-cols-3 gap-8 text-main xl:px-14"
  >
    <template v-if="schedule">
      <div class="col-span-3">
        <the-live-video-player
          :video-src="schedule.schedule.distributionUrl"
          :title="schedule.schedule.title"
          :start-at="schedule.schedule.startAt"
          :description="schedule.schedule.description"
          :is-archive="false"
          marche-name=""
          inventory:address=""
          cn-img-src=""
          cn-name=""
        />
        <the-live-timeline class="mt-4" :items="liveTimeLineItems" />
      </div>
    </template>

    <!-- PC画面のみ表示する右サイドバー -->
    <!--
    <div class="col-span-1 hidden lg:block">
      <div class="flex h-[450px] flex-col bg-white p-4">
        <div class="grow overflow-scroll"></div>
        <div class="relative text-[12px] text-main">
          <input
            type="text"
            placeholder="コメントを投稿する"
            class="w-full rounded-[20px] border border-main py-2 pl-4 pr-[56px]"
          />
          <button class="absolute right-4 top-2 min-w-max">送信</button>
        </div>
      </div>
    </div>
    -->
  </div>
</template>
