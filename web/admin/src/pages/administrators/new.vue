<script lang="ts" setup>
import { convertJapaneseToI18nPhoneNumber } from '~/lib/formatter'
import { useAlert } from '~/lib/hooks'
import { useAdministratorStore, useCommonStore } from '~/store'
import type { CreateAdministratorRequest } from '~/types/api'

const router = useRouter()
const commonStore = useCommonStore()
const administratorStore = useAdministratorStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const loading = ref<boolean>(false)
const formData = ref<CreateAdministratorRequest>({
  lastname: '',
  firstname: '',
  lastnameKana: '',
  firstnameKana: '',
  email: '',
  phoneNumber: ''
})

const handleSubmit = async (): Promise<void> => {
  try {
    loading.value = true
    const req: CreateAdministratorRequest = {
      ...formData.value,
      phoneNumber: convertJapaneseToI18nPhoneNumber(formData.value.phoneNumber)
    }
    await administratorStore.createAdministrator(req)
    commonStore.addSnackbar({
      message: `${formData.value.lastname} ${formData.value.firstname}を作成しました。`,
      color: 'info'
    })
    router.push('/administrators')
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
  <templates-administrator-new
    v-model:form-data="formData"
    :loading="loading"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    @submit="handleSubmit"
  />
</template>
