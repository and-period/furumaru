<template>
  <div>
    <v-card outlined class="mx-auto" :max-width="isMobile ? 360 : 440">
      <v-card-text class="pa-md-12 pa-sm-4">
        <form>
          <v-text-field
            type="tel"
            :label="$t('auth.signUp.tel')"
            outlined
            dense
            required
          />
          <v-text-field
            type="email"
            :label="$t('auth.signUp.email')"
            outlined
            dense
            required
          />
          <v-text-field
            type="password"
            :label="$t('auth.signUp.password')"
            outlined
            dense
            required
          />
          <v-text-field
            type="password"
            :label="$t('auth.signUp.passwordConfirm')"
            outlined
            dense
            required
          />
          <v-btn
            class="rounded-xl"
            color="accent"
            block
            type="submit"
            elevation="0"
            :height="isMobile ? 40 : 64"
            :x-large="!isMobile"
            >{{ $t('auth.signUp.signUp') }}</v-btn
          >
        </form>
      </v-card-text>
    </v-card>
    <div class="text-center mt-10">
      <nuxt-link to="/">{{ $t('auth.signUp.alreadyHas') }}</nuxt-link>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, useContext, computed } from '@nuxtjs/composition-api'

export default defineComponent({
  layout: 'auth',
  setup() {
    const { $vuetify } = useContext()
    const isMobile = computed(() =>
      ['sm', 'xs'].includes($vuetify.breakpoint.name)
    )

    return {
      isMobile,
    }
  },
})
</script>
