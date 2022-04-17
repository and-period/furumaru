<template>
  <div>
    <v-card outlined class="mx-auto" :max-width="isMobile ? 360 : 440">
      <v-card-text class="pa-md-12 pa-sm-4">
        <form @submit.prevent="handleSubmit">
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
          <the-submit-button :is-mobile="isMobile">
            {{ $t('auth.signUp.signUp') }}
          </the-submit-button>
        </form>
      </v-card-text>
    </v-card>
    <div class="text-center mt-10">
      <nuxt-link to="/">{{ $t('auth.signUp.alreadyHas') }}</nuxt-link>
    </div>
  </div>
</template>

<script lang="ts">
import {
  defineComponent,
  useContext,
  computed,
  useRouter,
} from '@nuxtjs/composition-api'
import TheSubmitButton from '~/components/molecules/TheSubmitButton.vue'

export default defineComponent({
  components: {
    TheSubmitButton,
  },
  layout: 'auth',
  setup() {
    const { $vuetify } = useContext()
    const router = useRouter()
    const isMobile = computed(() =>
      ['sm', 'xs'].includes($vuetify.breakpoint.name)
    )

    const handleSubmit = () => {
      router.push('/verify')
    }

    return {
      isMobile,
      handleSubmit,
    }
  },
})
</script>
