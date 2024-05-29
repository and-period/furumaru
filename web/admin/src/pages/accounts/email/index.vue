<script lang="ts" setup>
import { useAlert } from '~/lib/hooks'
import { useAuthStore, useCommonStore } from '~/store'
import type { UpdateAuthEmailRequest } from '~/types/api'

const router = useRouter()
const commonStore = useCommonStore()
const authStore = useAuthStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const loading = ref<boolean>(false)
const formData = ref<UpdateAuthEmailRequest>({
  email: '',
})

const handleSubmit = async (): Promise<void> => {
  try {
    loading.value = true
    await authStore.updateEmail(formData.value)
    commonStore.addSnackbar({
      message: '認証コードを送信しました。',
      color: 'info',
    })
    router.push({
      name: 'accounts-email-verification',
      params: { email: formData.value.email },
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
</script>

<template>
  <templates-auth-email-edit
    v-model:form-data="formData"
    :loading="loading"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    @submit="handleSubmit"
  />
</template>
