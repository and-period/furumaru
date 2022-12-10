<template>
  <v-app>
    <v-main class="d-flex justify-center align-center bg-color">
      <v-container>
        <Nuxt />
      </v-container>
    </v-main>
  </v-app>
</template>

<script lang="ts">
import { defineComponent } from '@vue/composition-api'

export default defineComponent({
  setup() {},
})
</script>

<style lang="scss" scoped>
.bg-color {
  background-color: #eef5f9;
}
</style>
