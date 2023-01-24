<template>
  <div>
    <atoms-the-marche-logo class="text-center mb-6" />
    <v-card elevation="0" class="mx-auto bg-transparent">
      <v-card-text class="pa-md-12 pa-sm-4">
        <form @submit.prevent="handleSubmit">
          <atoms-the-stack :space="24" class="text-center">
            <p class="text-h6">
              {{ t('message') }}
            </p>
            <molecules-the-pincode-input value="" :length="6" />
            <molecules-the-submit-button>
              {{ t('btnText') }}
            </molecules-the-submit-button>
          </atoms-the-stack>
        </form>
      </v-card-text>
    </v-card>
  </div>
</template>

<script lang="ts" setup>

import { I18n } from '~/types/locales'

definePageMeta({
  layout: 'auth'
})

const { $i18n } = useNuxtApp()

const t = (str: keyof I18n['auth']['verify']) => {
  return $i18n.t(`auth.verify.${str}`)
}

const handleSubmit = () => {
  console.log('submit not implemented')
}
</script>
