<script lang="ts" setup>
import { storeToRefs } from 'pinia'
import { convertI18nToJapanesePhoneNumber, convertJapaneseToI18nPhoneNumber } from '~/lib/formatter'
import { useAlert } from '~/lib/hooks'
import { useAdministratorStore, useCommonStore } from '~/store'
import type { UpdateAdministratorRequest } from '~/types/api/v1'

const route = useRoute()
const router = useRouter()
const commonStore = useCommonStore()
const administratorStore = useAdministratorStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const administratorId = route.params.id as string

const { administrator } = storeToRefs(administratorStore)

const loading = ref<boolean>(false)
const formData = ref<UpdateAdministratorRequest>({
  lastname: '',
  firstname: '',
  lastnameKana: '',
  firstnameKana: '',
  phoneNumber: '',
})

const fetchState = useAsyncData('administrator-detail', async (): Promise<void> => {
  try {
    await administratorStore.getAdministrator(administratorId)
    formData.value = {
      ...administrator.value,
      phoneNumber: convertI18nToJapanesePhoneNumber(administrator.value.phoneNumber),
    }
  }
  catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
})

const isLoading = (): boolean => {
  return fetchState?.pending?.value || loading.value
}

const handleSubmit = async (): Promise<void> => {
  try {
    loading.value = true
    const req: UpdateAdministratorRequest = {
      ...formData.value,
      phoneNumber: convertJapaneseToI18nPhoneNumber(formData.value.phoneNumber),
    }
    await administratorStore.updateAdministrator(administratorId, req)
    commonStore.addSnackbar({
      message: '管理者情報の更新が完了しました。',
      color: 'info',
    })
    router.push('/administrators')
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

try {
  await fetchState.execute()
}
catch (err) {
  console.log('failed to setup', err)
}
</script>

<template>
  <templates-administrator-edit
    v-model:form-data="formData"
    :loading="isLoading()"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    :administrator="administrator"
    @submit="handleSubmit"
  />
</template>
