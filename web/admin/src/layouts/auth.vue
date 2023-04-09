<style lang="scss" scoped>
.bg-color {
  background-color: #eef5f9;
}
</style>

<template>
  <v-app>
    <v-main class="d-flex justify-center align-center bg-color">
      <v-container>
        <slot />
      </v-container>
    </v-main>
  </v-app>
</template>
