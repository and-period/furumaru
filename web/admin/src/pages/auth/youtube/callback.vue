<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { unix } from 'dayjs'
import { useAlert } from '~/lib/hooks'
import type {
  CreateYoutubeBroadcastRequest,
  CallbackAuthYoutubeBroadcastRequest,
} from '~/types/api/v1'
import { useBroadcastStore } from '~/store'

const route = useRoute()
const router = useRouter()
const broadcastStore = useBroadcastStore()
const { show } = useAlert('error')

const state = route.query.state as string
const code = route.query.code as string
const { guestBroadcast } = storeToRefs(broadcastStore)

const items = [
  { title: '公開', value: true },
  { title: '限定公開', value: false },
]

const loading = ref<boolean>(false)
const formData = reactive<CreateYoutubeBroadcastRequest>({
  title: '',
  description: '',
  _public: false,
})

const parseTime = (unixtime: number): string => {
  return unix(unixtime).format('YYYY/MM/DD HH:mm')
}

const fetchBroadcast = async () => {
  try {
    loading.value = true
    await broadcastStore.getGuestBroadcast()
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

const connectYoutube = async () => {
  const req: CallbackAuthYoutubeBroadcastRequest = {
    state,
    authCode: code,
  }
  try {
    loading.value = true
    await broadcastStore.connectYouTube(req)
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

const handleSubmit = async () => {
  try {
    loading.value = true
    await broadcastStore.createYoutubeLive(formData)
    router.push('/auth/youtube/complete')
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
  if (state && code) {
    await connectYoutube()
  }
  else {
    await fetchBroadcast()
  }
  formData.title = guestBroadcast.value.title
  formData.description = guestBroadcast.value.description
}
catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <div>
    <h1 class="mb-4">
      YouTube連携情報
    </h1>

    <v-card>
      <v-card-title> ライブ配信情報 </v-card-title>
      <v-card-text>
        <v-table v-if="guestBroadcast">
          <tbody>
            <tr>
              <td>タイトル</td>
              <td>{{ guestBroadcast.title }}</td>
            </tr>
            <tr>
              <td>説明</td>
              <td>{{ guestBroadcast.description }}</td>
            </tr>
            <tr>
              <td>配信時間</td>
              <td>
                {{ parseTime(guestBroadcast.startAt) }}〜{{
                  parseTime(guestBroadcast.endAt)
                }}
              </td>
            </tr>
            <tr>
              <td>配信担当者（マルシェ名）</td>
              <td>{{ guestBroadcast.coordinatorMarche }}</td>
            </tr>
            <tr>
              <td>配信担当者（コーディネータ名）</td>
              <td>{{ guestBroadcast.coordinatorName }}</td>
            </tr>
          </tbody>
        </v-table>
      </v-card-text>
    </v-card>

    <v-card class="mt-4">
      <v-card-title> YouTube 配信情報</v-card-title>
      <v-card-text>
        <form @submit.prevent="handleSubmit">
          <v-text-field
            v-model="formData.title"
            label="タイトル"
            variant="outlined"
            primary
            required
          />
          <v-textarea
            v-model="formData.description"
            label="説明"
            variant="outlined"
          />
          <v-select
            v-model.boolean="formData.public"
            variant="outlined"
            :items="items"
            label="公開設定"
            chips
          />
          <v-btn
            :loading="loading"
            type="submit"
            block
            variant="outlined"
            color="primary"
          >
            送信
          </v-btn>
        </form>
      </v-card-text>
    </v-card>
  </div>
</template>
