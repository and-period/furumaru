<template>
  <div>
    <v-card outlined>
      <v-card-text class="pa-12">
        <form>
          <v-text-field
            type="tel"
            :label="$t('auth.signUp.tel')"
            outlined
            dense
          />
          <v-text-field
            type="email"
            :label="$t('auth.signUp.email')"
            outlined
            dense
          />
          <v-text-field
            type="password"
            :label="$t('auth.signUp.password')"
            outlined
            dense
          />
          <v-text-field
            type="password"
            :label="$t('auth.signUp.passwordConfirm')"
            outlined
            dense
          />
          <v-btn block type="submit">{{ $t('auth.signUp.signUp') }}</v-btn>
        </form>
      </v-card-text>
    </v-card>
    <div class="text-center mt-10">
      <nuxt-link to="/">{{ $t('auth.signUp.alreadyHas') }}</nuxt-link>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from '@nuxtjs/composition-api'

export default defineComponent({
  setup() {},
})
</script>
