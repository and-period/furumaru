<script lang="ts" setup>
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
    <v-btn color="primary" outlined class="mx-2" @click="handleClickTestButton">
      テスト配信
    </v-btn>
    <v-btn color="primary" outlined @click="handleClickAddButton">
      <v-icon>mdi-plus</v-icon>
      ライブマルシェ登録
    </v-btn>
  </div>
</template>
