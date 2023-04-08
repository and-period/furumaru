<script setup>
import HLS from 'hls.js'

const router = useRouter()
const hls = new HLS()

const formData = reactive({
  playbackUrl: '',
})

const handleClickStreaming = () => {
  router.push('/livestreaming')
}

const startWatching = () => {
  if (formData.playbackUrl === '') {
    alert('playback url is required!')
    return
  }

  try {
    hls.loadSource(formData.playbackUrl)
    hls.attachMedia(document.getElementById('video'))
  } catch (err) {
    console.error(err)
  }
}
</script>

<style scoped>
#video {
  margin-bottom: 1.5em;
  width: 100%;
}
</style>

<template>
  <v-card>
    <v-card-title>視聴テスト用モック</v-card-title>
    <v-card-subtitle>
      <p>Playback URLはAWSコンソールより取得すること。</p>
    </v-card-subtitle>

    <v-container>
      <video id="video" controls></video>
    </v-container>

    <v-card-text>
      <v-text-field v-model="formData.playbackUrl" placeholder="Playback URL" />
    </v-card-text>

    <v-card-actions>
      <v-btn @click="startWatching">Start Watching</v-btn>
      <v-btn @click="handleClickStreaming">Live Streaming</v-btn>
    </v-card-actions>
  </v-card>
</template>
