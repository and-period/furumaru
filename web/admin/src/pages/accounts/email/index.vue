<script lang="ts" setup>
import { useAlert } from '~/lib/hooks'
import { useAuthStore } from '~/store'
import { UpdateAuthEmailRequest } from '~/types/api'

const router = useRouter()
const authStore = useAuthStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const formData = reactive<UpdateAuthEmailRequest>({
  email: ''
})

const handleSubmit = async (): Promise<void> => {
  try {
    await authStore.emailUpdate(formData)
    router.push({
      name: 'accounts-email-verification',
      params: { email: formData.email }
    })
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log('failed to update email', err)
  }
}
</script>

<template>
  <templates-auth-edit-email
    v-model:form-data="formData"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    @submit="handleSubmit"
  />
</template>
