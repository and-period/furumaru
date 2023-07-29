<script lang="ts" setup>
import { I18n } from '~/types/locales'

definePageMeta({
  layout: 'auth',
})

const i18n = useI18n()

const t = (str: keyof I18n['auth']['verify']) => {
  return i18n.t(`auth.verify.${str}`)
}

const handleSubmit = () => {
  console.log('submit not implemented')
}

const code = ref<string>('')
</script>

<template>
  <the-verify-code-page
    v-model:code="code"
    :page-name="t('pageName')"
    :button-text="t('btnText')"
    :message="t('message')"
    @submit="handleSubmit"
  />
</template>
