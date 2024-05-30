<script lang="ts" setup>
import { useCommonStore } from '~/store'

const commonStore = useCommonStore()

const snackbars = computed(() => {
  return commonStore.snackbars.filter(item => item.isOpen)
})

const calcStyle = (i: number) => {
  if (i > 0) {
    return `top: ${60 * i}px;`
  }
}
</script>

<template>
  <v-app>
    <v-snackbar
      v-for="(snackbar, i) in snackbars"
      :key="i"
      v-model="snackbar.isOpen"
      :color="snackbar.color"
      location="top"
      variant="elevated"
      :timeout="snackbar.timeout"
      :style="calcStyle(i)"
    >
      {{ snackbar.message }}
      <template #actions>
        <v-btn
          variant="text"
          color="white"
          @click="commonStore.hideSnackbar(i)"
        >
          閉じる
        </v-btn>
      </template>
    </v-snackbar>

    <v-main class="d-flex justify-center align-center bg-color">
      <v-container>
        <slot />
      </v-container>
    </v-main>
  </v-app>
</template>

<style lang="scss" scoped>
.bg-color {
  background-color: #eef5f9;
}
</style>
