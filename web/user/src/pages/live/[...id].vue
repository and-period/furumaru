<script setup lang="ts">
import dayjs from 'dayjs'
import { useScheduleStore } from '~/store/schedule'
import { useShoppingCartStore } from '~/store/shopping'
import type { ScheduleResponse } from '~/types/api'
import type { Snackbar } from '~/types/props'
import type { LiveTimeLineItem } from '~/types/props/schedule'

const scheduleStore = useScheduleStore()
const { getSchedule } = scheduleStore

const shoppingCartStore = useShoppingCartStore()
const { addCart } = shoppingCartStore

const snackbarItems = ref<Snackbar[]>([])

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
        // 生産者情報のマッピング
        const producer = schedule.value?.producers.find(
          (p) => p.id === live.producerId,
        )
        // 商品のマッピング
        const products = live.productIds.map((id) => {
          return schedule.value?.products.find((p) => p.id === id)
        })
        // コーディネーターのマッピング
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

const isLiveStreaming = computed<boolean>(() => {
  if (schedule.value) {
    return dayjs().isAfter(schedule.value.schedule.startAt)
  } else {
    return false
  }
})

const handleClickAddCart = (name: string, id: string, quantity: number) => {
  addCart({ productId: id, quantity })
  snackbarItems.value.push({
    text: `買い物カゴに「${name}」を追加しました`,
    isShow: true,
  })
}

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
  <template v-for="(snackbarItem, i) in snackbarItems" :key="i">
    <the-snackbar
      v-model:is-show="snackbarItem.isShow"
      :text="snackbarItem.text"
    />
  </template>

  <div
    class="mx-auto grid max-w-[1440px] grid-flow-col auto-rows-max grid-cols-3 gap-8 text-main xl:px-14"
  >
    <template v-if="schedule">
      <div class="col-span-3">
        <the-live-video-player
          :video-src="schedule.schedule.distributionUrl"
          :title="schedule.schedule.title"
          :start-at="schedule.schedule.startAt"
          :end-at="schedule.schedule.endAt"
          :description="schedule.schedule.description"
          :is-live-streaming="isLiveStreaming"
          :marche-name="schedule.coordinator.marcheName"
          :address="schedule.coordinator.city"
          :cn-name="schedule.coordinator.username"
          :cn-img-src="schedule.coordinator.thumbnailUrl"
        />
        <the-live-timeline
          class="mt-4"
          :items="liveTimeLineItems"
          @click:add-cart="handleClickAddCart"
        />
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
