<script lang="ts" setup>
import { useAlert } from '~/lib/hooks'
import { useAuthStore } from '~/store'
import { UpdateAuthEmailRequest, VerifyAuthEmailRequest } from '~/types/api'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const { alertType, isShow, alertText, show } = useAlert('error')

// todo: emailの受け渡しを変更する
const email = route.params.email as string

const formData = reactive<VerifyAuthEmailRequest>({
  verifyCode: route.params.verifyCode as string
})

const handleClickResendEmail = async (): Promise<void> => {
  try {
    const req: UpdateAuthEmailRequest = { email }
    await authStore.emailUpdate(req)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}

const handleSubmit = async (): Promise<void> => {
  try {
    await authStore.codeVerify(formData)
    router.push('/')
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}
</script>

<template>
  <templates-auth-verify-email
    v-model:form-data="formData"
    v-model:is-alert="isShow"
    :email="email"
    :alert-type="alertType"
    :alert-text="alertText"
    @submit="handleSubmit"
    @click:resend-email="handleClickResendEmail"
  />
</template>
