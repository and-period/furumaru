<script lang="ts" setup>
import { useAlert } from '~/lib/hooks'
import { useAuthStore, useCommonStore } from '~/store'
import type { UpdateAuthEmailRequest, VerifyAuthEmailRequest } from '~/types/api/v1'

const route = useRoute()
const router = useRouter()
const commonStore = useCommonStore()
const authStore = useAuthStore()
const { alertType, isShow, alertText, show } = useAlert('error')

// todo: emailの受け渡しを変更する
const email = route.params.email as string

const loading = ref<boolean>(false)
const formData = ref<VerifyAuthEmailRequest>({
  verifyCode: '',
})

watch(formData.value, (): void => {
  if (formData.value.verifyCode.length !== 6) {
    return
  }
  handleSubmit()
})

const handleClickResendEmail = async (): Promise<void> => {
  try {
    loading.value = true
    const req: UpdateAuthEmailRequest = { email }
    await authStore.updateEmail(req)
    commonStore.addSnackbar({
      message: '認証コードを送信しました。',
      color: 'info',
    })
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    loading.value = false
  }
}

const handleSubmit = async (): Promise<void> => {
  try {
    loading.value = true
    await authStore.verifyEmail(formData.value)
    commonStore.addSnackbar({
      message: 'メールアドレスが変更されました。',
      color: 'info',
    })
    router.push('/')
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
  finally {
    formData.value.verifyCode = ''
    loading.value = false
  }
}
</script>

<template>
  <templates-auth-email-verify
    v-model:form-data="formData"
    :loading="loading"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :email="email"
    @submit="handleSubmit"
    @click:resend-email="handleClickResendEmail"
  />
</template>
