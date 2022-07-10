<template>
  <div>
    <v-card-title>生産者管理</v-card-title>
    <div class="d-flex">
      <v-spacer />
      <v-btn outlined @click="handleClickAddButton">
        <v-icon left>mdi-plus</v-icon>
        追加
      </v-btn>
    </div>
    {{ producers }}
  </div>
</template>

<script lang="ts">
import { defineComponent, useFetch, useRouter } from '@nuxtjs/composition-api'

import { useProducerStore } from '~/store/producer'

export default defineComponent({
  setup() {
    const router = useRouter()
    const { producers, fetchProducers } = useProducerStore()

    const handleClickAddButton = () => {
      router.push('/producers/add')
    }

    useFetch(async () => {
      try {
        await fetchProducers()
      } catch (err) {
        console.log(err)
      }
    })

    return {
      handleClickAddButton,
      producers,
    }
  },
})
</script>
