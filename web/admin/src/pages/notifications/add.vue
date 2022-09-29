<template>
  <the-notification-create-form-page
    :form-data="formData"
    @submit="handleSubmit"
    />
</template>

<script lang="ts">
import { useRouter, defineComponent, reactive } from '@nuxtjs/composition-api'
import dayjs from 'dayjs'

import { useNotificationStore } from '~/store/notification'
import { CreateNotificationRequest } from '~/types/api'

export default defineComponent({
  setup() {
    const router = useRouter()
    const notificationStore = useNotificationStore()

    const formData = reactive<CreateNotificationRequest>({
      title: '',
            body: '',
            targets: [0],
            public: false,
            publishedAt: dayjs().unix(),
    })

    const handleSubmit = async () => {
      try {
        await notificationStore.createNotification(formData)
        router.push('/notifications')
      } catch (error) {
        console.log(error)
      }
    }

    return {
      formData,
      handleSubmit,
    }
  },
})
</script>


