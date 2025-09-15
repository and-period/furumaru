<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useAlert, usePagination } from '~/lib/hooks'
import { useAuthStore, useCommonStore, useVideoStore } from '~/store'
import type { Video } from '~/types/api/v1'

const router = useRouter()
const commonStore = useCommonStore()
const authStore = useAuthStore()
const videoStore = useVideoStore()
const pagination = usePagination()
const { alertType, isShow, alertText, show } = useAlert('error')

const { adminType } = storeToRefs(authStore)
const { videoResponse } = storeToRefs(videoStore)

const loading = ref<boolean>(false)
const deleteDialog = ref<boolean>(false)

const fetchState = useAsyncData(async (): Promise<void> => {
  await fetchVideos()
})

watch(pagination.itemsPerPage, (): void => {
  fetchVideos()
})

const fetchVideos = async (): Promise<void> => {
  try {
    await videoStore.fetchVideos(pagination.itemsPerPage.value, pagination.offset.value)
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const isLoading = (): boolean => {
  return fetchState?.pending?.value || loading.value
}

const handleUpdatePage = async (page: number): Promise<void> => {
  pagination.updateCurrentPage(page)
  await fetchVideos()
}

const handleClickAdd = (): void => {
  router.push('/videos/new')
}

const handleClickRow = (videoId: string): void => {
  router.push(`/videos/${videoId}`)
}

const handleClickDelete = async (videoId: string): Promise<void> => {
  try {
    loading.value = true
    const video = videoResponse.value?.videos.find((video: Video): boolean => {
      return video.id === videoId
    })
    if (!video) {
      throw new Error(`failed to find video. videoId=${videoId}`)
    }
    await videoStore.deleteVideo(videoId)
    commonStore.addSnackbar({
      message: `${video.title}を削除しました。`,
      color: 'info',
    })
    deleteDialog.value = false
    fetchState.execute()
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
  }
}

try {
  await fetchState.execute()
}
catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-video-list
    v-model:delete-dialog="deleteDialog"
    :loading="isLoading()"
    :admin-type="adminType"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :videos="videoResponse?.videos || []"
    :table-items-per-page="pagination.itemsPerPage.value"
    :table-items-total="videoResponse?.total || 0"
    @click:row="handleClickRow"
    @click:add="handleClickAdd"
    @click:delete="handleClickDelete"
    @click:update-page="handleUpdatePage"
    @click:update-items-per-page="pagination.handleUpdateItemsPerPage"
  />
</template>
