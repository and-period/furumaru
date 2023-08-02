<script lang="ts" setup>
import Hls from 'hls.js'
import { Broadcast, BroadcastStatus } from '~/types/api'

const props = defineProps({
  selectedTabItem: {
    type: String,
    default: 'schedule'
  },
  broadcast: {
    type: Object as PropType<Broadcast>,
    default: (): Broadcast => ({
      id: '',
      scheduleId: '',
      status: BroadcastStatus.UNKNOWN,
      inputUrl: '',
      outputUrl: '',
      createdAt: 0,
      updatedAt: 0
    })
  }
})

const emit = defineEmits<{
  (e: 'update:broadcast', broadcast: Broadcast): void
}>()

const statuses = [
  { title: 'リソース未作成', value: BroadcastStatus.DISABLED },
  { title: 'リソース作成中', value: BroadcastStatus.WAITING },
  { title: '配信停止中', value: BroadcastStatus.IDLE },
  { title: '配信中', value: BroadcastStatus.ACTIVE },
  { title: '不明', value: BroadcastStatus.UNKNOWN }
]

const videoRef = ref<HTMLVideoElement>()

const broadcastValue = computed({
  get: (): Broadcast => props.broadcast,
  set: (broadcast: Broadcast): void => emit('update:broadcast', broadcast)
})

watch((): string => props.selectedTabItem, (): void => {
  if (props.selectedTabItem !== 'streaming') {
    return
  }
  handleClickVideo()
})

const handleClickVideo = (): void => {
  if (!videoRef.value) {
    return
  }

  const video = videoRef.value

  if (Hls.isSupported()) {
    const hls = new Hls({ enableWorker: false })
    hls.loadSource(props.broadcast.outputUrl)
    hls.attachMedia(video)
    video.play()
  } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
    video.src = props.broadcast.outputUrl
    video.play()
  }
}
</script>

<template>
  <v-row>
    <v-col sm="12" md="12" lg="8">
      <v-card>
        <v-card-text>
          <v-container>
            <video id="video" ref="videoRef" controls />
            <v-btn @click="handleClickVideo">映像の更新</v-btn>
          </v-container>
        </v-card-text>
      </v-card>
    </v-col>

    <v-col sm="12" md="12" lg="4">
      <v-card>
        <v-card-text>
          <v-select
            v-model="broadcastValue.status"
            label="配信状況"
            :items="statuses"
            item-title="title"
            item-value="value"
            variant="plain"
            readonly
          />
          <v-text-field
            v-model="broadcastValue.inputUrl"
            label="配信エンドポイント：入力側"
            readonly
          />
          <v-text-field
            v-model="broadcastValue.outputUrl"
            label="配信エンドポイント：出力側"
            readonly
          />
        </v-card-text>
      </v-card>
    </v-col>
  </v-row>
</template>

<style scoped>
#video {
  margin-bottom: 1.5em;
  width: 100%;
}
</style>
