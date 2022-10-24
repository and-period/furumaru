<template>
  <the-notification-edit-form-page
    :form-data="formData"
    :time-data="timeData"
    :form-data-loading="fetchState.pending"
    @submit="handleSubmit"
  />
</template>

<script lang="ts">
import { useFetch, useRoute, useRouter } from '@nuxtjs/composition-api'
import { defineComponent, reactive } from '@vue/composition-api'
import dayjs from 'dayjs'

import { useNotificationStore } from '~/store/notification'
import { NotificationResponse } from '~/types/api'
import { NotificationTime } from '~/types/props'

export default defineComponent({
  setup() {
    const router = useRouter()
    const route = useRoute()
    const id = route.value.params.id

    const { getNotification, editNotification } = useNotificationStore()

    const formData = reactive<NotificationResponse>({
      id,
      createdBy: '',
      creatorName: '',
      updatedBy: '',
      title: '',
      body: '',
      targets: [],
      public: false,
      publishedAt: -1,
      createdAt: -1,
      updatedAt: -1,
    })

    const timeData = reactive<NotificationTime>({
      publishedDate: '',
      publishedTime: '',
    })

    const { fetchState } = useFetch(async () => {
      const notification = await getNotification(id)
      formData.title = notification.title
      formData.body = notification.body
      formData.targets = notification.targets
      formData.public = notification.public
      formData.publishedAt = notification.publishedAt
      timeData.publishedDate = dayjs
        .unix(notification.publishedAt)
        .format('YYYY-MM-DD')
      timeData.publishedTime = dayjs
        .unix(notification.publishedAt)
        .format('HH:mm')
    })

    const handleSubmit = async () => {
      try {
        await editNotification(id, formData)
        router.push('/notifications')
      } catch (error) {
        console.log(error)
      }
    }

    return {
      fetchState,
      formData,
      timeData,
      handleSubmit,
    }
  },
})
</script>
