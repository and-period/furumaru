<template>
  <v-row justify="center" align="center">
    <v-col cols="12" sm="8" md="6">
      <v-card>
        <v-card-title class="headline">
          Welcome to the Vuetify + Nuxt.js template
        </v-card-title>
        <v-card-text>
          <div class="cube-area">
            <the-cube v-show="!cubeIsOpen" @click="handleCubeClick" />
            <transition name="fade">
              <the-cube-open v-show="cubeIsOpen" @click="handleCubeClick">
                <v-avatar size="28">
                  <img
                    src="https://4.bp.blogspot.com/-PYqNTVOP9hs/UgSMR6zIbdI/AAAAAAAAXAk/tDJMCfemfwk/s800/fruit_orange.png"
                    alt="sample"
                  />
                </v-avatar>
              </the-cube-open>
            </transition>
          </div>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn color="primary" nuxt to="/inspire"> Continue </v-btn>
        </v-card-actions>
      </v-card>
    </v-col>
  </v-row>
</template>

<script lang="ts">
import { defineComponent, ref } from '@nuxtjs/composition-api'

import TheCube from '~/components/atoms/TheCube.vue'

export default defineComponent({
  name: 'IndexPage',
  components: { TheCube },
  setup() {
    const cubeIsOpen = ref<boolean>(false)

    const handleCubeClick = () => {
      console.log('page')
      cubeIsOpen.value = !cubeIsOpen.value
    }

    return {
      handleCubeClick,
      cubeIsOpen,
    }
  },
})
</script>

<style lang="scss" scoped>
.cube-area {
  height: 160px;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s;
}
.fade-enter, .fade-leave-to /* .fade-leave-active below version 2.1.8 */ {
  opacity: 0;
}
</style>
