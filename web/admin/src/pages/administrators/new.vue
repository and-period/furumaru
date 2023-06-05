<script lang="ts" setup>
import { convertJapaneseToI18nPhoneNumber } from '~/lib/formatter'
import { useAlert } from '~/lib/hooks'
import { useAdministratorStore } from '~/store'
import { CreateAdministratorRequest } from '~/types/api'

const router = useRouter()
const administratorStore = useAdministratorStore()
const { alertType, isShow, alertText, show } = useAlert('error')

const formData = reactive<CreateAdministratorRequest>({
  lastname: '',
  firstname: '',
  lastnameKana: '',
  firstnameKana: '',
  email: '',
  phoneNumber: ''
})

const handleSubmit = async (): Promise<void> => {
  try {
    const req: CreateAdministratorRequest = {
      ...formData,
      phoneNumber: convertJapaneseToI18nPhoneNumber(formData.phoneNumber)
    }
    await administratorStore.createAdministrator(req)
    router.push('/administrators')
  } catch (err) {
    if (err instanceof Error) {
      show(err.message)
    }
    console.log(err)
  }
}
</script>

<template>
  <templates-administrator-new
    v-model:form-data="formData"
    :is-alert="isShow"
    :alert-type="alertType"
    :alert-text="alertText"
    @submit="handleSubmit"
  />
</template>
