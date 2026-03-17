<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { useAlert } from '~/lib/hooks'
import { useAuthStore, useCommonStore, useFeatureRequestStore } from '~/store'
import { FeatureRequestStatus } from '~/types/feature-request'
import type { UpdateFeatureRequestInput } from '~/types/feature-request'

const route = useRoute()
const router = useRouter()
const commonStore = useCommonStore()
const authStore = useAuthStore()
const featureRequestStore = useFeatureRequestStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const id = Array.isArray(route.params.id) ? route.params.id[0] : route.params.id
const { featureRequest } = storeToRefs(featureRequestStore)
const { adminType } = storeToRefs(authStore)

const loading = ref<boolean>(false)
const formData = ref<UpdateFeatureRequestInput>({
  status: FeatureRequestStatus.Waiting,
  note: '',
})

const fetchState = useAsyncData('feature-request-detail', async (): Promise<void> => {
  try {
    await featureRequestStore.getFeatureRequest(id)
    if (featureRequest.value) {
      formData.value = {
        status: featureRequest.value.status,
        note: featureRequest.value.note,
      }
    }
  }
  catch (err) {
    if (err instanceof Error)
      show(err.message)
  }
})

const isLoading = (): boolean => {
  return fetchState?.pending?.value || loading.value
}

const handleSubmit = async (): Promise<void> => {
  try {
    loading.value = true
    await featureRequestStore.updateFeatureRequest(id, formData.value)
    commonStore.addSnackbar({
      message: '要望リクエストを更新しました。',
      color: 'info',
    })
    router.push('/feature-requests')
  }
  catch (err) {
    if (err instanceof Error)
      show(err.message)
    if (import.meta.client)
      window.scrollTo({ top: 0, behavior: 'smooth' })
  }
  finally {
    loading.value = false
  }
}

const handleClickDelete = async (): Promise<void> => {
  try {
    loading.value = true
    await featureRequestStore.deleteFeatureRequest(id)
    commonStore.addSnackbar({
      message: '要望リクエストを削除しました。',
      color: 'info',
    })
    router.push('/feature-requests')
  }
  catch (err) {
    if (err instanceof Error)
      show(err.message)
  }
  finally {
    loading.value = false
  }
}

try {
  await fetchState.execute()
}
catch (err) {
  console.error('failed to setup', err)
  if (err instanceof Error)
    show(err.message)
}
</script>

<template>
  <templates-feature-request-edit
    v-model:form-data="formData"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :admin-type="adminType"
    :feature-request="featureRequest ?? undefined"
    @submit="handleSubmit"
    @click:delete="handleClickDelete"
  />
</template>
