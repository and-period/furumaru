<script lang="ts" setup>
import { useAlert } from '~/lib/hooks'
import { useAuthStore } from '~/store'
import { UpdateAuthEmailRequest, VerifyAuthEmailRequest } from '~/types/api'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const { alertType, isShow, alertText, show } = useAlert('error')

// todo: emailの受け渡しを変更する
const email = route.params.email as string

const loading = ref<boolean>(false)
const formData = ref<VerifyAuthEmailRequest>({
  verifyCode: route.params.verifyCode as string
})

const handleClickResendEmail = async (): Promise<void> => {
  try {
    loading.value = true
    const req: UpdateAuthEmailRequest = { email }
    await authStore.updateEmail(req)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  } finally {
    loading.value = false
  }
}

const handleSubmit = async (): Promise<void> => {
  try {
    loading.value = true
    await authStore.verifyEmail(formData.value)
    router.push('/')
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
