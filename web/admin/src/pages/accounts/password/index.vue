<script lang="ts" setup>
import { useAlert } from '~/lib/hooks'
import { useAuthStore } from '~/store'
import { UpdateAuthPasswordRequest } from '~/types/api'

const router = useRouter()
const authStore = useAuthStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const formData = reactive<UpdateAuthPasswordRequest>({
  oldPassword: '',
  newPassword: '',
  passwordConfirmation: ''
})

const handleSubmit = async (): Promise<void> => {
  try {
    await authStore.passwordUpdate(formData)
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
  <templates-auth-edit-password
    v-model:form-data="formData"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    @submit="handleSubmit"
  />
</template>
