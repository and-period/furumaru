<template>
  <div>
    <v-card
      elevation="0"
      class="mx-auto bg-transparent"
      :max-width="isMobile ? 360 : 440"
    >
      <v-card-text class="pa-md-12 pa-sm-4">
        <form>
          <stack :space="24" class="text-center">
            <p class="text-h6">{{ $t('auth.verify.message') }}</p>
            <the-pincode-input value="" :length="6" class="mt-6" />
            <the-submit-button :is-mobile="isMobile">{{
              $t('auth.verify.btnText')
            }}</the-submit-button>
          </stack>
        </form>
      </v-card-text>
    </v-card>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, useContext } from '@nuxtjs/composition-api'
import Stack from '~/components/atoms/Stack.vue'
import ThePincodeInput from '~/components/molecules/ThePincodeInput.vue'
import TheSubmitButton from '~/components/molecules/TheSubmitButton.vue'

export default defineComponent({
  components: {
    ThePincodeInput,
    TheSubmitButton,
    Stack,
  },
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
