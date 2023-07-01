<script lang="ts" setup>
import { useAlert } from '~/lib/hooks'
import { useAuthStore } from '~/store'
import { SignInRequest } from '~/types/api'

definePageMeta({
  layout: 'auth'
})

const router = useRouter()
const authStore = useAuthStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const formData = reactive<SignInRequest>({
  username: '',
  password: ''
})

const handleSubmit = async () => {
  try {
    const path = await authStore.signIn(formData)
    router.push(path)
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}
</script>

<template>
  <templates-auth-sign-in
    v-model:form-data="formData"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    @submit="handleSubmit"
  />
</template>
