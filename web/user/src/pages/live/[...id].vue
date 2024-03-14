<script setup lang="ts">
import { useAuthStore } from '~/store/auth'
import { useScheduleStore } from '~/store/schedule'
import { useShoppingCartStore } from '~/store/shopping'
import {
  ScheduleStatus,
  type ScheduleResponse,
  type LiveComment,
} from '~/types/api'
import type { Snackbar } from '~/types/props'
import type { LiveTimeLineItem } from '~/types/props/schedule'

const authStore = useAuthStore()
const { isAuthenticated } = storeToRefs(authStore)

const scheduleStore = useScheduleStore()
const { getSchedule, postComment, getComments } = scheduleStore

const shoppingCartStore = useShoppingCartStore()
const { addCart } = shoppingCartStore

const snackbarItems = ref<Snackbar[]>([])

const router = useRouter()
const route = useRoute()

const isLoading = ref<boolean>(false)
const schedule = ref<ScheduleResponse | undefined>(undefined)

const comments = ref<LiveComment[]>([])
const commentFormData = ref<string>('')
const commentIsSending = ref<boolean>(false)

const scheduleId = computed<string>(() => {
  return route.params.id as string
})

const handleSubmitComment = async () => {
  try {
    commentIsSending.value = true
    await postComment(scheduleId.value, commentFormData.value)
    commentFormData.value = ''
    const res = await getComments(scheduleId.value)
    comments.value = res.comments
  } catch (e) {
    snackbarItems.value.push({
      text: 'コメントの送信に失敗しました。',
      isShow: true,
    })
  } finally {
    commentIsSending.value = false
  }
}

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
    return schedule.value.schedule.status === ScheduleStatus.LIVE
  } else {
    return false
  }
})

const isArchive = computed<boolean>(() => {
  if (schedule.value) {
    return schedule.value.schedule.status === ScheduleStatus.ARCHIVED
  } else {
    return false
  }
})

const liveRef = ref<{ videoRef: HTMLVideoElement | null }>({ videoRef: null })

const livePlayerHeight = computed(() => {
  if (liveRef.value.videoRef) {
    if (liveRef.value.videoRef.offsetWidth >= 768) {
      return 0
    }
    return liveRef.value.videoRef.offsetHeight
  }
  return 0
})

const selectedTab = ref<'product' | 'comment'>('product')

const clickTab = (tab: 'product' | 'comment') => {
  selectedTab.value = tab
}

const handleClickItem = (productId: string) => {
  router.push(`/items/${productId}`)
}

const handleClickAddCart = (name: string, id: string, quantity: number) => {
  addCart({ productId: id, quantity })
  snackbarItems.value.push({
    text: `買い物カゴに「${name}」を追加しました`,
    isShow: true,
  })
}

const handleCLickCoordinator = (id: string) => {
  router.push(`/coordinator/${id}`)
}

const fetchComments = async () => {
  const res = await getComments(scheduleId.value)
  comments.value = res.comments
}

useAsyncData(`schedule-${scheduleId.value}`, async () => {
  isLoading.value = true
  const scheduleRes = await getSchedule(scheduleId.value)
  schedule.value = scheduleRes
  isLoading.value = false
})

onMounted(() => {
  if (isLiveStreaming.value) {
    fetchComments()
    const interval = setInterval(() => {
      fetchComments()
    }, 5000)

    onUnmounted(() => {
      clearInterval(interval)
    })
  } else {
    fetchComments()
  }
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
          ref="liveRef"
          :video-src="schedule.schedule.distributionUrl"
          :is-archive="isArchive"
          class="fixed z-[20] w-full md:static"
        />
        <div :style="{ 'padding-top': `${livePlayerHeight}px` }">
          <the-live-description
            :title="schedule.schedule.title"
            :description="schedule.schedule.description"
            :is-archive="isArchive"
            :is-live-streaming="isLiveStreaming"
            :start-at="schedule.schedule.startAt"
            :marche-name="schedule.coordinator.marcheName"
            :coordinator-id="schedule.coordinator.id"
            :coordinator-name="schedule.coordinator.username"
            :coordinator-img-src="schedule.coordinator.thumbnailUrl"
            :coordinator-address="schedule.coordinator.city"
            @click:coordinator="handleCLickCoordinator"
          />
          <div class="mt-4 flex w-full flex-col rounded">
            <div
              class="grid w-full grid-cols-2 gap-2 text-[12px] font-bold text-main"
            >
              <button
                class="rounded-t-xl p-4 text-center"
                :class="{
                  'bg-white': selectedTab === 'product',
                  'bg-main text-white hover:bg-main/80':
                    selectedTab !== 'product',
                }"
                @click="clickTab('product')"
              >
                このマルシェの商品
              </button>
              <button
                class="rounded-t-xl p-4 text-center"
                :class="{
                  'bg-white': selectedTab === 'comment',
                  'bg-main text-white hover:bg-main/80':
                    selectedTab !== 'comment',
                }"
                @click="clickTab('comment')"
              >
                コメント
              </button>
            </div>
            <template v-if="selectedTab === 'product'">
              <the-live-timeline
                :items="liveTimeLineItems"
                @click:item="handleClickItem"
                @click:add-cart="handleClickAddCart"
              />
            </template>
            <template v-else>
              <div class="flex flex-col bg-white p-4">
                <the-live-comment-input-form
                  v-model="commentFormData"
                  class="order-1 md:order-[0]"
                  :is-authenticated="isAuthenticated"
                  :is-sending="commentIsSending"
                  @submit="handleSubmitComment"
                />
                <div class="flex flex-col gap-4 py-8">
                  <div v-if="comments.length === 0" class="text-typography">
                    コメントがありません。
                  </div>
                  <div
                    v-for="(item, i) in comments"
                    :key="i"
                    class="flex items-center gap-3"
                  >
                    <img
                      :src="
                        item.thumbnailUrl
                          ? item.thumbnailUrl
                          : '/img/furuneko.png'
                      "
                      class="h-8 w-8 rounded-full"
                    />
                    <div
                      class="flex flex-row items-center gap-2 tracking-[10%] md:flex-col md:items-baseline md:gap-0"
                    >
                      <div
                        class="whitespace-nowrap text-[14px] text-typography"
                      >
                        {{ item.username ? item.username : 'ゲスト' }}
                      </div>
                      <div class="line-clamp-2 text-main">
                        {{ item.comment }}
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </template>
          </div>
        </div>
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
