<script setup lang="ts">
import type { I18n } from '~/types/locales'
definePageMeta({
  layout: 'auth',
})

const i18n = useI18n()
const t = (str: keyof I18n['auth']['register']) => {
  return i18n.t(`auth.register.${str}`)
}

const name = ref<string>('')
const username = ref<string>('')
</script>

<template>
  <the-user-register-page
    v-model:name="name"
    v-model:username="username"
    :page-name="t('pageName')"
    :name-label="t('name')"
    name-error-message=""
    :name-placeholder="t('name')"
    :username-label="t('username')"
    :username-placeholder="t('username')"
    username-error-message=""
    :button-text="t('btnText')"
    :caution-text="
      i18n.t('auth.register.cautionText', {
        privacyPolicyLinkText: `<a class='underline' href='#''>${t(
          'privacyPolicyLinkText',
        )}</a>`,
        termsOfService: `<a class='underline' href='#''>${t(
          'termsOfServiceLink',
        )}</a>`,
      })
    "
  />
</template>
