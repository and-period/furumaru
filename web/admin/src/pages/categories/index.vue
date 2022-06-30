<template>
  <div>
    <p class="text-h6">カテゴリー・品目設定</p>
    <v-tabs v-model="selector" grow color="dark">
      <v-tabs-slider color="accent"></v-tabs-slider>
      <v-tab v-for="item in items" :key="item.value" :href="`tab-${item.value}`">
        {{ item }}
      </v-tab>
    </v-tabs>
  </div>
</template>

<script lang="ts">
import { defineComponent, useRouter } from '@nuxtjs/composition-api'

export default defineComponent({
  setup() {
    const router = useRouter()

    const items = ['カテゴリー', '品目']

    const handleClickAddButton = () => {
      router.push('/categories/add')
    }

    return {
      items,
      handleClickAddButton,
    }
  },
})
</script>
