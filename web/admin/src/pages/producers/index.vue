<template>
  <div>
    <p class="text-h6">生産者管理</p>
    <div class="d-flex">
      <v-spacer />
      <v-btn outlined @click="handleClickAddButton">
        <v-icon left>mdi-plus</v-icon>
        追加
      </v-btn>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, useRouter } from '@nuxtjs/composition-api'

export default defineComponent({
  setup() {
    const router = useRouter()

    const handleClickAddButton = () => {
      router.push('/producers/add')
    }
    return {
      handleClickAddButton,
    }
  },
})
</script>
