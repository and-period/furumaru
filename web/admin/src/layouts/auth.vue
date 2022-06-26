<template>
  <v-app>
    <v-app-bar flat app color="primary" dark>
      <v-toolbar-title>Online Marche</v-toolbar-title>
    </v-app-bar>
    <v-main class="bg-color">
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
