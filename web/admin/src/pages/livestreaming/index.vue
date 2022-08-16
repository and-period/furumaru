<template>
  <v-card>
    <v-card-title>配信テスト用モック</v-card-title>
    <v-card-subtitle>
      <p>Ingest Endpoint, Stream KeyはAWSコンソールより取得すること。</p>
    </v-card-subtitle>

    <v-container>
      <canvas id="preview"></canvas>
    </v-container>

    <v-card-text>
      <v-select
        v-model="formData.videoDevice"
        :items="videoDevices"
        placeholder="Web Camera"
        @change="handleSelectVideo"
      />
      <v-select
        v-model="formData.audioDevice"
        :items="audioDevices"
        placeholder="Microphone"
        @change="handleSelectAudio"
      />
      <v-select
        v-model="formData.streamConfig"
        :items="channelConfigs"
        placeholder="Channel Config"
        @change="handleSelectChannelConfig"
      />
      <v-text-field
        v-model="formData.ingestEndpoint"
        placeholder="Ingest Endpoint"
      />
      <v-text-field v-model="formData.streamKey" placeholder="Stream Key" />
    </v-card-text>

    <v-card-actions>
      <v-btn @click="startBroadcast">Start Broadcast</v-btn>
      <v-btn @click="stopBroadcast">Stop Broadcast</v-btn>
      <v-btn @click="handleClickViewing">Live Viewing</v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>
import { defineComponent, onMounted, reactive, useRouter } from '@nuxtjs/composition-api'
import IVSBroadcastClient, {
  STANDARD_LANDSCAPE,
  STANDARD_PORTRAIT,
  LOG_LEVEL,
} from 'amazon-ivs-web-broadcast'

export default defineComponent({
  setup() {
    const router = useRouter()

    let client = IVSBroadcastClient.create({
      streamConfig: STANDARD_LANDSCAPE,
      ingestEndpoint: '',
      logLevel: LOG_LEVEL.DEBUG,
    })

    const formData = reactive({
      streamConfig: STANDARD_LANDSCAPE,
      ingestEndpoint: '',
      streamKey: '',
      videoDevice: undefined,
      audioDevice: undefined,
    })

    const channelConfigs = [
      { text: 'Standard: Landscape', value: STANDARD_LANDSCAPE },
      { text: 'Standard: Portrait', value: STANDARD_PORTRAIT },
    ]
    const videoDevices = reactive([])
    const audioDevices = reactive([])

    const handlePermissions = async () => {
      let permissions = {
        audio: false,
        video: false,
      }
      try {
        const stream = await navigator.mediaDevices.getUserMedia({
          video: true,
          audio: true,
        })
        for (const track of stream.getTracks()) {
          track.stop()
        }
        permissions = { video: true, audio: true }
      } catch (err) {
        permissions = { video: false, audio: false }
        console.error(err.message)
      }
      if (!permissions.video) {
        console.error('failed to get video permissions.')
      } else if (!permissions.audio) {
        console.error('failed to get audio permissions.')
      }
    }

    const getDevices = async () => {
      const devices = await navigator.mediaDevices.enumerateDevices()
      const video = devices.filter((d) => d.kind === 'videoinput')
      const audio = devices.filter((d) => d.kind === 'audioinput')

      video.forEach((val) => {
        videoDevices.push({ text: val.label, value: val.deviceId })
      })
      audio.forEach((val) => {
        audioDevices.push({ text: val.label, value: val.deviceId })
      })

      return { videoDevices: video, audioDevices: audio }
    }

    const attachPreview = async () => {
      await client.attachPreview(document.getElementById('preview'))
    }

    onMounted(async () => {
      await handlePermissions()
      await getDevices()
      await attachPreview()
    })

    const recreateClient = async () => {
      const config = {
        streamConfig: formData.streamConfig,
        ingestEndpoint: formData.ingestEndpoint,
        logLevel: LOG_LEVEL.DEBUG,
      }
      client = IVSBroadcastClient.create(config)
      client.on(
        IVSBroadcastClient.BroadcastClientEvents.ACTIVE_STATE_CHANGE,
        (active) => {
          onActiveStateChange(active)
        }
      )

      await handleSelectVideo(formData.videoDevice)
      await handleSelectAudio(formData.audioDevice)
      await attachPreview()
    }

    const handleSelectChannelConfig = async () => {
      await recreateClient()
    }

    const handleSelectVideo = async (deviceId) => {
      if (client.getVideoInputDevice('camera')) {
        client.removeVideoInputDevice('camera')
      }
      if (!formData.streamConfig || deviceId === '') {
        return
      }
      const { width, height } = formData.streamConfig.maxResolution
      const cameraStream = await navigator.mediaDevices.getUserMedia({
        video: { deviceId, width: { max: width }, height: { max: height } },
      })
      await client.addVideoInputDevice(cameraStream, 'camera', { index: 0 })
    }

    const handleSelectAudio = async (deviceId) => {
      if (client.getAudioInputDevice('microphone')) {
        client.removeAudioInputDevice('microphone')
      }
      if (!formData.streamConfig || deviceId === '') {
        return
      }
      const microphoneStream = await navigator.mediaDevices.getUserMedia({
        audio: { deviceId },
      })
      await client.addAudioInputDevice(microphoneStream, 'microphone')
    }

    const handleClickViewing = () => {
      router.push('/livestreaming/view')
    }

    const startBroadcast = async () => {
      client.config.ingestEndpoint = formData.ingestEndpoint
      await client
        .startBroadcast(formData.streamKey)
        .then((res) => {
          console.log('success to start broardcast', res)
        })
        .catch((err) => {
          alert(err)
        })
    }

    const stopBroadcast = () => {
      client.stopBroadcast()
    }

    return {
      formData,
      channelConfigs,
      videoDevices,
      audioDevices,
      handleSelectChannelConfig,
      handleSelectVideo,
      handleSelectAudio,
      handleClickViewing,
      startBroadcast,
      stopBroadcast,
    }
  },
})
</script>

<style scoped>
#preview {
  margin-bottom: 1.5rem;
  background: green;
  width: 100%;
}
</style>
