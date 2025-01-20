<script setup lang="ts">
import { useAuthStore } from '~/store/auth'
import { useVideoStore } from '~/store/video'
import { useShoppingCartStore } from '~/store/shopping'
import {
  type VideoComment,
  type VideoResponse,
} from '~/types/api'
import type { Snackbar } from '~/types/props'
import type { I18n } from '~/types/locales'

const authStore = useAuthStore()
const { isAuthenticated } = storeToRefs(authStore)

const videoStore = useVideoStore()
const { getVideo, postComment, getComments } = videoStore

const shoppingCartStore = useShoppingCartStore()
const { addCart } = shoppingCartStore

const snackbarItems = ref<Snackbar[]>([])

const i18n = useI18n()

const dt = (str: keyof I18n['lives']['details']) => {
  return i18n.t(`videos.details.${str}`)
}

const router = useRouter()
const route = useRoute()

const isLoading = ref<boolean>(false)
const video = ref<VideoResponse | undefined>(undefined)

const comments = ref<VideoComment[]>([])
const commentFormData = ref<string>('')
const commentIsSending = ref<boolean>(false)

const videoId = computed<string>(() => {
  return route.params.id as string
})

const handleSubmitComment = async () => {
  try {
    commentIsSending.value = true
    await postComment(videoId.value, commentFormData.value)
    commentFormData.value = ''
    const res = await getComments(videoId.value)
    comments.value = res.comments
  }
  catch (e) {
    snackbarItems.value.push({
      text: 'コメントの送信に失敗しました。',
      isShow: true,
    })
    console.log(e)
  }
  finally {
    commentIsSending.value = false
  }
}

const videoRef = ref<{ videoRef: HTMLVideoElement | null }>({ videoRef: null })

const videoPlayerHeight = computed(() => {
  if (videoRef.value.videoRef) {
    if (videoRef.value.videoRef.offsetWidth >= 768) {
      return 0
    }
    return videoRef.value.videoRef.offsetHeight
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
  const message = i18n.t('items.details.addCartSnackbarMessage', {
    itemName: name,
  })
  addCart({ productId: id, quantity })
  snackbarItems.value.push({
    text: message,
    isShow: true,
  })
}

const handleCLickCoordinator = (id: string) => {
  router.push(`/coordinator/${id}`)
}

const fetchComments = async () => {
  const res = await getComments(videoId.value)
  comments.value = res.comments
}

await useAsyncData(`schedule-${videoId.value}`, async () => {
  isLoading.value = true
  const videoRes = await getVideo(videoId.value)
  video.value = videoRes
  isLoading.value = false
})

onMounted(() => {
  fetchComments()
  const interval = setInterval(() => {
    fetchComments()
  }, 3000)

  onUnmounted(() => {
    clearInterval(interval)
  })
})

useSeoMeta({
  title: '動画',
})
</script>

<template>
  <template
    v-for="(snackbarItem, i) in snackbarItems"
    :key="i"
  >
    <the-snackbar
      v-model:is-show="snackbarItem.isShow"
      :text="snackbarItem.text"
    />
  </template>

  <div
    class="mx-auto grid max-w-[1440px] grid-flow-col auto-rows-max grid-cols-3 gap-8 text-main xl:px-14"
  >
    <template v-if="video">
      <div class="col-span-3">
        <the-video-player
          ref="liveRef"
          :video-src="video.video.videoUrl"
          class="fixed z-[20] w-full md:static"
        />
        <div :style="{ 'padding-top': `${videoPlayerHeight}px` }">
          <the-video-description
            :title="video.video.title"
            :description="video.video.description"
            :start-at="video.video.publishedAt"
            :marche-name="video.coordinator.marcheName"
            :coordinator-id="video.coordinator.id"
            :coordinator-name="video.coordinator.username"
            :coordinator-img-src="video.coordinator.thumbnailUrl"
            :coordinator-address="video.coordinator.city"
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
                {{ dt("itemsTabLabel") }}
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
                {{ dt("commentsTabLabel") }}
              </button>
            </div>
            <template v-if="selectedTab === 'product'">
              xxxx
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
                  <div
                    v-if="comments.length === 0"
                    class="text-typography"
                  >
                    {{ dt("noCommentsText") }}
                  </div>
                  <div
                    v-for="(item, i) in comments"
                    :key="i"
                    class="flex items-center gap-3"
                  >
                    <nuxt-img
                      provider="cloudFront"
                      width="32px"
                      height="32px"
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
                        {{
                          item.username ? item.username : dt("guestNameLabel")
                        }}
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
