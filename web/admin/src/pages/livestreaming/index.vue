<script lang="ts" setup>
import { mdiPlus } from '@mdi/js'

const router = useRouter()

const handleClickAddButton = () => {
  router.push('/livestreaming/add')
}

const handleClickTestButton = () => {
  router.push('/livestreaming/tests')
}
</script>

<template>
  <div class="d-flex mb-4">
    <v-spacer />
    <v-btn color="primary" variant="outlined" class="mx-2" @click="handleClickTestButton">
      テスト配信
    </v-btn>
    <v-btn color="primary" variant="outlined" @click="handleClickAddButton">
      <v-icon :icon="mdiPlus" />
      ライブマルシェ登録
    </v-btn>
  </div>
</template>
