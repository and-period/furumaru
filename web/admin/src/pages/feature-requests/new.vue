<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useAuthStore, useFeatureRequestStore } from '~/store'
import { FeatureRequestCategory, FeatureRequestPriority } from '~/types/feature-request'
import type { CreateFeatureRequestInput } from '~/types/feature-request'

const authStore = useAuthStore()
const featureRequestStore = useFeatureRequestStore()

const { adminId, user } = storeToRefs(authStore)

const formData = ref<CreateFeatureRequestInput>({
  title: '',
  description: '',
  category: FeatureRequestCategory.Feature,
  priority: FeatureRequestPriority.Medium,
  submittedBy: '',
  submitterName: '',
})

const { isLoading, isShow, alertType, alertText, handleSubmit } = useFormPage({
  submitFn: () => {
    formData.value.submittedBy = adminId.value
    formData.value.submitterName = user.value?.username ?? ''
    return featureRequestStore.createFeatureRequest(formData.value)
  },
  successMessage: '要望リクエストを提出しました。',
  redirectPath: '/feature-requests',
})
</script>

<template>
  <templates-feature-request-new
    v-model:form-data="formData"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    @submit="handleSubmit"
  />
</template>
