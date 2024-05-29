<script lang="ts" setup>
import { useAlert } from '~/lib/hooks'
import { useAuthStore, useCommonStore } from '~/store'
import type { UpdateAuthPasswordRequest } from '~/types/api'

const router = useRouter()
const commonStore = useCommonStore()
const authStore = useAuthStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const loading = ref<boolean>(false)
const formData = ref<UpdateAuthPasswordRequest>({
  oldPassword: '',
  newPassword: '',
  passwordConfirmation: '',
})

const handleSubmit = async (): Promise<void> => {
  try {
    loading.value = true
    await authStore.updatePassword(formData.value)
    commonStore.addSnackbar({
      message: 'パスワードを更新しました。',
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
    loading.value = false
  }
}
</script>

<template>
  <templates-auth-password-edit
    v-model:form-data="formData"
    :loading="loading"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    @submit="handleSubmit"
  />
</template>
