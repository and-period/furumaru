<script lang="ts" setup>
import { useAlert } from '~/lib/hooks'
import type { CreateYoutubeBroadcastRequest } from '~/types/api'
import { useBroadcastStore } from '~/store'

const route = useRoute()
const router = useRouter()
const broadcastStore = useBroadcastStore()
const { show } = useAlert('error')

const state = route.query.state as string
const code = route.query.code as string

const formData = reactive<CreateYoutubeBroadcastRequest>({
  state: '',
  authCode: '',
  public: false
})

const handleSubmit = async () => {
  try {
    formData.state = state
    formData.authCode = code
    await broadcastStore.connectYouTube(formData)

    router.push('/auth/youtube/complete')
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}
</script>

<template>
  <div>
    <h1>Youtube Callback</h1>
    <v-text-field v-model="formData.public" label="公開範囲" />
    <v-btn @click="handleSubmit">Submit</v-btn>
  </div>
</template>
