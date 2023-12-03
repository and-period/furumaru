<script lang="ts" setup>
import { useAlert } from '~/lib/hooks'
import { useAuthStore, useCommonStore } from '~/store'
import type { ForgotAuthPasswordRequest, ResetAuthPasswordRequest } from '~/types/api'

definePageMeta({
  layout: 'auth'
})

const router = useRouter()
const commonStore = useCommonStore()
const authStore = useAuthStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const loading = ref<boolean>(false)
const sentEmail = ref<boolean>(false)
const formData = ref<ResetAuthPasswordRequest>({
  email: '',
  verifyCode: '',
  password: '',
  passwordConfirmation: ''
})

const handleClickCancel = (): void => {
  router.back()
}

const handleSendEmail = async (): Promise<void> => {
  try {
    loading.value = true
    const req: ForgotAuthPasswordRequest = {
      email: formData.value.email
    }
    await authStore.forgotPassword(req)
    commonStore.addSnackbar({
      message: 'パスワードリセット用のメールをしました。',
      color: 'info'
    })
    sentEmail.value = true
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  } finally {
    loading.value = false
  }
}

const handleResetPassword = async (): Promise<void> => {
  try {
    loading.value = true
    await authStore.resetPassword(formData.value)
    commonStore.addSnackbar({
      message: 'パスワードをリセットしました。',
      color: 'info'
    })
    router.push('/signin')
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <templates-auth-password-reset
    v-if="sentEmail"
    v-model:form-data="formData"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    @click:cancel="handleClickCancel"
    @submit="handleResetPassword"
  />
  <templates-auth-password-forgot
    v-else
    v-model:form-data="formData"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    @click:cancel="handleClickCancel"
    @submit="handleSendEmail"
  />
</template>
