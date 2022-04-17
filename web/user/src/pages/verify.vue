<template>
  <div>
    <v-card
      elevation="0"
      class="mx-auto bg-transparent"
      :max-width="isMobile ? 360 : 440"
    >
      <v-card-text class="pa-md-12 pa-sm-4 text-center">
        <p>{{ $t('auth.verify.message') }}</p>
        <form>
          <the-pincode-input value="" :length="6" />
          <the-submit-button :is-mobile="isMobile">{{
            $t('auth.verify.btnText')
          }}</the-submit-button>
        </form>
      </v-card-text>
    </v-card>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, useContext } from '@nuxtjs/composition-api'
import ThePincodeInput from '~/components/molecules/ThePincodeInput.vue'
import TheSubmitButton from '~/components/molecules/TheSubmitButton.vue'

export default defineComponent({
  components: {
    ThePincodeInput,
    TheSubmitButton,
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
