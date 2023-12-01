<script lang="ts" setup>
import { useAuthStore } from '~/store/auth'
import { I18n } from '~/types/locales'

definePageMeta({
  layout: 'auth',
})

const authStore = useAuthStore()
const { verifyAuth } = authStore

const route = useRoute()

const i18n = useI18n()

const t = (str: keyof I18n['auth']['verify']) => {
  return i18n.t(`auth.verify.${str}`)
}

const id = computed<string>(() => {
  const id = route.query.id
  if (id) {
    return id as string
  } else {
    return ''
  }
})

const code = ref<string>('')

const handleSubmit = async () => {
  await verifyAuth({
    verifyCode: code.value,
    id: '',
  })
}
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
