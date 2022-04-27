<template>
  <div class="text-center">
    <the-marche-logo class="mb-6" />
    <v-form>
      <v-row justify="center">
        <v-col cols="12" sm="12" md="6">
          <v-text-field
            label="Email"
            placeholder="Email"
            outlined
            required
            type="email"
            dense
          ></v-text-field>
          <v-text-field
            label="Password"
            placeholder="Password"
            outlined
            required
            dense
          ></v-text-field>
        </v-col>
      </v-row>
      <div>
        <nuxt-link to="/"> パスワードを忘れた場合 </nuxt-link>
      </div>
      <v-row justify="center">
        <v-col cols="12" sm="12" md="6">
          <the-submit-button>ログイン</the-submit-button>
          <the-stack :space="32" class="my-10">
            <v-btn
              outlined
              block
              class="base rounded-xl google-btn"
              :height="isMobile ? 40 : 64"
              :x-large="!isMobile"
            >
              <v-icon left>fa-brands fa-google</v-icon>
              Log In with Google
            </v-btn>
            <v-btn
              block
              outlined
              class="facebook white--text rounded-xl"
              :height="isMobile ? 40 : 64"
              :x-large="!isMobile"
            >
              <v-icon left>fa-brands fa-facebook</v-icon>
              Log In with Facebook
            </v-btn>
            <v-btn
              block
              outlined
              class="line white--text rounded-xl"
              :height="isMobile ? 40 : 64"
              :x-large="!isMobile"
            >
              <v-icon left>fa-brands fa-line</v-icon>
              Log In with LINE
            </v-btn>
          </the-stack>
        </v-col>
      </v-row>
    </v-form>
    <p>
      アカウントをお持ちではないですか？
      <nuxt-link to="/"> 登録する </nuxt-link>
    </p>
  </div>
</template>

<script lang="ts">
import { defineComponent } from '@nuxtjs/composition-api'
import TheMarcheLogo from '~/components/atoms/TheMarcheLogo.vue'
import { useIsMobile } from '~/lib/hooks'

export default defineComponent({
  components: { TheMarcheLogo },
  layout: 'auth',
  setup() {
    const { isMobile } = useIsMobile()

    return {
      isMobile,
    }
  },
})
</script>

<style lang="scss" scoped>
.google-btn {
  border-color: rgba(0, 0, 0, 0.3) !important;
  color: rgba(0, 0, 0, 0.54);
}
</style>
